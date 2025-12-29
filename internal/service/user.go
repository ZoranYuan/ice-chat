package service

import (
	"ice-chat/config"
	"ice-chat/internal/constants"
	"ice-chat/internal/model/request"
	res "ice-chat/internal/model/response"
	"ice-chat/internal/redisService"
	"ice-chat/internal/repository"
	"ice-chat/pkg/jwtUtils"
	"strconv"

	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户业务服务
type UserService struct {
	userRedis redisService.UserReids
	userRepo  repository.UserRepository
	minio     *minio.Client
}

// NewUserService 构造函数：注入Redis操作接口
func NewUserService(userRedis redisService.UserReids, userRepo repository.UserRepository, minio *minio.Client) *UserService {
	return &UserService{
		userRedis: userRedis,
		userRepo:  userRepo,
		minio:     minio,
	}
}

func (us *UserService) Login(v request.Login) (*res.Login, error) {
	user, err := us.userRepo.FindUserByEmail(v.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(v.Password)); err != nil {
		return nil, err
	}

	j := jwtUtils.CreateJwtUtils(jwtUtils.Config{
		Secret: []byte(config.Conf.JWT.AccessTokenExpireDuration),
		Expire: config.Conf.JWT.GetAccessTokenExpireDuration(),
	})

	token, claims, err := j.Generate(user.UserId)
	if err != nil {
		return nil, err
	}

	// 将 token 存入到 redis 中
	accessKey := constants.ACCESSKEY + strconv.FormatUint(claims.UserId, 10) + ":" + claims.JTI
	us.userRedis.StoreAccessKey(accessKey)

	var res *res.Login = &res.Login{
		Token:    token,
		UserId:   user.UserId,
		UserName: user.Username,
		Avatar:   user.Avatar,
	}

	return res, nil
}
