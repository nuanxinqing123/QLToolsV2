package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/data/ent/user"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/nuanxinqing123/QLToolsV2/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

// NewAuthService 创建 AuthService
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register 用户注册
func (s *AuthService) Register(obj schema.RegisterRequest) error {
	ctx := context.Background()
	// 检查系统是否已存在用户（只允许注册一个用户）
	cnt, err := config.Ent.User.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}
	if cnt > 0 {
		return errors.New("系统已存在用户，不允许重复注册")
	}

	// 生成密码哈希（使用 bcrypt）
	hashed, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 使用 Ent 创建记录
	_, err = config.Ent.User.Create().
		SetUsername(obj.Username).
		SetPassword(string(hashed)).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

// Login 用户登录
func (s *AuthService) Login(obj schema.LoginRequest) (*schema.LoginResponse, error) {
	ctx := context.Background()
	// 按用户名查询用户
	u, err := config.Ent.User.Query().
		Where(user.UsernameEQ(obj.Username)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("用户不存在或查询失败: %w", err)
	}

	// 比对密码
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(obj.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT Token对
	jwtManager := utils.NewJWTManager()
	accessToken, refreshToken, err := jwtManager.GenerateTokenPair(u.ID)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	return &schema.LoginResponse{
		Message:      "登录成功",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout() error {
	jwtManager := utils.NewJWTManager()
	return jwtManager.RevokeToken()
}

// RefreshToken 刷新访问Token
func (s *AuthService) RefreshToken(obj schema.RefreshTokenRequest) (schema.RefreshTokenResponse, error) {
	jwtManager := utils.NewJWTManager()
	nat, err := jwtManager.RefreshAccessToken(obj.RefreshToken)
	return schema.RefreshTokenResponse{
		Message:     "token刷新成功",
		AccessToken: nat,
	}, err
}
