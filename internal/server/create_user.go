package server

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"rag/internal"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.logger.Error("json.NewDecoder(r.Body).Decode(&request)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := s.userService.CreateUser(request)
	if err != nil {
		s.logger.Error("json.NewDecoder(r.Body).Decode(&request)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		response := internal.GetErrResponse(err)
		if err = json.NewEncoder(w).Encode(&response); err != nil {
			s.logger.Error("json.NewEncoder(w).Encode(&response)", zap.Error(err))
		}
		return
	}
	response := map[string]int{
		"user_id": *userID,
	}
	if err = json.NewEncoder(w).Encode(&response); err != nil {
		s.logger.Error("json.NewEncoder(w).Encode(&response)", zap.Error(err))
	}
}
