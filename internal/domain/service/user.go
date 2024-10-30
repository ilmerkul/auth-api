package service

import (
	"auth-api/internal/domain/entity"
	v1 "auth-api/internal/transport/http/v1/dto"
	auth "auth-api/pkg/auth/jwt"
	"auth-api/pkg/cache"
	"auth-api/pkg/cache/mapCache"
	"auth-api/pkg/email"
	"auth-api/pkg/email/smtp"
	"auth-api/pkg/hash"
	"auth-api/pkg/rand"
	"fmt"
	"log"
	"time"
)

const (
	ttlSignUpUser = 60 * 15
)

type UserStorage interface {
	GetByID(id int) (user entity.User, err error)
	Create(user entity.User) (id int64, err error)
}

type userService struct {
	userStorage           UserStorage
	refreshSessionService refreshSessionService

	emailSender  email.Sender
	tokenManager auth.TokenManager
	cache        cache.Cache
	hasher       hash.PasswordHasher
}

func NewUserService(UserStorage UserStorage, senderConfig *smtp.SenderConfig, managerTokenConfig *auth.TokenManagerConfig) *userService {
	emailSender, err := smtp.NewSMTPSender(senderConfig)
	if err != nil {
		log.Fatal("error init email sender")
	}
	manager, err := auth.NewManager(managerTokenConfig.SigningKey)
	if err != nil {
		log.Fatal("error init token manager")
	}
	return &userService{userStorage: UserStorage, emailSender: emailSender, tokenManager: manager, cache: mapCache.NewMemoryCache(), hasher: hash.NewSHA1Hasher(managerTokenConfig.Salt)}
}

func (s *userService) SignUp(user v1.UserSignUp) error {
	code, err := rand.GetRandString()
	if err != nil {
		return err
	}
	user.Password, err = s.hasher.Hash(user.Password)
	if err != nil {
		return err
	}

	// TODO link to address
	msg := email.SendEmailInput{
		To:      user.Email,
		Subject: "Подтверждение почты",
		Body:    fmt.Sprintf("https://address//?code=%s", code),
	}

	if err = s.emailSender.Send(msg); err != nil {
		return err
	}

	return s.cache.Set(code, user, ttlSignUpUser)
}

func (s *userService) ConfirmEmail(code string) (id int64, err error) {
	val, err := s.cache.Get(code)
	if err != nil {
		return
	}

	userSignUp := val.(v1.UserSignUp)

	// TODO: new refreshSession

	user := entity.User{
		Email:     userSignUp.Email,
		Password:  userSignUp.Password,
		CreatedAt: time.Now().Format("02.01.2006 15:04:05"),
	}

	return s.userStorage.Create(user)
}

func (s *userService) NewRefreshSession() (v1.ResponseRefreshSession, error) {
	return nil, nil
}
