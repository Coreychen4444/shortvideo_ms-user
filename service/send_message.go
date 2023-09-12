package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-user/model"
)

// 发送消息
func (s *UserServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	token_user_id := req.GetTokenUserId()
	to_user_id := req.GetToUserId()
	content := req.GetContent()
	message := model.Message{
		FromUserID: token_user_id,
		ToUserID:   to_user_id,
		Content:    content,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := s.r.CreateMessage(&message)
	if err != nil {
		return nil, fmt.Errorf("发送消息失败")
	}
	return &pb.SendMessageResponse{}, nil
}
