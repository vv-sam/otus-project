package handlers

import (
	"encoding/json"
	"net/http"
)

type userService interface {
	Login(username, password string) bool
	GenerateToken(username string) (string, error)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type Auth struct {
	us userService
}

func NewAuth(us userService) *Auth {
	return &Auth{us: us}
}

// @Summary Login
// @Description Login and get a token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Login request"
// @Success 200 {object} loginResponse
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /api/auth/login [post]
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest loginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !a.us.Login(loginRequest.Username, loginRequest.Password) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := a.us.GenerateToken(loginRequest.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := loginResponse{
		Token: token,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
