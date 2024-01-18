package models

import (
	"encoding/json"
	"errors"
	//"fmt"
	//"log"

	//"strconv"
	"strings"

	"github.com/gorilla/mux"

	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	//"lib2.0/api/controllers"
	//"lib2.0/api/main"
)

type Book struct {
	ID      uint   `gorm:"primarykey"`
	Subject string `gorm:"size:50;not null" json:"subject"`
	//
	ISBN      uint  `gorm:"size:20;unique" json:"isbn"`
	IsRead    *bool `gorm:"size:50"          json:"isread"`
	StudentID uint  `json :"StudentID"` //`gorm:"foreignKey:StudentID"`
	TeacherID uint  `json: "TeacherID"`
	Available bool  `json: "available"`
	//AvailableDate time.Time `json:"availableDate"`
}

type JsonResponse struct {
	Data   []Book `json:"data"`
	Source string `json:"source"`
}
type BookHistory struct {
	StudentID  uint       `json:"userId"`
	BookID     uint       `json:"bookId"`
	IssueDate  *time.Time `json:"issueDate"`
	ReturnDate *time.Time `json:"returnDate"`
	Returned   bool       `json:"returned"`
}

type BookSubjects struct {
	Subject string `gorm:references: "subject"`
	ISBN    uint   `gorm:references: "isbn"`
}

// Prepare strips off white spaces
func (b *Book) Prepare() {
	b.Subject = strings.TrimSpace(b.Subject)

}

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Validate Book input
func (b *Book) Validate(action string) error {
	switch strings.ToLower("action") {
	case "createbook":
		if b.Subject == "" {
			return errors.New("please input subject")
		}
		return nil

		//isbn:= parseint(b.ISBN)
		//if b.ISBN ==strings.Contains(isbn,  )

	case "updatebook":
		if b.Subject == "" {
			return errors.New("please input subject")
		}
		return nil

	case "issuebook":
		if b.StudentID == 0 {
			return errors.New("please input student id")
		}
		return nil

	}
	return nil

}

// func CreateBook creates book
func (b *Book) CreateBook(db *gorm.DB) (*Book, error) {

	err := db.Debug().Create(&b).Error
	if err != nil {
		return &Book{}, err
	}
	return b, nil
}

//func

// func GetBook gets book from database
func (b *Book) GetBookById(isbn int, db *gorm.DB) (*Book, error) {
	book := &Book{}
	if err := db.Debug().Table("books").Where("isbn= ?", isbn).First(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

// func PopulateBooks adds books assigned to student
func (b *Book) PopulateBooks(studentID int, db *gorm.DB) (*Student, error) {
	students := Student{}
	if err := db.Debug().Raw("select students.id,students.created_at,students.updated_at,students.deleted_at,students.username,students.phone_number,students.email,students.password, count(*) as books FROM students inner join books on books.student_id= students.id WHERE students.id= ? AND students.deleted_at IS NULL group by students.id", studentID).Scan(&students).Error; err != nil {
		return nil, err
	}
	return &students, nil
}

func (b *Book) GetBookss(db *gorm.DB) ([]Book, error) {

	books := []Book{}
	if err := db.Debug().Table("books").Find(&books).Error; err != nil {
		return []Book{}, err
	}
	return books, nil
}

// func GetBooks gets all books from database
func (b *Book) GetBooks(db *gorm.DB) (*JsonResponse, error) {
	//var a *App
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cachedProducts, err := redisClient.Get("book").Bytes()
	if err == nil {
        // Data found in Redis
        books := []Book{}
        if err := json.Unmarshal(cachedProducts, &books); err != nil {
            return nil, err
        }
        return &JsonResponse{Data: books, Source: "Redis Cache"}, nil
    }

	// Data not found in Redis, fetch from the database
    books, err := b.GetBookss(db)
    if err != nil {
        return nil, err
    }

    // Fetch from the database and cache the result in Redis
    dbProducts := db.Find(&books)
    if dbProducts.Error != nil {
        return nil, dbProducts.Error
    }

    cachedProducts, err = json.Marshal(books)
    if err != nil {
        return nil, err
    }

    // Set cache in Redis
    err = redisClient.Set("book", string(cachedProducts), 20*time.Second).Err()
    if err != nil {
        return nil, err
    }

    return &JsonResponse{Data: books, Source: "Postgres"}, nil
}

// func UpdateBook updates book subject and student assigned
func (b *Book) UpdateBook(id int, db *gorm.DB) (*Book, error) {
	BookSubject := b.Subject
	StudentId := b.StudentID
	BookId := b.ID

	notif := Notification{}
	if err := db.Debug().Table("books").Where("isbn =?", b.ISBN).Updates(Book{
		//Subject: b.Subject,

		//StudentID: b.StudentID,
		IsRead: b.IsRead,
		//TeacherID: b.TeacherID,
		Subject: b.Subject,
		ISBN:    b.ISBN}).Error; err != nil {
		return &Book{}, err
	}
	if *b.IsRead {
		if db.Debug().Table("notifications").Where("book_subject= ? AND student_id=?", BookSubject, StudentId).First(&notif).RowsAffected == 0 {
			db.Create(&Notification{
				Message:     "Book is read",
				StudentID:   b.StudentID,
				BookSubject: BookSubject,
				BookId:      BookId,
				TeacherID:   b.TeacherID,
			})
		}
	}

	return b, nil
}

// func UpdateReadingProgress sets reading progress to true and creates notif
func (b *Book) UpdateReadingProgress(id int, db *gorm.DB) (*Book, error) {
	if err := db.Debug().Table("books").Where("id =?", b.ID).Updates(Book{
		IsRead: b.IsRead}).Error; err != nil {
		return &Book{}, err
	}
	return b, nil
}

// func return book returns book from student
func (b *Book) ReturnBook(id int, db *gorm.DB) (*Book, error) {

	if err := db.Debug().Table("books").Where("isbn=?", b.ISBN).Updates(map[string]interface{}{"is_read": false, "student_id": 0, "available": true, "available_date": time.Now()}).Error; err != nil {
		return &Book{}, err
	}

	return b, nil

}

// func AssignedBooks shows unavailable books
func (b *Book) AssignedBooks(db *gorm.DB) (*[]BookSubjects, error) {
	AssignedBooks := &[]BookSubjects{}
	if err := db.Debug().Table("books").Select("subject, isbn").Where("student_id>?", 0).Find(AssignedBooks).Error; err != nil {
		return nil, err
	}
	return AssignedBooks, nil
}

// func UnassignesBooks shows available books
func (b *Book) UnassignedBooks(db *gorm.DB) (*[]BookSubjects, error) {
	UnassignedBooks := &[]BookSubjects{}
	if err := db.Debug().Table("books").Select("subject, isbn").Where("student_id<?", 1).Find(UnassignedBooks).Error; err != nil {
		return nil, err
	}
	return UnassignedBooks, nil

}

// func IssueBook assigns book to student
func (b *Book) IssueBook(id int, db *gorm.DB) (*Book, error) {
	currentTime := time.Now()
	returnDate := currentTime.AddDate(0, 0, 15)
	if err := db.Debug().Table("books").Where("isbn=?", b.ISBN).Updates(map[string]interface{}{
		"is_read":        "false",
		"student_id":     b.StudentID,
		"teacher_id":     b.TeacherID,
		"available":      "false",
		"available_date": returnDate}).Error; err != nil {
		return &Book{}, err
	}
	return b, nil
}

// func DeleteBook removes books from db
func (b *Book) DeleteBookByISBN(isbn int, db *gorm.DB) *Book {
	if err := db.Debug().Where("isbn=?", isbn).Delete(&Book{}).Error; err != nil {
		return &Book{} //fmt.Sprintf("Could not delete Book with id %v", isbn)
	}
	return b //fmt.Sprintf("Deleted Book with id %v", isbn)
}
