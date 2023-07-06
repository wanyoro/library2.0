package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"lib2.0/api/models"
	"lib2.0/api/responses"
)

// StudentSignUp creates new student
func (a *App) StudentSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "student registered successfully"}

	student := models.Student{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &student)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	stu, _ := student.GetStudent(a.DB)
	if stu != nil {
		resp["status"] = "failure"
		resp["message"] = "student already exists with this email address. Login"
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	student.Prepare()

	err = student.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	studentCreated, err := student.SaveStudent(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["student"] = studentCreated
	responses.JSON(w, http.StatusCreated, resp)
	return

}

// func Login signs in user
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "logged in successfully"}

	student := &models.Student{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &student)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	student.Prepare()

	err = student.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	stu, err := student.GetStudent(a.DB)
	if 

}
