package quickapi

import (
	"errors"

	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/helper/nuts"
	"github.com/Meduzz/quickapi/http"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/quickapi/rpc"
	"github.com/Meduzz/quickapi/storage"
	arepece "github.com/Meduzz/rpc"
	"github.com/Meduzz/rpc/encoding"
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
		slice.ForEach(entities, func(entity model.Entity) {
			http.For(db, &engine.RouterGroup, entity)
		})

		return engine.Run(":8080")
	}

	return cmd
}

func RpcStarter(db *gorm.DB, entities ...model.Entity) *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "start"
	cmd.Short = "start a quickapi over rpc"
	cmd.Flags().String("prefix", "", "prefix to use when creating topic")
	cmd.Flags().String("queue", "", "queue to use as queueGroup when subscribint to topic")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		prefix, _ := cmd.Flags().GetString("prefix")
		queue, _ := cmd.Flags().GetString("queue")

		nc, err := nuts.Connect()

		if err != nil {
			return err
		}

		err = Migrate(db, entities...)

		if err != nil {
			return err
		}

		srv := arepece.NewRpc(nc, encoding.Json())

		// iterate entities and create their api
		slice.ForEach(entities, func(entity model.Entity) {
			rpc.For(db, srv, prefix, queue, entity)
		})

		srv.Run()

		return nil
	}

	return cmd
}

func Migrate(db *gorm.DB, entities ...model.Entity) error {
	errorz := slice.Map(entities, func(e model.Entity) error {
		if e.Kind() == model.KindNormal {
			return db.Table(e.Name()).AutoMigrate(e.Create())
		} else {
			return db.Table(e.Name()).AutoMigrate(&storage.JsonTable{})
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
