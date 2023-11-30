package shared

import (
	"fmt"
)

var (
	// general
	ErrInvalidBodySchema = NewCustomError(BadRequest, "Invalid Body Schema")

	// auth
	ErrUsernameAlreadyTaken     = NewCustomError(BadRequest, "Username already taken")
	ErrPasswordContainsUsername = NewCustomError(BadRequest, "Password contains username")
	ErrEmailNotAvailable        = NewCustomError(BadRequest, "Email not available")
	ErrEmailAlreadyUsed         = NewCustomError(BadRequest, "Email already used")
	ErrUsernameNotAvailable     = NewCustomError(BadRequest, "Username not available")
	ErrPasswordNotMatchRegex    = NewCustomError(BadRequest, "Password must contain at least 1 lowercase letter, 1 uppercase letter, 1 number, and with a minimum of 8 characters")
	ErrInvalidEmailOrPassword   = NewCustomError(BadRequest, "Invalid email or password")
	ErrInvalidPassword          = NewCustomError(BadRequest, "Invalid password")
	ErrInvalidAuthHeader        = NewCustomError(Unauthorized, "Invalid header")
	ErrInvalidToken             = NewCustomError(Unauthorized, "Invalid Token")
	ErrAccessTokenExpired       = NewCustomError(Unauthorized, "AccessTokenExpired")
	ErrRefreshTokenExpired      = NewCustomError(Unauthorized, "RefreshTokenExpired")
	ErrStepUpTokenExpired       = NewCustomError(Unauthorized, "StepUpTokenExpired")
	ErrInvalidTokenType         = NewCustomError(Unauthorized, "Invalid token type")
	ErrUserAlreadyLogout        = NewCustomError(BadRequest, "User already logged out")
	ErrUserDetailNotFound       = NewCustomError(BadRequest, "User not found")
	ErrFailedCreateWallet       = NewCustomError(InternalServer, "Failed create wallet")
	ErrSamePassword             = NewCustomError(BadRequest, "Password must be different")

	// user
	ErrFailedGetLocation             = NewCustomError(InternalServer, "Failed getting location")
	ErrFailedGetAddress              = NewCustomError(InternalServer, "Failed getting account address")
	ErrWrongPostalCode               = NewCustomError(BadRequest, "Wrong postal code")
	ErrWrongProvinceId               = NewCustomError(BadRequest, "Wrong province id")
	ErrWrongDistrictId               = NewCustomError(BadRequest, "Wrong district id")
	ErrDistrictNotFound              = NewCustomError(BadRequest, "District not found")
	ErrCreateAddress                 = NewCustomError(InternalServer, "Failed creating address")
	ErrFailedUpdateDefaultAddress    = NewCustomError(InternalServer, "Failed update default address")
	ErrNoAddress                     = NewCustomError(BadRequest, "Not yet filled in address")
	ErrAccountNotFound               = NewCustomError(NotFound, "user not found")
	ErrEmailNotFound                 = NewCustomError(BadRequest, "User email not found")
	ErrAddressNotFound               = NewCustomError(NotFound, "address not found")
	ErrAddressNotBelongToCurrentUser = NewCustomError(BadRequest, "address not belong to this user")
	ErrResetPasswordCodeExpired      = NewCustomError(BadRequest, "ResetPasswordCodeExpired")
	ErrChangePasswordExist           = NewCustomError(BadRequest, "User already request to change password")
	ErrChangePasswordCodeExpired     = NewCustomError(BadRequest, "ChangePasswordCodeExpired")
	ErrUnknownVerifCode              = NewCustomError(BadRequest, "Unknown verification code")

	// shop
	ErrNoShop                        = NewCustomError(NotFound, "Not yet registered as seller")
	ErrAlreadyHaveShop               = NewCustomError(BadRequest, "Already have shop")
	ErrFailedCreateShop              = NewCustomError(InternalServer, "Failed create shop")
	ErrFindShop                      = NewCustomError(InternalServer, "Failed find shop")
	ErrSameShopName                  = NewCustomError(BadRequest, "Shop name already used last time")
	ErrFailedGetShop                 = NewCustomError(InternalServer, "Failed getting shop")
	ErrShopNameAlreadyTaken          = NewCustomError(BadRequest, "Shop name already taken")
	ErrFailedUpdateShopName          = NewCustomError(InternalServer, "Failed update shop name")
	ErrFailedUpdateShopAddress       = NewCustomError(InternalServer, "Failed update shop address")
	ErrFailedUpdateShopCourier       = NewCustomError(InternalServer, "Failed update shop courier")
	ErrInvalidAddressId              = NewCustomError(BadRequest, "Invalid address id")
	ErrAlreadyDefaultShopAddress     = NewCustomError(BadRequest, "Address Id already chosen as default shop address")
	ErrCourierNotFull                = NewCustomError(BadRequest, "Courier condition not fully filled")
	ErrCourierNotAvailable           = NewCustomError(BadRequest, "Courier not available")
	ErrInvalidCourier                = NewCustomError(BadRequest, "Courier condition more than available courier")
	ErrCourierNotBelongToCurrentShop = NewCustomError(BadRequest, "Courier doesn't belong to current shop")
	ErrFailedFindShopCourier         = NewCustomError(InternalServer, "Failed find shop courier")
	ErrFailedActivateShopWallet      = NewCustomError(InternalServer, "Failed activate shop wallet")
	ErrNoCategory                    = NewCustomError(BadRequest, "Invalid product category")
	ErrInvalidWeight                 = NewCustomError(BadRequest, "Weight equal or less than 0")
	ErrInvalidPrice                  = NewCustomError(BadRequest, "Price equal or less than 0")
	ErrInvalidStock                  = NewCustomError(BadRequest, "Stock less than 0")
	ErrUpdateProduct                 = NewCustomError(InternalServer, "Failed update product details")
	ErrFindProductDetail             = NewCustomError(InternalServer, "Failed find product detail")
	ErrShopNotFound                  = NewCustomError(NotFound, "Shop did not exist")
	ErrFindShopDetail                = NewCustomError(InternalServer, "Failed find shop detail")
	ErrFindProduct                   = NewCustomError(InternalServer, "Failed find product")
	ErrInvalidPage                   = NewCustomError(BadRequest, "Invalid page")
	ErrFindCategory                  = NewCustomError(InternalServer, "Failed find category")
	ErrCountProduct                  = NewCustomError(InternalServer, "Failed count product")
	ErrFindProductDiscount           = NewCustomError(InternalServer, "Failed find product discount")
	ErrInvalidProductCode            = NewCustomError(BadRequest, "Invalid Product Code")
	ErrUpdateProductDiscount         = NewCustomError(InternalServer, "Failed update product discount")
	ErrShopNameIsNull                = NewCustomError(BadRequest, "Seller not yet set shop name")
	ErrProductNotFromSeller          = NewCustomError(BadRequest, "Seller do not own this product")
	ErrDeleteProduct                 = NewCustomError(InternalServer, "Failed delete product")

	// cart
	ErrDifferentSeller       = NewCustomError(BadRequest, "Product from the shop not found")
	ErrOwnSellerProduct      = NewCustomError(BadRequest, "Seller cannot add their own product to cart")
	ErrQuantityMoreThanStock = NewCustomError(BadRequest, "Cannot add quantity more than stock")
	ErrCartNotFound          = NewCustomError(NotFound, "Product not found in cart")
	ErrFindCart              = NewCustomError(NotFound, "Cart not found")
	ErrFailedDeleteInCart    = NewCustomError(InternalServer, "Failed delete item in cart")
	ErrNoCheckedCart         = NewCustomError(BadRequest, "No checked cart")

	// wallet
	ErrFindWallet           = NewCustomError(InternalServer, "Failed find wallet")
	ErrWalletNotActivated   = NewCustomError(BadRequest, "Wallet is not activated")
	ErrWalletPinIsNotSet    = NewCustomError(BadRequest, "Wallet pin is not set")
	ErrWrongWalletPin       = NewCustomError(BadRequest, "Wallet pin is wrong")
	ErrInsufficientBalance  = NewCustomError(BadRequest, "Insufficient balance")
	ErrTransferWallet       = NewCustomError(InternalServer, "Failed transfer wallet user to temp")
	ErrUpdateInactiveWallet = NewCustomError(BadRequest, "trying to update inactive wallet")
	ErrSameWalletPin        = NewCustomError(BadRequest, "Wallet pin cannot be same as previous pin")
	ErrWalletIsLocked       = NewCustomError(BadRequest, "Wallet is temporarily locked")

	// order
	ErrCreateOrder               = NewCustomError(InternalServer, "Failed create order")
	ErrWrongInitialStatus        = NewCustomError(BadRequest, "Wrong initial status")
	ErrCreateOrderProductVariant = NewCustomError(InternalServer, "Failed create order product variant")
	ErrFindOrder                 = NewCustomError(InternalServer, "Failed find order")
	ErrOrderIDNotFount           = NewCustomError(NotFound, "Order not found")
	ErrDuplicateVariantType      = NewCustomError(BadRequest, "Duplicate variant type")
	ErrUnauthorizedUser          = NewCustomError(Unauthorized, "User does not have access")

	// wishlist
	ErrWishlistAlreadyExist = NewCustomError(BadRequest, "Wishlist already exist")
	ErrProductNotFound      = NewCustomError(BadRequest, "Product does not exist")
	ErrWishlistNotFound     = NewCustomError(BadRequest, "Wishlist does not exist")

	// district
	ErrDistrictNotBelongToProvince        = NewCustomError(BadRequest, "District does not belong to current province")
	ErrModifyProvinceShouldModifyDistrict = NewCustomError(BadRequest, "Modifying province should also modify district")

	// promotion
	ErrPercentAndPriceNull = NewCustomError(BadRequest, "Percentage and Exact Price cannot both be 0")
	ErrExpiredBeforeStart  = NewCustomError(BadRequest, "Cannot set expired date before promo start date")
	ErrPromotionNotFound   = NewCustomError(NotFound, "Promotion did not exist")
	ErrFindPromotion       = NewCustomError(InternalServer, "Failed find promotion")
	ErrPromoNotFound       = NewCustomError(NotFound, "Shop promotion not found")

	// internal server
	ErrInternalServer = NewCustomError(InternalServer, "Internal server error")
)

func GenerateErrQueryParamRequired(param string) *CustomError {

	return NewCustomError(BadRequest, fmt.Sprintf("Query param %s is required", param))
}

func GenerateErrQueryParamInvalid(param string) *CustomError {

	return NewCustomError(BadRequest, fmt.Sprintf("Query param %s is invalid", param))
}

func GenerateErrPathParamInvalid(param string) *CustomError {

	return NewCustomError(BadRequest, fmt.Sprintf("Path param %s is invalid", param))
}
