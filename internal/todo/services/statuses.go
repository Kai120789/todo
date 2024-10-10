package services

import (
	"fmt"
	"strconv"
	"todo/internal/todo/dto"

	"go.uber.org/zap"
)

type StatusesService struct {
	storage StatusesStorager
}

type StatusesStorager interface {
	SetStatus(body dto.PostStatusDto) error
	DeleteStatus(id uint) error
}

func NewStatusesService(stor StatusesStorager, logger *zap.Logger) *StatusesService {
	return &StatusesService{
		storage: stor,
	}
}

func (t *StatusesService) SetStatus(body dto.PostStatusDto) error {
	err := t.storage.SetStatus(body)
	if err != nil {
		return err
	}

	return nil
}

func (t *StatusesService) DeleteStatus(id string) error {
	if id == "" {
		return fmt.Errorf("status ID is required")
	}

	Uintid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	err = t.storage.DeleteStatus(uint(Uintid))
	if err != nil {
		return err
	}

	return nil
}
