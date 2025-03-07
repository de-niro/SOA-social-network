package main

import (
	"context"
	"database/sql"
	api "users/users_api"
)

// UsersAPIServer Corresponds to the StrictServerInterface in gen.go
type UsersAPIServer struct {
	db *sql.DB
}

func NewUsersAPIServer(db *sql.DB) UsersAPIServer {
	return UsersAPIServer{db: db}
}

// GetUsersUsername Fetch user from database
func (s *UsersAPIServer) GetUsersUsername(ctx context.Context, request api.GetUsersUsernameRequestObject) (api.GetUsersUsernameResponseObject, error) {
	// We assume that lib/pq is thread-safe for now
	info, err := getUserInfoByUsername(s.db, request.Username)
	if err != nil {
		return api.GetUsersUsername404Response{}, nil
	}

	verified, err := isUserEmailVerified(s.db, info.id)
	if err != nil {
		return api.GetUsersUsername404Response{}, nil
	}

	status := AccountStatusString(info.account_status)

	return api.GetUsersUsername200JSONResponse{Username: &info.username, Id: &info.id, AccountStatus: &status, Bio: &info.bio, Birthday: &info.bday, Email: &info.email, EmailVerified: &verified, FullName: &info.full_name, Phone: &info.phone, RegistrationString: &info.registration, UpdateString: &info.user_update}, nil
}
