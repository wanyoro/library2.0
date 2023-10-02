package controllers

import (
	"encoding/json"
	//"fmt"

	"io"
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

	body, err := io.ReadAll(r.Body)
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

	body, err := io.ReadAll(r.Body)
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

	body, err := io.ReadAll(r.Body)
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

// func ReturnBook updates values to default
func (a *App) ReturnBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{"Status": "Successful", "Message": "Book returned successfully"}
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	book := models.Book{}

	BookGot, err := book.GetBookById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	ReturnedBook, err := BookGot.ReturnBook(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["returned_book"] = ReturnedBook
	responses.JSON(w, http.StatusOK, resp)
}

func (a *App) AssignedBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books := models.Book{}

	assignedBooks, err := books.AssignedBooks(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response := map[string]*[]models.BookSubjects{"Assigned Books": assignedBooks}
	responses.JSON(w, http.StatusOK, response)
	//responses.JSON(w, http.StatusOK, availableBooks) // TODO: change this later on
}

func (a *App) UnassignedBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	books := models.Book{}
	//ISBN := books.ISBN
	unassignedBooks, err := books.UnassignedBooks(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response := map[string]*[]models.BookSubjects{"UnassignedBooks": unassignedBooks}
	responses.JSON(w, http.StatusOK, response)
}

// func AssignBook issues book to student
func (a *App) AssignBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{"Status": "Success", "Message": "Book Issued Succesfully"}
	var response = map[string]interface{}{"Status": "Failed", "Message": "this book is already assigned to student"}
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	Book := models.Book{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	bookGotten, err := Book.GetBookById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if bookGotten.StudentID != 0 {
		//err= errors.New("this book is already assigned to student")
		// response["message"] = "failed"
		// response["message"] = "this book is already assigned to student"
		responses.JSON(w, http.StatusBadRequest, response)
		return
	}

	err = json.Unmarshal(body, &bookGotten)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	Book.Validate("issuebook")

	issuedBook, err := bookGotten.IssueBook(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "check book"

	}

	resp["book"] = issuedBook
	responses.JSON(w, http.StatusCreated, resp)

}
