package main

import (
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
		model.NewEntity[Person]("person"),
		model.NewEntity[Pet]("pet"))

	/*
		// dont forget you can provide --prefix and --queue flags here
		start := quickapi.RpcStarter(db,
			model.NewEntity[Person]("person"),
			model.NewEntity[Pet]("pet"))
	*/

	err = start.Execute()

	if err != nil {
		panic(err)
	}
}
