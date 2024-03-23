package quickapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Run runs db migrations for T, creates a gin server
// setups the routing for T (crud + list/search) (in the root /)
// and starts the server on 8080.
func Run(db *gorm.DB, entities ...Entity) error {
	// start gin
	engine := gin.Default()
	// create a place to store entities to migrate
	migrations := make([]any, 0)

	// iterate entities and create their api
	for _, entity := range entities {
		migrations = append(migrations, entity.Create())

		if entity.Name() == "" {
			For(db, &engine.RouterGroup, entity)
		} else {
			rg := engine.Group(fmt.Sprintf("/%s", entity.Name()))
			For(db, rg, entity)
		}
	}

	// run migration
	err := db.AutoMigrate(migrations...)

	if err != nil {
		return err
	}

	return engine.Run(":8080")
}

// For sets up routing for T in the provided router group
// but leaves up to you to deal with the server and run migrations.
func For(db *gorm.DB, api *gin.RouterGroup, entity Entity) {
	r := newRouter(entity)

	// setup REST endpoints
	api.POST("/", r.Create(db))      // create
	api.GET("/:id", r.Read(db))      // read
	api.PUT("/:id", r.Update(db))    // update
	api.DELETE("/:id", r.Delete(db)) // delete
	api.GET("/", r.Search(db))       // list/search
	api.PATCH("/:id", r.Patch(db))   // patch
}
