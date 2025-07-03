package sql

import "database/sql"

// getDB 获取数据库连接
func getDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./prompts.db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
