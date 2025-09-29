package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/aibekfatkhulla/shop/internal/mocks"
	"github.com/aibekfatkhulla/shop/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbErr := errors.New("db error")

	tests := []struct {
		name        string
		user        *domain.User
		repository  service.Repository
		expectedErr error
	}{
		{
			"success case",
			&domain.User{
				Name:     "arnur",
				Password: "123",
				Email:    "qwe@qwe.qwe",
				Number:   "123123123",
				Address:  "the capella",
				Balance:  100,
			},
			func() service.Repository {
				r := mocks.NewMockRepository(ctrl)

				r.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
				r.EXPECT().GetByEmail(gomock.Any(), "qwe@qwe.qwe").Return(nil, domain.ErrorUserNotFound)

				return r
			}(),
			nil,
		},
		{
			"user already exists",
			&domain.User{
				Name:     "arnur",
				Password: "123",
				Email:    "qwe@qwe.qwe",
				Number:   "123123123",
				Address:  "the capella",
				Balance:  100,
			},
			func() service.Repository {
				r := mocks.NewMockRepository(ctrl)
				r.EXPECT().GetByEmail(gomock.Any(), "qwe@qwe.qwe").Return(&domain.User{
					ID:        "123",
					Name:      "arnur",
					Password:  "123",
					Email:     "qwe@qwe.qwe",
					Number:    "123123123",
					Address:   "the capella",
					Balance:   100,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)

				return r
			}(),
			domain.ErrorUserAlreadyExists,
		},
		{
			"db error getting user by email",
			&domain.User{
				Name:     "arnur",
				Password: "123",
				Email:    "qwe@qwe.qwe",
				Number:   "123123123",
				Address:  "the capella",
				Balance:  100,
			},
			func() service.Repository {
				r := mocks.NewMockRepository(ctrl)
				r.EXPECT().GetByEmail(gomock.Any(), "qwe@qwe.qwe").Return(nil, dbErr)

				return r
			}(),
			dbErr,
		},
		{
			"db error creating user",
			&domain.User{
				Name:     "arnur",
				Password: "123",
				Email:    "qwe@qwe.qwe",
				Number:   "123123123",
				Address:  "the capella",
				Balance:  100,
			},
			func() service.Repository {
				r := mocks.NewMockRepository(ctrl)
				r.EXPECT().GetByEmail(gomock.Any(), "qwe@qwe.qwe").Return(nil, domain.ErrorUserNotFound)
				r.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(dbErr)

				return r
			}(),
			dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewService(tt.repository)

			err := s.CreateUser(t.Context(), tt.user)
			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestUpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbErr := errors.New("db error")

	order := &domain.Order{
		ID:     "123",
		UserID: "999",
		Status: domain.StatusPending,
	}

	tests := []struct {
		name        string
		inputOrder  *domain.Order
		mockSetup   func() service.Repository
		expectedErr error
	}{
		{
			name:       "success",
			inputOrder: order,
			mockSetup: func() service.Repository {
				r := mocks.NewMockRepository(ctrl)
				r.EXPECT().GetOrderByID(gomock.Any(), order.ID).Return(order, nil)
				r.EXPECT().UpdateOrder(gomock.Any(), order).Return(nil)
				return r
			},
			expectedErr: nil,
		},
		{
			name:       "get order error",
			inputOrder: order,
			mockSetup: func() service.Repository {
				r := mocks.NewMockRepository(ctrl)
				r.EXPECT().GetOrderByID(gomock.Any(), order.ID).Return(nil, dbErr)
				return r
			},
			expectedErr: dbErr,
		},
		{
			name:       "update order error",
			inputOrder: order,
			mockSetup: func() service.Repository {
				r := mocks.NewMockRepository(ctrl)
				r.EXPECT().GetOrderByID(gomock.Any(), order.ID).Return(order, nil)
				r.EXPECT().UpdateOrder(gomock.Any(), order).Return(dbErr)
				return r
			},
			expectedErr: dbErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockSetup()
			s := service.NewService(repo)

			err := s.UpdateOrder(t.Context(), tt.inputOrder)
			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
