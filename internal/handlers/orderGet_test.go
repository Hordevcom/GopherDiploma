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

type mockGetOrder struct {
	getOrderFunc func(ctx context.Context, user string) ([]models.Order, error)
}

func (m *mockGetOrder) GetUserOrders(ctx context.Context, user string) ([]models.Order, error) {
	return m.getOrderFunc(ctx, user)
}

type mockUserBalance struct {
	getUserBalance func(ctx context.Context, user string) (models.UserBalance, error)
}

func (m *mockUserBalance) GetUserBalance(ctx context.Context, user string) (models.UserBalance, error) {
	return m.getUserBalance(ctx, user)
}

func TestNewOrderGet(t *testing.T) {
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
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{}, nil
						},
					},
					GetOrder: &mockGetOrder{
						getOrderFunc: func(ctx context.Context, user string) ([]models.Order, error) {
							return []models.Order{}, nil
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
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{}, nil
						},
					},
					GetOrder: &mockGetOrder{
						getOrderFunc: func(ctx context.Context, user string) ([]models.Order, error) {
							return []models.Order{
								models.Order{
									Number:   "1",
									Status:   "A",
									Accrual:  200,
									UploadAt: time.Now(),
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
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{}, errors.New("U")
						},
					},
					GetOrder: &mockGetOrder{
						getOrderFunc: func(ctx context.Context, user string) ([]models.Order, error) {
							return []models.Order{}, nil
						},
					},
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewOrderGet(tt.args.serv)

			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
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
