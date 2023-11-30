package usecase

import (
	"context"

	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/repository"
)

type (
	DropdownUsecase interface {
		ListProvince(ctx context.Context) ([]dto.DropdownValue, error)
		ListDistrict(ctx context.Context, payload dto.ListDistrictPayload) ([]dto.DropdownValue, error)
		ListShopCourier(ctx context.Context, payload dto.ListShopCourierPayload) ([]dto.DropdownValue, error)
		ListTopCategory(ctx context.Context) ([]dto.DropdownValue, error)
		ListCategoryByParentID(ctx context.Context, payload dto.ListCategoryByParentIDPayload) ([]dto.DropdownValue, error)
	}
	dropdownUsecase struct {
		pr  repository.ProvinceRepository
		dr  repository.DistrictRepository
		scr repository.ShopCourierRepository
		cr  repository.CategoryRepository
	}
)

// ListCategoryByParentID implements DropdownUsecase.
func (uc *dropdownUsecase) ListCategoryByParentID(ctx context.Context, payload dto.ListCategoryByParentIDPayload) ([]dto.DropdownValue, error) {
	cats, err := uc.cr.FindByParentCategoryID(ctx, payload.ParentCategoryID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.DropdownValue, 0)
	for _, c := range cats {
		temp := dto.DropdownValue{
			Label: c.Name,
			Value: c.ID,
		}

		res = append(res, temp)
	}

	return res, nil
}

// ListTopCategory implements DropdownUsecase.
func (uc *dropdownUsecase) ListTopCategory(ctx context.Context) ([]dto.DropdownValue, error) {
	cats, err := uc.cr.FindByLevel(ctx, 1)
	if err != nil {
		return nil, err
	}

	res := make([]dto.DropdownValue, 0)
	for _, c := range cats {
		temp := dto.DropdownValue{
			Label: c.Name,
			Value: c.ID,
		}

		res = append(res, temp)
	}

	return res, nil
}

// ListShopCourier implements DropdownUsecase.
func (uc *dropdownUsecase) ListShopCourier(ctx context.Context, payload dto.ListShopCourierPayload) ([]dto.DropdownValue, error) {
	couriers, err := uc.scr.FindShopCourierByShopId(ctx, payload.ShopID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.DropdownValue, 0)
	for _, courier := range couriers {
		if courier.IsAvailable {
			temp := dto.DropdownValue{
				Label: courier.CourierName,
				Value: courier.ShopCourierID,
			}

			res = append(res, temp)
		}
	}

	return res, nil
}

// ListDistrict implements DropdownUsecase.
func (uc *dropdownUsecase) ListDistrict(ctx context.Context, payload dto.ListDistrictPayload) ([]dto.DropdownValue, error) {
	districts, err := uc.dr.FindByProvinceID(ctx, payload.ProvinceID)
	if err != nil {
		return nil, err
	}

	resPayload := make([]dto.DropdownValue, 0)
	for _, p := range districts {
		temp := dto.DropdownValue{
			Label: p.Name,
			Value: p.ID,
		}

		resPayload = append(resPayload, temp)
	}

	return resPayload, nil
}

// ListProvince implements DropdownUsecase.
func (uc *dropdownUsecase) ListProvince(ctx context.Context) ([]dto.DropdownValue, error) {
	provinces, err := uc.pr.Find(ctx)
	if err != nil {
		return nil, err
	}

	resPayload := make([]dto.DropdownValue, 0)
	for _, p := range provinces {
		temp := dto.DropdownValue{
			Label: p.Name,
			Value: p.ID,
		}

		resPayload = append(resPayload, temp)
	}

	return resPayload, nil

}

func NewDropdownUsecase(
	pr repository.ProvinceRepository,
	dr repository.DistrictRepository,
	scr repository.ShopCourierRepository,
	cr repository.CategoryRepository,
) DropdownUsecase {
	return &dropdownUsecase{
		pr:  pr,
		dr:  dr,
		scr: scr,
		cr:  cr,
	}
}
