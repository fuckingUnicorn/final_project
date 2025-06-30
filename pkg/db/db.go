package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// db — глобальная переменная для доступа к базе данных
var DB *sql.DB

// schema содержит SQL-команды для создания таблицы scheduler и индекса по дате
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// InitSchema выполняет SQL-команды для создания таблицы и индекса, если их нет
// Эту функцию можно вызывать отдельно для любого экземпляра БД
func InitSchema(db *sql.DB) error {
	// Выполняем schema — она безопасна из-за IF NOT EXISTS
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("ошибка инициализации схемы: %w", err)
	}
	return nil
}

// Init открывает (или создаёт) БД и таблицу scheduler, если её нет
func Init(dbFile string) error {
	// Открываем базу данных (создаст файл, если его нет)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Инициализируем схему
	if err := InitSchema(db); err != nil {
		db.Close()
		return err
	}

	DB = db
	return nil
}
