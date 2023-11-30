package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/lil-oren/rest/internal/constant"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthUsecase interface {
		RegisterUser(ctx context.Context, payload dto.RegisterUserRequestPayload) error
		Login(ctx context.Context, payload dto.LoginRequestPayload) (*dto.LoginResponsePayload, error)
		LoginWithGoogle(ctx context.Context, googleUser *dto.GoogleResponse) (*dto.LoginResponsePayload, error)
		Logout(ctx context.Context, payload dto.LogoutPayload) error
		RefreshToken(ctx context.Context, payload dto.RefreshTokenPayload) (*dto.RefreshTokenResponsePayload, error)
		GetUserDetail(ctx context.Context, payload dto.GetUserDetailPayload) (*dto.GetUserDetailResponsePayload, error)
		GetPaymentToken(ctx context.Context, payload dto.GetStepUpTokenPayload) (*dto.GetStepUpTokenResponse, error)
		ChangeEmail(ctx context.Context, payload dto.ChangeEmailPayload) error
		ForgotPassword(ctx context.Context, payload dto.ForgotPasswordPayload) error
		ResetPassword(ctx context.Context, payload dto.ResetPasswordPayload) error
		RequestChangePassword(ctx context.Context, userId int64) error
		ChangePassword(ctx context.Context, payload dto.ChangePasswordPayload, userID int64) error
	}
	authUsecase struct {
		ar       repository.AccountRepository
		cr       repository.CacheRepository
		cartRepo repository.CartRepository
		er       repository.WalletRepository
		cer      repository.ChangedEmailRepository
		sr       repository.ShopRepository
		cfg      dependency.Config
	}
)

// ChangeEmail implements AuthUsecase.
func (uc *authUsecase) ChangeEmail(ctx context.Context, payload dto.ChangeEmailPayload) error {
	email, _ := uc.cer.FirstByEmail(ctx, payload.Email)

	if email != nil {
		return shared.ErrEmailAlreadyUsed
	}

	user, err := uc.ar.FirstById(ctx, payload.UserID)
	if err != nil {
		return err
	}

	if user.Email == payload.Email {
		return shared.ErrEmailAlreadyUsed
	}

	if err := uc.cer.Create(ctx, payload.UserID, payload.Email, user.Email); err != nil {
		return err
	}

	return nil
}

// GetPaymentToken implements AuthUsecase.
func (uc *authUsecase) GetPaymentToken(ctx context.Context, payload dto.GetStepUpTokenPayload) (*dto.GetStepUpTokenResponse, error) {
	userID := payload.UserID
	user, err := uc.ar.FirstById(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.PinHash.Valid {
		return nil, shared.ErrWalletPinIsNotSet
	}

	boo, err := uc.cr.GetLockedWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if boo != nil {
		return nil, shared.ErrWalletIsLocked
	}

	hashedPinBytes := []byte(user.PinHash.String)
	pinBytes := []byte(payload.WalletPin)
	if err = bcrypt.CompareHashAndPassword(hashedPinBytes, pinBytes); err != nil {
		ctr, err := uc.cr.GetCountWrongPinWalletByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if ctr == nil {
			if err := uc.cr.SetCountWrongPinWalletForUserID(ctx, userID); err != nil {
				return nil, err
			}
			return nil, shared.ErrWrongWalletPin
		}
		if *ctr == 1 {
			if err := uc.cr.UpdateCountWrongPinWalletForUserID(ctx, userID); err != nil {
				return nil, err
			}
			return nil, shared.ErrWrongWalletPin
		}
		if err := uc.cr.DeleteCountWrongPinWalletByUserID(ctx, userID); err != nil {
			return nil, err
		}
		if err := uc.cr.SetLockedWalletForUserID(ctx, userID); err != nil {
			return nil, err
		}
		return nil, shared.ErrWrongWalletPin
	}
	if err := uc.cr.DeleteCountWrongPinWalletByUserID(ctx, userID); err != nil {
		return nil, err
	}

	token, err := shared.SignStepUpToken(uc.cfg)
	if err != nil {
		return nil, err
	}

	resPayload := dto.GetStepUpTokenResponse{
		StepUpToken: *token,
	}

	if err := uc.cr.SetPaymentTokenForUserID(ctx, *token, userID); err != nil {
		return nil, err
	}

	return &resPayload, nil
}

// GetUserDetail implements AuthUsecase.
func (uc *authUsecase) GetUserDetail(ctx context.Context, payload dto.GetUserDetailPayload) (*dto.GetUserDetailResponsePayload, error) {

	account, err := uc.ar.FirstById(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrUserDetailNotFound
		}

		return nil, err
	}

	count, err := uc.cartRepo.CountCartByAccountID(ctx, payload.UserID)
	if err != nil {
		return nil, err
	}

	resPayload := dto.GetUserDetailResponsePayload{
		UserID:    account.ID,
		Username:  account.Username,
		Email:     account.Email,
		IsSeller:  account.IsSeller,
		CartCount: *count,
		IsPinSet:  account.PinHash.Valid,
	}

	if account.ProfilePictureURL.Valid {
		resPayload.ImageURL = account.ProfilePictureURL.String
	}

	if account.IsSeller {
		shop, err := uc.sr.FirstShopById(ctx, int(payload.UserID))
		if err != nil {
			return nil, shared.ErrFindShop
		}
		if !shop.Name.Valid {
			return nil, shared.ErrShopNameIsNull
		}
		resPayload.ShopName = shop.Name.String
	}

	return &resPayload, nil
}

// Logout implements AuthUsecase.
func (uc *authUsecase) Logout(ctx context.Context, payload dto.LogoutPayload) error {
	if err := uc.cr.DeleteRefreshToken(ctx, payload.RefreshToken); err != nil {
		return err
	}

	return nil
}

// RefreshToken implements AuthUsecase.
func (uc *authUsecase) RefreshToken(ctx context.Context, payload dto.RefreshTokenPayload) (*dto.RefreshTokenResponsePayload, error) {
	_, err := shared.ValidateRefreshToken(payload.RefreshToken, uc.cfg)
	if err != nil {
		return nil, err
	}

	userID, err := uc.cr.GetUserIdByRefreshToken(ctx, payload.RefreshToken)
	if err != nil {
		return nil, shared.ErrRefreshTokenExpired
	}

	user, err := uc.ar.FirstById(ctx, *userID)
	if err != nil {
		return nil, err
	}

	accessTokenPayload := shared.SignAccessTokenPayload{
		UserID:   user.ID,
		IsSeller: user.IsSeller,
	}

	accessToken, err := shared.GenerateAccessToken(accessTokenPayload, uc.cfg)
	if err != nil {
		return nil, err
	}

	resPayload := dto.RefreshTokenResponsePayload{
		AccessToken: *accessToken,
	}

	return &resPayload, nil

}

// Login implements AuthUsecase.
func (uc *authUsecase) Login(ctx context.Context, payload dto.LoginRequestPayload) (*dto.LoginResponsePayload, error) {
	account, err := uc.ar.FirstByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrInvalidEmailOrPassword
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash.String), []byte(payload.Password))
	if err != nil {
		return nil, shared.ErrInvalidEmailOrPassword
	}

	refreshToken, err := shared.GenerateRefreshToken(uc.cfg)
	if err != nil {
		return nil, err
	}

	accessTokenSignPayload := shared.SignAccessTokenPayload{
		UserID:   account.ID,
		IsSeller: account.IsSeller,
	}
	accessToken, err := shared.GenerateAccessToken(accessTokenSignPayload, uc.cfg)
	if err != nil {
		return nil, err
	}

	responsePayload := &dto.LoginResponsePayload{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	if err := uc.cr.SetRefreshToken(ctx, *refreshToken, account.ID); err != nil {
		return nil, err
	}

	return responsePayload, nil
}

func (uc *authUsecase) LoginWithGoogle(ctx context.Context, googleUser *dto.GoogleResponse) (*dto.LoginResponsePayload, error) {
	email := strings.ToLower(googleUser.Email)
	account, err := uc.ar.FirstByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newEntity := model.Account{
				Username: email,
				Email:    email,
			}

			if err = uc.ar.Create(ctx, newEntity); err != nil {
				return nil, err
			}
			account, err = uc.ar.FirstByEmail(ctx, email)
			if err != nil {
				return nil, err
			}
		}
	}

	refreshToken, err := shared.GenerateRefreshToken(uc.cfg)
	if err != nil {
		return nil, err
	}

	accessTokenSignPayload := shared.SignAccessTokenPayload{
		UserID:   account.ID,
		IsSeller: account.IsSeller,
	}
	accessToken, err := shared.GenerateAccessToken(accessTokenSignPayload, uc.cfg)
	if err != nil {
		return nil, err
	}

	responsePayload := &dto.LoginResponsePayload{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	if err := uc.cr.SetRefreshToken(ctx, *refreshToken, account.ID); err != nil {
		return nil, err
	}

	return responsePayload, nil
}

// RegisterUser implements AuthUsecase.
func (uc *authUsecase) RegisterUser(ctx context.Context, payload dto.RegisterUserRequestPayload) error {
	e, _ := uc.ar.FirstByUsername(ctx, payload.Username)
	if e != nil {
		return shared.ErrUsernameAlreadyTaken
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newEntity := model.Account{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: sql.NullString{String: string(hashedBytes), Valid: true},
	}

	if err = uc.ar.Create(ctx, newEntity); err != nil {
		strErr := err.Error()
		if strings.Contains(strErr, "duplicate") && strings.Contains(strErr, "email") {
			return shared.ErrEmailNotAvailable
		}

		if strings.Contains(strErr, "duplicate") && strings.Contains(strErr, "username") {
			return shared.ErrUsernameNotAvailable
		}
		return err
	}

	return nil
}

func (uc *authUsecase) ForgotPassword(ctx context.Context, payload dto.ForgotPasswordPayload) error {
	user, err := uc.ar.FirstByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.ErrEmailNotFound
		}
		return err
	}

	resetCode := shared.GenerateUUID()
	if err := uc.cr.SetResetPasswordCode(ctx, resetCode, user.ID); err != nil {
		return err
	}

	link := ""
	if uc.cfg.App.OriginDomain == "localhost" {
		link = fmt.Sprintf(constant.ResetPasswordLinkTemplate, "http://localhost", resetCode)
	} else {
		link = fmt.Sprintf(constant.ResetPasswordLinkTemplate, fmt.Sprintf("https://%s/vm1", uc.cfg.App.OriginDomain), resetCode)
	}

	content := fmt.Sprintf(constant.ForgotPwContentTemplate, user.Username, link)

	mail := shared.MakeEmail(uc.cfg.EmailSender.Name, uc.cfg.EmailSender.Address, constant.ForgotPwSubject, content, payload.Email)

	smtpAuth := smtp.PlainAuth("", uc.cfg.EmailSender.Address, uc.cfg.EmailSender.Password, constant.SmtpAuthAddress)
	if err := mail.Send(constant.SmtpServerAddress, smtpAuth); err != nil {
		uc.cr.DeleteResetPasswordCode(ctx, resetCode)
		return err
	}
	return nil
}

func (uc *authUsecase) ResetPassword(ctx context.Context, payload dto.ResetPasswordPayload) error {
	userId, err := uc.cr.GetUserIdByResetPasswordCode(ctx, payload.ResetCode)
	if err != nil {
		return err
	}
	if userId == nil {
		return shared.ErrResetPasswordCodeExpired
	}
	user, err := uc.ar.FirstById(ctx, *userId)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(payload.Password))
	if err == nil {
		return shared.ErrSamePassword
	}

	if strings.Contains(strings.ToLower(payload.Password), strings.ToLower(user.Username)) {
		return shared.ErrPasswordContainsUsername
	}
	re := regexp2.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[a-zA-Z]).{8,}$`, regexp2.None)
	passMatch, err := re.MatchString(payload.Password)
	if err != nil {
		return err
	}
	if !passMatch {
		return shared.ErrPasswordNotMatchRegex
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedBytes)
	if err := uc.ar.UpdatePassword(ctx, *userId, hashedPassword); err != nil {
		return err
	}
	if err := uc.cr.DeleteResetPasswordCode(ctx, payload.ResetCode); err != nil {
		return err
	}
	return nil
}

func (uc *authUsecase) RequestChangePassword(ctx context.Context, userId int64) error {
	user, err := uc.ar.FirstById(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shared.ErrEmailNotFound
		}
		return err
	}
	code, err := uc.cr.GetUserChangePasswordCodeByUserID(ctx, userId)
	if err != nil {
		return err
	}
	if code != nil {
		return shared.ErrChangePasswordExist
	}

	verifCode := shared.GenerateNanoID()
	if err := uc.cr.SetChangePasswordCode(ctx, verifCode, user.ID); err != nil {
		return err
	}

	content := fmt.Sprintf(constant.ChangePwContentTemplate, user.Username, verifCode)
	mail := shared.MakeEmail(uc.cfg.EmailSender.Name, uc.cfg.EmailSender.Address, constant.ChangePwSubject, content, user.Email)

	smtpAuth := smtp.PlainAuth("", uc.cfg.EmailSender.Address, uc.cfg.EmailSender.Password, constant.SmtpAuthAddress)
	if err := mail.Send(constant.SmtpServerAddress, smtpAuth); err != nil {
		uc.cr.DeleteResetPasswordCode(ctx, verifCode)
		return err
	}
	return nil
}

func (uc *authUsecase) ChangePassword(ctx context.Context, payload dto.ChangePasswordPayload, userID int64) error {
	code, err := uc.cr.GetUserChangePasswordCodeByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if code == nil || *code != payload.VerifCode {
		return shared.ErrUnknownVerifCode
	}

	user, err := uc.ar.FirstById(ctx, userID)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(payload.Password))
	if err == nil {
		return shared.ErrSamePassword
	}

	if strings.Contains(strings.ToLower(payload.Password), strings.ToLower(user.Username)) {
		return shared.ErrPasswordContainsUsername
	}
	re := regexp2.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[a-zA-Z]).{8,}$`, regexp2.None)
	passMatch, err := re.MatchString(payload.Password)
	if err != nil {
		return err
	}
	if !passMatch {
		return shared.ErrPasswordNotMatchRegex
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedBytes)
	if err := uc.ar.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return err
	}
	if err := uc.cr.DeleteChangePasswordCode(ctx, userID); err != nil {
		return err
	}
	return nil
}

func NewAuthUsecase(
	ar repository.AccountRepository,
	cr repository.CacheRepository,
	cartRepo repository.CartRepository,
	er repository.WalletRepository,
	cer repository.ChangedEmailRepository,
	sr repository.ShopRepository,
	cfg dependency.Config,
) AuthUsecase {
	return &authUsecase{
		ar:       ar,
		cfg:      cfg,
		cr:       cr,
		cartRepo: cartRepo,
		er:       er,
		sr:       sr,
		cer:      cer,
	}
}
