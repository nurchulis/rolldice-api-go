package user_test

import (
	"context"
	"testing"
	"time"

	"rolldice-go-api/internal/user"
	mock_user "rolldice-go-api/internal/user/mock"

	"rolldice-go-api/pkg/log"

	"rolldice-go-api/pkg/model"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
)

func TestUserService_ListUsers(t *testing.T) {
	validate := validator.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_user.NewMockRepository(ctrl)

	user1 := model.User{
		Firstname:  "Isak",
		Lastname:   "Rickyanto",
		Email:      "isak@ricky.com",
		DateJoined: time.Now(),
		LastLogin:  time.Now(),
	}
	user2 := model.User{
		Firstname:  "Fafa",
		Lastname:   "Tjan",
		Email:      "fafa@tjan.com",
		DateJoined: time.Now(),
		LastLogin:  time.Now(),
	}

	mockListUsers := make([]model.User, 0)
	mockListUsers = append(mockListUsers, user1, user2)

	pagination := model.Pagination{Page: 1, Limit: 2, Sort: "asc"}
	req := &user.ListUsersRequest{Pagination: pagination}
	ctx := context.Background()

	repo.EXPECT().
		List(ctx, gomock.Eq(req)).
		Return(mockListUsers, nil).AnyTimes()

	logger := log.New()
	service := user.NewService(repo, validate, logger)

	t.Log("Test UserService")
	{
		cases := []struct {
			name      string
			page      int
			limit     int
			sort      string
			wantErr   bool
			wantCount int
		}{
			{
				name:      "list with valid parameters",
				page:      1,
				limit:     2,
				sort:      "asc",
				wantErr:   false,
				wantCount: 2,
			},
			{
				name:      "list with invalid page",
				page:      100000,
				limit:     2,
				sort:      "asc",
				wantErr:   true,
				wantCount: 0,
			},
			{
				name:      "list with invalid limit",
				page:      1,
				limit:     100000,
				sort:      "asc",
				wantErr:   true,
				wantCount: 0,
			},
			{
				name:      "list with invalid sort",
				page:      1,
				limit:     100000,
				sort:      "abc",
				wantErr:   true,
				wantCount: 0,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				req.Limit = tt.limit
				req.Page = tt.page
				req.Sort = tt.sort
				users, err := service.ListUsers(ctx, req)

				if tt.wantCount != len(users) {
					t.Errorf("Failed to get correct number of list, want: %d, get: %d, with page : %d, limit: %d, sort: %s",
						tt.wantCount,
						len(users),
						tt.page,
						tt.limit,
						tt.sort)
				}

				if tt.wantErr && err == nil {
					t.Errorf("Failed to validate list users with page : %d, limit: %d, sort: %s", tt.page, tt.limit, tt.sort)
				}
			})
		}

	}

}

func TestUserService_GetUserbyUsername(t *testing.T) {
	validate := validator.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_user.NewMockRepository(ctrl)

	const username = "isak"
	user1 := model.User{
		Username:   "isak",
		Firstname:  "Isak",
		Lastname:   "Rickyanto",
		Email:      "isak@ricky.com",
		DateJoined: time.Now(),
		LastLogin:  time.Now(),
	}

	ctx := context.Background()

	repo.EXPECT().
		GetByUsername(ctx, username).
		Return(&user1, nil).AnyTimes()

	logger := log.New()
	service := user.NewService(repo, validate, logger)

	user, err := service.GetByUsername(ctx, username)

	if err != nil {
		t.Errorf("Failed to get user by username: %s, %s", username, err)
	}

	if user.Username != username {
		t.Errorf("Failed to get user, got : %s, expected :%s", user.Username, username)
	}

}
