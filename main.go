package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"ther.cool/micro-go-course-user/dao"
	"ther.cool/micro-go-course-user/endpoint"
	"ther.cool/micro-go-course-user/redis"
	"ther.cool/micro-go-course-user/service"
	"ther.cool/micro-go-course-user/transport"
)

func main() {
	var servicePort = flag.Int("service.port", 10086, "service port")
	flag.Parse()
	ctx := context.Background()
	errChan := make(chan error)
	err := dao.InitMysql("127.0.0.1", "3306", "root", "root", "user")
	if err != nil {
		log.Fatal(err)
	}
	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		log.Fatal(err)
	}

	userService := service.NewServiceImpl(&dao.UserDAOImpl{})
	userEndpoints := &endpoint.UserEndpoints{
		RegisterEndpoint: endpoint.NewRegisterEndpoint(userService),
		LoginEndpoint:    endpoint.NewLoginEndpoint(userService),
	}
	handler := transport.NewHttpHandler(ctx, userEndpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	serverError := <-errChan
	log.Println(serverError)
}
