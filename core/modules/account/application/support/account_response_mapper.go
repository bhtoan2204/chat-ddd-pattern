package support

import (
	"time"

	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/domain/entity"
)

func ToGetProfileResponse(account *entity.Account) *out.GetProfileResponse {
	if account == nil {
		return nil
	}

	return &out.GetProfileResponse{
		ID:                account.ID,
		DisplayName:       account.DisplayName,
		Email:             account.Email.Value(),
		Username:          stringValue(account.Username),
		AvatarObjectKey:   stringValue(account.AvatarObjectKey),
		Status:            account.Status,
		EmailVerified:     account.EmailVerifiedAt != nil,
		EmailVerifiedAt:   formatOptionalTime(account.EmailVerifiedAt),
		LastLoginAt:       formatOptionalTime(account.LastLoginAt),
		PasswordChangedAt: formatOptionalTime(account.PasswordChangedAt),
		CreatedAt:         account.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:         account.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func ToUpdateProfileResponse(account *entity.Account) *out.UpdateProfileResponse {
	if account == nil {
		return nil
	}

	return &out.UpdateProfileResponse{
		ID:                account.ID,
		DisplayName:       account.DisplayName,
		Email:             account.Email.Value(),
		Username:          stringValue(account.Username),
		AvatarObjectKey:   stringValue(account.AvatarObjectKey),
		Status:            account.Status,
		EmailVerified:     account.EmailVerifiedAt != nil,
		EmailVerifiedAt:   formatOptionalTime(account.EmailVerifiedAt),
		LastLoginAt:       formatOptionalTime(account.LastLoginAt),
		PasswordChangedAt: formatOptionalTime(account.PasswordChangedAt),
		CreatedAt:         account.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:         account.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func formatOptionalTime(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
