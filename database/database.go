package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Driver de conexao com o mysql
)

// Abre conexao com o banco de dados
func Conect() (*sql.DB, error) {

	conn := "hennangadelha:capivara@/devbook?charset=utf8&parseTime=True&Local"

	db, erro := sql.Open("mysql", conn)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}
