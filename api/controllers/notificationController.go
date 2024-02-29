package controllers

import (
	"fmt"
	
	"net/http"
	"strconv"
	

	"github.com/gorilla/mux"
	
	
	"lib2.0/api/models"
	"lib2.0/api/responses"
	//"lib2.0/api/models"
)

// Func CreateNotif creates notification after reading progress is set to true
func (a *App) CreateNotif(w http.ResponseWriter, r *http.Request) {
	// 	var resp = map[string]interface {}{"status":"successful", "message":"notification created successfully"}
	// 	notif := models.Notification{}
}

// func GetNotifs gets all notifications from table
func (a *App) GetNotifs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	notif := models.Notification{}
	notifs, err := notif.GetNotifs(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, notifs)
}

// func BooksREadByStudent
func (a *App) BooksReadByStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")

	params := mux.Vars(r)
	notif := models.Notification{}

	student_id, err := strconv.Atoi(params["student_id"])
	if err != nil {
		responses.JSON(w, http.StatusUnprocessableEntity, err)
	}

	notifs, err := notif.GetReadBooks(student_id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, notifs)

}

func (a *App) DeleteAllNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	notif := models.Notification{}
	err := notif.DeleteAllNotifications(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Couldn't delete the notifications"))
		return
	}
	responses.JSON(w, http.StatusNoContent, fmt.Sprint("Deleted Successfully!"))
}


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Maintain a list of connected clients
	clients := make(map[*websocket.Conn]bool)

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}

		// Process the received message
		// Check if it's related to a book being marked as read
		isRead, ok := msg["isRead"].(bool)
		if ok && isRead {
			// Broadcast a notification to all connected clients (teachers)
			for client := range clients {
				err := client.WriteJSON(map[string]interface{}{"message": "Book marked as read"})
				if err != nil {
					log.Println(err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
