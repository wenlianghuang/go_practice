package data

import (
	"context"
	"database/sql"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	const createTable = `
CREATE TABLE IF NOT EXISTS books (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    title  TEXT NOT NULL,
    author TEXT NOT NULL,
    price  REAL NOT NULL
);`
	if _, err := db.ExecContext(ctx, createTable); err != nil {
		return err
	}

	var cnt int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM books").Scan(&cnt); err != nil {
		return err
	}
	if cnt == 0 {
		const seed = `
INSERT INTO books (title, author, price) VALUES
('1984','George Orwell',9.99),
('Brave New World','Aldous Huxley',8.99),
('To Kill a Mockingbird','Harper Lee',12.99);`
		if _, err := db.ExecContext(ctx, seed); err != nil {
			return err
		}
	}
	return nil
}
