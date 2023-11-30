package dto

type (
	RegisterUserRequestPayload struct {
		Username string
		Email    string
		Password string
	}
	RegisterUserRequestBody struct {
		Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	LoginRequestPayload struct {
		Email    string
		Password string
	}
	LoginRequestBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	LoginResponsePayload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	RefreshTokenPayload struct {
		RefreshToken string
	}
	RefreshTokenResponsePayload struct {
		AccessToken string `json:"access_token"`
	}
	LogoutPayload struct {
		RefreshToken string
	}
	GetUserDetailPayload struct {
		UserID int64
	}
	GetUserDetailResponsePayload struct {
		UserID    int64  `json:"user_id"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		IsSeller  bool   `json:"is_seller"`
		ShopName  string `json:"shop_name,omitempty"`
		IsPinSet  bool   `json:"is_pin_set"`
		CartCount int64  `json:"cart_count"`
		ImageURL  string `json:"profile_picture_url"`
	}
	GetStepUpTokenRequestBody struct {
		WalletPin string `json:"wallet_pin" validate:"required,numeric,len=6"`
	}
	GetStepUpTokenPayload struct {
		WalletPin string
		UserID    int64
	}
	GetStepUpTokenResponse struct {
		StepUpToken string `json:"step_up_token"`
	}
	ChangeEmailRequestBody struct {
		Email string `json:"email" validate:"required,email"`
	}
	ChangeEmailPayload struct {
		UserID int64
		Email  string
	}
	ForgotPasswordRequestBody struct {
		Email string `json:"email" validate:"required,email"`
	}
	ForgotPasswordPayload struct {
		Email string
	}
	ResetPasswordRequestBody struct {
		ResetCode string `json:"reset_code" validate:"required"`
		Password  string `json:"password" validate:"required"`
	}
	ResetPasswordPayload struct {
		ResetCode string
		Password  string
	}
	ChangePasswordRequestBody struct {
		VerifCode string `json:"verif_code" validate:"required"`
		Password  string `json:"password" validate:"required"`
	}
	ChangePasswordPayload struct {
		VerifCode string
		Password  string
	}
)
