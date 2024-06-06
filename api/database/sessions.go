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
	b64       string
	user_id   string
	user_name string
	expires   time.Time
}

func (sess *Session) SetID(key string) {
	sess.user_id = key
}

func (sess *Session) GetID() string {
	return sess.user_id
}

func (sess *Session) SetName(key string) {
	sess.user_name = key
}

func (sess *Session) GetName() string {
	return sess.user_name
}

func CreateSession(w http.ResponseWriter, r *http.Request) (s *Session) {
	s = new(Session)
	sessid := utils.GenerateBase64ID(16)
	s.b64 = sessid
	cookie := http.Cookie{
		Name:     cookie_name,
		Value:    sessid,
		Expires:  time.Now().Add(session_timeout),
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)
	private_store.mx.Lock()
	s.expires = cookie.Expires
	private_store.sessions[sessid] = s
	private_store.mx.Unlock()
	return s
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

func (sess *Session) End() {
	private_store.mx.Lock()
	delete(private_store.sessions, sess.b64)
	private_store.mx.Unlock()
}

type session_store struct {
	mx       sync.RWMutex
	sessions map[string]*Session
}

func (st *session_store) timeout_cycle() {
	for {
		time.Sleep(session_timeout)
		for key, sess := range st.sessions {
			if sess.expires.Before(time.Now()) {
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
