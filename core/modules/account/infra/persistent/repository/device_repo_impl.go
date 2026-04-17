package repos

import (
	"context"
	"fmt"

	"wechat-clone/core/modules/account/domain/aggregate"
	"wechat-clone/core/modules/account/domain/entity"
	accountrepos "wechat-clone/core/modules/account/domain/repos"
	"wechat-clone/core/modules/account/infra/persistent/models"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/utils"

	"gorm.io/gorm"
)

type deviceRepoImpl struct {
	db *gorm.DB
}

func NewDeviceRepoImpl(db *gorm.DB) accountrepos.DeviceRepository {
	return &deviceRepoImpl{db: db}
}

func (r *deviceRepoImpl) FindByAccountAndUID(ctx context.Context, accountID string, deviceUID string) (*aggregate.DeviceAggregate, error) {
	var model models.DeviceModel
	if err := r.db.WithContext(ctx).
		Where("account_id = ? AND device_uid = ?", accountID, deviceUID).
		First(&model).Error; err != nil {
		return nil, stackErr.Error(err)
	}

	device, err := r.toEntity(&model)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return r.toAggregate(device)
}

func (r *deviceRepoImpl) GetByAccountAndID(ctx context.Context, accountID string, deviceID string) (*aggregate.DeviceAggregate, error) {
	var model models.DeviceModel
	if err := r.db.WithContext(ctx).
		Where("account_id = ? AND id = ?", accountID, deviceID).
		First(&model).Error; err != nil {
		return nil, stackErr.Error(err)
	}

	device, err := r.toEntity(&model)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return r.toAggregate(device)
}

func (r *deviceRepoImpl) Save(ctx context.Context, device *aggregate.DeviceAggregate) error {
	if device == nil {
		return stackErr.Error(fmt.Errorf("device is nil"))
	}

	snapshot, err := device.Snapshot()
	if err != nil {
		return stackErr.Error(err)
	}

	if err := r.db.WithContext(ctx).Save(r.toModel(snapshot)).Error; err != nil {
		return stackErr.Error(err)
	}
	return nil
}

func (r *deviceRepoImpl) toAggregate(device *entity.Device) (*aggregate.DeviceAggregate, error) {
	agg, err := aggregate.NewDeviceAggregate(device.ID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if err := agg.Restore(device); err != nil {
		return nil, stackErr.Error(err)
	}
	return agg, nil
}

func (r *deviceRepoImpl) toEntity(model *models.DeviceModel) (*entity.Device, error) {
	if model == nil {
		return nil, stackErr.Error(fmt.Errorf("device model is nil"))
	}

	return &entity.Device{
		ID:            model.ID,
		AccountID:     model.AccountID,
		DeviceUID:     model.DeviceUID,
		DeviceName:    utils.ClonePtr(model.DeviceName),
		DeviceType:    entity.DeviceType(model.DeviceType),
		OSName:        utils.ClonePtr(model.OSName),
		OSVersion:     utils.ClonePtr(model.OSVersion),
		AppVersion:    utils.ClonePtr(model.AppVersion),
		UserAgent:     utils.ClonePtr(model.UserAgent),
		LastIPAddress: utils.ClonePtr(model.LastIPAddress),
		LastSeenAt:    utils.ClonePtr(model.LastSeenAt),
		IsTrusted:     model.IsTrusted == 1,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *deviceRepoImpl) toModel(device *entity.Device) *models.DeviceModel {
	isTrusted := int8(0)
	if device.IsTrusted {
		isTrusted = 1
	}

	return &models.DeviceModel{
		ID:            device.ID,
		AccountID:     device.AccountID,
		DeviceUID:     device.DeviceUID,
		DeviceName:    utils.ClonePtr(device.DeviceName),
		DeviceType:    device.DeviceType.String(),
		OSName:        utils.ClonePtr(device.OSName),
		OSVersion:     utils.ClonePtr(device.OSVersion),
		AppVersion:    utils.ClonePtr(device.AppVersion),
		UserAgent:     utils.ClonePtr(device.UserAgent),
		LastIPAddress: utils.ClonePtr(device.LastIPAddress),
		LastSeenAt:    utils.ClonePtr(device.LastSeenAt),
		IsTrusted:     isTrusted,
		CreatedAt:     device.CreatedAt,
		UpdatedAt:     device.UpdatedAt,
	}
}
