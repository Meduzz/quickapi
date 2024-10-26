package http

import (
	"fmt"

	"github.com/Meduzz/quickapi/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// For sets up routing for T in the provided router group
// but leaves up to you to deal with the server and run migrations.
func For(db *gorm.DB, e *gin.RouterGroup, entity model.Entity) {
	r := newRouter(db, entity)

	if entity.Name() != "" {
		api := e.Group(fmt.Sprintf("/%s", entity.Name()))

		// setup REST endpoints
		api.POST("/", r.Create)      // create
		api.GET("/:id", r.Read)      // read
		api.PUT("/:id", r.Update)    // update
		api.DELETE("/:id", r.Delete) // delete
		api.GET("/", r.Search)       // list/search
		api.PATCH("/:id", r.Patch)   // patch
	} else {
		// setup REST endpoints
		e.POST("/", r.Create)      // create
		e.GET("/:id", r.Read)      // read
		e.PUT("/:id", r.Update)    // update
		e.DELETE("/:id", r.Delete) // delete
		e.GET("/", r.Search)       // list/search
		e.PATCH("/:id", r.Patch)   // patch
	}
}

func createScopes(ctx *gin.Context, filters []*model.NamedFilter) []model.Hook {
	if len(filters) == 0 {
		return nil
	}

	scopes := []model.Hook{}

	for _, filter := range filters {
		data, ok := ctx.GetQueryMap(filter.Name)

		if ok {
			scopes = append(scopes, filter.Scope(data))
		}
	}

	return scopes
}
