---
title: Паттерны HTMX
date: 2026-04-07
summary: Основные паттерны работы с HTMX — lazy load, infinite scroll, формы.
tags: htmx, frontend
---

## hx-get и hx-target

Самый базовый паттерн — загрузить фрагмент и вставить его в элемент:

```html
<button hx-get="/api/data" hx-target="#result">
    Загрузить
</button>
<div id="result"></div>
```

## hx-push-url

Чтобы адрес в браузере менялся при навигации:

```html
<a hx-get="/posts/hello" hx-target="#content" hx-push-url="true">
    Читать
</a>
```

## hx-swap

Управляет тем, как вставляется контент:

- `innerHTML` (по умолчанию) — заменяет содержимое
- `outerHTML` — заменяет сам элемент
- `beforeend` — добавляет в конец (infinite scroll)
- `afterbegin` — добавляет в начало

## Серверная сторона

Сервер должен различать обычные запросы и HTMX-запросы:

```go
func isHTMX(r *http.Request) bool {
    return r.Header.Get("HX-Request") == "true"
}
```

Если запрос от HTMX — отдаём только фрагмент.
Если обычный — отдаём полную страницу с layout.

## Итог

HTMX работает по принципу **HTML over the wire** — сервер возвращает HTML, а не JSON. Это упрощает архитектуру и убирает необходимость в отдельном SPA-фреймворке.
