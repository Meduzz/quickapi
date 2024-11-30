package main

import (
	"github.com/Meduzz/quickapi"
	"github.com/Meduzz/quickapi/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Pets []*Pet `json:"pets,omitempty"`
	}

	Pet struct {
		Name  string `json:"name"`
		Alive bool   `json:"alive"`
	}
)

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	start := quickapi.GinStarter(db,
		model.NewJsonEntity[Person]("person"))

	err = start.Execute()

	if err != nil {
		panic(err)
	}
}
