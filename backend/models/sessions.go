package models

type Session struct {
	Values map[interface{}]interface{}
}

type SessionStore struct {
	Sessions map[string]Session
}

func NewStore() (store *SessionStore) {
	store = new(SessionStore)
	return store
}
