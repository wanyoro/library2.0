package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"lib2.0/api/models"
	"lib2.0/api/responses"
)

// func CreateBook adds a book to the database
func (a *App) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Application/json")
	var resp = map[string]interface{}{"status": "success", "message": "book created successfully"}

	book := models.Book{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &book)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// params := mux.Vars(r)
	// id, err := strconv.Atoi(params["id"])
	// dbbook, _ := book.GetBookById(id, a.DB)
	// if dbbook != nil {
	// 	resp["status"] = "failed"
	// 	resp["message"] = "book number already exists"
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	book.Prepare()

	book.Validate("createbook")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	bookCreated, err := book.CreateBook(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["book"] = bookCreated
	responses.JSON(w, http.StatusCreated, resp)
}

// func GetBookById gets book with specific id
func (a *App) GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	book := models.Book{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	bookGotten, err := book.GetBookById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, bookGotten)

}

// func GetBooks get all books from database
func (a *App) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	booksGotten := models.Book{}
	books, err := booksGotten.GetBooks(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, books)
}

// func UpdateBook sets reading progress to true
func (a *App) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{"status": "successful", "message": "book updated successfully"}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	book := models.Book{}
	BookGotten, err := book.GetBookById(id, a.DB)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "book does not exist in database"
		return
	}

	err = json.Unmarshal(body, &BookGotten)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	book.Validate("updatereadingprogress")

	UpdatedBook, err := BookGotten.UpdateBook(id, a.DB)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "check book"

	} else {
		resp["status"] = "Success"
		resp["message"] = "Book updated successfully"
	}
	resp["book"] = UpdatedBook
	responses.JSON(w, http.StatusCreated, resp)
}

// func UpdateReadingProgress sets reading progress to true or false
func (a *App) UpdateReadingProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{"status": "Successful", "message": "progress updated successfully"}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	book := models.Book{}
	BookGot, err := book.GetBookById(id, a.DB)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "book does not exist in database"
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &BookGot)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	UpdatedBook, err := BookGot.UpdateReadingProgress(id, a.DB)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "check book"

	} else {
		resp["status"] = "Success"
		resp["message"] = "Book updated successfully"
	}
	resp["book"] = UpdatedBook
	responses.JSON(w, http.StatusCreated, resp)
}
