package post

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"example/internal/pkg/known"
	"example/internal/pkg/model"
)

// 模拟PostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx interface{}, post *model.PostM) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) Get(ctx interface{}, postID string) (*model.PostM, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PostM), args.Error(1)
}

func (m *MockPostRepository) Update(ctx interface{}, post *model.PostM) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(ctx interface{}, postID string) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockPostRepository) List(ctx interface{}, offset, limit int) ([]*model.PostM, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*model.PostM), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) ListByUser(ctx interface{}, username string, offset, limit int) ([]*model.PostM, int64, error) {
	args := m.Called(ctx, username, offset, limit)
	return args.Get(0).([]*model.PostM), args.Get(1).(int64), args.Error(2)
}

// 模拟UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Get(ctx interface{}, username string) (*model.UserM, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserM), args.Error(1)
}

// 其他方法实现(省略)...

// 模拟Authz
type MockAuthz struct {
	mock.Mock
}

func (m *MockAuthz) Authorize(sub, obj, act string) (bool, error) {
	args := m.Called(sub, obj, act)
	return args.Bool(0), args.Error(1)
}

// 其他方法实现(省略)...

func TestPostController_Get(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name           string
		postID         string
		setupAuth      func(c *gin.Context)
		setupMocks     func(*MockPostRepository, *MockUserRepository, *MockAuthz)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "成功获取文章",
			postID: "test-post-id",
			setupAuth: func(c *gin.Context) {
				c.Set(known.XUsernameKey, "testuser")
				c.Set(known.XRoleKey, "reader")
			},
			setupMocks: func(mockPostRepo *MockPostRepository, mockUserRepo *MockUserRepository, mockAuthz *MockAuthz) {
				createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

				post := &model.PostM{
					ID:        1,
					Username:  "testuser",
					PostID:    "test-post-id",
					Title:     "测试文章",
					Content:   "这是测试内容",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				}

				mockPostRepo.On("Get", mock.Anything, "test-post-id").Return(post, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":0,"message":"","data":{"post_id":"test-post-id","title":"测试文章","content":"这是测试内容","username":"testuser","created_at":"2023-01-01 12:00:00","updated_at":"2023-01-01 12:00:00"}}`,
		},
		{
			name:   "文章不存在",
			postID: "nonexistent-post-id",
			setupAuth: func(c *gin.Context) {
				c.Set(known.XUsernameKey, "testuser")
				c.Set(known.XRoleKey, "reader")
			},
			setupMocks: func(mockPostRepo *MockPostRepository, mockUserRepo *MockUserRepository, mockAuthz *MockAuthz) {
				mockPostRepo.On("Get", mock.Anything, "nonexistent-post-id").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"code":20200,"message":"文章不存在","data":null}`,
		},
		{
			name:   "无权限访问",
			postID: "other-user-post-id",
			setupAuth: func(c *gin.Context) {
				c.Set(known.XUsernameKey, "testuser")
				c.Set(known.XRoleKey, "reader")
			},
			setupMocks: func(mockPostRepo *MockPostRepository, mockUserRepo *MockUserRepository, mockAuthz *MockAuthz) {
				createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

				post := &model.PostM{
					ID:        2,
					Username:  "otheruser", // 不是当前用户的文章
					PostID:    "other-user-post-id",
					Title:     "其他用户的文章",
					Content:   "这是其他用户的文章内容",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				}

				mockPostRepo.On("Get", mock.Anything, "other-user-post-id").Return(post, nil)
				mockAuthz.On("Authorize", "testuser", "/v1/posts/other-user-post-id", "GET").Return(false, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"code":10007,"message":"禁止访问","data":null}`,
		},
		{
			name:   "缺少文章ID参数",
			postID: "",
			setupAuth: func(c *gin.Context) {
				c.Set(known.XUsernameKey, "testuser")
				c.Set(known.XRoleKey, "reader")
			},
			setupMocks: func(mockPostRepo *MockPostRepository, mockUserRepo *MockUserRepository, mockAuthz *MockAuthz) {
				// 不需要设置模拟，因为参数错误会提前返回
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":10001,"message":"参数错误","data":null}`,
		},
		{
			name:   "数据库错误",
			postID: "test-post-id",
			setupAuth: func(c *gin.Context) {
				c.Set(known.XUsernameKey, "testuser")
				c.Set(known.XRoleKey, "reader")
			},
			setupMocks: func(mockPostRepo *MockPostRepository, mockUserRepo *MockUserRepository, mockAuthz *MockAuthz) {
				mockPostRepo.On("Get", mock.Anything, "test-post-id").Return(nil, errors.New("数据库连接错误"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":11000,"message":"数据库错误","data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初始化Gin
			gin.SetMode(gin.TestMode)

			// 创建模拟对象
			mockPostRepo := new(MockPostRepository)
			mockUserRepo := new(MockUserRepository)
			mockAuthz := new(MockAuthz)

			// 设置模拟行为
			if tt.setupMocks != nil {
				tt.setupMocks(mockPostRepo, mockUserRepo, mockAuthz)
			}

			// 创建控制器
			ctrl := &Controller{
				db:       nil, // 不需要真实数据库
				authz:    mockAuthz,
				postRepo: mockPostRepo,
				userRepo: mockUserRepo,
			}

			// 创建HTTP请求上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置认证信息
			if tt.setupAuth != nil {
				tt.setupAuth(c)
			}

			// 设置请求参数
			if tt.postID != "" {
				c.AddParam("postID", tt.postID)
			}

			// 执行被测试函数
			ctrl.Get(c)

			// 验证状态码
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 验证响应体
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// 验证模拟对象的期望被满足
			mockPostRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockAuthz.AssertExpectations(t)
		})
	}
}
