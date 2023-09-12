package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
	"gorm.io/gorm"
)

func (s *UserServer) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	id := req.GetId()
	token_user_id := req.GetTokenUserId()
	user, err := s.r.GetUserById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("该用户不存在")
		}
		return nil, fmt.Errorf("查找用户时出错")
	}
	// 判断是否关注该用户
	if id == token_user_id {
		return &pb.UserInfoResponse{User: user}, nil
	}
	isFollow, err := s.r.IsFollow(id, token_user_id)
	if err != nil {
		return nil, fmt.Errorf("查找用户时出错")
	}
	user.IsFollow = isFollow
	return &pb.UserInfoResponse{User: user}, nil
}
