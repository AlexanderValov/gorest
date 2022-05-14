package server

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"rag/internal"
)

func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.userService.GetAllUsers()
	if err != nil {
		s.logger.Error("s.userService.GetAllUsers()", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, internal.ErrUsernameExist) {
			response := internal.GetErrResponse(err)
			if err = json.NewEncoder(w).Encode(&response); err != nil {
				s.logger.Error("json.NewEncoder(w).Encode(&response)", zap.Error(err))
			}
		}
		return
	}
	if err := json.NewEncoder(w).Encode(&users); err != nil {
		s.logger.Error("json.NewEncoder(w).Encode(&users)", zap.Error(err))
	}
}
