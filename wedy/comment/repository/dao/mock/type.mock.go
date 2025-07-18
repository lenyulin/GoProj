// Code generated by MockGen. DO NOT EDIT.
// Source: ./dao/comment/type.go
//
// Generated by this command:
//
//	mockgen --package=mockDAOComment --source=./dao/comment/type.go --destination=./dao/mock/type.mock.go
//

// Package mockDAOComment is a generated GoMock package.
package mockDAOComment

import (
	comment "GoProj/wedy/comment/repository/dao"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCommentDAO is a mock of CommentDAO interface.
type MockCommentDAO struct {
	ctrl     *gomock.Controller
	recorder *MockCommentDAOMockRecorder
	isgomock struct{}
}

// MockCommentDAOMockRecorder is the mock recorder for MockCommentDAO.
type MockCommentDAOMockRecorder struct {
	mock *MockCommentDAO
}

// NewMockCommentDAO creates a new mock instance.
func NewMockCommentDAO(ctrl *gomock.Controller) *MockCommentDAO {
	mock := &MockCommentDAO{ctrl: ctrl}
	mock.recorder = &MockCommentDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentDAO) EXPECT() *MockCommentDAOMockRecorder {
	return m.recorder
}

// FindById mocks base method.
func (m *MockCommentDAO) FindById(ctx context.Context, id int64) ([]comment.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, id)
	ret0, _ := ret[0].([]comment.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockCommentDAOMockRecorder) FindById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockCommentDAO)(nil).FindById), ctx, id)
}

// Insert mocks base method.
func (m *MockCommentDAO) Insert(ctx context.Context, comment comment.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCommentDAOMockRecorder) Insert(ctx, comment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCommentDAO)(nil).Insert), ctx, comment)
}

// Update mocks base method.
func (m *MockCommentDAO) Update(ctx context.Context, comment comment.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCommentDAOMockRecorder) Update(ctx, comment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCommentDAO)(nil).Update), ctx, comment)
}
