package quickapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type router struct {
	entity Entity
}

func newRouter(entity Entity) *router {
	return &router{entity}
}

func (r *router) Create(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := r.entity.Create()
		err := ctx.BindJSON(entity)

		if err != nil {
			println("binding body threw error", err.Error())
			ctx.AbortWithStatus(400)
			return
		}

		err = db.Create(entity).Error

		if err != nil {
			// TODO here be dragons
			println("creating row threw error", err.Error())
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(http.StatusCreated, entity)
	}
}

func (r *router) Read(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := r.entity.Create()

		id := ctx.Param("id")
		query := db

		scopes := createScopes(ctx, r.entity.Filters())

		if scopes != nil {
			query = query.Scopes(scopes...)
		}

		err := query.First(entity, id).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatus(404)
				return
			}

			// TODO here be dragons
			println("reading row threw error", err.Error())
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, entity)
	}
}

func (r *router) Update(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := r.entity.Create()
		// id := ctx.Param("id") // not used in logic, only in routing

		err := ctx.BindJSON(entity)

		if err != nil {
			println("binding body threw error", err.Error())
			ctx.AbortWithStatus(400)
			return
		}

		query := db.Session(&gorm.Session{FullSaveAssociations: true})

		scopes := createScopes(ctx, r.entity.Filters())

		if scopes != nil {
			query = query.Scopes(scopes...)
		}

		err = query.Save(entity).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatus(404)
				return
			}

			// TODO here be dragons
			println("updating row threw error", err.Error())
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, entity)
	}
}

func (r *router) Delete(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := r.entity.Create()
		id := ctx.Param("id")

		query := db.Select(clause.Associations)

		scopes := createScopes(ctx, r.entity.Filters())

		if scopes != nil {
			query = query.Scopes(scopes...)
		}

		err := query.Delete(entity, id).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatus(404)
				return
			}

			// TOOD here be dragons
			println("deleting row threw error", err.Error())
			ctx.AbortWithStatus(500)
			return
		}

		ctx.Status(200)
	}
}

func (r *router) Search(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sSkip := ctx.DefaultQuery("skip", "0")
		sTake := ctx.DefaultQuery("take", "25")
		where, ok := ctx.GetQueryMap("where")

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

		data := r.entity.CreateArray()

		query := db.
			Offset(iSkip).
			Limit(iTake)

		if ok {
			query = query.Where(where)
		}

		scopes := createScopes(ctx, r.entity.Filters())

		if scopes != nil {
			query = query.Scopes(scopes...)
		}

		err = query.Find(&data).Error

		if err != nil {
			// TODO here be dragons
			println("searching for data threw error", err.Error())
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, data)
	}
}

func (r *router) Patch(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := r.entity.Create()
		id := ctx.Param("id")

		data := make(map[string]any)
		err := ctx.BindJSON(&data)

		if err != nil {
			println("binding request threw error", err.Error())
			ctx.AbortWithStatus(400)
			return
		}

		err = db.Model(entity).
			Where("id = ?", id).
			Updates(data).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatus(404)
				return
			}

			// TODO here be dragons
			println("patching data threw error", err.Error())
			ctx.AbortWithStatus(500)
			return
		}

		query := db

		scopes := createScopes(ctx, r.entity.Filters())

		if scopes != nil {
			query = query.Scopes(scopes...)
		}

		err = query.Find(entity, id).Error

		if err != nil {
			// TODO here be dragons too
			println("loading updated data threw error", err.Error())
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, entity)
	}
}
