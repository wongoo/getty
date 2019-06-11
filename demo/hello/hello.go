package hello

import (
	"github.com/dubbogo/getty"
	"math/rand"
	"sync"
	"syscall"
	"time"
)

var (
	Sessions []getty.Session
	lock     sync.RWMutex
)

func HelloClientRequest() {
	for {
		if selectSession() != nil {
			break
		}
		time.Sleep(time.Second)
	}
	echoTimes := 10

	counter := getty.CountWatch{}
	counter.Start()
	for i := 0; i < echoTimes; i++ {
		session := selectSession()
		err := session.WritePkg("hello", WritePkgTimeout)
		if err != nil {
			log.Infof("session.WritePkg(session{%s}, error{%v}", session.Stat(), err)
			session.Close()
			removeSession(session)
		}
	}
	cost := counter.Count()
	log.Infof("after loop %d times, echo cost %d ms", echoTimes, cost/1e6)
}

func selectSession() getty.Session {
	lock.RLock()
	defer lock.RUnlock()
	count := len(Sessions)
	if count == 0 {
		log.Infof("client session array is nil...")
		return nil
	}

	return Sessions[rand.Int31n(int32(count))]
}

func removeSession(session getty.Session) {
	if session == nil {
		return
	}
	lock.Lock()
	for i, s := range Sessions {
		if s == session {
			Sessions = append(Sessions[:i], Sessions[i+1:]...)
			log.Infof("delete session{%s}, its index{%d}", session.Stat(), i)
			break
		}
	}
	log.Infof("after remove session{%s}, left session number:%d", session.Stat(), len(Sessions))
	lock.Unlock()
}

func SetLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Infof("set cur limit: %d", rLimit.Cur)
}
