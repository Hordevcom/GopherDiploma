package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hordevcom/GopherDiploma/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type mockPasswordGetter struct {
	getUserPasswordFunc func(ctx context.Context, username string) string
}

func (m *mockPasswordGetter) GetUserPassword(ctx context.Context, username string) string {
	return m.getUserPasswordFunc(ctx, username)
}

// type mockCheckUserLogin struct {
// 	checkUsernameLoginFunc func(ctx context.Context, username string) bool
// 	addUserToDBFunc        func(ctx context.Context, username, password string) error
// }

// func (m *mockCheckUserLogin) CheckUsernameLogin(ctx context.Context, username string) bool {
// 	return m.checkUsernameLoginFunc(ctx, username)
// }

// func (m *mockCheckUserLogin) AddUserToDB(ctx context.Context, username, password string) error {
// 	return m.addUserToDBFunc(ctx, username, password)
// }

// --- Тест --- //

func TestNewUserLogin_Unauthorized(t *testing.T) {
	type args struct {
		serv service.Service
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)

	tests := []struct {
		name               string
		args               args
		requestBody        string
		expectedStatusCode int
	}{
		{
			name: "unauthorized login attempt",
			args: args{
				serv: service.Service{
					PassGetter: &mockPasswordGetter{
						getUserPasswordFunc: func(ctx context.Context, username string) string {
							return string(hashedPassword)
						},
					},
					// UserChecker: &mockCheckUserLogin{
					// 	checkUsernameLoginFunc: func(ctx context.Context, username string) bool {
					// 		return true
					// 	},
					// 	addUserToDBFunc: func(ctx context.Context, username, password string) error {
					// 		return nil
					// 	},
					// },
				},
			},
			requestBody:        `{"login":"user1", "password":"wrong_password"}`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "authorized login attempt",
			args: args{
				serv: service.Service{
					PassGetter: &mockPasswordGetter{
						getUserPasswordFunc: func(ctx context.Context, username string) string {
							hash, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
							return string(hash)
						},
					},
					// UserChecker: &mockCheckUserLogin{
					// 	checkUsernameLoginFunc: func(ctx context.Context, username string) bool {
					// 		return true
					// 	},
					// 	addUserToDBFunc: func(ctx context.Context, username, password string) error {
					// 		return nil
					// 	},
					// },
				},
			},
			requestBody:        `{"login":"user1", "password":"correct_password"}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewUserLogin(tt.args.serv)

			req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			fmt.Println("StatusCode: ", resp.StatusCode)

			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}
