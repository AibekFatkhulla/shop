package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aibekfatkhulla/shop/internal/domain"
	internalMock "github.com/aibekfatkhulla/shop/internal/mocks"
	"github.com/aibekfatkhulla/shop/internal/server"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServer_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		user         server.UserDTO
		svc          server.Service
		expectedCode int
		expectedBody []byte
	}{
		{
			"success case",
			server.UserDTO{
				Name:     "arnur",
				Password: "qwe",
				Email:    "qwe@qwe.qwe",
				Number:   "123",
				Address:  "123",
				Balance:  100,
			},
			func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)

				return s
			}(),
			http.StatusOK,
			nil,
		}, {
			"user already exists",
			server.UserDTO{
				Name:     "arnur",
				Password: "qwe",
				Email:    "qwe@qwe.qwe",
				Number:   "123",
				Address:  "123",
				Balance:  100,
			},
			func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(domain.ErrorUserAlreadyExists)

				return s
			}(),
			http.StatusConflict,
			func() []byte {
				b, err := json.Marshal(map[string]string{"error": domain.ErrorUserAlreadyExists.Error()})
				assert.NoError(t, err)

				return b
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqJSON, err := json.Marshal(tt.user)
			assert.NoError(t, err)

			s := server.NewServer(tt.svc)

			r := s.SetupRouter()

			w := httptest.NewRecorder()

			body := bytes.NewBuffer(reqJSON)

			req, err := http.NewRequest("POST", "/users", body)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.Bytes())
		})
	}
}

func TestServer_AddProductToCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		categoryID   string
		productID    string
		svc          server.Service
		expectedCode int
		expectedBody []byte
	}{
		{
			name:       "success case",
			categoryID: "cat123",
			productID:  "prod456",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					AddProductToCategory(gomock.Any(), "cat123", "prod456").
					Return(nil)
				return s
			}(),
			expectedCode: http.StatusOK,
			expectedBody: func() []byte {
				b, _ := json.Marshal(map[string]string{"message": "product added to category"})
				return b
			}(),
		},
		{
			name:       "category not found",
			categoryID: "notfound",
			productID:  "prod456",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					AddProductToCategory(gomock.Any(), "notfound", "prod456").
					Return(domain.ErrorCategoryNotFound)
				return s
			}(),
			expectedCode: http.StatusNotFound,
			expectedBody: func() []byte {
				b, _ := json.Marshal(map[string]string{"error": "category not found"})
				return b
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.NewServer(tt.svc)
			r := s.SetupRouter()

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST",
				"/categories/"+tt.categoryID+"/products/"+tt.productID, nil)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.Bytes())
		})
	}
}

func TestServer_RemoveProductFromCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		categoryID   string
		productID    string
		svc          server.Service
		expectedCode int
		expectedBody []byte
	}{
		{
			name:       "success case",
			categoryID: "cat123",
			productID:  "prod456",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					RemoveProductFromCategory(gomock.Any(), "cat123", "prod456").
					Return(nil)
				return s
			}(),
			expectedCode: http.StatusOK,
			expectedBody: []byte(`{"message":"product removed from category"}`),
		},
		{
			name:       "product not found in category",
			categoryID: "cat123",
			productID:  "notfound",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					RemoveProductFromCategory(gomock.Any(), "cat123", "notfound").
					Return(domain.ErrorProductNotFound)
				return s
			}(),
			expectedCode: http.StatusNotFound,
			expectedBody: []byte(`{"error":"product not found"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.NewServer(tt.svc)
			r := s.SetupRouter()

			w := httptest.NewRecorder()
			req, err := http.NewRequest("DELETE",
				"/categories/"+tt.categoryID+"/products/"+tt.productID, nil)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.Bytes())
		})
	}
}

func TestServer_GetSupplierByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		supplierID   string
		svc          server.Service
		expectedCode int
		expectedBody []byte
	}{
		{
			name:       "success case",
			supplierID: "sup123",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					GetSupplierByID(gomock.Any(), "sup123").
					Return(&domain.Supplier{
						"sup123",
						"Test Supplier"},
						nil)
				return s
			}(),
			expectedCode: http.StatusOK,
			expectedBody: []byte(`{"id":"sup123","name":"Test Supplier"}`),
		},
		{
			name:       "supplier not found",
			supplierID: "notfound",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					GetSupplierByID(gomock.Any(), "notfound").
					Return(nil, domain.ErrorSupplierNotFound)
				return s
			}(),
			expectedCode: http.StatusNotFound,
			expectedBody: []byte(`{"error":"supplier not found"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.NewServer(tt.svc)
			r := s.SetupRouter()

			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/supplier/"+tt.supplierID, nil)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, string(tt.expectedBody), w.Body.String())
		})
	}
}

func TestServer_DeleteSupplierByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		supplierID   string
		svc          server.Service
		expectedCode int
		expectedBody []byte
	}{
		{
			name:       "success case",
			supplierID: "sup123",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					DeleteSupplierByID(gomock.Any(), "sup123").
					Return(nil)
				return s
			}(),
			expectedCode: http.StatusOK,
			expectedBody: []byte(`{"message":"supplier deleted"}`),
		},
		{
			name:       "supplier not found",
			supplierID: "notfound",
			svc: func() server.Service {
				s := internalMock.NewMockService(ctrl)
				s.EXPECT().
					DeleteSupplierByID(gomock.Any(), "notfound").
					Return(domain.ErrorSupplierNotFound)
				return s
			}(),
			expectedCode: http.StatusNotFound,
			expectedBody: []byte(`{"error":"supplier not found"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.NewServer(tt.svc)
			r := s.SetupRouter()

			w := httptest.NewRecorder()
			req, err := http.NewRequest("DELETE", "/supplier/"+tt.supplierID, nil)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, string(tt.expectedBody), w.Body.String())
		})
	}
}
