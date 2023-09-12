package service

import (
	"context"
	"fmt"
	pb "github.com/Coreychen4444/shortvideo"
)

func (s *UserServer) FollowOrCancel(ctx context.Context, req *pb.FollowOrCancelRequest) (*pb.FollowOrCancelResponse, error) {
	author_id := req.GetId()
	token_user_id := req.GetTokenUserId()
	action_type := req.GetActionType()
	if action_type == "1" {
		err := s.r.CreateFollow(author_id, token_user_id)
		if err != nil {
			return nil, fmt.Errorf("关注失败")
		}
	}
	if action_type == "2" {
		err := s.r.DeleteFollow(author_id, token_user_id)
		if err != nil {
			return nil, fmt.Errorf("取消关注失败")
		}
	}
	return &pb.FollowOrCancelResponse{}, nil
}
