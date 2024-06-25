package util

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	E "github.com/bmg-c/product-diary/errorhandler"
	"github.com/bmg-c/product-diary/localization"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

func ComponentBytes(component templ.Component, r *http.Request) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := component.Render(r.Context(), buf)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func RenderComponent(writeBytes *[]byte, component templ.Component, r *http.Request) error {
	if writeBytes == nil {
		writeBytes = &[]byte{}
	}
	out, err := ComponentBytes(component, r)
	if err != nil {
		return err
	}
	*writeBytes = append(*writeBytes, out...)
	return nil
}

func RespondHTTP(w http.ResponseWriter, code *int, out *[]byte) {
	if code == nil {
		panic("Code should not be nil")
	}
	if *code != http.StatusOK {
		w.WriteHeader(*code)
	}
	w.Write(*out)
}

func InitHTMLHandler(w http.ResponseWriter, r *http.Request) (l *localization.Localizer) {
	w.Header().Set("Content-Type", "text/html")
	locale, _ := GetLocaleCookieValue(r)
	l = localization.NewLocilizer(locale)
	return l
}

func SetUserSessionCookie(w http.ResponseWriter, sessionUUID uuid.UUID) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    sessionUUID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func GetUserSessionCookieValue(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	sessionUUID, err := getUserSessionCookieValue(r)
	if err != nil {
		switch err {
		case E.ErrNotFound:
			return uuid.UUID{}, err
		case E.ErrUnprocessableEntity:
			DeleteUserSessionCookie(w)
			return uuid.UUID{}, E.ErrNotFound
		default:
			return uuid.UUID{}, E.ErrInternalServer
		}
	}
	return sessionUUID, nil
}

func getUserSessionCookieValue(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return uuid.UUID{}, E.ErrNotFound
		}
		return uuid.UUID{}, E.ErrInternalServer
	}
	sessionUUID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return uuid.UUID{}, E.ErrUnprocessableEntity
	}
	return sessionUUID, nil
}

func DeleteUserSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, &cookie)
}

func SetLocaleCookie(w http.ResponseWriter, locale localization.Locale) {
	cookie := http.Cookie{
		Name:     "locale",
		Value:    fmt.Sprint(locale),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func GetLocaleCookieValue(r *http.Request) (localization.Locale, error) {
	cookie, err := r.Cookie("locale")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return localization.LocaleEnUS, E.ErrNotFound
		}
		return localization.LocaleEnUS, E.ErrInternalServer
	}
	locale64, err := strconv.ParseUint(cookie.Value, 10, 8)
	if err != nil {
		return localization.LocaleEnUS, E.ErrUnprocessableEntity
	}
	locale := localization.Locale(locale64)
	return locale, nil
}

func DeleteLocaleCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, &cookie)
}

func IsErrorSQL(err error, sqlErr error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if sqliteErr.Code == sqlErr {
			return true
		}
	}
	return false
}
