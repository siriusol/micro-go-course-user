package service

import (
	"context"
	"errors"
	"log"
	"ther.cool/micro-go-course-user/dao"
)

type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserVO struct {
	Username string
	Password string
	Email    string
}

var (
	ErrorUserExisted = errors.New("user is existed")
	ErrorPassword    = errors.New("email and password are not match")
	ErrorRegistering = errors.New("email is registering")
)

type UserService interface {
	// login interface
	Login(ctx context.Context, email, password string) (*UserInfoDTO, error)

	// register interface
	Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error)
}

type UserServiceImpl struct {
	userDAO dao.UserDAO
}

func NewServiceImpl(userDao dao.UserDAO) UserService {
	return &UserServiceImpl{
		userDAO: userDao,
	}
}

func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error) {
	user, err := userService.userDAO.SelectByEmail(email)
	if err != nil {
		log.Printf("[Login] select by email error:%s", err)
		return nil, err
	}
	if user.Password == password {
		return &UserInfoDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}, nil
	}
}

func (userService UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error) {

}
