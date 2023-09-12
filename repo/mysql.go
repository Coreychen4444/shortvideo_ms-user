package repo

import (
	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-user/model"
	"gorm.io/gorm"
	"sync"
)

// 创建用户
func (r *DbRepository) CreateUsers(user *model.User) (int64, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return -1, err
	}
	return user.ID, nil
}

// 根据用户名获取用户
func (r *DbRepository) GetUserByName(username string) (*pb.User, string, error) {
	var user pb.User
	err := r.db.Model(&model.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, "", err
	}
	var passwordHash string
	err = r.db.Model(&model.User{}).Where("username = ?", username).Pluck("password_hash", &passwordHash).Error
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

// 根据用户id获取用户
func (r *DbRepository) GetUserById(id int64) (*pb.User, error) {
	var user pb.User
	err := r.db.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// 判断是否关注
func (r *DbRepository) IsFollow(authorID, fansID int64) (bool, error) {
	var relation model.Relation
	err := r.db.Where("author_id = ? and fans_id = ?", authorID, fansID).First(&relation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 关注
func (r *DbRepository) CreateFollow(authorID, fansID int64) error {
	relation := model.Relation{
		AuthorID: authorID,
		FansID:   fansID,
	}
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	// 添加关注记录
	if err := r.db.Create(&relation).Error; err != nil {
		txErr = err
		return err
	}
	// 更新作者的粉丝数
	if err := r.db.Model(&model.User{}).Where("id = ?", authorID).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户的关注数
	if err := r.db.Model(&model.User{}).Where("id = ?", fansID).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 取消关注
func (r *DbRepository) DeleteFollow(authorID, fansID int64) error {
	var relation model.Relation
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()

	// 删除关注记录
	if err := r.db.Where("author_id = ? and fans_id = ?", authorID, fansID).Delete(&relation).Error; err != nil {
		txErr = err
		return err
	}
	// 更新作者的粉丝数
	if err := r.db.Model(&model.User{}).Where("id = ?", authorID).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户的关注数
	if err := r.db.Model(&model.User{}).Where("id = ?", fansID).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 获取关注列表
func (r *DbRepository) GetFollowList(userID int64) ([]*pb.User, error) {
	var following_id []int64
	err := r.db.Model(&model.Relation{}).Where("fans_id = ?", userID).Pluck("author_id", &following_id).Error
	if err != nil {
		return nil, err
	}
	var followings []*pb.User
	err = r.db.Model(&model.User{}).Where("id in (?)", following_id).Find(&followings).Error
	if err != nil {
		return nil, err
	}
	return followings, nil
}

// 获取粉丝列表
func (r *DbRepository) GetFansList(userID int64) ([]*pb.User, error) {
	var follower_id []int64
	err := r.db.Model(&model.Relation{}).Where("author_id = ?", userID).Pluck("fans_id", &follower_id).Error
	if err != nil {
		return nil, err
	}
	var followers []*pb.User
	err = r.db.Model(&model.User{}).Where("id in (?)", follower_id).Find(&followers).Error
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// 获取用户好友列表
// 由于涉及多次数据库查询，使用协程并发查询，并且使用redis缓存好友id,减少数据库IO
func (r *DbRepository) GetFriendList(userID int64) ([]*pb.User, error) {
	var following, follower []int64
	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err := r.db.Model(&model.Relation{}).Where("fans_id = ?", userID).Pluck("author_id", &following).Error
		if err != nil {
			errChan <- err
		}
		wg.Done()
	}()
	go func() {
		err := r.db.Model(&model.Relation{}).Where("author_id = ?", userID).Pluck("fans_id", &follower).Error
		if err != nil {
			errChan <- err
		}
		wg.Done()
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	friendID := GetFriID(following, follower)
	friends, err := r.GetUserListByIds(friendID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// 获取好友ID
func GetFriID(following, follower []int64) []int64 {
	var friendID []int64
	friendmap := make(map[int64]bool)
	for _, v := range following {
		friendmap[v] = true
	}
	for _, v := range follower {
		if friendmap[v] {
			friendID = append(friendID, v)
		}
	}
	return friendID
}

// 根据用户id获取用户列表
func (r *DbRepository) GetUserListByIds(ids []int64) ([]*pb.User, error) {
	var users []*pb.User
	err := r.db.Model(&model.User{}).Where("id in (?)", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 获取聊天记录
func (r *DbRepository) GetMessages(user_id, to_user_id int64, pre_msg_time string) ([]*pb.Message, error) {
	var messages []*pb.Message
	err := r.db.Model(&model.Message{}).Where(
		"(create_time > ?) and "+
			"((from_user_id = ? and to_user_id = ?) "+
			"or (from_user_id = ? and to_user_id = ?))",
		pre_msg_time, user_id, to_user_id, to_user_id, user_id).Order("create_time asc").Limit(20).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// 创建消息记录
func (r *DbRepository) CreateMessage(message *model.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		return err
	}
	return nil
}
