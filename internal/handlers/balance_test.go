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

func TestNewBalance(t *testing.T) {
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
			name: "test status 200",
			args: args{
				serv: service.Service{
					GetBalance: &mockUserBalance{
						getUserBalance: func(ctx context.Context, user string) (models.UserBalance, error) {
							return models.UserBalance{}, nil
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
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewBalance(tt.args.serv)

			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", bytes.NewBufferString(tt.requestBody))

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
