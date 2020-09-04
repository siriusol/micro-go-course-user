package service

import (
	"context"
	"testing"
	"ther.cool/micro-go-course-user/dao"
	"ther.cool/micro-go-course-user/redis"
)

func initMysql() error {
	return dao.InitMysql("127.0.0.1", "3306", "root", "root", "user")
}

func initRedis() error {
	return redis.InitRedis("127.0.0.1", "6379", "")
}

func TestUserServiceImpl_Login(t *testing.T) {
	if err := initMysql(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := initRedis(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}
	user, err := userService.Login(context.Background(), "ther@ther.cool", "ther")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("[TestUserServiceImpl_Login] user login success:%v", user)
}

func TestUserServiceImpl_Register(t *testing.T) {
	if err := initMysql(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := initRedis(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}
	user, err := userService.Register(context.Background(), &RegisterUserVO{
		Username: "Ther2",
		Password: "ther2",
		Email:    "ther2@ther.cool",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("[TestUserServiceImpl_Register] user register success:%v", user)
}
