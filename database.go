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
	Name     int
	Quantity int
	Status   int
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

func AddItems(name, status, quantity int) error {
	_, err := db.Exec(`
		INSERT INTO items(name, status, quantity)
		VALUES (?, ?, ?)`,
		name, status, quantity)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func GetItems(name int, filters map[int]bool) ([]ItemQuery, error) {
	var rows *sql.Rows
	var err error
	if name == -1 {
		rows, err = db.Query(`
			SELECT id, name, quantity, status
			FROM items`)
	} else {
		rows, err = db.Query(`
			SELECT id, name, quantity, status
			FROM items
			WHERE name = ?`,
			name)
	}

	if err != nil {
		log.Print(err)
		return []ItemQuery{}, err
	}

	var queries []ItemQuery
	defer rows.Close()
	for rows.Next() {
		var query ItemQuery
		rows.Scan(&query.Id, &query.Name, &query.Quantity, &query.Status)
		if _, ok := filters[query.Status]; !ok {
			continue
		}

		queries = append(queries, query)
	}

	return queries, nil
}

func AddItemName(name string) error {
	_, err := db.Exec(`
		INSERT INTO item_names(name)
		VALUES (?)`,
		name)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func GetItemName(id int) (string, error) {
	row := db.QueryRow(`
		SELECT name
		FROM item_names
		WHERE item_names.id = ?`,
		id)

	err := row.Err()
	if err != nil {
		log.Println("GetItemName:", err)
		return "", err
	}

	var name string
	row.Scan(&name)

	return name, nil
}

func GetItemNameId(name string) (int, error) {
	row := db.QueryRow(`
		SELECT id
		FROM item_names
		WHERE item_names.name = ?`,
		name)

	err := row.Err()
	if err != nil {
		log.Println("GetItemNameId:", err)
		return -1, err
	}

	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Println("GetItemNameId:", err)
		return -1, err
	}

	return id, nil
}

func AddItemStatus(name string) error {
	_, err := db.Exec(`
		INSERT INTO item_statuses(name)
		VALUES (?)`,
		name)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func GetItemStatus(id int) (string, error) {
	row := db.QueryRow(`
		SELECT name
		FROM item_statuses
		WHERE item_statuses.id = ?`,
		id)

	err := row.Err()
	if err != nil {
		log.Println(err)
		return "", err
	}

	var status string
	row.Scan(&status)

	return status, nil
}

func AddUserItem(id, quantity, user int) error {
	_, err := db.Exec(`
		INSERT INTO items_users(id, quantity, user)
		VALUES(?, ?, ?)`,
		id, quantity, user)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func GetUserItems(user int, name string, filters map[int]bool) ([]ItemQuery, error) {
	rows, err := db.Query(`
		SELECT items.id, items.name, items.quantity, items.status
		FROM items, items_names 
		WHERE (items.name = ? AND ? <> '' OR ? = '' )
		AND items_names.user = ?
		AND items_names.id = items.id`,
		name, user)

	if err != nil {
		log.Print(err)
		return []ItemQuery{}, err
	}

	var queries []ItemQuery
	defer rows.Close()
	for rows.Next() {
		var query ItemQuery
		rows.Scan(&query.Id, &query.Name, &query.Quantity, &query.Status)
		if _, ok := filters[query.Status]; !ok {
			continue
		}

		queries = append(queries, query)
	}

	return queries, nil
}
