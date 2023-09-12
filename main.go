package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-user/repo"
	"github.com/Coreychen4444/shortvideo_ms-user/service"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	mysql_db := repo.InitMysql()
	redis_db := repo.InitRedis()
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, service.NewUserServer(repo.NewDbRepository(mysql_db, redis_db)))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
