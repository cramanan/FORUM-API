package database

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	private_store  = new_store()
	CookieValue    = "session-id"
	SessionTimeout = 3 * time.Hour
)

func generate_id_64(lenght int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-")
	s := make([]rune, lenght)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

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

func NewSession(w http.ResponseWriter, r *http.Request) (s *Session) {
	s = new(Session)
	sessid := generate_id_64(16)
	s.cookie = http.Cookie{
		Name:     CookieValue,
		Value:    sessid,
		Expires:  time.Now().Add(SessionTimeout),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
		HttpOnly: false,
	}
	s.values = make(map[string]any)

	cookie, _ := r.Cookie(CookieValue)
	if cookie != nil {
		private_store.mx.RLock()
		sx, ok := private_store.sessions[cookie.Value]
		private_store.mx.RUnlock()
		if ok {
			return sx
		}

	}
	http.SetCookie(w, &s.cookie)
	private_store.mx.Lock()
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
		time.Sleep(SessionTimeout)
		for key, sess := range st.sessions {
			if sess.cookie.Expires.Before(time.Now()) {
				st.mx.Lock()
				delete(st.sessions, key)
				fmt.Println("Deleted", key)
				st.mx.Unlock()
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
