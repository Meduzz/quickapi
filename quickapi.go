package quickapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Run runs db migrations for T, creates a gin server
// setups the routing for T (crud + list/search) (in the root /)
// and starts the server on 8080.
func Run[T any](db *gorm.DB) error {
	// run migration
	mig := db.Migrator()
	err := mig.AutoMigrate(new(T))

	if err != nil {
		return err
	}

	// start gin
	engine := gin.Default()

	For[T](db, &engine.RouterGroup)

	return engine.Run(":8080")
}

// For sets up routing for T in the provided router group
// but leaves up to you to deal with the server and run migrations.
func For[T any](db *gorm.DB, api *gin.RouterGroup) {
	// setup REST endpoints
	api.POST("/", create[T](db))      // create
	api.GET("/:id", read[T](db))      // read
	api.PUT("/:id", update[T](db))    // update
	api.DELETE("/:id", remove[T](db)) // delete
	api.GET("/", search[T](db))       // list/search
	api.PATCH("/:id", patch[T](db))   // patch
}

func create[T any](db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := new(T)
		err := ctx.BindJSON(entity)

		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		err = db.Create(entity).Error

		if err != nil {
			// TODO here be dragons
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(http.StatusCreated, entity)
	}
}

func read[T any](db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := new(T)

		id := ctx.Param("id")

		err := db.First(entity, id).Error

		if err != nil {
			// TODO here be dragons
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, entity)
	}
}

func update[T any](db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := new(T)
		// id := ctx.Param("id") // not used in logic, only in routing

		err := ctx.BindJSON(entity)

		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		err = db.Save(entity).Error

		if err != nil {
			// TODO here be dragons
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, entity)
	}
}

func remove[T any](db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := new(T)
		id := ctx.Param("id")

		err := db.Delete(entity, id).Error

		if err != nil {
			// TOOD here be dragons
			ctx.AbortWithStatus(500)
			return
		}

		ctx.Status(200)
	}
}

func search[T any](db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		sSkip := ctx.DefaultQuery("skip", "0")
		sTake := ctx.DefaultQuery("take", "25")
		// TODO figure out more advanced query capabilities
		where, ok := ctx.GetQueryMap("where")

		iSkip, err := strconv.Atoi(sSkip)

		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		iTake, err := strconv.Atoi(sTake)

		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		data := make([]*T, 0)

		query := db.
			Offset(iSkip).
			Limit(iTake)

		if ok {
			query = query.Where(where)
		}

		err = query.Find(&data).Error

		if err != nil {
			// TODO here be dragons
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, data)
	}
}

func patch[T any](db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		entity := new(T)
		id := ctx.Param("id")

		data := make(map[string]any)
		err := ctx.BindJSON(&data)

		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		err = db.Model(entity).
			Where("id = ?", id).
			Updates(data).Error

		if err != nil {
			// TODO here be dragons
			ctx.AbortWithStatus(500)
			return
		}

		err = db.Find(entity, id).Error

		if err != nil {
			// TODO here be dragons too
			ctx.AbortWithStatus(500)
			return
		}

		ctx.JSON(200, entity)
	}
}
