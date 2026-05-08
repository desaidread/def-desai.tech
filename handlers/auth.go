package handlers

import (
	"net/http"
	"os"
	"site/templates"
	"strings"
)

const cookieName = "admin_tok"

func adminKey() string {
	k := strings.TrimSpace(os.Getenv("ADMIN_KEY"))
	if k != "" {
		return k
	}
	return "changeme"
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil || c.Value != adminKey() {
			if isHTMX(r) {
				w.Header().Set("HX-Redirect", "/admin/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		input := strings.TrimSpace(r.FormValue("key"))
		if input == adminKey() {
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    adminKey(),
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   86400 * 7,
			})
			w.Header().Set("HX-Redirect", "/admin")
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		if isHTMX(r) {
			templates.AdminLoginForm(true).Render(r.Context(), w)
		} else {
			templates.AdminLoginPage(true).Render(r.Context(), w)
		}
		return
	}
	if isHTMX(r) {
		templates.AdminLoginForm(false).Render(r.Context(), w)
		return
	}
	templates.AdminLoginPage(false).Render(r.Context(), w)
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
