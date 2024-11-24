package main

import (
	"strconv"

	"github.com/Meduzz/quickapi"
	"github.com/Meduzz/quickapi/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Person struct {
		ID   int64  `gorm:"autoIncrement" json:"id,omitempty"`
		Name string `gorm:"size:32" json:"name" binding:"required"`
		Age  int    `json:"age" binding:"gt=-1"`
		Pets []*Pet `json:"pets,omitempty"` // gorm:"constraint:OnDelete:CASCADE" works in PG but not sqlite.
	}

	Pet struct {
		ID       int64  `gorm:"autoIncrement" json:"id,omitempty"`
		Name     string `gorm:"size:32" json:"name" binding:"required"`
		PersonID int64  `json:"-"`
		Alive    bool   `json:"alive"`
	}
)

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	start := quickapi.GinStarter(db,
		model.NewEntity[Person]("person", personPreload, model.NewFilter("asdf", preloadPets())),
		model.NewEntity[Pet]("pet", nil))

	/*
		// dont forget you should provide --prefix and --queue flags here
		start := quickapi.RpcStarter(db,
			model.NewEntity[Person]("person", personPreload),
			model.NewEntity[Pet]("pet", nil))
	*/

	err = start.Execute()

	if err != nil {
		panic(err)
	}
}

var preload = map[string]map[string]*model.PreloadConfig{
	"status": {
		"Pets": {
			Condition: "alive = ?",
			Converter: func(s string) any {
				it, _ := strconv.ParseBool(s)
				return it
			},
		},
	},
	"naming": {
		"Pets": {
			Condition: "name = ?",
			Converter: nil,
		},
	},
}

func personPreload(name string) map[string]*model.PreloadConfig {
	return preload[name]
}

func preloadPets() model.Scope {
	return func(m map[string]string) model.Hook {
		return func(d *gorm.DB) *gorm.DB {
			alive, ok := m["alive"]

			if ok {
				isAlive, err := strconv.ParseBool(alive)

				if err != nil {
					println("parseBool threw error", err.Error())
					isAlive = false
				}

				return d.Preload("Pets", "alive = ?", isAlive)
			}

			return d.Preload("Pets")
		}
	}
}
