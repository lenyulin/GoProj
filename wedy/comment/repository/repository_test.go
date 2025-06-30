package repository

import (
	"GoProj/wedy/comment/domain"
	mockRepoComm "GoProj/wedy/comment/repository/cache/comment/mock"
	"context"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCommentSetRepository(t *testing.T) {
	testcases := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantRes []domain.Comment
		key     string
		offset  int64
	}{
		{
			name:   "success",
			key:    "123456",
			offset: 0,
			wantRes: []domain.Comment{
				{Id: 12345,
					VId:     12345,
					Content: "this is the first comment content",
					User: domain.User{
						Id:        1234567,
						Name:      "1234567",
						AvatarURL: "",
					}},
				{Id: 123456,
					VId:     12345,
					Content: "this is the 2rd comment content",
					User: domain.User{
						Id:        12345678,
						Name:      "12345678",
						AvatarURL: "",
					}},
				{Id: 123457,
					VId:     12345,
					Content: "this is the 3th comment content",
					User: domain.User{
						Id:        123456789,
						Name:      "123456789",
						AvatarURL: "",
					}},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			primaryCache := mockRepoComm.NewMockCache(ctrl)
			primaryCache.EXPECT().Get(context.Background(), tc.key, tc.offset).Return(
				[]domain.Comment{
					{Id: 12345,
						VId:     12345,
						Content: "this is the first comment content",
						User: domain.User{
							Id:        1234567,
							Name:      "1234567",
							AvatarURL: "",
						}},
					{Id: 123456,
						VId:     12345,
						Content: "this is the 2rd comment content",
						User: domain.User{
							Id:        12345678,
							Name:      "12345678",
							AvatarURL: "",
						}},
					{Id: 123457,
						VId:     12345,
						Content: "this is the 3th comment content",
						User: domain.User{
							Id:        123456789,
							Name:      "123456789",
							AvatarURL: "",
						}},
				})
			secondaryCache := mockRepoComm.NewMockCache(ctrl)
			secondaryCache.EXPECT().Get(context.Background(), tc.key, tc.offset).Return(
				[]domain.Comment{
					{Id: 12345,
						VId:     12345,
						Content: "this is the first comment content",
						User: domain.User{
							Id:        1234567,
							Name:      "1234567",
							AvatarURL: "",
						}},
					{Id: 123456,
						VId:     12345,
						Content: "this is the 2rd comment content",
						User: domain.User{
							Id:        12345678,
							Name:      "12345678",
							AvatarURL: "",
						}},
					{Id: 123457,
						VId:     12345,
						Content: "this is the 3th comment content",
						User: domain.User{
							Id:        123456789,
							Name:      "123456789",
							AvatarURL: "",
						}},
				})
			commentRepo := NewCommentRepository(primaryCache, secondaryCache)
			commentRepo.Get(context.Background(), "12345", 0)
		})
	}
}
