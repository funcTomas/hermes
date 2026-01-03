package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/funcTomas/hermes/common"
	"github.com/funcTomas/hermes/service"
	"github.com/funcTomas/hermes/tool"
)

type UserHandler struct {
	factory service.Factory
}

func NewUserHandler(factory service.Factory) UserHandler {
	return UserHandler{factory: factory}
}

func (uh *UserHandler) UserAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	phone := strings.TrimSpace(r.FormValue("phone"))
	uniqId := strings.TrimSpace(r.FormValue("uniqId"))
	channelId, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("channelId")))

	if phone == "" || uniqId == "" || channelId == 0 {
		http.Error(w, "Missing required parameters: phone, uniqId, channelId", http.StatusBadRequest)
		return
	}
	putDate := tool.GetNowDate()
	userService := uh.factory.GetUserService()
	if err := userService.SendMqAddUser(r.Context(), phone, uniqId, channelId, putDate); err != nil {
		http.Error(w, "send rocketmq failed "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data := struct{}{}
	if err := json.NewEncoder(w).Encode(common.SuccessRet(data)); err != nil {
		log.Printf("Error encoding response: %v\n", err)
	}
}

func (uh *UserHandler) UserEnterGroup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	uniqId := strings.TrimSpace(r.FormValue("uniqId"))

	if uniqId == "" {
		http.Error(w, "Missing required parameters: uniqId", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"success": "true", "message": "Request processed successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
	}
}
