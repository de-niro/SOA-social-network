package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	api "users/users_api"
)

// UsersAPIServer Corresponds to the StrictServerInterface in gen.go
type UsersAPIServer struct {
	db *sql.DB
}

func NewUsersAPIServer(db *sql.DB) UsersAPIServer {
	return UsersAPIServer{db: db}
}

func execLoopIgnition(db *sql.DB, addr string) {
	server := NewUsersAPIServer(db)
	r := gin.Default()

	sh := api.NewStrictHandler(server, nil)
	api.RegisterHandlers(r, sh)
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Print(s.ListenAndServe())
}

// GetUsersUsername Fetch user from database
func (UsersAPIServer) GetUsersUsername(ctx context.Context, request api.GetUsersUsernameRequestObject) (api.GetUsersUsernameResponseObject, error) {
	// We assume that lib/pq is thread-safe for now
	db := ctx.Value("db").(*sql.DB)
	info, err := getUserInfoByUsername(db, request.Username)
	if err != nil {
		return api.GetUsersUsername404Response{}, nil
	}

	verified, err := isUserEmailVerified(db, info.id)
	if err != nil {
		return api.GetUsersUsername404Response{}, nil
	}

	status := AccountStatusString(info.account_status)
	bday := dateFormatShort(info.bday)
	regDate := dateFormatShort(info.registration)
	updDate := dateFormatShort(info.user_update)

	return api.GetUsersUsername200JSONResponse{Username: &info.username, Id: &info.id, AccountStatus: &status, Bio: &info.bio, Birthday: &bday, Email: &info.email, EmailVerified: &verified, FullName: &info.full_name, Phone: &info.phone, RegistrationString: &regDate, UpdateString: &updDate}, nil
}

func (UsersAPIServer) PostLogin(ctx context.Context, request api.PostLoginRequestObject) (api.PostLoginResponseObject, error) {
	db := ctx.Value("db").(*sql.DB)
	id, err := getUserIDByCredentials(db, UserCredentials{phone: *request.Body.Phone, email: *request.Body.Email, username: *request.Body.Username})
	if err != nil {
		errText := "NO_SUCH_USER"
		return api.PostLogindefaultJSONResponse{Body: api.Error{Text: &errText}}, nil
	}

	cred, err := getUserCredentials(db, id)
	if err != nil {
		errText := "NO_SUCH_USER"
		return api.PostLogindefaultJSONResponse{Body: api.Error{Text: &errText}}, nil
	}

	if comparePasswords(cred.passwd_hash, request.Body.Passwd) {
		return api.PostLogin201JSONResponse{}, nil
	}

	errText := "PASSWORD_INVALID"
	return api.PostLogindefaultJSONResponse{Body: api.Error{Text: &errText}}, nil
}

func (UsersAPIServer) PostRegister(ctx context.Context, request api.PostRegisterRequestObject) (api.PostRegisterResponseObject, error) {
	db := ctx.Value("db").(*sql.DB)
	if checkIfUserExists(db, UserCredentials{phone: request.Body.Phone, email: request.Body.Email, username: request.Body.Username}) {
		errText := "USER_ALREADY_EXISTS"
		return api.PostRegisterdefaultJSONResponse{Body: api.Error{Text: &errText}}, nil
	}

	hash, err := passwd2hash(request.Body.Passwd)
	if err != nil {
		return nil, err
	}

	birthday, err := parseDateShort(*request.Body.Birthday)
	if err != nil {
		return nil, err
	}

	u := UserInstance{id: 0, registration: time.Now(), user_update: time.Now(), bday: birthday, username: request.Body.Username, full_name: *request.Body.FullName, email: request.Body.Email, phone: request.Body.Phone, bio: "", passwd_hash: string(hash), account_status: AccountStatusActive}
	err = createUser(db, &u)
	if err != nil {
		return nil, err
	}
	return api.PostRegister201Response{}, nil
}

func (UsersAPIServer) PostEditInfo(ctx context.Context, request api.PostEditInfoRequestObject) (api.PostEditInfoResponseObject, error) {
	db := ctx.Value("db").(*sql.DB)
	info, err := getUserInfoByID(db, *request.Body.Id)
	if err != nil {
		text := "NO_SUCH_USER"
		return api.PostEditInfodefaultJSONResponse{Body: api.Error{Text: &text}}, nil
	}

	// Crappy struct fields update
	if *request.Body.Birthday != "" {
		info.bday, err = parseDateShort(*request.Body.Birthday)
		if err != nil {
			return nil, err
		}
	}

	if *request.Body.FullName != "" {
		info.full_name = *request.Body.FullName
	}

	if *request.Body.Bio != "" {
		info.bio = *request.Body.Bio
	}

	if *request.Body.AccountStatus != "" {
		info.account_status, err = AccountStatusEnum(*request.Body.AccountStatus)
		if err != nil {
			badAccS := "BAD_ACCOUNT_STATUS"
			return api.PostEditInfodefaultJSONResponse{Body: api.Error{Text: &badAccS}}, nil
		}
	}

	err = editUserInfo(db, &info)
	if err != nil {
		return nil, err
	}
	return api.PostEditInfo201Response{}, nil
}

func (UsersAPIServer) PostEditCredentials(ctx context.Context, request api.PostEditCredentialsRequestObject) (api.PostEditCredentialsResponseObject, error) {
	return nil, errors.New("NOT_IMPLEMENTED")
}

func (UsersAPIServer) GetEmailVerify(ctx context.Context, request api.GetEmailVerifyRequestObject) (api.GetEmailVerifyResponseObject, error) {
	return nil, errors.New("NOT_IMPLEMENTED")
}

func (UsersAPIServer) PostEmailVerify(ctx context.Context, request api.PostEmailVerifyRequestObject) (api.PostEmailVerifyResponseObject, error) {
	return nil, errors.New("NOT_IMPLEMENTED")
}
