package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/go-redis/redis"
)

var ctx = context.Background()

// 在redis 中获取好友列表id
func (r *DbRepository) GetFriendListByRedis(userId int64) ([]int64, error) {
	cacheKey := fmt.Sprintf("user:%d:friends_list", userId)
	ids, err := r.rdb.SMembers(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(ids) > 0 {
		var idList []int64
		for _, id := range ids {
			userId, err := strconv.ParseInt(id, 10, 64)
			if err == nil {
				idList = append(idList, userId)
			}
		}
		return idList, nil
	}
	return []int64{}, nil
}

// 将好友列表id存入redis
func (r *DbRepository) AddFriendList(userId int64, friends []*pb.User) error {
	cacheKey := fmt.Sprintf("user:%d:friends_list", userId)
	var ids []interface{}
	for _, friend := range friends {
		ids = append(ids, friend.Id)
	}
	_, err := r.rdb.SAdd(ctx, cacheKey, ids...).Result()
	if err != nil {
		return err
	}
	r.rdb.Expire(ctx, cacheKey, time.Hour)
	return nil
}
