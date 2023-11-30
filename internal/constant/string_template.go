package constant

const (
	ConnectionStringTemplate = "host=%s user=%s password=%s dbname=%s port=%s timezone=Asia/Jakarta sslmode=disable"
	RedisConnectionTemplate  = "%s:%s"

	RedisRefreshTokenTemplate       = "refresh_token:%s"
	RedisPaymentTokenTemplate       = "payment_token:%s"
	RedisResetPwCodeTemplate        = "reset_password:%s"
	RedisChangePwCodeTemplate       = "change_password:%d"
	RedisWrongPinTemplate           = "wrong_pin:%d"
	RedisLockedWalletTemplate       = "locked_wallet:%d"
	RedisRecommendedProductTemplate = "recommended_product"

	VerifCodeAlphaNum = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
)
