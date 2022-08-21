package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"task/common/config"
	"task/discovery"
	"task/internal/handler"
	"task/internal/repository"
	"task/internal/service"
)

func main() {
	config.InitConfig()
	repository.InitDB()

	etcdAddresses := []string{viper.GetString("etcd.address")}
	user := viper.GetString("etcd.user")
	password := viper.GetString("etcd.password")
	etcdRegister := discovery.NewRegister(etcdAddresses, user, password, logrus.New())
	grpcAddress := viper.GetString("server.grpcAddress")
	userNode := discovery.Server{
		Address: grpcAddress,
		Name:    viper.GetString("server.domain"),
	}
	server := grpc.NewServer()
	defer server.Stop()
	service.RegisterTaskServiceServer(server, handler.NewTaskService())
	listen, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}

	if _, err = etcdRegister.Register(userNode, 10); err != nil {
		panic(err)
	}
	log.Printf("server listening at %v", listen.Addr())
	if err = server.Serve(listen); err != nil {
		panic(err)
	}
}
