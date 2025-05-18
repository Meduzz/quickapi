package http

import (
	"net/http"
	"strconv"
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
	}
)

func newRouter(db *gorm.DB, entity model.Entity) (*router, error) {
	store, err := storage.CreateStorage(db, entity)

	if err != nil {
		return nil, err
	}

	return &router{store, entity}, nil
}

func (r *router) Create(ctx *gin.Context) {
	entity := r.entity.Create()
	err := ctx.BindJSON(entity)

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
	id := ctx.Param("id")
	preload := ctx.QueryMap("preload")

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
	id := ctx.Param("id")
	entity := r.entity.Create()

	err := ctx.BindJSON(entity)

	if err != nil {
		println("binding body threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	req := api.NewUpate(id, entity)
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
	id := ctx.Param("id")

	req := api.NewDelete(id)
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
	sSkip := ctx.DefaultQuery("skip", "0")
	sTake := ctx.DefaultQuery("take", "25")
	where, ok := ctx.GetQueryMap("where")

	if !ok {
		where = make(map[string]string)
	}

	sort, ok := ctx.GetQueryMap("sort")

	if !ok {
		sort = make(map[string]string)
	}

	iSkip, err := strconv.Atoi(sSkip)

	if err != nil {
		println("parsing query parameter 'skip' threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	iTake, err := strconv.Atoi(sTake)

	if err != nil {
		println("parsing query parameter 'take' threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	preload := ctx.QueryMap("preload")
	hooks := make([]model.Hook, 0)

	scopeSupport, ok := r.entity.(model.ScopeSupport)

	if ok {
		hooks = createScopes(ctx, scopeSupport.Scopes())
	}

	req := api.NewSearch(iSkip, iTake, where, sort, preload, hooks)

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
	id := ctx.Param("id")
	data := make(map[string]any)

	err := ctx.BindJSON(&data)

	if err != nil {
		println("binding request threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	preload := ctx.QueryMap("preload")

	req := api.NewPatch(id, data, preload)
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
