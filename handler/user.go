package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/funcTomas/hermes/common"
	"github.com/funcTomas/hermes/service"
	"github.com/funcTomas/hermes/tool"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(srv service.UserService) *UserHandler {
	return &UserHandler{service: srv}
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
	if err := uh.service.SendMqAddUser(r.Context(), phone, uniqId, channelId, putDate); err != nil {
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

	idStr := strings.TrimSpace(r.FormValue("id"))
	putDateStr := strings.TrimSpace(r.FormValue("putDate"))
	tsStr := strings.TrimSpace(r.FormValue("timestamp"))
	id, _ := strconv.ParseInt(idStr, 10, 64)
	putDate, _ := strconv.Atoi(putDateStr)
	ts, _ := strconv.ParseInt(tsStr, 10, 64)

	if id == 0 || putDate == 0 || ts == 0 {
		http.Error(w, fmt.Sprintf("Missing required parameters: id %s putDate %s timestamp %s",
			idStr, putDateStr, tsStr), http.StatusBadRequest)
		return
	}
	if err := uh.service.SendMqEnterGroup(r.Context(), id, putDate, ts); err != nil {
		http.Error(w, "send rocketmq failed "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"success": "true", "message": "Request processed successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
	}
}
