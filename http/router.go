package http

import (
	"net/http"
	"strings"

	"github.com/Meduzz/helper/http/herror"
	"github.com/Meduzz/quickapi/api"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/quickapi/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	router struct {
		storage storage.Storage
		entity  model.Entity
		config  *Config
	}
)

func newRouter(db *gorm.DB, config *Config, entity model.Entity) *router {
	store := storage.CreateStorage(db, entity)

	return &router{store, entity, config}
}

func (r *router) Create(ctx *gin.Context) {
	entity, err := r.config.Body(r.entity.Create(), ctx)

	if err != nil {
		println("binding body threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	req := api.NewCreate(entity)
	entity, err = r.storage.Create(req)

	if err != nil {
		println("creating row threw error", err.Error())
		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.JSON(http.StatusCreated, entity)
}

func (r *router) Read(ctx *gin.Context) {
	id := r.config.ID(ctx)
	preload := r.config.Preload(ctx)

	req := api.NewRead(id, preload)
	entity, err := r.storage.Read(req)

	if err != nil {
		println("reading row threw error", err.Error())
		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.JSON(200, entity)
}

func (r *router) Update(ctx *gin.Context) {
	id := r.config.ID(ctx)
	entity, err := r.config.Body(r.entity.Create(), ctx)

	if err != nil {
		println("binding body threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	hooks := CreateHooks(r.entity, ctx)

	req := api.NewUpate(id, entity, hooks)
	entity, err = r.storage.Update(req)

	if err != nil {
		println("updating row threw error", err.Error())
		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.JSON(200, entity)
}

func (r *router) Delete(ctx *gin.Context) {
	id := r.config.ID(ctx)
	hooks := CreateHooks(r.entity, ctx)

	req := api.NewDelete(id, hooks)
	err := r.storage.Delete(req)

	if err != nil {
		println("deleting row threw error", err.Error())
		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.Status(200)
}

func (r *router) Search(ctx *gin.Context) {
	take := r.config.Take(ctx)
	skip := r.config.Skip(ctx)
	where := r.config.Where(ctx)
	sort := r.config.Sorting(ctx)

	preload := r.config.Preload(ctx)
	hooks := CreateHooks(r.entity, ctx)

	req := api.NewSearch(skip, take, where, sort, preload, hooks)
	data, err := r.storage.Search(req)

	if err != nil {
		println("searching for data threw error", err.Error())

		if strings.Contains(err.Error(), "syntax error") || strings.Contains(err.Error(), "no such column") {
			ctx.AbortWithStatus(400)
			return
		}

		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.JSON(200, data)
}

func (r *router) Patch(ctx *gin.Context) {
	id := r.config.ID(ctx)
	data := make(map[string]any)
	err := ctx.BindJSON(&data)

	if err != nil {
		println("binding request threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	preload := r.config.Preload(ctx)
	hooks := CreateHooks(r.entity, ctx)

	req := api.NewPatch(id, data, preload, hooks)
	entity, err := r.storage.Patch(req)

	if err != nil {
		println("patching data threw error", err.Error())

		if strings.Contains(err.Error(), "syntax error") {
			ctx.AbortWithStatus(400)
			return
		}

		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.JSON(200, entity)
}
