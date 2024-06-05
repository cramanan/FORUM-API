package database

import (
	"errors"
	"fmt"
	"net/http"
	"real-time-forum/api/utils"
	"sync"
	"time"
)

type Session struct {
	mx     sync.RWMutex
	cookie http.Cookie
	values map[string]any
}

func (sess *Session) Get(key string) (value any, ok bool) {
	sess.mx.RLock()
	defer sess.mx.RUnlock()
	value, ok = sess.values[key]
	return value, ok
}

func (sess *Session) Set(key string, value any) {
	sess.mx.Lock()
	defer sess.mx.Unlock()
	sess.values[key] = value
}

func (sess *Session) Values() map[string]any {
	return sess.values
}

func GetSession(w http.ResponseWriter, r *http.Request) (s *Session, err error) {
	cookie, err := r.Cookie(cookie_name)
	if err != nil {
		return nil, err
	}

	s, ok := private_store.sessions[cookie.Value]
	if !ok {
		return nil, errors.New("no session found")
	}
	return s, nil
}

func CreateSession(w http.ResponseWriter, r *http.Request) (s *Session) {
	s = new(Session)
	s.values = make(map[string]any)
	sessid := utils.GenerateBase64ID(16)
	cookie := http.Cookie{
		Name:     cookie_name,
		Value:    sessid,
		Expires:  time.Now().Add(session_timeout),
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	private_store.mx.Lock()
	s.cookie = cookie
	private_store.sessions[sessid] = s
	private_store.mx.Unlock()
	return s
}

func (sess *Session) End() {
	sess.cookie.Expires = time.Now()
}

type session_store struct {
	mx       sync.RWMutex
	sessions map[string]*Session
}

func (st *session_store) timeout_cycle() {
	for {
		time.Sleep(session_timeout)
		for key, sess := range st.sessions {
			if sess.cookie.Expires.Before(time.Now()) {
				st.mx.Lock()
				delete(st.sessions, key)
				st.mx.Unlock()
				fmt.Println("Deleted", key)
			}
		}
	}
}

func new_store() *session_store {
	store := new(session_store)
	store.sessions = make(map[string]*Session)
	go store.timeout_cycle()
	return store
}

var private_store = new_store()

const (
	cookie_name     = "session-id"
	session_timeout = time.Minute * 10
)
