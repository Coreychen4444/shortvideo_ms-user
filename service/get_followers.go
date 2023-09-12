package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
)

func (s *UserServer) GetFollowers(ctx context.Context, req *pb.GetFollowersRequest) (*pb.GetFollowersResponse, error) {
	followers, err := s.r.GetFansList(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("获取粉丝列表失败")
	}
	return &pb.GetFollowersResponse{UserList: followers}, nil
}
