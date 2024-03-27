package models

import "math/rand"

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
	Values map[interface{}]interface{}
}

func NewSession() (session *Session) {
	session = new(Session)
	session.ID = generate_id_64(16)
	return session
}

type SessionStore struct {
	Sessions map[string]*Session
}

func NewStore() (store *SessionStore) {
	store = new(SessionStore)
	store.Sessions = make(map[string]*Session)
	return store
}
