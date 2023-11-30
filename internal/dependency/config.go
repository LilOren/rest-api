package dependency

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		App          app
		Rest         rest
		PostgreDB    postgreDB
		RedisCache   redisconfig
		Jwt          jwt
		ThirdParty   thirdParty
		EmailSender  emailSender
		ResetPW      resetPW
		ChangePW     changePW
		LockedWallet lockedWallet
		GOauth       gOauth
	}

	app struct {
		AppName         string `env:"APP_NAME"`
		GracefulTimeout uint   `env:"GRACEFUL_TIMEOUT"`
		OriginDomain    string `env:"ORIGIN_DOMAIN"`
	}

	rest struct {
		RequestTimeout uint `env:"REQUEST_TIMEOUT"`
		Port           uint `env:"REST_PORT"`
	}

	postgreDB struct {
		DBHost string `env:"DB_HOST"`
		DBPort string `env:"DB_PORT"`
		DBName string `env:"DB_NAME"`
		DBUser string `env:"DB_USER"`
		DBPass string `env:"DB_PASS"`
		DBTz   string `env:"DB_TZ" env-default:"Asia/Jakarta"`
	}

	redisconfig struct {
		HOST     string `env:"REDIS_HOST"`
		PORT     string `env:"REDIS_PORT"`
		Password string `env:"REDIS_PASSWORD"`
	}

	jwt struct {
		JWTSecret              string `env:"JWT_SECRET"`
		AccessTokenExpiration  uint   `env:"ACCESS_TOKEN_EXPIRATION"`
		RefreshTokenExpiration uint   `env:"REFRESH_TOKEN_EXPIRATION"`
		StepUpTokenExpiration  int    `env:"STEP_UP_TOKEN_EXPIRATION"`
	}

	thirdParty struct {
		RajaOngkirBaseURL string `env:"RO_BASE_URL"`
		RajaOngkirAPIKey  string `env:"RO_API_KEY"`
	}

	emailSender struct {
		Name     string `env:"EMAIL_SENDER_NAME"`
		Address  string `env:"EMAIL_SENDER_ADDRESS"`
		Password string `env:"EMAIL_SENDER_PASSWORD"`
	}

	resetPW struct {
		ResetPWCodeExpiration uint `env:"RESET_PW_CODE_EXPIRATION"`
	}

	changePW struct {
		ChangePWCodeExpiration uint `env:"CHANGE_PW_CODE_EXPIRATION"`
	}

	lockedWallet struct {
		LockedWalletExpiration uint `env:"LOCKED_WALLET_EXPIRATION"`
	}

	gOauth struct {
		ClientID     string `env:"GOOGLE_OAUTH_CLIENT_ID"`
		ClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET"`
		RedirectURL  string `env:"GOOGLE_OAUTH_REDIRECT_URL"`
		RedirectFE   string `env:"GOOGLE_REDIRECT_FRONTEND"`
	}
)

func NewConfig(logger Logger) (*Config, error) {
	config := new(Config)

	err := cleanenv.ReadEnv(config)
	if err != nil {
		logger.Fatalf("Failed to load config")
		return nil, err
	}

	logger.Infof("Successfully load config", nil)

	return config, err
}
