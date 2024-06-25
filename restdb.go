package restdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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


// FromJSON decodes a serialized JSON record - User{}
func (p *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(r)
}

// ToJSON encodes a User JSON record
func (p *User) ToJSON (r io.Writer) error {
	e := json.NewEncoder(r)
	return e.Encode(r)
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

// FindUserUsername is for returning a user record defined by a username
func FindUserUsername(username string) User {
	db := ConnectPostgres()

	if db == nil {
		log.Println("Cannot connect to PostgreSQL!")
		db.Close()
		return User{}
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where Username = $1 \n", username)

	if err != nil {
		log.Println("FindUserUsername Query:", err)
		db.Close()
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

// ReturnLoggedUsers is for returning all logged in users
func ReturnLoggedUsers() []User {
	db := ConnectPostgres()

	if db == nil {
		log.Println("Cannot connect to PostgreSQL!")
		db.Close()
		return []User{}

	}

	defer db.Close()

	rows, err := db.Query("SELECT * from users WHERE Active=1 \n")
	if err != nil {
		log.Println(err)
		db.Close()
		return []User{}
	}

	defer rows.Close()

	all := []User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6 int

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)

		if err != nil {
			log.Println(err)
			db.Close()
			return []User{}
		}

		temp := User{c1, c2, c3, c4, c5, c6}
		log.Println("temp:", all)
		all = append(all, temp)
	}

	return all

}

// IsUserAdmin determines whether a user is
// an administrator or not

func IsUserAdmin(u User) bool {
	db := ConnectPostgres()

	if db == nil {
		log.Println("Cannot connect to PostgreSQL!")
		db.Close()
		return false
	}

	db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE Username=$1 \n", u.Username)

	if err != nil {
		log.Println(err)
		db.Close()
		return false
	}

	temp := User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6 int

	// If there exist multiple users with the same username,
	// we will get the FIRST ONE only.
	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6)
		if err != nil {
			log.Println(err)
			db.Close()
			return false
		}

		temp = User{c1, c2, c3, c4, c5, c6}

	}

	defer rows.Close()

	if u.Username == temp.Username && u.Password == temp.Password && u.Admin == 1 {
		return true
	}

	return false

}

// UpdateUser allows you to update user name
func UpdateUser(u User) bool {

	db := ConnectPostgres()

	if db == nil {
		fmt.Println("Cannot connect to PostgreSQL!")
		db.Close()
		return false
	}

	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET Username=$1, Password=$2, Admin=$3, Active=$4 WHERE id=$5")

	if err != nil {
		log.Println("Adduser:", err)
		return false
	}

	res, err := stmt.Exec(u.Username, u.Password, u.Admin, u.Active, u.ID)

	if err != nil {
		log.Println("UpdateUser failed:", err)
		return false
	}

	affect, err := res.RowsAffected()

	if err != nil {
		log.Println("RowsAffected() failed:", err)
		return false
	}

	log.Println("Affected:", affect)
	return true
}

func InsertUser(u User) bool {
	db := ConnectPostgres()

	if db == nil {
		fmt.Println("Cannot connect to PostgreSQL!")
		return false
	}

	defer db.Close()

	if IsUserValid(u) {
		log.Println("User", u.Username, "already exists!")
		return false
	}

	stmt, err := db.Prepare("INSERT INTO users (Username, Password, LastLogin, Admin, Active) values($1, $2, $3, $4, $5)")
	if err != nil {
		log.Println("Adduser:", err)
		return false
	}

	stmt.Exec(u.Username, u.Password, u.LastLogin, u.Admin, u.Active)
	return true

}
