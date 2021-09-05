package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/i18n"
)

type UserFakeStore struct{}

func (s *UserFakeStore) First(m *database.User) error {
	return nil
}
func (s *UserFakeStore) Find(m *[]database.User) error {
	return nil
}
func (s *UserFakeStore) Create(m *database.User) error {
	return nil
}
func (s *UserFakeStore) Ping() error {
	return nil
}

type FakeHasher struct{}

func (c *FakeHasher) Generate(password string) (string, error) {
	return string("hash-password"), nil
}

func (c *FakeHasher) Verify(hash, passowrd string) bool {
	return true
}

func TestSignupGet(t *testing.T) {
	a := e.app.Echo.Group("/auth")
	a.GET("/signup", SignupGet)

	req := httptest.NewRequest("GET", "/auth/signup", nil)
	rec := httptest.NewRecorder()

	e.app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("could not visit '/auth/signup' page, status want %d got %d", http.StatusOK, rec.Code)
	}
}

func TestSignupPost(t *testing.T) {
	a := e.app.Echo.Group("/auth")
	a.POST("/signup", SignupPost)

	cc := core.AppContext{
		App:   e.app,
		Loc:   i18n.New(),
		Store: &database.Store{User: &UserFakeStore{}},
	}

	e.app.Echo.Use(core.AppCtxMiddleware(&cc))

	body := `name=john&email=j.doe@example.com&password=secret`
	req := httptest.NewRequest("POST", "/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	e.app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("could not visit '/auth/signup' page, status want %d got %d", http.StatusOK, rec.Code)
	}
}
