package httpadapter

import (
	"encoding/json"
	"hexa-user/internal/domain"
	"net/http"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest CreateUserRequest

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserListResponse struct {
	Users    []*UserResponse `json:"users"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Total    int             `json:"total"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

func toUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        string(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := toUserResponse(user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10

	users, err := h.service.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &UserListResponse{
		Users:    make([]*UserResponse, len(users)),
		Page:     page,
		PageSize: pageSize,
		Total:    len(users),
	}

	for i, user := range users {
		response.Users[i] = toUserResponse(user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
