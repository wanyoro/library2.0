package controllers

import (
	
	"encoding/json"
	

	//"errors"

	"fmt"
	"io"
	"net/http"

	//"strconv"
	"strings"

	"github.com/gorilla/mux"
	"lib2.0/api/models"
	"lib2.0/api/responses"
	"lib2.0/utils"
)

// func TeacherSignUp adds teacher to db
func (a *App) TeacherSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "teacher created successfully"}

	teacher := models.Teacher{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &teacher)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teach, _ := teacher.GetTeacherByUsername(teacher.Username, a.DB)
	if teach != nil {
		
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("%v already exists", teacher.Username))
		return
	}

	teacher.Prepare()
	teacher.BeforeSave()


	teacherCreated, err := teacher.SaveTeacher(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["teacher"] = teacherCreated
	responses.JSON(w, http.StatusCreated, resp)

}

// func TeacherLogIn logins teacher to app
func (a *App) TeacherLogIn(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "Success", "message": "teacher logged in successfully"}

	teacher := models.Teacher{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &teacher)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teacher.Prepare()

	err = teacher.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	teach, err := teacher.GetTeacher(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if teach == nil {
		resp["status"] = "failed"
		resp["message"] = "teacher not found please sign up"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(teacher.Password, teach.Password)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "login failed check password"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}
	token, err := utils.EncodeAuthTokenTeacher(teacher.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)

}

func (a *App) GetTeachers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{
		"status": "success",
	}
	teacher := models.Teacher{}
	teachers, _ := teacher.FindAllTeachers(a.DB)
	resp["teachers"] = teachers
	responses.JSON(w, http.StatusOK, resp)

}

func (a *App) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	var resp = map[string]interface{}{
		"status": "success",
	}
	params := mux.Vars(r)
	username := (params["username"])
	account := models.Teacher{}
	teacherGot, _ := account.GetTeacherByUsername(username, a.DB)
	if teacherGot == nil {
		resp["status"] = "failed"
		resp["message"] = fmt.Sprintf("%v does not exist", username)
		return
	}
	//teacher := models.Teacher{}
	teachers, _ := teacherGot.RemoveTeacher(username, a.DB)
	resp["teacher successfully removed"] = strings.Split(teachers, "REMOVED")
	responses.JSON(w, http.StatusOK, resp)
}

//func ChangePasswd changes password for user
// func (a *App) ChangePasswd(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-Type", "Application/json")
// 	var resp = map[string]interface{}{}
// 	params := mux.Vars(r)
// 	email := (params["email"])
// 	account := models.Teacher{}
// 	teacherGot, _ := account.GetTeacherByUsername(email, a.DB)
// 	if teacherGot == nil {
// 		resp["status"] = "failed"
// 		resp["message"] = fmt.Sprintf("%v does not exist", email)
// 		return
// 	}
	
// 	var teacher models.Teacher
// 	var err error
// 	err = json.NewDecoder(r.Body).Decode(&teacher)
// 	if err != nil {
// 		resp["status"] = "failed"
// 		resp["message"] = err.Error()
// 		return
// 	}
// 	newPassword,err := teacherGot.ChangePasswd(email, a.DB)
// 	if err!= nil{
// 		resp["status"] = "failed"
// 		resp["message"] = err.Error()
// 		return
	
// 	}


	
// }

//RequestPasswordReset generates a reset token and sends an email
