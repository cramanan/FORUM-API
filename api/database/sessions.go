package database

import (
	"errors"
	"fmt"
	"math/rand"
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
}

func NewSessionStore() *SessionStore {
	store := new(SessionStore)
	store.sessions = make(map[string]*Session)
	store.timeoutCycle()
	return store
}

func (store *SessionStore) timeoutCycle() {
	go func() {
		for {
			time.Sleep(session_timeout)
			for key, sess := range store.sessions {
				if sess.Expires.Before(time.Now()) {
					store.mx.Lock()
					delete(store.sessions, key)
					store.mx.Unlock()
					fmt.Println("Deleted", key)
				}
			}
		}
	}()
}

func generateBase64ID(lenght int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+-")
	s := make([]rune, lenght)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

type Session struct {
	ID      string
	User    *models.User
	Expires time.Time
}

func (store *SessionStore) NewSession(w http.ResponseWriter, r *http.Request) *Session {
	session := new(Session)
	session.ID = generateBase64ID(16)
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
