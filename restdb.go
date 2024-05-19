package restdb

import (
	"database/sql"
	"fmt"
	"log"
)

var Hostname = "localhost"
var Port = 5432
var Username = "postgres"
var Password = "postgres"
var Database = "restapi"

type User struct {
	ID        int
	Username  string
	Password  string
	LastLogin int64
	Admin     int
	Active    int
}

func ConnectPostgres() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}

// DeleteUser is for deleting a user defined by ID
func DeleteUser(ID int) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostgreSQL")
		db.Close()
	}
	defer db.Close()

	t := FindUserID(ID)

	if t.ID == 0 {
		log.Printf("User with ID %d does not exist", ID)
		return false
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE ID=$1")
	if err != nil {
		log.Println("Delete user", err)
		return false
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		log.Println("Delete user", err)
		return false
	}

	return true
}

// FindUserID is for returning a user record defined by ID
func FindUserID(ID int) User {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostgreSQL")
		db.Close()
		return User{}
	}

	defer db.Close()

	rows, err := db.Query("SELECT FROM users WHERE ID=$1\n", ID)
	if err != nil {
		log.Println("Query", err)
		return User{}
	}

	defer rows.Close()

	u := User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6 int

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
		if err != nil {
			log.Println(err)
			return User{}
		}

		u = User{c1, c2, c3, c4, c5, c6}
		log.Println("Found user:", u)
	}

	return u

}

func IsUserValid(u User) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostgreSQL")
		db.Close()
		return false
	}

	db.Close()

	rows, err := db.Query("SELECT FROM users WHERE Username = $1 \n", u.Username)
	if err != nil {
		log.Println(err)
		return false
	}

	temp := User{}

	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6 int

	for rows.Next() {
		err := rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
		if err != nil {
			log.Println(err)
			return false
		}
		temp = User{c1, c2, c3, c4, c5, c6}
	}

	if u.Username == temp.Username && u.Password == temp.Password {
		return true
	}

	return false
}

func ListAllUsers() []User {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostgreSQL")
		db.Close()
		return []User{}
	}

	db.Close()

	rows, err := db.Query("SELECT * FROM users \n")
	if err != nil {
		log.Println("Query:", err)
		return []User{}
	}

	var all = []User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6 int

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
		if err != nil {
			log.Println("Rows.Scan:", err)
		}
		temp := User{c1, c2, c3, c4, c5, c6}
		all = append(all, temp)
	}

	log.Println(all)
	return all

}
