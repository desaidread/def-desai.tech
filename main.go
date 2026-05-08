package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"site/handlers"
	"syscall"
	"time"
)

func main() {
	if os.Getenv("ADMIN_KEY") == "" {
		log.Println("⚠  ADMIN_KEY не задан, используется ключ по умолчанию: changeme")
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/posts/", handlers.PostHandler)
	mux.HandleFunc("/about", handlers.AboutHandler)

	mux.HandleFunc("/admin/login", handlers.AdminLogin)
	mux.HandleFunc("/admin/logout", handlers.AdminLogout)
	mux.HandleFunc("/admin", handlers.RequireAdmin(handlers.AdminIndex))
	mux.HandleFunc("/admin/new", handlers.RequireAdmin(handlers.AdminNew))
	mux.HandleFunc("/admin/edit/", handlers.RequireAdmin(handlers.AdminEdit))
	mux.HandleFunc("/admin/delete/", handlers.RequireAdmin(handlers.AdminDelete))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Сервер запущен на http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Ошибка завершения:", err)
	}
	log.Println("Сервер остановлен.")
}
