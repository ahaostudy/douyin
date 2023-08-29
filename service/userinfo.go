package service

import (
	"context"
	"main/config"
	"main/dao"
	"main/middleware/redis"
	"main/model"
	"reflect"
	"strconv"
	"sync"
)

// GetUserByID 通过ID获取用户
// id 目标用户ID
// curID 当前登录的用户ID
func GetUserByID(id, curID uint) (*model.User, bool) {
	ctx, cancel := redis.WithTimeoutContextBySecond(3)
	defer cancel()

	// 由于获取用户信息接口是高频接口
	// 所以对每步操作都采用并发执行，保证接口的性能
	var userinfo map[string]string
	var isFollow bool
	wg, status := sync.WaitGroup{}, true
	wg.Add(2)

	// 1. 从redis中获取用户信息
	go func() {
		defer wg.Done()
		info, ok := GetCacheUserInfo(ctx, id)
		if !ok {
			status = false
		}
		userinfo = info
	}()

	// 2. 获取 is_follow 字段
	//    因为 is_follow 是一个用户相对于另一个用户的值，而不是某个用户绝对的值，所以需要实时获取
	go func() {
		defer wg.Done()
		f, ok := IsFollow(curID, id)
		if !ok {
			status = false
		}
		isFollow = f
	}()

	// 3. 将 is_follow 字段添加到 userinfo 中
	wg.Wait()
	if !status {
		return nil, false
	}
	userinfo["is_follow"] = strconv.FormatBool(isFollow)

	// 4. 使用反射将用户信息转换为用户模型对象
	user, ok := infoToUser(userinfo)
	if !ok {
		return nil, false
	}

	return user, ok
}

// ExistsUserInfo 判断redis中是否有参数id对应的用户信息
func ExistsUserInfo(ctx context.Context, id uint) bool {
	key := redis.GenerateUserKey(id)
	n, err := redis.RdbUser.Exists(ctx, key).Result()
	return n != 0 && err == nil
}

// GetCacheUserInfo 获取缓存用户信息
func GetCacheUserInfo(ctx context.Context, id uint) (map[string]string, bool) {
	key := redis.GenerateUserKey(id)

	// 更新过期时间
	defer func() {
		go func() { redis.RdbUser.Expire(ctx, key, config.RedisKeyTTL) }()
	}()

	// 获取用户信息，如果获取不到时从MySQL中读取到redis
	info, err := redis.RdbUser.HGetAll(ctx, key).Result()
	if err != nil || info == nil || len(info) == 0 {
		// 获取到用户对象（不包含is_follow字段，因为is_follow字段是相对的值，而不是绝对的值）
		user, err := dao.GetUserInfoByID(id)
		if err != nil {
			return nil, false
		}

		// 使用反射遍历用户对象的每个字段，并将内容写入到redis
		// 为提高性能，每个字段都使用并发来写入
		typ, val := reflect.TypeOf(*user), reflect.ValueOf(*user)
		wg := sync.WaitGroup{}
		for i := 0; i < val.NumField(); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				tag, value := typ.Field(i).Tag.Get("json"), val.Field(i).Interface()
				if len(tag) == 0 || tag == "-" {
					return
				}
				if e := redis.RdbUser.HSet(ctx, key, tag, value).Err(); e != nil {
					err = e
				}
			}(i)
		}
		wg.Wait()

		// 如果过程中某一字段写入失败，则直接将key删除，防止缓存数据不全
		if err != nil {
			redis.RdbUser.Del(ctx, key)
			return nil, false
		}

		// 重新获取数据并返回
		info, err = redis.RdbUser.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, false
		}
	}

	return info, true
}

// 将map类型的info转换为用户模型对象
func infoToUser(info map[string]string) (*model.User, bool) {
	user := new(model.User)

	// 通过反射获取类型和值
	typ := reflect.TypeOf(user).Elem()
	val := reflect.ValueOf(user).Elem()

	// 遍历每个字段
	// 获取每个字段的json tag，因为redis存储时使用的是json tag作为key，这里解析也用json tag获取
	// 以tag值为key获取map中的值，通过反射将值写入用户对象的对应字段
	for i := 0; i < val.NumField(); i++ {
		// 获取tag和字段值对象
		tag := typ.Field(i).Tag.Get("json")
		field := val.Field(i)

		// 跳过忽略的tag
		if len(tag) == 0 || tag == "-" {
			continue
		}

		// 判断map中是否存在该key
		if v, ok := info[tag]; ok {
			// 判断字段的类型，将map获取出来的值解析为对应类型
			switch field.Kind() {
			case reflect.Int:
				intVal, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, false
				}
				field.SetInt(intVal)
			case reflect.Uint:
				uintVal, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return nil, false
				}
				field.SetUint(uintVal)
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(v)
				if err != nil {
					return nil, false
				}
				field.SetBool(boolVal)
			case reflect.String:
				if reflect.TypeOf(v).AssignableTo(field.Type()) {
					field.Set(reflect.ValueOf(v))
				} else {
					return nil, false
				}
			}
		} else {
			return nil, false
		}
	}

	return user, true
}
