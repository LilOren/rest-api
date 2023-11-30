package usecase

import (
	"context"
	"errors"
	"math"

	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
)

type (
	ReviewUsecase interface {
		GetAllReviewOfProduct(ctx context.Context, code string, params *dto.GetAllReviewParams) (*dto.GetAllReviewResponse, error)
		AddReviewOfProduct(ctx context.Context, payload *dto.AddReviewPayload) error
	}
	reviewUsecase struct {
		pr repository.ProductRepository
		rr repository.ReviewRepository
	}
)

func (ruc *reviewUsecase) GetAllReviewOfProduct(ctx context.Context, code string, params *dto.GetAllReviewParams) (*dto.GetAllReviewResponse, error) {
	reviews, err := ruc.rr.FindByProductCode(ctx, code, params)
	if err != nil {
		return nil, err
	}
	userReview := make([]dto.GetAllReviewUserResponse, 0)
	for _, review := range reviews {
		userReview = append(userReview, dto.GetAllReviewUserResponse{
			Rating:    review.Rating,
			Comment:   review.Comment,
			AccountID: review.AccountID,
			Username:  review.Username,
			ImageUrls: review.ImageUrls,
			CreatedAt: review.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	metadata, err := ruc.rr.FindByProductCodeMetadata(ctx, code, params)
	if err != nil {
		return nil, err
	}
	totalData := len(metadata)
	totalPage := math.Ceil(float64(totalData) / constant.ProductReviewDefaultItems)
	resp := &dto.GetAllReviewResponse{
		UserReview:  userReview,
		TotalReview: totalData,
		TotalPage:   int(totalPage),
		CurrentPage: params.Page,
	}
	return resp, nil
}

func (ruc *reviewUsecase) AddReviewOfProduct(ctx context.Context, payload *dto.AddReviewPayload) error {
	_, err := ruc.pr.FirstProductByCode(ctx, payload.ProductCode)
	if err != nil {
		if errors.Is(err, shared.ErrProductNotFound) {
			return shared.ErrProductNotFound
		}
		return err
	}
	review := &model.Review{
		Rating:      payload.Rating,
		Comment:     payload.Comment,
		ProductCode: payload.ProductCode,
		AccountID:   payload.AccountID,
	}
	if err := ruc.rr.Create(ctx, review, payload.ImageUrls); err != nil {
		return err
	}
	return nil
}

func NewReviewUsecase(rr repository.ReviewRepository, pr repository.ProductRepository) ReviewUsecase {
	return &reviewUsecase{
		rr: rr,
		pr: pr,
	}
}
