package main

import (
	"github.com/Meduzz/quickapi"
	"github.com/Meduzz/quickapi/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Person struct {
		FullName string `json:"name"`
		Age      int    `json:"age"`
		Pets     []*Pet `json:"pets,omitempty"`
	}

	Pet struct {
		FullName string `json:"name"`
		Alive    bool   `json:"alive"`
	}
)

var (
	_ model.Entity = Person{}
)

// Implementing the Entity interface for Person
func (p Person) Name() string {
	return "persons"
}

func (p Person) Create() any {
	return &Person{}
}

func (p Person) CreateArray() any {
	return []*Person{}
}

func (p Person) Kind() model.EntityKind {
	return model.JsonKind
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	start := quickapi.GinStarter(db, Person{})

	err = start.Execute()

	if err != nil {
		panic(err)
	}
}
