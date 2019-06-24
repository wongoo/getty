// +build !linux

/******************************************************
# MAINTAINER : wongoo
# LICENCE    : Apache License 2.0
# EMAIL      : gelnyang@163.com
# MOD        : 2019-06-11
******************************************************/

package getty

import (
	"github.com/pkg/errors"
)

var (
	unsupportedError = errors.New("unsupported")
)

type epollOther struct {
}

func NewEpoller(opts ...EpollOption) (Epoller, error) {
	return nil, unsupportedError
}

func (e *epollOther) Add(ss *session) error {
	return unsupportedError
}

func (e *epollOther) Remove(ss *session) error {
	return unsupportedError
}

func (e *epollOther) Wait() error {
	return unsupportedError
}

func (e *epollOther) Start() {
}
