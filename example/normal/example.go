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
		ID       int64  `gorm:"autoIncrement" json:"id,omitempty"`
		FullName string `gorm:"size:32" json:"name" binding:"required"`
		Age      int    `json:"age" binding:"gt=-1"`
		Pets     []*Pet `json:"pets,omitempty"` // gorm:"constraint:OnDelete:CASCADE" works in PG but not sqlite.
	}

	Pet struct {
		ID       int64  `gorm:"autoIncrement" json:"id,omitempty"`
		FullName string `gorm:"size:32" json:"name" binding:"required"`
		PersonID int64  `json:"-"`
		Alive    bool   `json:"alive"`
	}
)

var (
	_ model.Entity         = Person{}
	_ model.Entity         = Pet{}
	_ model.PreloadSupport = Person{}
)

func (p Person) Kind() model.EntityKind {
	return model.NormalKind
}

func (p Person) Name() string {
	return "persons"
}

func (p Person) Create() any {
	return &Person{}
}

func (p Person) CreateArray() any {
	return make([]*Person, 0)
}

func (p Person) Preload(key string) map[string]*model.PreloadConfig {
	it, ok := preload[key]

	if !ok {
		return nil
	}

	return it
}

func (p Pet) Kind() model.EntityKind {
	return model.NormalKind
}

func (p Pet) Name() string {
	return "pets"
}

func (p Pet) Create() any {
	return &Pet{}
}

func (p Pet) CreateArray() any {
	return make([]*Pet, 0)
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	start := quickapi.GinStarter(db, Person{}, Pet{})

	// model.NewEntity[Person]("person", personPreload, model.NewFilter("asdf", preloadPets())),

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
			Condition: "full_name = ?",
			Converter: nil,
		},
	},
}

/*
func personPreload(name string) map[string]*model.PreloadConfig {
	return preload[name]
}
*/

/*
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
*/
