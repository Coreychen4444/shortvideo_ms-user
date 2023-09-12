package service

import (
	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-user/repo"
)

type UserServer struct {
	r *repo.DbRepository
	pb.UnimplementedUserServiceServer
}

func NewUserServer(r *repo.DbRepository) *UserServer {
	return &UserServer{r: r}
}
