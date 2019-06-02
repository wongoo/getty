// +build !linux

package getty

import "github.com/pkg/errors"

var (
	unsupportedError = errors.New("unsupported")
)

type epollOther struct {
}

func createEpoll() (epoller, error) {
	return nil, unsupportedError
}

func (e *epollOther) add(ss *session) error {
	return unsupportedError
}

func (e *epollOther) remove(ss *session) error {
	return unsupportedError
}

func (e *epollOther) wait() ([]*session, error) {
	return nil, unsupportedError
}

func (e *epollOther) start() {
}
