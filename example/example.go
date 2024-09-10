package main

import (
	"strconv"

	"github.com/Meduzz/quickapi"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Person struct {
		ID   int64  `gorm:"autoIncrement" json:"id,omitempty"`
		Name string `gorm:"size:32" json:"name" binding:"required"`
		Age  int    `json:"age" binding:"gt=-1"`
		Pets []*Pet `json:"pets,omitempty" gorm:"constraint:OnDelete:CASCADE"`
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

	err = quickapi.Run(db,
		quickapi.NewEntity[Person]("", quickapi.NewFilter("pets", preloadPets())),
		quickapi.NewEntity[Pet]("pet"))

	if err != nil {
		panic(err)
	}
}

// preloadPets creates a named filter that preload the Pets collection.
// if a filter on alive column was provided, then that is used in the
// preload.
func preloadPets() quickapi.Scope {
	return func(m map[string]string) func(*gorm.DB) *gorm.DB {
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
