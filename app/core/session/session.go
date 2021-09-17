package session

import (
	"encoding/gob"

	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func init() {
	gob.Register(database.User{})
	gob.Register(Flash{})
}

type (
	OptionFunc       = func(sess *sessions.Session)
	SessionStoreFunc = func(c echo.Context) *SessionStore
)

var (
	// FlashError is a bootstrap 5 class
	FlashError = "danger"
	// FlashSuccess is a bootstrap 5 class
	FlashSuccess = "success"
	// FlashInfo is a bootstrap 5 class
	FlashInfo = "info"
	// FlashWarning is a bootstrap 5 class
	FlashWarning = "warning"
)

type Flash struct {
	Message string
	Type    string
}

type SessionStore struct {
	*sessions.Session
	ctx echo.Context
}

func (s *SessionStore) Save() error {
	return s.Session.Save(s.ctx.Request(), s.ctx.Response())
}

func (s *SessionStore) Put(key interface{}, val interface{}) {
	s.Session.Values[key] = val
}

func (s *SessionStore) Get(key interface{}) interface{} {
	return s.Session.Values[key]
}

func (s *SessionStore) GetBool(key interface{}) bool {
	if _, ok := s.Get(key).(bool); ok {
		return true
	}

	return false
}

func (s *SessionStore) FlashError(msg string) {
	s.Session.AddFlash(Flash{
		Message: msg,
		Type:    FlashError,
	})
	s.Save()
}

func (s *SessionStore) FlashWarning(msg string) {
	s.Session.AddFlash(Flash{
		Message: msg,
		Type:    FlashWarning,
	})
	s.Save()
}

func (s *SessionStore) FlashSuccess(msg string) {
	s.Session.AddFlash(Flash{
		Message: msg,
		Type:    FlashSuccess,
	})
	s.Save()
}

func (s *SessionStore) Flashes() []Flash {
	flashes := s.Session.Flashes()
	fm := make([]Flash, len(flashes))
	if len(flashes) > 0 {
		for i, f := range flashes {
			switch f.(type) {
			case Flash:
				fm[i] = f.(Flash)
			default:
				fm[i] = Flash{Message: f.(string), Type: FlashInfo}
			}
		}
	}
	s.Save()
	return fm
}

func (s *SessionStore) Authenticate(user database.User, opts ...OptionFunc) error {
	for _, opt := range opts {
		opt(s.Session)
	}

	s.Session.Values["Auth"] = true
	s.Session.Values["User"] = user

	return s.Save()
}

func (s *SessionStore) Logout() {
	delete(s.Values, "Auth")
	delete(s.Values, "User")
	s.Save()
}

func NewSessionStore(ctx echo.Context) *SessionStore {
	sess, err := session.Get("session", ctx)
	if err != nil {
		ctx.Logger().Errorf("could not get session: %v", err)
	}
	return &SessionStore{sess, ctx}
}
