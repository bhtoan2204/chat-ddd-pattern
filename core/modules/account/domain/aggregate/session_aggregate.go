package aggregate

import (
	"errors"
	"strings"
	"time"

	"wechat-clone/core/modules/account/domain/entity"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/utils"
)

var ErrSessionAggregateNotInitialized = errors.New("session aggregate is not initialized")

type SessionAggregate struct {
	session *entity.Session
}

func NewSessionAggregate(sessionID string) (*SessionAggregate, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, stackErr.Error(ErrSessionAggregateNotInitialized)
	}
	return &SessionAggregate{
		session: &entity.Session{ID: sessionID},
	}, nil
}

func (a *SessionAggregate) Restore(snapshot *entity.Session) error {
	if snapshot == nil {
		return stackErr.Error(ErrSessionAggregateNotInitialized)
	}

	cloned := *snapshot
	cloned.IPAddress = utils.ClonePtr(snapshot.IPAddress)
	cloned.UserAgent = utils.ClonePtr(snapshot.UserAgent)
	cloned.LastActivityAt = utils.ClonePtr(snapshot.LastActivityAt)
	cloned.RevokedAt = utils.ClonePtr(snapshot.RevokedAt)
	cloned.RevokedReason = utils.ClonePtr(snapshot.RevokedReason)
	a.session = &cloned
	return nil
}

func (a *SessionAggregate) Create(
	accountID string,
	deviceID string,
	refreshTokenHash string,
	expiresAt time.Time,
	now time.Time,
	ipAddress string,
	userAgent string,
) error {
	if a == nil || a.session == nil || strings.TrimSpace(a.session.ID) == "" {
		return stackErr.Error(ErrSessionAggregateNotInitialized)
	}

	session, err := entity.NewSession(
		a.session.ID,
		accountID,
		deviceID,
		refreshTokenHash,
		expiresAt,
		now,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return stackErr.Error(err)
	}
	a.session = session
	return nil
}

func (a *SessionAggregate) EnsureRefreshAllowed(now time.Time) error {
	if a == nil || a.session == nil {
		return stackErr.Error(ErrSessionAggregateNotInitialized)
	}
	return stackErr.Error(a.session.EnsureRefreshAllowed(now))
}

func (a *SessionAggregate) Rotate(
	refreshTokenHash string,
	expiresAt time.Time,
	now time.Time,
	ipAddress string,
	userAgent string,
) error {
	if a == nil || a.session == nil {
		return stackErr.Error(ErrSessionAggregateNotInitialized)
	}
	return stackErr.Error(a.session.Rotate(refreshTokenHash, expiresAt, now, ipAddress, userAgent))
}

func (a *SessionAggregate) Revoke(reason string, now time.Time) (bool, error) {
	if a == nil || a.session == nil {
		return false, stackErr.Error(ErrSessionAggregateNotInitialized)
	}
	return a.session.Revoke(reason, now), nil
}

func (a *SessionAggregate) MarkExpired(now time.Time) bool {
	if a == nil || a.session == nil {
		return false
	}
	return a.session.MarkExpired(now)
}

func (a *SessionAggregate) Snapshot() (*entity.Session, error) {
	if a == nil || a.session == nil {
		return nil, stackErr.Error(ErrSessionAggregateNotInitialized)
	}

	cloned := *a.session
	cloned.IPAddress = utils.ClonePtr(a.session.IPAddress)
	cloned.UserAgent = utils.ClonePtr(a.session.UserAgent)
	cloned.LastActivityAt = utils.ClonePtr(a.session.LastActivityAt)
	cloned.RevokedAt = utils.ClonePtr(a.session.RevokedAt)
	cloned.RevokedReason = utils.ClonePtr(a.session.RevokedReason)
	return &cloned, nil
}

func (a *SessionAggregate) SessionID() string {
	if a == nil || a.session == nil {
		return ""
	}
	return a.session.ID
}

func (a *SessionAggregate) AccountID() string {
	if a == nil || a.session == nil {
		return ""
	}
	return a.session.AccountID
}

func (a *SessionAggregate) DeviceID() string {
	if a == nil || a.session == nil {
		return ""
	}
	return a.session.DeviceID
}

func (a *SessionAggregate) RefreshTokenHash() string {
	if a == nil || a.session == nil {
		return ""
	}
	return a.session.RefreshTokenHash
}
