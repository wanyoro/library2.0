package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"lib2.0/api/models"
	"lib2.0/api/responses"
	"lib2.0/utils"

	"github.com/gorilla/mux"
)

// StudentSignUp creates new student
func (a *App) StudentSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "student registered successfully"}

	student := models.Student{}
	body, err := io.ReadAll(r.Body)
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

}

// func Login signs in user
func (a *App) StudentLogin(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "logged in successfully"}

	student := &models.Student{}
	body, err := io.ReadAll(r.Body)
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
	token, err := utils.EncodeAuthTokenStudent(stu.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)

}

// func UpdateStudent updates student info and reading progress
func (a *App) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "student updated successfully"}
	w.Header().Set("Content-Type", "Aplication/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	student := models.Student{}
	studentGotten, err := student.GetStudentById(id, a.DB)
	if err != nil {
		resp["message"] = "student not found"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	//newStu := &models.Student{}

	err = json.Unmarshal(body, &studentGotten)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	student.BeforeSave()
	updatedStudent, err := studentGotten.UpdateStudent(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	token, err := utils.EncodeAuthTokenStudent(updatedStudent.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)

}

// func GetStudents gets all students from db
func (a *App) GetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	//allStudents := models.Student{}

	students, err := models.GetStudents(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	responses.JSON(w, http.StatusOK, students)

}

// func GetStudentsAndBooks get students with assigned books
func (a *App) GetStudentsAndBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	snb, err := models.GetStudentsAndBooks(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)

	}
	responses.JSON(w, http.StatusOK, snb)

}

// func GetStudentBookCount count number of assigned books
func (a *App) GetStudentBookCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	bookcounted := models.Student{}
	bookcount, err := bookcounted.CountBooks(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	responses.JSON(w, http.StatusOK, bookcount)
}

// func GetStudentById get specific student id
func (a *App) GetStudentById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	studentGot := models.Student{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.JSON(w, http.StatusBadRequest, err)
	}

	studentbyid, err := studentGot.GetStudentById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, studentbyid)
}

func (a *App) PopulateBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["student_id"])
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}
	books := models.Student{}
	loadedbooks, err := books.PopulateBooks(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusOK, loadedbooks)

	//books:=models.PopulateBooks()
}

func (a *App) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	students := models.Student{}
	//books := models.Book{}
	student, _ := students.PopulateBooks(id, a.DB)
	if student.BooksAssigned != 0 {
		responses.ERROR(w, http.StatusBadRequest, fmt.Errorf("%v owns %v books . Return book first to delete", student.Username, student.BooksAssigned))
		return
	}

	studentDeleted := student.DeleteStudent(id, a.DB)
	responses.JSON(w, http.StatusOK, fmt.Sprintf("Student of username %v and ID %v deleted successfully", studentDeleted.Username, studentDeleted.ID))
}

func (a *App) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	params := mux.Vars(r)
	email := params["email"]
	user := models.Student{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responses.ERROR(w, http.StatusNotAcceptable, errors.New("Error while decoding the request body"))
		return
	}

	student := user.GetStudentByEmail(email, a.DB)
	resetedUser := student.ResetPassword(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to reset password"))
		return
	}
	responses.JSON(w, http.StatusCreated, resetedUser)

}

