package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"lib2.0/api/models"
	"lib2.0/api/responses"
	"lib2.0/utils"
)

// func TeacherSignUp adds teacher to db
func (a *App) TeacherSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "teacher created successfully"}

	teacher := models.Teacher{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &teacher)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teach, _ := teacher.GetTeacher(a.DB)
	if teach != nil {
		resp["status"] = "unsusccesful"
		resp["message"] = "teacher already exists please login"
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teacher.Prepare()

	err = teacher.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teacherCreated, err := teacher.SaveTeacher(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["teacher"] = teacherCreated
	responses.JSON(w, http.StatusCreated, resp)
	

}

//func TeacherLogIn logins teacher to app
func (a *App) TeacherLogIn(w http.ResponseWriter, r *http.Request){
	var resp = map[string]interface{}{"status":"Success", "message":"teacher logged in successfully"}

	teacher := models.Teacher{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err= json.Unmarshal(body, &teacher)
	if err!= nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teacher.Prepare()

	err= teacher.Validate("login")
	if err != nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teach, err := teacher.GetTeacher(a.DB)
	if err!= nil{
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if teach == nil{
		resp["status"]= "failed"
		resp["message"]= "teacher not found please sign up"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(teacher.Password, teach.Password)
	if err != nil{
		resp["status"]= "failed"
		resp["message"]= "login failed check password"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}
	token, err := utils.EncodeAuthToken(teacher.ID)
	if err!= nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["token"]= token
	responses.JSON(w, http.StatusOK, resp)
	
}
