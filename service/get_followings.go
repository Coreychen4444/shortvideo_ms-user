package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
)

func (s *UserServer) GetFollowings(ctx context.Context, req *pb.GetFollowingsRequest) (*pb.GetFollowingsResponse, error) {

	followings, err := s.r.GetFollowList(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("获取关注列表失败")
	}
	return &pb.GetFollowingsResponse{UserList: followings}, nil
}
