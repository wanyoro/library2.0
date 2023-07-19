package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"lib2.0/api/models"
	"lib2.0/api/responses"
	"lib2.0/utils"
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
		resp["message"] = "student already exists with this email address please login"
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	student.BeforeSave()
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
func (a *App) StudentLogin(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if stu == nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed please signup"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(student.Password, stu.Password)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "login failed check password"
		responses.JSON(w, http.StatusForbidden, resp)
		return

	}
	token, err := utils.EncodeAuthToken(stu.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
	return

}

//func UpdateStudent updates student info and reading progress
func (a *App) UpdateStudent (w http.ResponseWriter, r*http.Request)
var resp = map[string]interface{}{"status":"success", "message":"student updated successfully"}

//params := mux.vars(r)
id, err:= strconv.Atoi(params["id"])
if err!= nil {
	responses.ERROR(w, http.StatusBadRequest, errors.New("invalid id"))
	return
}

body, err := ioutil.ReadAll(r.Body)
if err!= nil{
	responses.ERROR(w, http.StatusBadRequest, err)
	return
}

studentGotten, err := student.GetStudentById(id, a.DB)
if studentGotten == nil{
	resp["message"]= "student not found"
	responses.JSON(w, http.StatusNotfound, resp)
	return
}

newStu := &models.Student{}

err = json.Unmarshal(body, &newStu)
if err != nil{
	responses.ERROR(w, http.StatusBadRequest, err)
	return
}

_, err= newStu.UpdateStudent(id, a.DB)
if err!= nil{
	responses.ERROR(w, http.StatusBadRequest, err)
	return
}

resp["new student"]= newstu
responses.JSON(w, http.StatusOk, resp)
return
