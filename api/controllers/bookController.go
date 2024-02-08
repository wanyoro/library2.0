package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"log"
	"time"

	//"time"

	//"fmt"

	"io"
	"net/http"
	"strconv"

	"github.com/jung-kurt/gofpdf"

	//"gorm.io/gorm"

	"errors"

	"github.com/gorilla/mux"
	"lib2.0/api/models"
	"lib2.0/api/responses"
	//"lib2.0/redis"
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
	//bookisbn, err := book.GetBookById()

	err = json.Unmarshal(body, &book)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	book.Prepare()

	book.Validate("createbook")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	book.AvailableDate = time.Now()
	book.Available = true
	bookCreated, err := book.CreateBook(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["book"] = bookCreated
	responses.JSON(w, http.StatusCreated, resp)
}

// func GetBooks get all books from database
func (a *App) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	//cachedBooks, err := redis.RedisClient.Get(context.Background(), )

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
	id, err := strconv.Atoi(params["isbn"])
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
	resp["book"], _ = UpdatedBook.GetBookById(id, a.DB)
	responses.JSON(w, http.StatusCreated, resp)
}

// func UpdateReadingProgress sets reading progress to true or false
func (a *App) UpdateReadingProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{"status": "Successful", "message": "progress updated successfully"}
	params := mux.Vars(r)
	isbn, err := strconv.Atoi(params["isbn"])
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

	BookGot, err := book.GetBookById(isbn, a.DB)
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
	UpdatedBook, err := BookGot.UpdateReadingProgress(isbn, a.DB)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "check book"

	} else {
		resp["status"] = "Success"
		resp["message"] = "Book updated successfully"
	}
	resp["book"], _ = UpdatedBook.GetBookById(isbn, a.DB)
	responses.JSON(w, http.StatusCreated, resp)
}

// func ReturnBook updates values to default
func (a *App) ReturnBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{"Status": "Successful", "Message": "Book returned successfully"}
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["isbn"])

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
	resp["returned_book"], _ = ReturnedBook.GetBookById(id, a.DB)
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

}

func (a *App) UnassignedBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	books := models.Book{}
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
	id, _ := strconv.Atoi(params["isbn"])
	Book := models.Book{}
	//var db *gorm.DB
	//overdueDays := models.OverdueDays{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	err = json.Unmarshal(body, &Book)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	fmt.Printf("Student id:%v", Book.StudentID)
	bookGotten, err := Book.GetBookById(id, a.DB)

	fmt.Println(bookGotten.StudentID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	//fmt.Println(bookGotten.ISBN)
	if bookGotten.StudentID != 0 {
		responses.JSON(w, http.StatusBadRequest, response)
		return
	}
	overdue, err := Book.GetOverDueDaysPerStudent(Book.StudentID, a.DB)

	if overdue.OverdueDays != 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("Student has overdue book of isbn %v and subject %s. Return overdue book first", overdue.ISBN, overdue.Subject))
		fmt.Println(err)
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

	resp["book"], _ = issuedBook.GetBookById(id, a.DB)
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

func (a *App) DeleteBook(w http.ResponseWriter, r *http.Request) {
	// var resp = map[string]interface{}{
	// 	"status":"Successful",
	// }
	w.Header().Set("Content-Type", "Application/json")
	var params = mux.Vars(r)
	isbn, err := strconv.Atoi(params["isbn"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("Invalid ISBN"))
		return
	}

	books := models.Book{}
	getBook, err := books.GetBookById(isbn, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("No Book Found With This ID %v ", isbn))
		return
	}
	if !getBook.Available {
		responses.ERROR(w, http.StatusConflict, fmt.Errorf("Book of isbn %v is issued please return book first", getBook.ISBN))
		return
	}

	deletedbook := getBook.DeleteBookByISBN(isbn, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}
	responses.JSON(w, http.StatusOK, fmt.Sprintf("Book of isbn %v, successfully deleted", deletedbook.ISBN))
}

func (a *App) GetBoooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	//cachedBooks, err := redis.RedisClient.Get(context.Background(), )

	booksGotten := models.Book{}
	books, err := booksGotten.GetBookss(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, books)
}

// func GetDefaulters gets defaulters from DB
func (a *App) GetBookDefaulters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	//defaulters:= models.BookDefaulters{}
	book := models.Book{}

	def, err := book.GetDefaultedBooks(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return

	}
	responses.JSON(w, http.StatusOK, def)

}

// func GetOverdueDays on books
func (a *App) GetOverdueDays(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	//var days models.OverdueDays
	var book models.Book
	books, err := book.GetOverdueDays(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, books)
}

func (a *App) RateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn, err := strconv.Atoi(vars["isbn"])
	if err != nil {
		return
	}

	var ratingData struct {
		Score float64 `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ratingData); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var book models.Book
	result := a.DB.First(&book, "isbn=?", isbn)
	if result.Error != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Book not found"))
		return
	}
	//newRating := models.Rating{Score: ratingData.Score}
	book.AvgRating = (book.AvgRating + ratingData.Score) / 2
	a.DB.Save(&book)

	//averageRating := book.CalculateAverageRating()
	//.Rating = averageRating
	//a.DB.Save(&book)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book rated successfully.")
}

func (a *App) ExportBooksToPDF(w http.ResponseWriter, r *http.Request) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)
	//fetch books
	var books []models.Book
	a.DB.Find(&books)
	for _, book := range books {
		bookInfo := fmt.Sprintf("Subject: %s\nISBN: %d\n\n", book.Subject, book.ISBN)
		pdf.MultiCell(0, 10, bookInfo, "", "0", false)
	}
	w.Header().Set("Content-Disposition", "attachement; filename=books_export.pdf")
	w.Header().Set("Content-Type", "application/pdf")

	err := pdf.Output(w)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	fmt.Println("PDF exported successfully")
}

func (a *App) CompleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// studentID, err := strconv.Atoi(vars["ID"])
	// if err!= nil{
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	isbn, err := strconv.Atoi(vars["isbn"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	book := models.Book{}
	result, err := book.GetBookById(isbn, a.DB)
	if result.ISBN == 0 {
		log.Printf("book of isbn %s does not exist", book.ISBN)
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Println("Request Body:", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// var student models.Student
	// if result :=a.DB.First(&student, studentID);result.Error!= nil{
	// 	responses.ERROR(w, http.StatusNotFound, err)
	// 	return
	// }
	var completionStatus struct {
		IsRead bool `json:"isRead"`
	}
	if err := json.NewDecoder(r.Body).Decode(&completionStatus); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	book.IsRead = &completionStatus.IsRead
	a.DB.Save(&book)

	if *book.IsRead {
		var student models.Student
		if result := a.DB.First(&student, "id=?", book.StudentID); result.Error != nil {
			responses.ERROR(w, http.StatusNotFound, err)
			return
		}
		//var student models.Student
		completedBook := models.CompletedBook{
			BookID:      book.ID,
			StudentID:   student.ID,
			CompletedAt: time.Now(),
		}
		student.CompletedBooks = append(student.CompletedBooks, completedBook)
		a.DB.Save(&student)
	}
	//responses.JSON(w, http.StatusOK, fmt.Fprintf())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book marked as completed with ID %d", book.ID)

}
