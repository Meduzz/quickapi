package quickapi

import (
	"errors"
	"fmt"

	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/quickapi/http"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/quickapi/storage"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func GinStarter(db *gorm.DB, entities ...model.Entity) *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "start"
	cmd.Short = "start a quickapi over gin"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		// start gin
		engine := gin.Default()

		err := Migrate(db, entities...)

		if err != nil {
			return err
		}

		// iterate entities and create their api
		err = http.For(db, &engine.RouterGroup, entities...)

		if err != nil {
			return err
		}

		return engine.Run(":8080")
	}

	return cmd
}

func Migrate(db *gorm.DB, entities ...model.Entity) error {
	errorz := slice.Map(entities, func(e model.Entity) error {
		if e.Kind() == model.NormalKind {
			return db.Table(e.Name()).AutoMigrate(e.Create())
		} else if e.Kind() == model.JsonKind {
			return db.Table(e.Name()).AutoMigrate(&storage.JsonTable{})
		} else {
			return fmt.Errorf("unknown entity kind: %s", e.Kind())
		}
	})

	return slice.Fold(errorz, nil, func(err, agg error) error {
		if err != nil {
			if agg != nil {
				return errors.Join(agg, err)
			}

			return err
		}

		return nil
	})
}
