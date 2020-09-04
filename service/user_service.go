package service

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"ther.cool/micro-go-course-user/dao"
	"ther.cool/micro-go-course-user/redis"
	"time"
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
		return nil, err
	}
	if user.Password == password {
		return &UserInfoDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}, nil
	} else {
		return nil, ErrorPassword
	}
}

func (userService UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error) {
	lock := redis.GetRedisLock(vo.Email, time.Duration(5)*time.Second)
	if err := lock.Lock(); err != nil {
		log.Printf("[Register] get redis lock error:%s", err)
		return nil, ErrorRegistering
	}
	defer lock.Unlock() // TODO 写在这里合适？

	existUser, err := userService.userDAO.SelectByEmail(vo.Email)
	if err == nil && existUser != nil {
		return nil, ErrorUserExisted
	}
	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		// TODO 前一项什么时候发生？
		newUser := &dao.UserEntity{
			Username: vo.Username,
			Password: vo.Password,
			Email:    vo.Email,
		}
		// TODO err 陷阱，注意 err 已经为 ErrRecordNotFound 的情况
		err = userService.userDAO.Save(newUser)
		if err == nil {
			return &UserInfoDTO{
				ID:       newUser.ID,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	return nil, err
}
