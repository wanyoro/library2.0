package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lib2.0/api/middleware"
	"lib2.0/api/models"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// connect to db
func (a *App) Initialize() {
	var err error
	const DNS = "postgres://postgres:nakaasana@localhost/lib?sslmode=disable"

	a.DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("Connected to database lib")
	}
	a.DB.Debug().AutoMigrate(&models.Student{}, &models.Book{}, &models.Teacher{}, &models.Notification{}, &models.StudentAndBooks{}, &models.BookDefaulters{}, &models.OverdueDays{}, &models.CompletedBook{})

	a.Router = mux.NewRouter().StrictSlash(true)
	a.InitializeRoutes()
}

// Starts server connection
func (a *App) RunServer() {
	log.Printf("\nServer starting on port 8001")
	log.Fatal(http.ListenAndServe(":8001", a.Router))
}

// func InitializeRoutes sets http routes
func (a *App) InitializeRoutes() {
	//a.Router.HandleFunc("/")
	a.Router.Use(middleware.SetContentTypeMiddleware)
	s := a.Router.PathPrefix("/api").Subrouter()
	s.Use(middleware.AuthJwtVerify)
	a.Router.HandleFunc("/teachersignup", a.TeacherSignUp).Methods("POST")
	a.Router.HandleFunc("/studentsignup", a.StudentSignUp).Methods("POST")
	a.Router.HandleFunc("/teacherlogin", a.TeacherLogIn).Methods("POST")
	a.Router.HandleFunc("/studentlogin", a.StudentLogin).Methods("POST")
	a.Router.HandleFunc("/createbook", a.CreateBook).Methods("POST")
	a.Router.HandleFunc("/getbook/{id}", a.GetBookById).Methods("GET")
	s.HandleFunc("/getbooks", a.GetBooks).Methods("GET")
	a.Router.HandleFunc("/updatestudent/{id}", a.UpdateStudent).Methods("PUT")
	s.HandleFunc("/getstudents", a.GetStudents).Methods("GET")
	a.Router.HandleFunc("/getstudentsandbooks", a.GetStudentsAndBooks).Methods("GET")
	s.HandleFunc("/updatebook/{isbn}", a.UpdateBook).Methods("PUT")
	a.Router.HandleFunc("/getstudentbookcount", a.GetStudentBookCount).Methods("GET")
	a.Router.HandleFunc("/createnotification", a.CreateNotif).Methods("POST")
	a.Router.HandleFunc("/updatereadingprogress/{isbn}", a.UpdateReadingProgress).Methods("PUT")
	a.Router.HandleFunc("/getstudent/{id}", a.GetStudentById).Methods("GET")
	a.Router.HandleFunc("/getnotifs", a.GetNotifs).Methods("GET")
	s.HandleFunc("/returnbook/{isbn}", a.ReturnBook).Methods("PUT")
	a.Router.HandleFunc("/booksread/{student_id}", a.BooksReadByStudent).Methods("GET")
	a.Router.HandleFunc("/assignedbooks", a.AssignedBooks).Methods("GET")
	a.Router.HandleFunc("/unassignedbooks", a.UnassignedBooks).Methods("GET")
	s.HandleFunc("/assignbook/{isbn}", a.AssignBook).Methods("PUT")
	s.HandleFunc("/getteachers", a.GetTeachers).Methods("GET")
	a.Router.HandleFunc("/deleteteacher/{username}", a.DeleteTeacher).Methods("DELETE")
	a.Router.HandleFunc("/deletebook/{isbn}", a.DeleteBook).Methods("DELETE")
	a.Router.HandleFunc("/deletenotifications", a.DeleteAllNotifications).Methods("DELETE")
	a.Router.HandleFunc("/deletestudent/{id}", a.DeleteStudent).Methods("DELETE")
	a.Router.HandleFunc("/resetpassword/{email}", a.ResetPassword).Methods("PUT")
	a.Router.HandleFunc("/populatebooks/{student_id}", a.PopulateBooks).Methods("GET")
	//s.HandleFunc("/forgotpassword/{email}", a.ForgotPassword).Methods("POST")
	//s.HandleFunc("/resetpassword/:resetToken", a.ResetPassword).Methods("PATCH")
	a.Router.HandleFunc("/getbookss", a.GetBoooks).Methods("GET")
	a.Router.HandleFunc("/getdefaulters", a.GetBookDefaulters).Methods("GET")
	s.HandleFunc("/getoverduedays", a.GetOverdueDays).Methods("GET")
	a.Router.HandleFunc("/exportbooks", a.exportBooksHandler).Methods("GET")
	a.Router.HandleFunc("/ratebook/{isbn}", a.RateBook).Methods("PUT")
	a.Router.HandleFunc("/exportBooksPDF", a.ExportBooksToPDF).Methods("GET")
	//a.Router.HandleFunc("/completeBook/{studentID}/{bookID}",a.CompleteBook).Methods("POST")
	a.Router.HandleFunc("/completedBook/{isbn}", a.CompleteBook).Methods("POST")
	//a.Router.HandleFunc("/averagerating/{isbn}", a.GetAverageRating).Methods("GET")
	a.Router.HandleFunc("/ws", a.handleWebSocket)
	a.Router.HandleFunc("/teacherassignbook/{isbn}", a.AssignTeacherBook).Methods("PUT")
}
