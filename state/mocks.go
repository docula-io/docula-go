// Code generated by MockGen. DO NOT EDIT.
// Source: dependencies.go

// Package state is a generated GoMock package.
package state

import (
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// mockFileSystem is a mock of FileSystem interface.
type mockFileSystem struct {
	ctrl     *gomock.Controller
	recorder *mockFileSystemMockRecorder
}

// mockFileSystemMockRecorder is the mock recorder for mockFileSystem.
type mockFileSystemMockRecorder struct {
	mock *mockFileSystem
}

// NewmockFileSystem creates a new mock instance.
func NewmockFileSystem(ctrl *gomock.Controller) *mockFileSystem {
	mock := &mockFileSystem{ctrl: ctrl}
	mock.recorder = &mockFileSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *mockFileSystem) EXPECT() *mockFileSystemMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *mockFileSystem) Create(name string) (File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *mockFileSystemMockRecorder) Create(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*mockFileSystem)(nil).Create), name)
}

// Getwd mocks base method.
func (m *mockFileSystem) Getwd() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Getwd")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Getwd indicates an expected call of Getwd.
func (mr *mockFileSystemMockRecorder) Getwd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Getwd", reflect.TypeOf((*mockFileSystem)(nil).Getwd))
}

// ReadFile mocks base method.
func (m *mockFileSystem) ReadFile(name string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", name)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *mockFileSystemMockRecorder) ReadFile(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*mockFileSystem)(nil).ReadFile), name)
}

// Remove mocks base method.
func (m *mockFileSystem) Remove(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *mockFileSystemMockRecorder) Remove(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*mockFileSystem)(nil).Remove), name)
}

// Rename mocks base method.
func (m *mockFileSystem) Rename(oldpath, newpath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rename", oldpath, newpath)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rename indicates an expected call of Rename.
func (mr *mockFileSystemMockRecorder) Rename(oldpath, newpath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rename", reflect.TypeOf((*mockFileSystem)(nil).Rename), oldpath, newpath)
}

// Stat mocks base method.
func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", name)
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *mockFileSystemMockRecorder) Stat(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*mockFileSystem)(nil).Stat), name)
}

// mockFile is a mock of File interface.
type mockFile struct {
	ctrl     *gomock.Controller
	recorder *mockFileMockRecorder
}

// mockFileMockRecorder is the mock recorder for mockFile.
type mockFileMockRecorder struct {
	mock *mockFile
}

// NewmockFile creates a new mock instance.
func NewmockFile(ctrl *gomock.Controller) *mockFile {
	mock := &mockFile{ctrl: ctrl}
	mock.recorder = &mockFileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *mockFile) EXPECT() *mockFileMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *mockFile) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *mockFileMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*mockFile)(nil).Close))
}

// Write mocks base method.
func (m *mockFile) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *mockFileMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*mockFile)(nil).Write), arg0)
}