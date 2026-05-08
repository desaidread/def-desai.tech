package handlers

import (
	"net/http"
	"site/models"
	"site/templates"
	"strings"
)

const postsDir = "data/posts"

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if isHTMX(r) {
		templates.NotFound().Render(r.Context(), w)
		return
	}
	templates.NotFoundPage().Render(r.Context(), w)
}

func buildIndexData(all []models.Post, activeTag string) templates.IndexData {
	posts := all
	if activeTag != "" {
		var filtered []models.Post
		for _, p := range all {
			for _, t := range p.Tags {
				if t == activeTag {
					filtered = append(filtered, p)
					break
				}
			}
		}
		posts = filtered
	}

	data := templates.IndexData{
		AllTags:    models.CollectTags(all),
		ActiveTag:  activeTag,
		TotalPosts: len(all),
	}

	n := len(all)
	if n > 5 {
		n = 5
	}
	data.Recent = all[:n]

	if len(posts) == 0 {
		return data
	}

	f := posts[0]
	data.Featured = &f

	rest := posts[1:]
	if len(rest) >= 4 {
		data.Grid = rest[:4]
		data.List = rest[4:]
	} else {
		data.Grid = rest
	}
	return data
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound(w, r)
		return
	}
	all, _ := models.LoadPosts(postsDir)
	tag := r.URL.Query().Get("tag")
	d := buildIndexData(all, tag)

	if isHTMX(r) {
		templates.Index(d).Render(r.Context(), w)
		return
	}
	templates.IndexPage(d).Render(r.Context(), w)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	slug = strings.Trim(slug, "/")
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	post, ok := models.GetPost(postsDir, slug)
	if !ok {
		notFound(w, r)
		return
	}
	if isHTMX(r) {
		templates.PostPage(post).Render(r.Context(), w)
		return
	}
	templates.PostFullPage(post).Render(r.Context(), w)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	if isHTMX(r) {
		templates.About().Render(r.Context(), w)
		return
	}
	templates.AboutPage().Render(r.Context(), w)
}
