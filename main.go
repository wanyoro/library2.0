package main

import (
	//"log"
	//"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	//"github.com/joho/godotenv"
	"lib2.0/api/controllers"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	app := controllers.App{}
	app.Initialize()
	app.RunServer()
}
