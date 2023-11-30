package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
)

type (
	CacheRepository interface {
		SetRefreshToken(ctx context.Context, refreshToken string, userID int64) error
		GetUserIdByRefreshToken(ctx context.Context, refreshToken string) (*int64, error)
		DeleteRefreshToken(ctx context.Context, refreshToken string) error
		SetPaymentTokenForUserID(ctx context.Context, paymentToken string, userID int64) error
		GetUserIDByPaymentToken(ctx context.Context, paymentToken string) (*int64, error)
		SetResetPasswordCode(ctx context.Context, resetPwCode string, userID int64) error
		GetUserIdByResetPasswordCode(ctx context.Context, resetPwCode string) (*int64, error)
		DeleteResetPasswordCode(ctx context.Context, resetPwCode string) error
		SetChangePasswordCode(ctx context.Context, changePwCode string, userID int64) error
		GetUserChangePasswordCodeByUserID(ctx context.Context, userID int64) (*string, error)
		DeleteChangePasswordCode(ctx context.Context, userID int64) error
		SetCountWrongPinWalletForUserID(ctx context.Context, userID int64) error
		UpdateCountWrongPinWalletForUserID(ctx context.Context, userID int64) error
		GetCountWrongPinWalletByUserID(ctx context.Context, userID int64) (*int, error)
		DeleteCountWrongPinWalletByUserID(ctx context.Context, userID int64) error
		SetLockedWalletForUserID(ctx context.Context, userID int64) error
		GetLockedWalletByUserID(ctx context.Context, userID int64) (*bool, error)
		GetRecommendedProduct(ctx context.Context) ([]dto.HomePageProductResponseBody, error)
	}
	cacheRepository struct {
		rd  *redis.Client
		cfg dependency.Config
	}
)

// GetUserIDByPaymentToken implements CacheRepository.
func (r *cacheRepository) GetUserIDByPaymentToken(ctx context.Context, paymentToken string) (*int64, error) {
	key := fmt.Sprintf(constant.RedisPaymentTokenTemplate, paymentToken)

	cmd := r.rd.Get(ctx, key)
	if cmd == nil {
		return nil, nil
	}

	if err := cmd.Err(); err != nil {
		return nil, err
	}

	userID, err := cmd.Int64()
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

// SetPaymentTokenForUserID implements CacheRepository.
func (r *cacheRepository) SetPaymentTokenForUserID(ctx context.Context, paymentToken string, userID int64) error {
	expiration := time.Duration(r.cfg.Jwt.StepUpTokenExpiration) * time.Minute

	key := fmt.Sprintf(constant.RedisPaymentTokenTemplate, paymentToken)
	cmd := r.rd.SetEX(ctx, key, userID, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

// DeleteRefreshToken implements CacheRepository.
func (r *cacheRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf(constant.RedisRefreshTokenTemplate, refreshToken)
	cmd := r.rd.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

// GetUserIdByRefreshToken implements CacheRepository.
func (r *cacheRepository) GetUserIdByRefreshToken(ctx context.Context, refreshToken string) (*int64, error) {
	key := fmt.Sprintf(constant.RedisRefreshTokenTemplate, refreshToken)

	cmd := r.rd.Get(ctx, key)
	if cmd == nil {
		return nil, nil
	}

	if err := cmd.Err(); err != nil {
		return nil, err
	}

	str, err := cmd.Int64()
	if err != nil {
		return nil, err
	}

	return &str, nil
}

// SetRefreshToken implements CacheRepository.
func (r *cacheRepository) SetRefreshToken(ctx context.Context, refreshToken string, userID int64) error {
	key := fmt.Sprintf(constant.RedisRefreshTokenTemplate, refreshToken)
	expiration := time.Duration(r.cfg.Jwt.RefreshTokenExpiration) * time.Minute

	cmd := r.rd.SetEX(ctx, key, userID, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) SetResetPasswordCode(ctx context.Context, resetPwCode string, userID int64) error {
	expiration := time.Duration(r.cfg.ResetPW.ResetPWCodeExpiration) * time.Minute

	key := fmt.Sprintf(constant.RedisResetPwCodeTemplate, resetPwCode)
	cmd := r.rd.SetEX(ctx, key, userID, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) GetUserIdByResetPasswordCode(ctx context.Context, resetPwCode string) (*int64, error) {
	key := fmt.Sprintf(constant.RedisResetPwCodeTemplate, resetPwCode)

	cmd := r.rd.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	str, err := cmd.Int64()
	if err != nil {
		return nil, err
	}

	return &str, nil
}

func (r *cacheRepository) DeleteResetPasswordCode(ctx context.Context, resetPwCode string) error {
	key := fmt.Sprintf(constant.RedisResetPwCodeTemplate, resetPwCode)
	cmd := r.rd.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) SetChangePasswordCode(ctx context.Context, resetPwCode string, userID int64) error {
	expiration := time.Duration(r.cfg.ChangePW.ChangePWCodeExpiration) * time.Minute

	key := fmt.Sprintf(constant.RedisChangePwCodeTemplate, userID)
	cmd := r.rd.SetEX(ctx, key, resetPwCode, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
func (r *cacheRepository) GetUserChangePasswordCodeByUserID(ctx context.Context, userID int64) (*string, error) {
	key := fmt.Sprintf(constant.RedisChangePwCodeTemplate, userID)

	cmd := r.rd.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	str, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return &str, nil
}
func (r *cacheRepository) DeleteChangePasswordCode(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(constant.RedisChangePwCodeTemplate, userID)
	cmd := r.rd.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) SetCountWrongPinWalletForUserID(ctx context.Context, userID int64) error {
	expiration := time.Duration(r.cfg.LockedWallet.LockedWalletExpiration) * time.Minute

	key := fmt.Sprintf(constant.RedisWrongPinTemplate, userID)
	cmd := r.rd.SetEX(ctx, key, 1, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) UpdateCountWrongPinWalletForUserID(ctx context.Context, userID int64) error {
	expiration := time.Duration(r.cfg.LockedWallet.LockedWalletExpiration) * time.Minute

	key := fmt.Sprintf(constant.RedisWrongPinTemplate, userID)
	cmd := r.rd.SetEX(ctx, key, 2, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) GetCountWrongPinWalletByUserID(ctx context.Context, userID int64) (*int, error) {
	key := fmt.Sprintf(constant.RedisWrongPinTemplate, userID)

	cmd := r.rd.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	ctr, err := cmd.Int()
	if err != nil {
		return nil, err
	}

	return &ctr, nil
}

func (r *cacheRepository) DeleteCountWrongPinWalletByUserID(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(constant.RedisWrongPinTemplate, userID)
	cmd := r.rd.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) SetLockedWalletForUserID(ctx context.Context, userID int64) error {
	expiration := time.Duration(r.cfg.LockedWallet.LockedWalletExpiration) * time.Minute

	key := fmt.Sprintf(constant.RedisLockedWalletTemplate, userID)
	cmd := r.rd.SetEX(ctx, key, true, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
func (r *cacheRepository) GetLockedWalletByUserID(ctx context.Context, userID int64) (*bool, error) {
	key := fmt.Sprintf(constant.RedisLockedWalletTemplate, userID)

	cmd := r.rd.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	boo, err := cmd.Bool()
	if err != nil {
		return nil, err
	}

	return &boo, nil
}

func (r *cacheRepository) GetRecommendedProduct(ctx context.Context) ([]dto.HomePageProductResponseBody, error) {
	resProducts := make([]dto.HomePageProductResponseBody, 0)
	products := r.rd.HGetAll(ctx, constant.RedisRecommendedProductTemplate)
	val := products.Val()
	err := json.Unmarshal([]byte(val[constant.RedisRecommendedProductTemplate]), &resProducts)
	if err != nil {
		return nil, err
	}
	return resProducts, nil
}

func NewCacheRepository(rd *redis.Client, cfg dependency.Config) CacheRepository {
	return &cacheRepository{
		rd:  rd,
		cfg: cfg,
	}
}
