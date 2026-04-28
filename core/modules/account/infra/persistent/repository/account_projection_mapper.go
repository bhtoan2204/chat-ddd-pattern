package repos

import (
	"wechat-clone/core/modules/account/domain/entity"
	valueobject "wechat-clone/core/modules/account/domain/value_object"
	"wechat-clone/core/modules/account/infra/persistent/models"
	accounttypes "wechat-clone/core/modules/account/types"
	"wechat-clone/core/shared/pkg/stackErr"
)

func projectionModelToAccount(m *models.AccountModel) (*entity.Account, error) {
	email, err := valueobject.NewEmail(m.Email)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	passwordHash, err := valueobject.NewHashedPassword(m.Password)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	status, err := accounttypes.ParseAccountStatus(m.Status)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return &entity.Account{
		ID:                m.ID,
		Email:             email,
		PasswordHash:      passwordHash,
		DisplayName:       m.DisplayName,
		Username:          m.Username,
		AvatarObjectKey:   m.AvatarObjectKey,
		Status:            status,
		EmailVerifiedAt:   m.EmailVerifiedAt,
		LastLoginAt:       m.LastLoginAt,
		PasswordChangedAt: m.PasswordChangedAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		BannedReason:      m.BannedReason,
		BannedUntil:       m.BannedUntil,
	}, nil
}

func accountToProjectionModel(e *entity.Account) *models.AccountModel {
	return &models.AccountModel{
		ID:                e.ID,
		Email:             e.Email.Value(),
		Password:          e.PasswordHash.Value(),
		DisplayName:       e.DisplayName,
		Username:          e.Username,
		AvatarObjectKey:   e.AvatarObjectKey,
		Status:            e.Status.String(),
		EmailVerifiedAt:   e.EmailVerifiedAt,
		LastLoginAt:       e.LastLoginAt,
		PasswordChangedAt: e.PasswordChangedAt,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
		BannedReason:      e.BannedReason,
		BannedUntil:       e.BannedUntil,
	}
}
