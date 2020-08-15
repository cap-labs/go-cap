package mock

import (
	"github.com/libs4go/scf4go"
	"github.com/libs4go/smf4go"
)

// NullLocker .
type mockLocker struct {
}

// Lock .
func (locker *mockLocker) Lock() {

}

// Unlock .
func (locker *mockLocker) Unlock() {

}

// NewLocker .
func NewLocker(config scf4go.Config) (smf4go.Service, error) {
	return &mockLocker{}, nil
}
