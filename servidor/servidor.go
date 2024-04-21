package servidor

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// struct de usuario, nome minusculo devido ao objeto ser usado apenas neste pacote
type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {

	bodyRequest, erro := io.ReadAll(r.Body)

	if erro != nil {
		w.Write([]byte("Fail reading body request"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(bodyRequest, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter usuario"))
		return
	}

	fmt.Println(usuario)

	db, erro := database.Conect()

	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}

	defer db.Close()

	statement, erro := db.Prepare("INSERT INTO usuarios (nome, email) VALUES (?, ?)")

	if erro != nil {
		w.Write([]byte("Erro ao criar statment"))
		return
	}

	defer statement.Close()

	insert, erro := statement.Exec(usuario.Nome, usuario.Email)

	if erro != nil {
		w.Write([]byte("Erro ao criar usuario"))
		return
	}

	idUsuario, erro := insert.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao recuperar id do usuario"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso! ID: %d", idUsuario)))

}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {

	db, erro := database.Conect()
	if erro != nil {
		w.Write([]byte("erro ao conectar com o banco de dados"))
		return
	}

	defer db.Close()

	rows, erro := db.Query("SELECT * FROM usuarios")

	if erro != nil {
		w.Write([]byte("erro ao buscar usuarios"))
		return
	}

	defer rows.Close()

	var usuarios []usuario

	for rows.Next() {
		var usuario usuario

		if erro := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {

			w.Write([]byte("erro ao capturar usuario"))
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)

	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.Write([]byte("erro na conversão de  usuario para json"))
		return
	}

}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	// -> ParseUint (parametro, base decimal e tamanho bytes)
	ID, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.Write([]byte("erro ao converter parametro"))
		return
	}

	db, erro := database.Conect()

	if erro != nil {
		w.Write([]byte("erro ao conectar ao bacno de dados"))
		return
	}

	defer db.Close()

	row, erro := db.Query("SELECT * FROM usuarios WHERE id = ?", ID)
	if erro != nil {
		w.Write([]byte("erro ao buscar usuario no banco de dados"))
		return
	}

	var usuario usuario

	if row.Next() {

		if erro := row.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("erro ao scan  usuario"))
			return
		}
	}

	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		if erro := row.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("erro ao converter o usuario para json "))
			return
		}
	}

	row.Close()
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	// -> ParseUint (parametro, base decimal e tamanho bytes)
	ID, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.Write([]byte("erro ao converter parametro"))
		return
	}

	var usuario usuario

	bodyRequest, erro := io.ReadAll(r.Body)

	if erro := json.Unmarshal(bodyRequest, &usuario); erro != nil {
		w.Write([]byte("erro ao mapear usuario em struct"))
		return
	}

	db, erro := database.Conect()

	if erro != nil {
		w.Write([]byte("erro ao conectar ao bacno de dados"))
		return
	}

	defer db.Close()

	statement, erro := db.Prepare("UPDATE usuarios SET nome = ?, email = ? WHERE id = ?")

	if erro != nil {
		w.Write([]byte("erro ao criar statement"))
		return
	}

	defer statement.Close()

	if _, erro := statement.Exec(usuario.Nome, usuario.Email, ID); erro != nil {
		w.Write([]byte("erro ao atualizar usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	// -> ParseUint (parametro, base decimal e tamanho bytes)
	ID, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.Write([]byte("erro ao converter parametro"))
		return
	}

	db, erro := database.Conect()

	if erro != nil {
		w.Write([]byte("erro ao conectar ao bacno de dados"))
		return
	}

	defer db.Close()

	statement, erro := db.Prepare("DELETE FROM usuarios WHERE id = ?")

	if erro != nil {
		w.Write([]byte("erro ao criar statement"))
		return
	}

	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("erro ao deletar usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
