package generator

import (
	"errors"

	"go-socket/scaffold/models"
)

func GenerateRegistry(spec *models.APISpec) (string, error) {
	if spec == nil {
		return "", errors.New("api spec is nil")
	}
	if len(spec.Endpoints) == 0 {
		return "", errors.New("no endpoints to generate registry")
	}
	return "skipped registry generation (routing is now per module)", nil
}
