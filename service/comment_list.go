package service

import (
	"context"
	"encoding/json"
	"fmt"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"sort"
	"sync"
)

// GetCommentList 获取评论列表
func GetCommentList(vid uint, uid uint) ([]*model.Comment, bool) {
	fmt.Println("enter comment list")
	var commentList []*model.Comment
	var err error
	// 先从redis中加载相关的视频评论信息
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()
	key := redis.GenerateCommentKey(vid)

	// 1. 加载评论列表到redis
	// 加载到redis失败时，直接返回从MySQL中查询的数据
	if commentList, err = LoadCommentList(ctx, vid, uid); err != nil {
		redis.RdbComment.Del(ctx, key)
		return commentList, true
	}

	commentHash, err := redis.RdbComment.HGetAll(ctx, key).Result()

	if err != nil {
		// 直接从MySQL返回
		redis.RdbComment.Del(ctx, key)
		commentList, err = dao.GetCommentList(vid, uid)
		if err != nil {
			return nil, false
		}
		return commentList, true
	}

	if len(commentHash) != 0 {
		// 定义一个等待组，用于等待所有协程完成
		var wg sync.WaitGroup
		// 定义一个互斥锁，用于保护共享数据的并发访问
		var mu sync.Mutex
		// 定义一个通道，用于接收协程的结果或错误
		ch := make(chan *model.Comment, len(commentHash))

		for _, v := range commentHash {
			wg.Add(1) // 增加等待组的计数器

			go func(commentData string) {
				defer wg.Done() // 减少等待组的计数器

				var commentStruct model.Comment
				if err := json.Unmarshal([]byte(commentData), &commentStruct); err != nil {
					// 将错误信息发送到通道
					ch <- nil
					return
				}

				curUser, ok := GetUserByID(commentStruct.UserID, uid)
				if !ok {
					// 将错误信息发送到通道
					ch <- nil
					return
				}

				commentStruct.User = *curUser

				// 由于并发操作，需要保证对 commentList 的安全访问
				// 使用互斥锁进行保护
				mu.Lock()
				commentList = append(commentList, &commentStruct)
				mu.Unlock()

				// 将评论结果发送到通道
				// ch <- &commentStruct
			}(v)
		}

		// 等待所有协程完成
		wg.Wait()

		// 关闭通道，表示不再发送数据
		close(ch)

		// 处理通道接收结果，如果中间有一步转换是错误的，直接从MySQL中返回
		for c := range ch {
			// 如果接收到 nil，则表示协程处理出错
			if c == nil {
				// 处理错误的逻辑
				redis.RdbComment.Del(ctx, key)
				commentList, err = dao.GetCommentList(vid, uid)
				if err != nil {
					return nil, false
				}
				return commentList, true
			}
		}
		// 处理正常的评论结果
		// 处理并发 HSet过程中导致的创建日期乱序情况
		// 排序评论列表
		sort.Slice(commentList, func(i, j int) bool {
			idI := commentList[i].ID
			idJ := commentList[j].ID
			return idI > idJ // 从小到大排序，若要从大到小排序可改为 idI > idJ
		})

	}

	return commentList, true
}

// LoadCommentList 加载评论列表到redis
// 评论是依赖于视频的，所以传入视频id vid
func LoadCommentList(ctx context.Context, vid, uid uint) ([]*model.Comment, error) {
	key := redis.GenerateCommentKey(vid)
	var comments []*model.Comment
	// 刷新过期时间
	defer func() {
		go func() { redis.RdbComment.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 判断redis里是否存在该视频的评论列表，不存在时从MySQL读取到Redis
	if n, err := redis.RdbComment.Exists(ctx, key).Result(); n == 0 && err == nil {
		if comments, err = updateRedisComments(ctx, key, vid, uid); err != nil {
			// 如果同步MySQL的redis信息有误，直接删除redis中存在的对应视频vid信息
			redis.RdbComment.Del(ctx, key)
			return comments, err
		}

	} else if err != nil {
		// 如果redis有错误，直接从MySQL中查询
		commentList, err := dao.GetCommentList(vid, uid)
		if err != nil {
			return nil, err
		}
		return commentList, nil
	}

	return comments, nil
}

// 从MySQL中查询当前视频vid对应的所有评论信息，更新到Redis中
func updateRedisComments(ctx context.Context, key string, vid, uid uint) ([]*model.Comment, error) {
	// 从数据库读取数据，如果数据库中的数据是可以正常查询出来的，就直接返回
	// 不能因为redis的更新有误就无法返回查询的评论列表
	commentList, err := dao.GetCommentList(vid, uid)
	if err != nil {
		return nil, err
	}

	// redis操作 先添加一个默认值，防止缓存穿透
	if err = redis.RdbComment.HSet(ctx, key, config.RedisValueOfNULL, config.RedisValueOfNULL).Err(); err == nil {
		// 并发写入redis
		wg := sync.WaitGroup{}
		errcCh := make(chan error)
		for _, comment := range commentList {
			wg.Add(1)
			go func(comment model.Comment) {
				defer wg.Done()
				commentJson, err := json.Marshal(comment)
				if err != nil {
					errcCh <- err
					return
				}
				if e := redis.RdbComment.HSet(ctx, key, comment.ID, commentJson).Err(); e != nil {
					errcCh <- e
					return
				}
			}(*comment)
		}
		wg.Wait()

		close(errcCh)

		//  处理并发过程中产生的错误
		if len(errcCh) > 0 {
			redis.RdbComment.Del(ctx, key)
			return commentList, nil
		}

		// 清除占位防止缓存穿透的NULL key
		var existNULL bool
		existNULL, err = redis.RdbComment.HExists(ctx, key, config.RedisValueOfNULL).Result()
		if err != nil {
			redis.RdbComment.Del(ctx, key)
			return commentList, nil
		}
		if existNULL {
			err = redis.RdbComment.HDel(ctx, key, config.RedisValueOfNULL).Err()
			if err != nil {
				redis.RdbComment.Del(ctx, key)
				return commentList, nil
			}
		}
		return commentList, nil
	} else {
		redis.RdbComment.Del(ctx, key)
		return commentList, nil
	}
}
