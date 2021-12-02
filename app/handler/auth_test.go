package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahmadwaleed/choreui/app/database/model"
)

var testuser = User{
	Name:  "john",
	Email: "j.doe@example.com",
}

func TestSignupGet(t *testing.T) {
	a := app.Echo.Group("/auth")
	a.GET("/signup", SignupGet)

	req := httptest.NewRequest("GET", "/auth/signup", nil)
	rec := httptest.NewRecorder()

	app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("could not visit '/auth/signup' page, status want %d got %d", http.StatusOK, rec.Code)
	}
}

func TestSignupPost(t *testing.T) {
	if err := migrateUp(); err != nil {
		t.Error(err)
	}

	a := app.Echo.Group("/auth")
	a.POST("/signup", SignupPost)

	body := fmt.Sprintf("name=%s&email=%s&password=secret", testuser.Name, testuser.Email)
	req := httptest.NewRequest("POST", "/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("could not signup new user, status want %d got %d", http.StatusOK, rec.Code)
	}

	if err := migrateDown(); err != nil {
		t.Error(err)
	}
}

func TestSignInGet(t *testing.T) {
	a := app.Echo.Group("/auth")
	a.GET("/signin", SignupGet)

	req := httptest.NewRequest("GET", "/auth/signin", nil)
	rec := httptest.NewRecorder()

	app.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("could not visit '/auth/signin' page, status want %d got %d", http.StatusOK, rec.Code)
	}
}

func TestSignInPost(t *testing.T) {
	if err := migrateUp(); err != nil {
		t.Error(err)
	}

	// create test user
	store := ctx.Store(ctx.App.DB())
	if err := store.User.Create(testuser.Name, testuser.Email, "$2y$10$9P7pi./SZRBmilkg3ELey.AgM8vYbUiDenWxYF2r6X8CcyUllNIDO"); err != nil {
		t.Errorf("could not create test user: %v", err)
	}

	a := app.Echo.Group("/auth")
	a.POST("/signin", SignInPost)

	body := fmt.Sprintf("email=%s&password=secret", testuser.Email)
	req := httptest.NewRequest("POST", "/auth/signin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	app.Echo.ServeHTTP(rec, req)
	sess := ctx.SessionStore(ctx.Context)

	if !sess.GetBool("Auth") {
		t.Error("could not get login value from session store")
	}

	login := sess.GetBool("Auth")
	if !sess.GetBool("Auth") {
		t.Errorf("unexpected session auth value, want: %t, got: %t", true, login)
	}

	user := (sess.Values["User"]).(model.User)
	if user.Email != testuser.Email {
		t.Errorf("unexpected session user email, want: %s, got: %s", testuser.Email, user.Email)
	}

	if rec.Code != http.StatusSeeOther {
		t.Errorf("could not login user, status want %d got %d", http.StatusOK, rec.Code)
	}

	if err := migrateDown(); err != nil {
		t.Error(err)
	}
}
