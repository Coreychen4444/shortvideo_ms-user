package service

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-user/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser 注册用户
func (s *UserServer) Register(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	//校验
	username := req.GetUsername()
	password := req.GetPassword()
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("用户名或密码不能为空,请重新输入")
	}
	if len(password) < 6 {
		return nil, errors.New("密码长度不能小于6位,请重新输入")
	}
	if len(username) > 32 || len(password) > 32 {
		return nil, errors.New("用户名或密码长度不能超过32位,请重新输入")
	}
	user, _, err := s.r.GetUserByName(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查找用户时出错")
	}
	//判断用户名是否存在
	if user != nil {
		return nil, errors.New("用户名已存在,请重新输入")
	}
	var newUser model.User
	//加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("设置的的密码格式有误")
	}
	newUser.Username = username
	newUser.PasswordHash = string(hashedPassword)
	//创建用户
	id, err := s.r.CreateUsers(&newUser)
	if err != nil {
		return nil, fmt.Errorf("创建用户时出错")
	}
	token, tknerr := GenerateToken(id)
	if tknerr != nil {
		return nil, fmt.Errorf("生成token时出错")
	}
	return &pb.UserResponse{Id: id, Token: token}, nil
}
