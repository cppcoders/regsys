package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func staffRegister(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "staffRegister.html")
		return
	}
	username := req.FormValue("username")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	email := req.FormValue("email")
	birthdate := req.FormValue("birthdate")
	gender := req.FormValue("gender")
	password := req.FormValue("password")
	course1 := req.FormValue("course1")
	course2 := req.FormValue("course2")
	course3 := req.FormValue("course3")

	var user string

	err := db.QueryRow("SELECT username FROM staff WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO staff(username ,firstname ,lastname ,email ,birthdate ,gender ,password ,course1 ,course2 ,course3) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", username, firstname, lastname, email, birthdate, gender, hashedPassword, course1, course2, course3)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Error(res, "User Already Exists.", 500)
	}
}
func courseRegister(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "courseRegister.html")
		return
	}
}
func studentRegister(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "studentRegister.html")
		return
	}

}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databaseFirstname string
	var databaseLastname string
	var databaseEmail string
	var databaseBirthdate string
	var databaseGender string
	var databasePassword string

	err := db.QueryRow("SELECT username ,firstname ,lastname ,email ,birthdate ,gender ,password FROM staff, students WHERE username=?", username).Scan(&databaseUsername, &databaseFirstname, &databaseLastname, &databaseEmail, &databaseBirthdate, &databaseGender, &databasePassword)

	if err != nil {
		http.Error(res, "User Doesn't Exist.", 500)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Error(res, "Wrong Password.", 500)
		return
	}

	res.Write([]byte("Hello " + databaseFirstname + " " + databaseLastname + "!\nHow are you doing today\nYour data is as Following\nUsername : " + databaseUsername + "\nEmail : " + databaseEmail + "\nBirthdate : " + databaseBirthdate + "\nGender : " + databaseGender))

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

func main() {
	if _, err := os.Stat("database.db"); err != nil {
		file, _ := os.Create("database.db")
		file.Close()
	}
	db, _ = sql.Open("sqlite3", "database.db")
	db.Exec("CREATE TABLE IF NOT EXISTS staff (id INTEGER PRIMARY KEY AUTOINCREMENT ,username VARCHAR(50) NOT NULL ,firstname VARCHAR(50) NOT NULL ,lastname VARCHAR(50) NOT NULL ,email VARCHAR(50) NOT NULL ,birthdate VARCHAR(50) NOT NULL ,gender VARCHAR(50) NOT NULL ,password VARCHAR(120) NOT NULL, course1 VARCHAR(50), course2 VARCHAR(50))")
	db.Exec("CREATE TABLE IF NOT EXISTS students (id INTEGER PRIMARY KEY AUTOINCREMENT ,username VARCHAR(50) NOT NULL ,firstname VARCHAR(50) NOT NULL ,lastname VARCHAR(50) NOT NULL ,email VARCHAR(50) NOT NULL ,birthdate VARCHAR(50) NOT NULL ,gender VARCHAR(50) NOT NULL ,password VARCHAR(120) NOT NULLL, course1 VARCHAR(50), course2 VARCHAR(50)L, course3 VARCHAR(50), course4 VARCHAR(50)L, course5 VARCHAR(50), course6 VARCHAR(50))")
	db.Exec("CREATE TABLE IF NOT EXISTS courses (id INTEGER PRIMARY KEY AUTOINCREMENT ,coursecode VARCHAR(50) NOT NULL ,coursename VARCHAR(50) NOT NULL)")
	defer db.Close()
	fmt.Println("http://localhost:5000/")
	http.HandleFunc("/staffRegister", staffRegister)
	http.HandleFunc("/studentRegister", studentRegister)
	http.HandleFunc("/courseRegister", courseRegister)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":5000", nil)

}
