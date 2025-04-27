package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/service"
)

type mockUserWithdrawnsGetter struct {
	getUserWithdrawnsFunc func(ctx context.Context, user string) ([]models.UserWithdrawal, error)
}

func (m *mockUserWithdrawnsGetter) GetUserWithdrawns(ctx context.Context, user string) ([]models.UserWithdrawal, error) {
	return m.getUserWithdrawnsFunc(ctx, user)
}

func TestNewWithdraw(t *testing.T) {
	type args struct {
		serv service.Service
	}
	tests := []struct {
		name               string
		args               args
		requestBody        string
		expectedStatusCode int
	}{
		{
			name: "test status 204",
			args: args{
				serv: service.Service{
					UserWithdrawnsGetter: &mockUserWithdrawnsGetter{
						getUserWithdrawnsFunc: func(ctx context.Context, user string) ([]models.UserWithdrawal, error) {
							return []models.UserWithdrawal{}, nil
						},
					},
				},
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "test status 200",
			args: args{
				serv: service.Service{
					UserWithdrawnsGetter: &mockUserWithdrawnsGetter{
						getUserWithdrawnsFunc: func(ctx context.Context, user string) ([]models.UserWithdrawal, error) {
							return []models.UserWithdrawal{
								models.UserWithdrawal{
									OrderNum:    "1",
									Sum:         23,
									ProcessedAt: time.Now(),
								},
							}, nil
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "test status 500",
			args: args{
				serv: service.Service{
					UserWithdrawnsGetter: &mockUserWithdrawnsGetter{
						getUserWithdrawnsFunc: func(ctx context.Context, user string) ([]models.UserWithdrawal, error) {
							return []models.UserWithdrawal{}, errors.New("U")
						},
					},
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewWithdraw(tt.args.serv)

			req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", bytes.NewBufferString(tt.requestBody))

			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: "dummy_token",
			})
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}
