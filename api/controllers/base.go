package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lib2.0/api/models"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// connect to db
func (a *App) Initialize() {
	var err error
	const DNS = "postgres://postgres@localhost/lib?sslmode=disable"

	a.DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("Connected to database lib")
	}
	a.DB.Debug().AutoMigrate(&models.Student{}, &models.Teacher{}, &models.Book{}, &models.Notification{})

	a.Router = mux.NewRouter().StrictSlash(true)
}

func (a *App) RunServer() {
	log.Printf("\nServer starting on port 8001")
	log.Fatal(http.ListenAndServe(":8001", a.Router))
}
