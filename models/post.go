package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Post struct {
	Slug     string
	Title    string
	Date     time.Time
	Summary  string
	Content  string
	Tags     []string
	Featured bool
}

func LoadPosts(dir string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		post := parsePost(strings.TrimSuffix(e.Name(), ".md"), string(data))
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Featured != posts[j].Featured {
			return posts[i].Featured
		}
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, nil
}

func GetPost(dir, slug string) (Post, bool) {
	data, err := os.ReadFile(filepath.Join(dir, slug+".md"))
	if err != nil {
		return Post{}, false
	}
	return parsePost(slug, string(data)), true
}

func SavePost(dir string, post Post) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, post.Slug+".md"), []byte(serialize(post)), 0644)
}

func DeletePost(dir, slug string) error {
	return os.Remove(filepath.Join(dir, slug+".md"))
}

func TitleToSlug(title string) string {
	cyr := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i",
		'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
		'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
		'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch",
		'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "e", 'ю': "yu", 'я': "ya",
	}
	title = strings.ToLower(title)
	var sb strings.Builder
	for _, r := range title {
		if tr, ok := cyr[r]; ok {
			sb.WriteString(tr)
		} else if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			sb.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' {
			sb.WriteRune('-')
		}
	}
	slug := strings.Trim(sb.String(), "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	if slug == "" {
		slug = fmt.Sprintf("post-%d", time.Now().Unix())
	}
	return slug
}

func CollectTags(posts []Post) []string {
	seen := map[string]bool{}
	var tags []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	return tags
}

func serialize(p Post) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString("title: " + p.Title + "\n")
	sb.WriteString("date: " + p.Date.Format("2006-01-02") + "\n")
	if p.Summary != "" {
		sb.WriteString("summary: " + p.Summary + "\n")
	}
	if len(p.Tags) > 0 {
		sb.WriteString("tags: " + strings.Join(p.Tags, ", ") + "\n")
	}
	if p.Featured {
		sb.WriteString("featured: true\n")
	}
	sb.WriteString("---\n\n")
	sb.WriteString(p.Content)
	return sb.String()
}

func parsePost(slug, raw string) Post {
	post := Post{Slug: slug, Date: time.Now()}
	raw = strings.TrimPrefix(raw, "\xef\xbb\xbf")
	if !strings.HasPrefix(raw, "---") {
		post.Content = raw
		return post
	}
	parts := strings.SplitN(raw, "---", 3)
	if len(parts) < 3 {
		post.Content = raw
		return post
	}
	for _, line := range strings.Split(parts[1], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		switch key {
		case "title":
			post.Title = val
		case "date":
			if t, err := time.Parse("2006-01-02", val); err == nil {
				post.Date = t
			}
		case "summary":
			post.Summary = val
		case "tags":
			for _, tag := range strings.Split(val, ",") {
				if t := strings.TrimSpace(tag); t != "" {
					post.Tags = append(post.Tags, t)
				}
			}
		case "featured":
			post.Featured = val == "true"
		}
	}
	post.Content = strings.TrimSpace(parts[2])
	if post.Title == "" {
		post.Title = slug
	}
	return post
}
