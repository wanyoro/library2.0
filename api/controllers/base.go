package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"lib2.0/api/models"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// connect to db
func (a *App) Initialize(DbHost, DbPort, DbUser, DbName string) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", DbHost, DbPort, DbUser, DbName)

	a.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		fmt.Printf("\n Cannot connect to database %s", DbName)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("Connected to database %s ", DbName)
	}
	a.DB.Debug().AutoMigrate(&models.Book{}, &models.Student{}, &models.Teacher{}, &models.Notification{})

	a.Router = mux.NewRouter().StrictSlash(true)
}

func (a *App) RunServer() {
	log.Printf("\nServer starting on port 8001")
	log.Fatal(http.ListenAndServe(":8001", a.Router))
}
