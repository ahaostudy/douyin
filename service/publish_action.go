package service

import (
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"path"
	"strconv"
	"strings"
	"time"
)

// InsertVideo 保存上传的视频数据到数据库
func InsertVideo(id uint, fileName string, title string) error {
	coverName := strings.Split(fileName, ".")[0] + ".jpg"
	video := model.Video{
		AuthorID:  id,
		Title:     title,
		PlayUrl:   path.Join("play", strconv.Itoa(int(id)), fileName),
		CoverUrl:  path.Join("cover", strconv.Itoa(int(id)), coverName),
		CreatedAt: time.Now(),
	}
	err := dao.InsertVideo(&video)
	return err
}

// PublishAction 发布视频操作
func PublishAction(video *model.Video) error {
	// 往数据库插入一条记录
	err := dao.InsertVideo(video)
	if err != nil {
		return err
	}

	// 更新redis，维护用户信息
	go func() {
		ctx, cancel := redis.WithTimeoutContextBySecond(2)
		defer cancel()

		if !ExistsUserInfo(ctx, video.AuthorID) {
			return
		}
		key := redis.GenerateUserKey(video.AuthorID)
		if redis.RdbUser.HIncrBy(ctx, key, "work_count", 1).Err() != nil {
			redis.RdbUser.Del(ctx, key)
		}
	}()

	return nil
}
