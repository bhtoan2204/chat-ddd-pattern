package repos

import gomock "go.uber.org/mock/gomock"

// MockDeviceRepository preserves compatibility for callers that use the legacy
// DeviceRepository alias while the aggregate repository remains the canonical
// persistence boundary.
type MockDeviceRepository = MockDeviceAggregateRepository

type MockDeviceRepositoryMockRecorder = MockDeviceAggregateRepositoryMockRecorder

func NewMockDeviceRepository(ctrl *gomock.Controller) *MockDeviceRepository {
	return NewMockDeviceAggregateRepository(ctrl)
}

// MockSessionRepository preserves compatibility for callers that use the legacy
// SessionRepository alias while the aggregate repository remains the canonical
// persistence boundary.
type MockSessionRepository = MockSessionAggregateRepository

type MockSessionRepositoryMockRecorder = MockSessionAggregateRepositoryMockRecorder

func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	return NewMockSessionAggregateRepository(ctrl)
}
