package controllers

import (
	"fmt"
	
	"net/http"
	"strconv"
	"log"
	//"encoding/json"
	

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	
	
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
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections = make(map[*websocket.Conn]bool)
var broadcast = make(chan Notification)

type Notification struct {
	Message string `json:"message"`
	//IsRead bool `json:"isRead"`
}

func(a *App) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	connections[conn] = true

	for {
		select {
		case notification := <-broadcast:
			// Send notification to all connected clients
			for c := range connections {
				err := c.WriteJSON(notification)
				if err != nil {
					log.Println(err)
					delete(connections, c)
					break
				}
			}
		}
	}
}

// Define a Message struct
type Message struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// sendMessage sends a message to the WebSocket connection
// func sendMessage(conn *websocket.Conn, message string) error {
// 	// Create a Message object
// 	msg := Message{
// 		Type:    1, // Set the message type as needed
// 		Content: message,
// 	}

// 	// Marshal the Message object to JSON
// 	jsonMessage, err := json.Marshal(msg)
// 	if err != nil {
// 		return err
// 	}

// 	// Send the JSON message to the WebSocket connection
// 	if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
// 		return err
// 	}

// 	return nil
// }