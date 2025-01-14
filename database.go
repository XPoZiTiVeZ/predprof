package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type User struct {
	Id              int
	Email           string
	Password        string
	IsAuthenticated bool
	IsActive        bool
	IsAdmin         bool
	IsSuperuser     bool
	LastLogin       string
	CreatedAt       string
}

type Status struct {
	Id   int
	Name string
}

type Inventory struct {
	Id       int
	Name     string
	Status   string
	Quantity int
}

type ItemQuery struct {
	Id       int
	Quantity int
	User     int
	Status   bool
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "db.db")
	if err != nil {
		os.WriteFile("db.sqlite3", []byte{}, 0777)
		db, _ = sql.Open("sqlite", "db.db")
	}

	_, err = db.Exec(`
    CREATE TABLE
    IF NOT EXISTS
    users (
        id          INTEGER PRIMARY KEY AUTOINCREMENT,
        email       TEXT UNIQUE NOT NULL,
        password    TEXT NOT NULL,
        isActive    BOOLEAN DEFAULT FALSE,
        isAdmin     BOOLEAN DEFAULT FALSE,
        isSuperuser BOOLEAN DEFAULT FALSE,
        last_login  DATETIME DEFAULT CURRENT_TIMESTAMP,
        created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
    );`)

	if err != nil {
		log.Printf("openDB: %s", err)
		return db, err
	}

	_, err = db.Exec(`
	CREATE TABLE
	IF NOT EXISTS
	item_names (
	    id      INTEGER PRIMARY KEY AUTOINCREMENT,
	    name    TEXT UNIQUE
	);`)

	if err != nil {
		log.Printf("openDB: %s", err)
		return db, err
	}

	_, err = db.Exec(`
    CREATE TABLE
    IF NOT EXISTS
    item_statuses (
        id      INTEGER PRIMARY KEY AUTOINCREMENT,
        name    TEXT UNIQUE
    );`)

	if err != nil {
		log.Printf("openDB: %s", err)
		return db, err
	}

	_, err = db.Exec(`
    CREATE TABLE
    IF NOT EXISTS
    items (
        id          INTEGER PRIMARY KEY AUTOINCREMENT,
        name        TEXT,
        status      TEXT,
        quantity    INTEGER,
        FOREIGN KEY(name)   REFERENCES item_names(name),
        FOREIGN KEY(status) REFERENCES item_statuses(name)
    );`)

	if err != nil {
		log.Printf("openDB: %s", err)
		return db, err
	}

	_, err = db.Exec(`
    CREATE TABLE
    IF NOT EXISTS
    items_users (
        id       INTEGER,
		quantity INTEGER,
		user     INTEGER,
		status   BOOLEAN DEFAULT FALSE,
		FOREIGN KEY(id)     REFERENCES items(id),
        FOREIGN KEY(user)   REFERENCES users(id)
    );`)

	if err != nil {
		log.Printf("openDB: %s", err)
		return db, err
	}

	log.Print("DB opened successfully")
	return db, nil
}

// func getUserById(id int) (User, error) {
//     row := db.QueryRow(`SELECT * FROM users WHERE id = ?`, id)

//     var user User
//     var LastLogin int64
//     var CreatedAt int64
//     err := row.Scan(&user.Email, &user.Password, &user.IsActive, &user.IsAdmin, &user.IsSuperUser, &LastLogin, &CreatedAt)
//     if err != nil { log.Print("getUserById", err); return *new(User), err }

//     user.LastLogin = time.UnixMilli(LastLogin)
//     user.CreatedAt = time.UnixMilli(CreatedAt)

//     return user, nil
// }

func AddUser(email, password string) (User, error) {
	_, err := db.Exec(`
		INSERT INTO users(email, password)
		VALUES (?, ?)`,
		email, password)

	if err != nil {
		return User{}, err
	}

	user, _ := GetUserByEmail(email)
	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	row := db.QueryRow(`SELECT * FROM users;`, email)
	var user User
	var LastLogin *string
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.IsActive, &user.IsAdmin, &user.IsSuperuser, &LastLogin, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}

	if LastLogin == nil {
		user.LastLogin = ""
	} else {
		user.LastLogin = *LastLogin
	}

	user.IsAuthenticated = true

	return user, nil
}

func AddItemQuery(id, quantity, user int) error {
	_, err := db.Exec(`
		INSERT INTO items_users(id, quantity, user)
		VALUES(?, ?, ?)`,
		id, quantity, user)

	return err
}

func GetItemQueries(user int) []ItemQuery {
	rows, err := db.Query(`
		SELECT id, quantity, user, status
		FROM items_names
		WHERE user = ?`,
		user)

	var queries []ItemQuery
	if err != nil {
		log.Print(err)
		return queries
	}

	defer rows.Close()
	for rows.Next() {
		var query ItemQuery
		rows.Scan(query.Id, query.Quantity, query.User, query.Status)
		queries = append(queries, query)
	}

	return queries
}
