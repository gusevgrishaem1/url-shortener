package server

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gusevgrishaem1/url-shortener/internal/config"
	"github.com/gusevgrishaem1/url-shortener/internal/model"
	"github.com/gusevgrishaem1/url-shortener/internal/storage"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/teris-io/shortid"
	"log"
	"net/http"
	"time"
)

// Storage интерфейс для работы с хранилищем URL.
type Storage interface {
	Save(url model.URL) error
	Get(short string) (original string, exists bool)
}

// Handler структура для обработки HTTP запросов.
type Handler struct {
	storage Storage
	config  *config.Config
}

// StartServer запускает HTTP сервер.
func StartServer(_ context.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.NewConfig()

	r := mux.NewRouter()

	// Установка параметров пула соединений
	maxOpenConns := 25
	maxIdleConns := 25
	connMaxLifetime := 5 * time.Minute

	postgresStorage, err := storage.NewPostgresStorage(cfg.Database, maxOpenConns, maxIdleConns, connMaxLifetime)
	if err != nil {
		return err
	}
	h := Handler{postgresStorage, cfg}

	// Используем middleware для логирования запросов.
	r.Use(loggingMiddleware)

	r.HandleFunc("/shorten", h.shortenURL).Methods("POST")
	r.HandleFunc("/{short}", h.getOriginalURL).Methods("GET")

	// Настройка CORS для всех источников
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
	})
	handler := c.Handler(r)

	addr := ":" + cfg.Port
	if cfg.EnableHTTPS {
		log.Printf("Starting HTTPS server on %s\n", addr)
		return http.ListenAndServeTLS(addr, "server.crt", "server.key", handler)
	} else {
		log.Printf("Starting HTTP server on %s\n", addr)
		return http.ListenAndServe(addr, handler)
	}
}

// loggingMiddleware логирует информацию о каждом запросе.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// shortenURL обрабатывает POST-запрос для сокращения URL.
func (h *Handler) shortenURL(w http.ResponseWriter, r *http.Request) {
	var u model.URL
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерируем уникальный короткий URL.
	short, err := shortid.Generate()
	if err != nil {
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		return
	}

	short = h.config.ServerURL + "/" + short

	// Сохраняем оригинальный URL в хранилище.
	err = h.storage.Save(model.URL{Original: u.Original, Short: short})
	if err != nil {
		http.Error(w, "Failed to save short URL", http.StatusInternalServerError)
		return
	}

	u.Short = short
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)
}

// getOriginalURL обрабатывает GET-запрос для получения оригинального URL.
func (h *Handler) getOriginalURL(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path

	short = h.config.ServerURL + short

	original, ok := h.storage.Get(short)
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}
