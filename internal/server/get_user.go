package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"rag/internal"
	"strconv"
)

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		s.logger.Error("strconv.Atoi(mux.Vars(r)[])", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := s.userService.GetUser(userID)
	if err != nil {
		s.logger.Error("s.userService.GetUser()", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		response := internal.GetErrResponse(err)
		if err = json.NewEncoder(w).Encode(&response); err != nil {
			s.logger.Error("json.NewEncoder(w).Encode(&response)", zap.Error(err))
		}
		return
	}
	if err = json.NewEncoder(w).Encode(&user); err != nil {
		s.logger.Error("json.NewEncoder(w).Encode(&user)", zap.Error(err))
	}
}
