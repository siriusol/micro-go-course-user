package redis

import "testing"

func TestGetRedisConn(t *testing.T) {
	// TODO 为什么这里填 3306 不报错
	InitRedis("127.0.0.1", "6379", "")
	conn, err := GetRedisConn()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	reply, _ := conn.Do("PING")
	t.Log("get connection success", reply)
}
