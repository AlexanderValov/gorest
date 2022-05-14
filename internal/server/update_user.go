package server

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"rag/internal"
)

type UpdateUserRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	IsActive bool   `json:"is_active"`
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var request UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.logger.Error("json.NewDecoder(r.Body).Decode(&request)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := s.userService.UpdateUser(request)
	if err != nil {
		s.logger.Error("json.NewDecoder(r.Body).Decode(&request)", zap.Error(err))
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
	return
}
