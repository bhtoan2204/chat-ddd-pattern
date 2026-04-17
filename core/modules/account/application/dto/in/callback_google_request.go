// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"go-socket/core/shared/pkg/stackErr"
	"strings"
)

type CallbackGoogleRequest struct {
	Code       string `json:"code" form:"code" binding:"required"`
	State      string `json:"state" form:"state" binding:"required"`
	DeviceUid  string `json:"device_uid" form:"device_uid" binding:"required"`
	DeviceName string `json:"device_name" form:"device_name"`
	DeviceType string `json:"device_type" form:"device_type"`
	OsName     string `json:"os_name" form:"os_name"`
	OsVersion  string `json:"os_version" form:"os_version"`
	AppVersion string `json:"app_version" form:"app_version"`
	UserAgent  string `json:"user_agent" form:"user_agent"`
	IpAddress  string `json:"ip_address" form:"ip_address"`
}

func (r *CallbackGoogleRequest) Normalize() {
	r.Code = strings.TrimSpace(r.Code)
	r.State = strings.TrimSpace(r.State)
	r.DeviceUid = strings.TrimSpace(r.DeviceUid)
	r.DeviceName = strings.TrimSpace(r.DeviceName)
	r.DeviceType = strings.TrimSpace(r.DeviceType)
	r.OsName = strings.TrimSpace(r.OsName)
	r.OsVersion = strings.TrimSpace(r.OsVersion)
	r.AppVersion = strings.TrimSpace(r.AppVersion)
	r.UserAgent = strings.TrimSpace(r.UserAgent)
	r.IpAddress = strings.TrimSpace(r.IpAddress)
}

func (r *CallbackGoogleRequest) Validate() error {
	r.Normalize()
	if r.Code == "" {
		return stackErr.Error(errors.New("code is required"))
	}
	if r.State == "" {
		return stackErr.Error(errors.New("state is required"))
	}
	if r.DeviceUid == "" {
		return stackErr.Error(errors.New("device_uid is required"))
	}
	return nil
}
