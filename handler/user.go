package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/funcTomas/hermes/common"
	"github.com/funcTomas/hermes/service"
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
	channelId := strings.TrimSpace(r.FormValue("channelId"))

	if phone == "" || uniqId == "" || channelId == "" {
		http.Error(w, "Missing required parameters: phone, uniqId, channelId", http.StatusBadRequest)
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
