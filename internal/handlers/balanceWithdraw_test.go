package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/service"
)

type mockBalanceUpdater struct {
	updateUserBalanceFunc func(ctx context.Context, user string, accrual, withdrawn float32) error
	setUserWithdrawnFunc  func(ctx context.Context, orderNum, user string, withdrawn float32) error
}

func (m *mockBalanceUpdater) UpdateUserBalance(ctx context.Context, user string, accrual, withdrawn float32) error {
	return m.updateUserBalanceFunc(ctx, user, accrual, withdrawn)
}

func (m *mockBalanceUpdater) SetUserWithdrawn(ctx context.Context, orderNum, user string, withdrawn float32) error {
	return m.setUserWithdrawnFunc(ctx, orderNum, user, withdrawn)
}

func TestNewBalanceWithdrawn(t *testing.T) {
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
			name: "test status 400",
			args: args{
				serv: service.Service{
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{}, nil
						},
					},
					BalanceUpdater: &mockBalanceUpdater{
						updateUserBalanceFunc: func(ctx context.Context, user string, accrual, withdrawn float32) error {
							return nil
						},
						setUserWithdrawnFunc: func(ctx context.Context, orderNum, user string, withdrawn float32) error {
							return nil
						},
					},
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "test status 402",
			args: args{
				serv: service.Service{
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{
								Current: 100,
							}, nil
						},
					},
					BalanceUpdater: &mockBalanceUpdater{
						updateUserBalanceFunc: func(ctx context.Context, user string, accrual, withdrawn float32) error {
							return nil
						},
						setUserWithdrawnFunc: func(ctx context.Context, orderNum, user string, withdrawn float32) error {
							return nil
						},
					},
				},
			},
			requestBody:        `{"order": "2131231", "sum":751}`,
			expectedStatusCode: http.StatusPaymentRequired,
		},
		{
			name: "test status 200",
			args: args{
				serv: service.Service{
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{
								Current: 1000,
							}, nil
						},
					},
					BalanceUpdater: &mockBalanceUpdater{
						updateUserBalanceFunc: func(ctx context.Context, user string, accrual, withdrawn float32) error {
							return nil
						},
						setUserWithdrawnFunc: func(ctx context.Context, orderNum, user string, withdrawn float32) error {
							return nil
						},
					},
				},
			},
			requestBody:        `{"order": "2131231", "sum":51}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "test status 500",
			args: args{
				serv: service.Service{
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{
								Current: 1000,
							}, errors.New("U")
						},
					},
					BalanceUpdater: &mockBalanceUpdater{
						updateUserBalanceFunc: func(ctx context.Context, user string, accrual, withdrawn float32) error {
							return nil
						},
						setUserWithdrawnFunc: func(ctx context.Context, orderNum, user string, withdrawn float32) error {
							return nil
						},
					},
				},
			},
			requestBody:        `{"order": "2131231", "sum":51}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBalanceWithdrawn(tt.args.serv)

			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(tt.requestBody))

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
