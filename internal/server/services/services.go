package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
)

func RegistService(r *http.Request, auth types.Authorization) (int, string, error) {
	var (
		userData types.UserData
		token    string
	)
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		return http.StatusBadRequest, token, fmt.Errorf("can't decode body %w", err)
	}
	if err := auth.CheckData(userData); err != nil {
		return http.StatusBadRequest, token, fmt.Errorf("no data provided: %w", err)
	}
	user, err := auth.RegisterUser(userData)
	if err != nil && !errors.Is(err, types.ErrInvalidData) {
		return http.StatusLoopDetected, token, fmt.Errorf("RegistHandler: %w", err)
	}
	if errors.Is(err, types.ErrInvalidData) {
		return http.StatusUnauthorized, token, fmt.Errorf("RegistHandler: %w", err)
	}
	token, err = auth.GenerateToken(user)
	if err != nil {
		return http.StatusInternalServerError, token, fmt.Errorf("RegistHandler: can't generate token %w", err)
	}

	return http.StatusOK, token, nil
}

func AuthService(r *http.Request, auth types.Authorization) (int, string, error) {
	var (
		userData types.UserData
		token    string
	)
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		return http.StatusBadRequest, token, err
	}
	if err := auth.CheckData(userData); err != nil {
		return http.StatusBadRequest, token, err
	}
	user, err := auth.LoginUser(userData)
	if err != nil && !errors.Is(err, types.ErrInvalidData) {
		return http.StatusInternalServerError, token, err
	}
	if errors.Is(err, types.ErrInvalidData) {
		return http.StatusUnauthorized, token, err
	}
	token, err = auth.GenerateToken(user)
	if err != nil {
		return http.StatusInternalServerError, token, err
	}

	return http.StatusOK, token, err
}
