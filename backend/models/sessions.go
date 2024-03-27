package models

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
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
	ID     string
	Client Client
	Cookie *http.Cookie
	Values map[interface{}]interface{}
}

func NewSession(c Client) (session *Session) {
	session = new(Session)
	session.ID = generate_id_64(16)
	session.Client = c
	session.Cookie = &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Expires:  time.Now().Add(3 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
		HttpOnly: false,
	}
	return session
}

type SessionStore struct {
	Sessions sync.Map // preferably map[string]*Session
	Timeout  time.Duration
}

func NewStore() (store *SessionStore) {
	store = new(SessionStore)
	store.Timeout = 10 * time.Second
	go store.TimeOutCycle()
	return store
}

func (store *SessionStore) TimeOutCycle() {
	for {
		time.Sleep(store.Timeout)
		store.Sessions.Range(func(key, value any) bool {
			session, ok := value.(*Session)
			if !ok {
				return false
			}

			if session.Cookie.Expires.Before(time.Now()) {
				store.Sessions.Delete(key)
			}
			fmt.Println(key)
			return true
		})

	}
}
