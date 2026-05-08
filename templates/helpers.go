package templates

import (
	"strings"
	"time"
)

func monthRu(m time.Month) string {
	months := []string{
		"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}
	idx := int(m) - 1
	if idx < 0 || idx >= 12 {
		return ""
	}
	return months[idx]
}

func stripeClass(i int) string {
	return []string{"", "green", "checker", "diag"}[i%4]
}

func tagClass(i int) string {
	return []string{"ptag-pink", "ptag-green", "ptag-white"}[i%3]
}

func renderMarkdown(md string) string {
	var sb strings.Builder
	lines := strings.Split(md, "\n")
	inCode := false
	inList := false
	var para []string

	flushPara := func() {
		if len(para) == 0 {
			return
		}
		sb.WriteString("<p>")
		sb.WriteString(inlineMarkdown(strings.Join(para, " ")))
		sb.WriteString("</p>\n")
		para = para[:0]
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCode {
				sb.WriteString("</code></pre>\n")
				inCode = false
			} else {
				flushPara()
				lang := strings.TrimPrefix(line, "```")
				if lang != "" {
					sb.WriteString(`<pre><code class="language-` + lang + `">`)
				} else {
					sb.WriteString("<pre><code>")
				}
				inCode = true
			}
			continue
		}
		if inCode {
			sb.WriteString(htmlEsc(line) + "\n")
			continue
		}

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "### ") {
			flushPara()
			if inList {
				sb.WriteString("</ul>\n")
				inList = false
			}
			sb.WriteString("<h3>" + inlineMarkdown(trimmed[4:]) + "</h3>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			flushPara()
			if inList {
				sb.WriteString("</ul>\n")
				inList = false
			}
			sb.WriteString("<h2>" + inlineMarkdown(trimmed[3:]) + "</h2>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			flushPara()
			if inList {
				sb.WriteString("</ul>\n")
				inList = false
			}
			sb.WriteString("<h1>" + inlineMarkdown(trimmed[2:]) + "</h1>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			flushPara()
			if !inList {
				sb.WriteString("<ul>\n")
				inList = true
			}
			sb.WriteString("<li>" + inlineMarkdown(trimmed[2:]) + "</li>\n")
			continue
		}
		if trimmed == "" {
			flushPara()
			if inList {
				sb.WriteString("</ul>\n")
				inList = false
			}
			continue
		}
		para = append(para, trimmed)
	}

	flushPara()
	if inList {
		sb.WriteString("</ul>\n")
	}
	if inCode {
		sb.WriteString("</code></pre>\n")
	}
	return sb.String()
}

func inlineMarkdown(s string) string {
	s = replaceDelim(s, "**", "<strong>", "</strong>")
	s = replaceDelim(s, "*", "<em>", "</em>")
	s = replaceDelim(s, "`", "<code>", "</code>")
	return s
}

func replaceDelim(s, delim, open, close string) string {
	var sb strings.Builder
	count := 0
	for {
		idx := strings.Index(s, delim)
		if idx == -1 {
			sb.WriteString(s)
			break
		}
		sb.WriteString(s[:idx])
		if count%2 == 0 {
			sb.WriteString(open)
		} else {
			sb.WriteString(close)
		}
		s = s[idx+len(delim):]
		count++
	}
	return sb.String()
}

func htmlEsc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
