package main

import (
	"api-gateway/common/config"
	"api-gateway/discovery"
	"api-gateway/internal/service"
	"api-gateway/router"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()
	startListen()
}

func startListen() {
	etcdAddress := []string{viper.GetString("etcd.address")}
	user := viper.GetString("etcd.user")
	password := viper.GetString("etcd.password")
	etcdRegister := discovery.NewResolver(etcdAddress, user, password, logrus.New())
	resolver.Register(etcdRegister)
	defer etcdRegister.Close()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	userAddress := fmt.Sprintf("%s:///%s", etcdRegister.Scheme(), viper.GetString("domain.user"))
	userConn, err := grpc.DialContext(ctx, userAddress, opts...)
	if err != nil {
		panic(err)
	}
	services := make(map[string]interface{})
	services["user"] = service.NewUserServiceClient(userConn)
	services["task"] = service.NewTaskServiceClient(userConn)
	ginRouter := router.NewRouter(services)
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")),
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("gateway启动失败, err: ", err)
		}
	}()
	fmt.Printf("gateway listen on %s\n", server.Addr)

	// 优雅关闭
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("closing http server gracefully ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("closing http server gracefully failed: ", err)
	}
}
