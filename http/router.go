package http

import (
	"net/http"
	"strconv"

	"github.com/Meduzz/helper/http/herror"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/quickapi/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	router struct {
		storer storage.Storer
		entity model.Entity
	}
)

func newRouter(db *gorm.DB, entity model.Entity) (*router, error) {
	storer, err := storage.CreateStorage(db, entity)

	if err != nil {
		return nil, err
	}

	return &router{storer, entity}, nil
}

func (r *router) Create(ctx *gin.Context) {
	entity := r.entity.Create()
	err := ctx.BindJSON(entity)

	if err != nil {
		println("binding body threw error", err.Error())
		ctx.AbortWithStatus(400)
		return
	}

	entity, err = r.storer.Create(entity)

	if err != nil {
		// TODO here be dragons
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

	entity, err := r.storer.Read(id, preload)

	if err != nil {
		// TODO here be dragons
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

	entity, err = r.storer.Update(id, entity)

	if err != nil {
		// TODO here be dragons
		println("updating row threw error", err.Error())
		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.JSON(200, entity)
}

func (r *router) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := r.storer.Delete(id)

	if err != nil {
		// TOOD here be dragons
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

	data, err := r.storer.Search(iSkip, iTake, where, sort, preload, hooks...)

	if err != nil {
		// TODO here be dragons
		println("searching for data threw error", err.Error())
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
	entity, err := r.storer.Patch(id, data, preload)

	if err != nil {
		// TODO here be dragons
		println("patching data threw error", err.Error())
		code := herror.CodeFromError(err)

		ctx.AbortWithStatus(code)
		return
	}

	ctx.JSON(200, entity)
}
