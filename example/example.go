package main

import (
	"github.com/Meduzz/quickapi"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Entity struct {
		ID   int64  `gorm:"autoIncrement" json:"id,omitempty"`
		Name string `gorm:"size:32" json:"name" binding:"required"`
		Age  int    `json:"age" binding:"gt=-1"`
	}
)

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = quickapi.Run[Entity](db)

	if err != nil {
		panic(err)
	}
}
