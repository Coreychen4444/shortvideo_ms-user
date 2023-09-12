package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
)

// 获取用户好友列表
func (s *UserServer) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	friendscache, _ := s.r.GetFriendListByRedis(req.GetId())
	if len(friendscache) > 0 {
		friends, err := s.r.GetUserListByIds(friendscache)
		if err == nil {
			return &pb.GetFriendsResponse{UserList: friends}, nil
		}
	}
	friends, err := s.r.GetFriendList(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("获取好友列表失败")
	}
	_ = s.r.AddFriendList(req.GetId(), friends)
	return &pb.GetFriendsResponse{UserList: friends}, nil
}
