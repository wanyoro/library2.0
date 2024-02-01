package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"encoding/csv"
	"encoding/json"

	"lib2.0/api/models"
	"lib2.0/api/responses"
)

var book models.Book

func (a *App) exportBooksHandler(w http.ResponseWriter, r *http.Request) {
	//book:= models.Book{}
	format := r.URL.Query().Get("format")
	if format == "" {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("Format parameter is required (csv or json)"))
		return
	}
	//Set the appropriate content type header
	var contentType string
	var fileExtension string

	switch format {
	case "csv":
		contentType = "test/csv"
		fileExtension = "csv"
	case "json":
		contentType = "application/json"
		fileExtension = "json"
	default:
		responses.ERROR(w, http.StatusNotFound, fmt.Errorf("Invalid Format provided"))
		return
	}
	// Set Content-Disposition header for the browser to suggest a filename
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=books.%s", fileExtension))
	w.Header().Set("Content-Type", contentType)

	//Export data based on the selected format
	switch format {
	case "csv":
		a.exportBooksCSV(w)
	case "json":
		a.exportBooksJSON(w)
	}

}
func (a *App) exportBooksCSV(w http.ResponseWriter) {
	// Set Content-Disposition header to suggest a filename
	w.Header().Set("Content-Disposition", "attachment; filename=books.csv")
	// Set the Content-Type header
	w.Header().Set("Content-Type", "text/csv")
	//Create a CSV writer
	books, err := book.GetBookss(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	//Write Header
	header := []string{"ID", "Subject", "ISBN", "IsRead", "StudentID", "TeacherID", "Available", "AvailableDate"}
	csvWriter.Write(header)

	//Write data
	for _, book := range books {
		var isReadstring string
		if book.IsRead!= nil{
			isReadstring = strconv.FormatBool(*book.IsRead)
		}else{
			isReadstring =""
		}
		row := []string{strconv.Itoa(int(book.ID)), book.Subject, strconv.Itoa(int(book.ISBN)), isReadstring , strconv.Itoa(int(book.StudentID)), strconv.Itoa(int(book.TeacherID)), strconv.FormatBool(book.Available), book.AvailableDate.Format(time.RFC3339)}
		csvWriter.Write(row)
	}
	fmt.Println("Books exported as CSV")
}

func (a *App) exportBooksJSON(w http.ResponseWriter) {
	// Encode books slice to JSON
	books, err := book.GetBookss(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	jsonData, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Write JSON data to the response writer
	w.Write(jsonData)
	fmt.Println("Books exported as JSON")
}
