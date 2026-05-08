package handlers

import (
	"net/http"
	"site/models"
	"site/templates"
	"strings"
	"time"
)

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	posts, _ := models.LoadPosts(postsDir)
	if isHTMX(r) {
		templates.AdminIndex(posts).Render(r.Context(), w)
		return
	}
	templates.AdminIndexPage(posts).Render(r.Context(), w)
}

func AdminNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		adminSave(w, r, "", true)
		return
	}
	blank := models.Post{Date: time.Now()}
	if isHTMX(r) {
		templates.AdminForm(blank, true).Render(r.Context(), w)
		return
	}
	templates.AdminFormPage(blank, true).Render(r.Context(), w)
}

func AdminEdit(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/admin/edit/")
	slug = strings.Trim(slug, "/")

	if r.Method == http.MethodPost {
		adminSave(w, r, slug, false)
		return
	}

	post, ok := models.GetPost(postsDir, slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if isHTMX(r) {
		templates.AdminForm(post, false).Render(r.Context(), w)
		return
	}
	templates.AdminFormPage(post, false).Render(r.Context(), w)
}

func AdminDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	slug := strings.TrimPrefix(r.URL.Path, "/admin/delete/")
	slug = strings.Trim(slug, "/")
	if err := models.DeletePost(postsDir, slug); err != nil {
		http.Error(w, "Failed to delete post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// HTMX: return empty — target element will be removed via outerHTML swap
	w.WriteHeader(http.StatusOK)
}

func adminSave(w http.ResponseWriter, r *http.Request, existingSlug string, isNew bool) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	summary := strings.TrimSpace(r.FormValue("summary"))
	content := strings.TrimSpace(r.FormValue("content"))
	dateStr := r.FormValue("date")
	tagsRaw := r.FormValue("tags")
	featured := r.FormValue("featured") == "true"

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	date := time.Now()
	if d, err := time.Parse("2006-01-02", dateStr); err == nil {
		date = d
	}

	var tags []string
	for _, t := range strings.Split(tagsRaw, ",") {
		if t = strings.TrimSpace(t); t != "" {
			tags = append(tags, t)
		}
	}

	slug := existingSlug
	if isNew {
		slug = strings.TrimSpace(r.FormValue("slug"))
		if slug == "" {
			slug = models.TitleToSlug(title)
		}
		// Ensure slug is unique
		if _, exists := models.GetPost(postsDir, slug); exists {
			slug = slug + "-" + time.Now().Format("20060102150405")
		}
	}

	post := models.Post{
		Slug:     slug,
		Title:    title,
		Date:     date,
		Summary:  summary,
		Content:  content,
		Tags:     tags,
		Featured: featured,
	}

	if err := models.SavePost(postsDir, post); err != nil {
		http.Error(w, "Failed to save post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// After save: redirect to admin (works for both HTMX and full page)
	w.Header().Set("HX-Redirect", "/admin")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
