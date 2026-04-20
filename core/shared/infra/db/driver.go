package db

import (
	"fmt"
	"strings"

	"wechat-clone/core/shared/pkg/stackErr"
)

const (
	DriverPostgres   = "postgres"
	DriverPostgreSQL = "postgresql"
)

func normalizedDriverNameOrError(driver string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "", DriverPostgres, DriverPostgreSQL:
		return DriverPostgres, nil
	default:
		return "", stackErr.Error(fmt.Errorf("unsupported db driver: %s", driver))
	}
}

func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToUpper(err.Error())
	return strings.Contains(msg, "SQLSTATE 23505") || strings.Contains(msg, "DUPLICATE KEY VALUE")
}

func isObjectExistsError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToUpper(err.Error())
	return strings.Contains(msg, "SQLSTATE 42P07") || strings.Contains(msg, "ALREADY EXISTS")
}
