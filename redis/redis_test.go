package redis

import "testing"

func TestGetRedisConn(t *testing.T) {
	InitRedis("127.0.0.1", "3306", "")
	conn, err := GetRedisConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	reply, _ := conn.Do("keys *")
	t.Log("get connection success", reply)
}
