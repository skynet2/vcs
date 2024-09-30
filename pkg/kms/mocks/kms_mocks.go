// Code generated by MockGen. DO NOT EDIT.
// Source: kms.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	jwk "github.com/trustbloc/kms-go/doc/jose/jwk"
	kms "github.com/trustbloc/kms-go/spi/kms"
	vc "github.com/trustbloc/vcs/pkg/doc/vc"
	verifiable "github.com/trustbloc/vcs/pkg/doc/verifiable"
)

// MockVCSKeyManager is a mock of VCSKeyManager interface.
type MockVCSKeyManager struct {
	ctrl     *gomock.Controller
	recorder *MockVCSKeyManagerMockRecorder
}

// MockVCSKeyManagerMockRecorder is the mock recorder for MockVCSKeyManager.
type MockVCSKeyManagerMockRecorder struct {
	mock *MockVCSKeyManager
}

// NewMockVCSKeyManager creates a new mock instance.
func NewMockVCSKeyManager(ctrl *gomock.Controller) *MockVCSKeyManager {
	mock := &MockVCSKeyManager{ctrl: ctrl}
	mock.recorder = &MockVCSKeyManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVCSKeyManager) EXPECT() *MockVCSKeyManagerMockRecorder {
	return m.recorder
}

// CreateCryptoKey mocks base method.
func (m *MockVCSKeyManager) CreateCryptoKey(keyType kms.KeyType) (string, interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCryptoKey", keyType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(interface{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateCryptoKey indicates an expected call of CreateCryptoKey.
func (mr *MockVCSKeyManagerMockRecorder) CreateCryptoKey(keyType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCryptoKey", reflect.TypeOf((*MockVCSKeyManager)(nil).CreateCryptoKey), keyType)
}

// CreateJWKKey mocks base method.
func (m *MockVCSKeyManager) CreateJWKKey(keyType kms.KeyType) (string, *jwk.JWK, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJWKKey", keyType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*jwk.JWK)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateJWKKey indicates an expected call of CreateJWKKey.
func (mr *MockVCSKeyManagerMockRecorder) CreateJWKKey(keyType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJWKKey", reflect.TypeOf((*MockVCSKeyManager)(nil).CreateJWKKey), keyType)
}

// NewVCSigner mocks base method.
func (m *MockVCSKeyManager) NewVCSigner(creator string, signatureType verifiable.SignatureType) (vc.SignerAlgorithm, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewVCSigner", creator, signatureType)
	ret0, _ := ret[0].(vc.SignerAlgorithm)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewVCSigner indicates an expected call of NewVCSigner.
func (mr *MockVCSKeyManagerMockRecorder) NewVCSigner(creator, signatureType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewVCSigner", reflect.TypeOf((*MockVCSKeyManager)(nil).NewVCSigner), creator, signatureType)
}

// SupportedKeyTypes mocks base method.
func (m *MockVCSKeyManager) SupportedKeyTypes() []kms.KeyType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportedKeyTypes")
	ret0, _ := ret[0].([]kms.KeyType)
	return ret0
}

// SupportedKeyTypes indicates an expected call of SupportedKeyTypes.
func (mr *MockVCSKeyManagerMockRecorder) SupportedKeyTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportedKeyTypes", reflect.TypeOf((*MockVCSKeyManager)(nil).SupportedKeyTypes))
}