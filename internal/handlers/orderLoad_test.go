package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hordevcom/GopherDiploma/internal/service"
)

type mockOrderGetter struct {
	getOrderFunc func(ctx context.Context, orderID string) (order, username string, err error)
}

func (m *mockOrderGetter) GetOrderAndUser(ctx context.Context, orderID string) (order, username string, err error) {
	return m.getOrderFunc(ctx, orderID)
}

type mockOrderAdder struct {
	setOrderFunc func(ctx context.Context, orderID, username string) error
}

func (m *mockOrderAdder) AddOrderToDB(ctx context.Context, orderID, username string) error {
	return m.setOrderFunc(ctx, orderID, username)
}

func TestNewOrderLoad(t *testing.T) {
	type args struct {
		accrualAddress string
		serv           service.Service
	}
	tests := []struct {
		name               string
		args               args
		requestBody        string
		expectedStatusCode int
	}{
		{
			name: "test status 409",
			args: args{
				accrualAddress: "localhost:8081",
				serv: service.Service{
					Adder: &mockOrderAdder{
						setOrderFunc: func(ctx context.Context, orderID, username string) error {
							return nil
						},
					},
					Getter: &mockOrderGetter{
						getOrderFunc: func(ctx context.Context, orderID string) (order, username string, err error) {
							return "12345678903", "alusha", nil
						},
					},
				},
			},
			requestBody:        `12345678903`,
			expectedStatusCode: http.StatusConflict,
		},
		{
			name: "test status 422",
			args: args{
				accrualAddress: "localhost:8081",
				serv: service.Service{
					Adder: &mockOrderAdder{
						setOrderFunc: func(ctx context.Context, orderID, username string) error {
							return nil
						},
					},
					Getter: &mockOrderGetter{
						getOrderFunc: func(ctx context.Context, orderID string) (order, username string, err error) {
							return "12345678903", "alusha", nil
						},
					},
				},
			},
			requestBody:        `12345678902`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "test status 500",
			args: args{
				accrualAddress: "localhost:8081",
				serv: service.Service{
					Adder: &mockOrderAdder{
						setOrderFunc: func(ctx context.Context, orderID, username string) error {
							return errors.New("something went wrong!")
						},
					},
					Getter: &mockOrderGetter{
						getOrderFunc: func(ctx context.Context, orderID string) (order, username string, err error) {
							return "12345678901", "alusha", nil
						},
					},
				},
			},
			requestBody:        `12345678903`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewOrderLoad(tt.args.accrualAddress, tt.args.serv)

			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(tt.requestBody))

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
