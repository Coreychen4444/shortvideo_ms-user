package service

import (
	"context"
	"fmt"
	pb "github.com/Coreychen4444/shortvideo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *UserServer) Login(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	//校验输入
	username := req.GetUsername()
	password := req.GetPassword()
	if len(username) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("用户名或密码不能为空,请重新输入")
	}
	if len(username) > 32 || len(password) > 32 {
		return nil, fmt.Errorf("用户名或密码长度不能超过32位,请重新输入")
	}
	//查找用户
	user, passwordHash, err := s.r.GetUserByName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名不存在,请重新输入")
		}
		return nil, fmt.Errorf("查找用户时出错")
	}
	//验证密码
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("密码错误,请重新输入")
		}
		return nil, fmt.Errorf("验证密码时出错")
	}
	token, tknerr := GenerateToken(user.Id)
	if tknerr != nil {
		return nil, fmt.Errorf("生成token时出错")
	}
	return &pb.UserResponse{Id: user.Id, Token: token}, nil
}
