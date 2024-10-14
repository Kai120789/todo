package services

import (
	"todo/internal/todo/dto"
	"todo/internal/todo/models"
	"todo/internal/todo/utils/hash"

	"go.uber.org/zap"
)

type UserService struct {
	storage UserStorager
}

type UserStorager interface {
	RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error)
	AuthorizateUser(body dto.PostUserDto) (*uint, *string, error)
	WriteRefreshToken(userId uint, refreshTokenValue string) error
	GetAuthUser(id uint) (*models.UserToken, error)
	UserLogout(id uint) error
	AddChatID(tgName string, chatID int64) error
}

func NewUserService(stor UserStorager, logger *zap.Logger) *UserService {
	return &UserService{
		storage: stor,
	}
}

func (t *UserService) RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error) {
	token, err := t.storage.RegisterNewUser(body)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *UserService) AuthorizateUser(body dto.PostUserDto) (*uint, error) {
	id, passwordHash, err := t.storage.AuthorizateUser(body)
	if err != nil {
		return nil, err
	}

	if !hash.CheckPasswordHash(body.PasswordHash, *passwordHash) {
		return nil, err
	}

	return id, nil
}

func (t *UserService) GetAuthUser(id uint) (*models.UserToken, error) {
	token, err := t.storage.GetAuthUser(uint(id))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *UserService) UserLogout(id uint) error {
	err := t.storage.UserLogout(uint(id))
	if err != nil {
		return err
	}

	return nil
}

func (t *UserService) WriteRefreshToken(userId uint, refreshTokenValue string) error {
	err := t.storage.WriteRefreshToken(userId, refreshTokenValue)
	if err != nil {
		return err
	}

	return nil
}

func (t *UserService) AddChatID(tgName string, chatID int64) error {
	err := t.storage.AddChatID(tgName, chatID)
	if err != nil {
		return err
	}

	return nil
}
