package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
)

type (
	ProfileUsecase interface {
		AddAddress(ctx context.Context, payload dto.AddAddressPayload, accountId int) error
		GetAddressDetailsByAccountId(ctx context.Context, accountId int) ([]dto.AccountDetailsAddressResponse, error)
		ChangeDefaultAddress(ctx context.Context, accountId int, defaultAddressId int) error
		UploadProfilePicture(ctx context.Context, accountID int64, photoURL string) error
		UpdateAddress(ctx context.Context, payload dto.UpdateAddressByIDPayload) error
		GetAddressDetailByID(ctx context.Context, payload dto.GetAddressByIDPayload) (*dto.GetAddressByIDResponse, error)
	}
	profileUsecase struct {
		aar repository.AccountAddressRepository
		ar  repository.AccountRepository
		dr  repository.DistrictRepository
		pr  repository.ProvinceRepository
	}
)

// GetAddressDetailByID implements ProfileUsecase.
func (uc *profileUsecase) GetAddressDetailByID(ctx context.Context, payload dto.GetAddressByIDPayload) (*dto.GetAddressByIDResponse, error) {
	address, err := uc.aar.FirstByID(ctx, payload.AddressID)
	if err != nil {
		return nil, err
	}

	if address.AccountId != payload.AccountID {
		return nil, shared.ErrAddressNotBelongToCurrentUser
	}

	res := dto.GetAddressByIDResponse{
		AddressID:           address.ID,
		ReceiverName:        address.ReceiverName,
		ReceiverPhoneNumber: address.ReceiverPhoneNumber,
		Address:             address.Detail,
		PostalCode:          address.PostalCode,
		ProvinceID:          address.ProvinceId,
		DistrictID:          address.DistrictId,
	}

	district, err := uc.dr.FirstByID(ctx, address.DistrictId)
	if err != nil {
		return nil, err
	}

	res.DistrictName = district.Name

	province, err := uc.pr.FirstByID(ctx, address.ProvinceId)
	if err != nil {
		return nil, err
	}

	res.ProvinceName = province.Name

	return &res, nil
}

// UpdateAddress implements ProfileUsecase.
func (uc *profileUsecase) UpdateAddress(ctx context.Context, payload dto.UpdateAddressByIDPayload) error {
	// check if address belongs to current user
	address, err := uc.aar.FirstByID(ctx, payload.AddressID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.ErrAddressNotFound
		}

		return err
	}

	if address.AccountId != payload.UserID {
		return shared.ErrAddressNotBelongToCurrentUser
	}

	if payload.ReceiverName != "" {
		address.ReceiverName = payload.ReceiverName
	}

	if payload.Address != "" {
		address.Detail = payload.Address
	}

	if payload.PostalCode != "" {
		address.PostalCode = payload.PostalCode
	}

	if payload.ReceiverPhoneNumber != "" {
		address.ReceiverPhoneNumber = payload.ReceiverPhoneNumber
	}

	if payload.ProvinceID > 0 {
		address.ProvinceId = payload.ProvinceID

		if payload.DistrictID == 0 {
			return shared.ErrModifyProvinceShouldModifyDistrict
		}
	}

	if payload.DistrictID > 0 {
		address.DistrictId = payload.DistrictID

		d, err := uc.dr.FirstByID(ctx, address.DistrictId)
		if err != nil {
			return err
		}

		if d.ProvinceID != address.ProvinceId {
			return shared.ErrDistrictNotBelongToProvince
		}
	}

	if err := uc.aar.UpdateAddressByID(ctx, *address); err != nil {
		return err
	}

	return nil
}

// UploadProfilePicture implements ProfileUsecase.
func (uc *profileUsecase) UploadProfilePicture(ctx context.Context, accountID int64, photoURL string) error {
	err := uc.ar.UpdateProfilePicture(ctx, accountID, photoURL)
	if err != nil {
		return err
	}

	return nil
}

func (uc *profileUsecase) AddAddress(ctx context.Context, payload dto.AddAddressPayload, accountId int) error {
	addressList, err := uc.aar.FindAddressById(ctx, accountId)
	if err != nil {
		return shared.ErrFailedGetAddress
	}

	length := len(addressList)

	address := model.AccountAddresses{
		ReceiverName:        payload.ReceiverName,
		ReceiverPhoneNumber: payload.ReceiverPhoneNumber,
		Detail:              fmt.Sprintf(payload.Address + ", " + payload.SubSubDistrict + ", " + payload.SubDistrict),
		ProvinceId:          int64(payload.ProvinceId),
		DistrictId:          int64(payload.CityId),
		PostalCode:          payload.PostalCode,
	}

	err = uc.aar.CreateAddress(ctx, address, length, accountId)
	if err != nil {
		if strings.Contains(err.Error(), "province") {
			return shared.ErrWrongProvinceId
		}
		if strings.Contains(err.Error(), "district") {
			return shared.ErrWrongDistrictId
		}
		return shared.ErrCreateAddress
	}

	return nil
}

func (uc *profileUsecase) GetAddressDetailsByAccountId(ctx context.Context, accountId int) ([]dto.AccountDetailsAddressResponse, error) {
	response := make([]dto.AccountDetailsAddressResponse, 0)
	addressList, err := uc.aar.FindDetailsAddressById(ctx, accountId)
	if err != nil {
		return nil, shared.ErrFailedGetAddress
	}

	for _, v := range addressList {
		r := dto.AccountDetailsAddressResponse{
			ID:                  int(v.ID),
			ReceiverName:        v.ReceiverName,
			Details:             v.Details,
			PostalCode:          v.PostalCode,
			ReceiverPhoneNumber: v.ReceiverPhoneNumber,
		}

		response = append(response, r)
	}

	return response, nil
}

func (uc *profileUsecase) ChangeDefaultAddress(ctx context.Context, accountId int, defaultAddressId int) error {

	address, err := uc.aar.FindAddressById(ctx, accountId)
	if err != nil {
		return shared.ErrFailedGetAddress
	}

	if defaultAddressId <= 0 {
		return shared.ErrInvalidAddressId
	}

	for idx, v := range address {
		if defaultAddressId == int(v.ID) {
			break
		}
		if idx == len(address)-1 {
			return shared.ErrInvalidAddressId
		}
	}

	err = uc.aar.UpdateDefaultAddress(ctx, accountId, defaultAddressId)
	if err != nil {
		return shared.ErrFailedUpdateDefaultAddress
	}

	return nil
}

func NewProfileUsecase(
	aar repository.AccountAddressRepository,
	ar repository.AccountRepository,
	dr repository.DistrictRepository,
	pr repository.ProvinceRepository,
) ProfileUsecase {
	return &profileUsecase{
		aar: aar,
		ar:  ar,
		dr:  dr,
		pr:  pr,
	}
}
