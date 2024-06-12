package database

import (
	"errors"
	"log"
	"net/http"
	"real-time-forum/api/models"
	"sync"
	"time"
)

const (
	session_timeout = time.Hour
	cookie_name     = "SESSION-ID"
)

type SessionStore struct {
	mx       sync.RWMutex
	sessions map[string]*Session
	timer    *time.Ticker
}

func NewSessionStore() *SessionStore {
	store := new(SessionStore)
	store.sessions = make(map[string]*Session)
	store.timeoutCycle()
	return store
}

func (store *SessionStore) timeoutCycle() {
	store.timer = time.NewTicker(session_timeout)
	go func() {
		for range store.timer.C {
			store.mx.Lock()
			for key, sess := range store.sessions {
				if sess.Expires.Before(time.Now()) {
					log.Println("Deleted", key)
					delete(store.sessions, key)
				}
			}
			store.mx.Unlock()
		}
	}()
}

type Session struct {
	ID      string
	User    models.User
	Expires time.Time
}

func (store *SessionStore) NewSession(w http.ResponseWriter, r *http.Request) *Session {
	session := new(Session)
	session.ID = GenerateB64(5)
	cookie := http.Cookie{
		Name:     cookie_name,
		Value:    session.ID,
		Expires:  time.Now().Add(session_timeout),
		Path:     "/",
		HttpOnly: false,
	}
	session.Expires = cookie.Expires
	http.SetCookie(w, &cookie)
	store.mx.Lock()
	store.sessions[session.ID] = session
	store.mx.Unlock()
	return session
}

func (store *SessionStore) GetSession(r *http.Request) (s *Session, err error) {
	cookie, err := r.Cookie(cookie_name)
	if err != nil {
		return nil, err
	}

	s, ok := store.sessions[cookie.Value]
	if !ok {
		return nil, errors.New("no session found")
	}
	return s, nil
}
