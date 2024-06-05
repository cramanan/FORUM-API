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
	sync.RWMutex
	cookie http.Cookie
	values map[string]any
}

func (sess *Session) Get(key string) (value any, ok bool) {
	sess.RLock()
	defer sess.RUnlock()
	value, ok = sess.values[key]
	return value, ok
}

func (sess *Session) Set(key string, value any) {
	sess.Lock()
	defer sess.Unlock()
	sess.values[key] = value
}

func (sess *Session) Values() map[string]any {
	return sess.values
}

func GetSession(w http.ResponseWriter, r *http.Request) (s *Session, err error) {
	cookie, err := r.Cookie(CookieName)
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
	sessid := utils.Generate_id_64(16)
	cookie := http.Cookie{
		Name:     CookieName,
		Value:    sessid,
		Expires:  time.Now().Add(SessionTimeout),
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	private_store.Lock()
	s.cookie = cookie
	private_store.sessions[sessid] = s
	private_store.Unlock()
	return s
}

func (sess *Session) End() {
	sess.cookie.Expires = time.Now()
}

type session_store struct {
	sync.RWMutex
	sessions map[string]*Session
}

func (st *session_store) timeout_cycle() {
	for {
		time.Sleep(SessionTimeout)
		for key, sess := range st.sessions {
			if sess.cookie.Expires.Before(time.Now()) {
				st.Lock()
				delete(st.sessions, key)
				fmt.Println("Deleted", key)
				st.Unlock()
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

var (
	private_store  = new_store()
	CookieName     = "session-id"
	SessionTimeout = time.Minute * 10
)
