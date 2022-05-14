package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		s.logger.Error("strconv.Atoi(mux.Vars(r)[])", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.userService.DeleteUser(userID); err != nil {
		s.logger.Error("strconv.Atoi(mux.Vars(r)[])", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("User with userID %d deleted", userID)
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		s.logger.Error("json.NewEncoder(w).Encode(&response)", zap.Error(err))
	}
}
