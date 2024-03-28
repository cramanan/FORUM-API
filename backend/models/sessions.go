package models

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	store            = new_session_store()
	SessionsIdLenght = 10
	SessionTimeout   = time.Hour
)

const (
	COOKIENAME = "session-id"
)

func generate_id_64(lenght int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-")
	s := make([]rune, lenght)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

type session_store struct {
	sync.Map // Preferably map[string]*session
	timeout  time.Duration
}

func new_session_store() (st *session_store) {
	st = new(session_store)
	st.timeout = SessionTimeout
	go st.timeout_cycle()
	return st
}

func (st *session_store) timeout_cycle() {
	for {
		time.Sleep(st.timeout)
		st.Map.Range(func(key, value any) bool {
			fmt.Println(key, value)
			return true
		})

	}
}

type Sessions struct {
	sync.Map
}

func Session(w http.ResponseWriter, r *http.Request) (sess *Sessions) {
	cookie, _ := r.Cookie(COOKIENAME)
	sess = new(Sessions)
	if cookie != nil {
		sess_any, ok := store.Map.Load(cookie.Value)
		if !ok {
			return sess
		}
		sess = sess_any.(*Sessions)
	} else {
		sessid := generate_id_64(SessionsIdLenght)
		http.SetCookie(w, &http.Cookie{
			Name:     COOKIENAME,
			Value:    sessid,
			Expires:  time.Now().Add(SessionTimeout),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Path:     "/",
			HttpOnly: false,
		})
		store.Map.Store(sessid, sess)
	}
	return sess
}

func (sess *Sessions) Store(key any, value any) {
	sess.Map.Store(key, value)
}

func (sess *Sessions) Load(key any) {
	sess.Map.Load(key)
}
