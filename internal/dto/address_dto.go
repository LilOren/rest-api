package dto

type (
	AddAddressRequestBody struct {
		ReceiverName        string `json:"receiver_name" validate:"required,max=20"`
		ReceiverPhoneNumber string `json:"receiver_phone_number" validate:"required,max=14"`
		Address             string `json:"address" validate:"required"`
		ProvinceId          int    `json:"province_id" validate:"required"`
		CityId              int    `json:"city_id" validate:"required"`
		SubDistrict         string `json:"sub_district" validate:"required"`
		SubSubDistrict      string `json:"sub_sub_district" validate:"required"`
		PostalCode          string `json:"postal_code" validate:"required"`
	}
	AddAddressPayload struct {
		ReceiverName        string
		ReceiverPhoneNumber string
		Address             string
		ProvinceId          int
		CityId              int
		SubDistrict         string
		SubSubDistrict      string
		PostalCode          string
	}
	AddAddressResponseBody struct {
		ID                  int    `json:"id"`
		ReceiverName        string `json:"receiver_name"`
		ReceiverPhoneNumber string `json:"receiver_phone_number"`
		Address             string `json:"address"`
		ProvinceId          int    `json:"province_id"`
		CityId              int    `json:"city_id"`
		SubDistrict         string `json:"sub_district"`
		SubSubDistrict      string `json:"sub_sub_district"`
		PostalCode          string `json:"postal_code"`
	}
	AccountDetailsAddressBody struct {
		ID int `json:"id" validate:"required"`
	}
	AccountDetailsAddressPayload struct {
		ID int
	}
	AccountDetailsAddress struct {
		ID                  int    `db:"id"`
		ReceiverName        string `db:"receiver_name"`
		Details             string `db:"detail"`
		PostalCode          string `db:"postal_code"`
		ReceiverPhoneNumber string `db:"receiver_phone_number"`
	}
	AccountDetailsAddressResponse struct {
		ID                  int    `json:"id"`
		ReceiverName        string `json:"receiver_name"`
		Details             string `json:"address"`
		PostalCode          string `json:"postal_code"`
		ReceiverPhoneNumber string `json:"receiver_phone_number"`
	}
	UpdateAddressByIDRequestBody struct {
		ReceiverName        string `json:"receiver_name" validate:"omitempty,min=4"`
		ReceiverPhoneNumber string `json:"receiver_phone_number" validate:"omitempty,e164"`
		Address             string `json:"address" validate:"omitempty,min=5"`
		PostalCode          string `json:"postal_code" validate:"omitempty,numeric,len=5"`
		ProvinceID          int64  `json:"province_id" validate:"omitempty,numeric"`
		DistrictID          int64  `json:"district_id" validate:"omitempty,numeric"`
	}
	UpdateAddressByIDPayload struct {
		UserID              int64
		AddressID           int64
		ReceiverName        string
		ReceiverPhoneNumber string
		Address             string
		PostalCode          string
		ProvinceID          int64
		DistrictID          int64
	}

	GetAddressByIDPayload struct {
		AccountID int64
		AddressID int64
	}
	GetAddressByIDResponse struct {
		AddressID           int64  `json:"address_id"`
		ReceiverName        string `json:"receiver_name"`
		ReceiverPhoneNumber string `json:"receiver_phone_number"`
		Address             string `json:"address"`
		PostalCode          string `json:"postal_code"`
		ProvinceID          int64  `json:"province_id"`
		ProvinceName        string `jsoN:"province_name"`
		DistrictID          int64  `json:"district_id"`
		DistrictName        string `json:"district_name"`
	}
)
