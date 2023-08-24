package controllers

import (
	"net/http"

	"lib2.0/api/models"
	"lib2.0/api/responses"
	//"lib2.0/api/models"
)

//Func CreateNotif creates notification after reading progress is set to true
func(a *App) CreateNotif (w http.ResponseWriter, r*http.Request){
// 	var resp = map[string]interface {}{"status":"successful", "message":"notification created successfully"}
// 	notif := models.Notification{} 
} 

//func GetNotifs gets all notifications from table
func (a *App) GetNotifs (w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "Application/json")
	notif := models.Notification{}
	notifs, err := notif.GetNotifs(a.DB)
	if err!= nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, notifs)
}