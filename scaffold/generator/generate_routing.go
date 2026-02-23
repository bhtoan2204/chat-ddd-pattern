package generator

import (
	"errors"

	"go-socket/scaffold/models"
)

func GenerateRouting(spec *models.APISpec) (string, error) {
	if spec == nil {
		return "", errors.New("api spec is nil")
	}
	if len(spec.Endpoints) == 0 {
		return "", errors.New("no endpoints to generate routing")
	}
	return "skipped routing generation (routes live in core/*/transport/http/routes.go)", nil
}
