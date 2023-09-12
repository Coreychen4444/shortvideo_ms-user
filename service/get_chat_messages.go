package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/Coreychen4444/shortvideo"
)

func (s *UserServer) GetChatMessages(ctx context.Context, req *pb.GetChatMessagesRequest) (*pb.GetChatMessagesResponse, error) {
	token_user_id := req.GetTokenUserId()
	to_user_id := req.GetToUserId()
	pre_msg_time := req.GetPreMsgTime()
	preMsgTime, err := strconv.ParseInt(pre_msg_time, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("时间戳格式错误")
	}
	pre_msg_time = time.Unix(preMsgTime, 0).Format("2006-01-02 15:04:05")
	messages, err := s.r.GetMessages(token_user_id, to_user_id, pre_msg_time)
	if err != nil {
		return nil, fmt.Errorf("获取聊天记录失败")
	}
	return &pb.GetChatMessagesResponse{ChatMessageList: messages}, nil
}
