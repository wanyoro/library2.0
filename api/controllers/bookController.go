package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"lib2.0/api/models"
	"lib2.0/api/responses"
)

//func CreateBook adds a book to the database
func (a *App)CreateBook (w http.ResponseWriter, r* http.Request){
	w.Header().Set("Content-type", "Application/json")
	var resp = map[string]interface{}{"status":"success", "message":"book created successfully"}

	book:= models.Book{}

	body, err := ioutil.ReadAll(r.Body)
	if err!= nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err= json.Unmarshal(body, &book)
	if err!= nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	dbbook,_ := book.GetBook(a.DB)
	if dbbook != nil {
		resp["status"]= "failed"
		resp["message"]= "book number already exists"
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} 

	book.Prepare()

	book.Validate("createbook")
	if err!= nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	bookCreated, err:= book.CreateBook(a.DB)
	if err!= nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp ["book"]= bookCreated
	responses.JSON(w, http.StatusCreated, resp)
}
