package dao

import "testing"

func initMysql() error {
	return InitMysql("127.0.0.1", "3306", "root", "root", "user")
}

func TestUserDAOImpl_Save(t *testing.T) {
	userDAO := &UserDAOImpl{}
	err := initMysql()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user := &UserEntity{
		Username: "Ther",
		Password: "ther",
		Email:    "ther@ther.cool",
	}
	err = userDAO.Save(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("new User ID is %d", user.ID)
}

func TestUserDAOImpl_SelectByEmail(t *testing.T) {
	userDAO := &UserDAOImpl{}
	if err := initMysql(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	user, err := userDAO.SelectByEmail("ther@ther.cool")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("result user is %v", user)
}
