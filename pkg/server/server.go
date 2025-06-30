package server

import (
	"log"
	"net/http"
	"os"

	"go1f/pkg/api"
)

// Run запускает веб-сервер для отдачи статических файлов фронтенда (директория web).
// По умолчанию используется порт 3988, но если задана переменная окружения TODO_PORT,
func Run() error {
	webDir := "web" // Директория с фронтендом

	// Порт по умолчанию
	port := "3988"
	// Если задана переменная окружения TODO_PORT — используем её значение
	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		port = envPort
	}

	// Регистрируем API-обработчики
	api.Init()

	// Регистрируем файловый сервер для отдачи фронтенда и статики
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// Запускаем сервер
	log.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
