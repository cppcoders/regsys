package main

import (
	"database/sql"
	"encoding/json"
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
			fmt.Println(err)
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO staff(username ,firstname ,lastname ,email ,birthdate ,gender ,password ,course1 ,course2 ,course3) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", username, firstname, lastname, email, birthdate, gender, hashedPassword, course1, course2, course3)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		fmt.Println(err)
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
	coursename := req.FormValue("coursename")
	coursecode := req.FormValue("coursecode")

	var course string

	err := db.QueryRow("SELECT coursename FROM courses WHERE coursename=?", coursename).Scan(&course)

	switch {
	case err == sql.ErrNoRows:

		_, err = db.Exec("INSERT INTO courses(coursename ,coursecode) VALUES(?, ?)", coursename, coursecode)
		if err != nil {
			http.Error(res, "Server error, unable to create your course.", 500)
			return
		}

		res.Write([]byte("Course created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your Course.", 500)
		return
	default:
		http.Error(res, "Course Already Exists.", 500)
	}
}
func studentRegister(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "studentRegister.html")
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
	course4 := req.FormValue("course4")
	course5 := req.FormValue("course5")
	course6 := req.FormValue("course6")

	var user string

	err := db.QueryRow("SELECT username FROM students WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO students(username ,firstname ,lastname ,email ,birthdate ,gender ,password ,course1 ,course2 ,course3, course4, course5, course6) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", username, firstname, lastname, email, birthdate, gender, hashedPassword, course1, course2, course3, course4, course5, course6)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		fmt.Println(err)
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Error(res, "User Already Exists.", 500)
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
	var databaseCourse1 string
	var databaseCourse2 string
	var databaseCourse3 string
	var databaseCourse4 string
	var databaseCourse5 string
	var databaseCourse6 string
	var student bool

	err := db.QueryRow("SELECT username FROM students WHERE username=?", username).Scan(&databaseUsername)
	if err == nil {
		student = true
		db.QueryRow("SELECT username ,firstname ,lastname ,email ,birthdate ,gender ,password ,course1 ,course2 ,course3 ,course4 ,course5 ,course6 FROM students WHERE username=?", username).Scan(&databaseUsername, &databaseFirstname, &databaseLastname, &databaseEmail, &databaseBirthdate, &databaseGender, &databasePassword, &databaseCourse1, &databaseCourse2, &databaseCourse3, &databaseCourse4, &databaseCourse5, &databaseCourse6)

	} else {
		err := db.QueryRow("SELECT username FROM staff WHERE username=?", username).Scan(&databaseUsername)
		if err != nil {
			http.Error(res, "User Doesn't Exist.", 500)
			return
		}
		db.QueryRow("SELECT username ,firstname ,lastname ,email ,birthdate ,gender ,password ,course1 ,course2 ,course3 FROM staff WHERE username=?", username).Scan(&databaseUsername, &databaseFirstname, &databaseLastname, &databaseEmail, &databaseBirthdate, &databaseGender, &databasePassword, &databaseCourse1, &databaseCourse2, &databaseCourse3)
	}

	errP := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if errP != nil {
		http.Error(res, "Wrong Password.", 500)
		return
	}
	if student {
		res.Write([]byte("Hello Student " + databaseFirstname + " " + databaseLastname + "!\nHow are you doing today\nYour data is as Following\nUsername : " + databaseUsername + "\nEmail : " + databaseEmail + "\nBirthdate : " + databaseBirthdate + "\nGender : " + databaseGender + "\nYou Have The Following Courses" + "\n\t" + databaseCourse1 + "\n\t" + databaseCourse2 + "\n\t" + databaseCourse3 + "\n\t" + databaseCourse4 + "\n\t" + databaseCourse5 + "\n\t" + databaseCourse6))
	} else {
		res.Write([]byte("Hello Dr. " + databaseFirstname + " " + databaseLastname + "!\nHow are you doing today\nYour data is as Following\nUsername : " + databaseUsername + "\nEmail : " + databaseEmail + "\nBirthdate : " + databaseBirthdate + "\nGender : " + databaseGender + "\nYou Teach The Following Courses" + "\n\t" + databaseCourse1 + "\n\t" + databaseCourse2 + "\n\t" + databaseCourse3))
	}

}

type person struct {
	name string
	age  int
}

func getCourses(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	var courses []string
	var course string
	rows, _ := db.Query("select coursename from courses")
	for rows.Next() {
		rows.Scan(&course)
		courses = append(courses, course)
	}

	json.NewEncoder(res).Encode(courses)
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
	db.Exec("CREATE TABLE IF NOT EXISTS staff (id INTEGER PRIMARY KEY AUTOINCREMENT ,username VARCHAR(50) NOT NULL ,firstname VARCHAR(50) NOT NULL ,lastname VARCHAR(50) NOT NULL ,email VARCHAR(50) NOT NULL ,birthdate VARCHAR(50) NOT NULL ,gender VARCHAR(50) NOT NULL ,password VARCHAR(120) NOT NULL, course1 VARCHAR(50), course2 VARCHAR(50), course3 VARCHAR(50))")
	db.Exec("CREATE TABLE IF NOT EXISTS students (id INTEGER PRIMARY KEY AUTOINCREMENT ,username VARCHAR(50) NOT NULL ,firstname VARCHAR(50) NOT NULL ,lastname VARCHAR(50) NOT NULL ,email VARCHAR(50) NOT NULL ,birthdate VARCHAR(50) NOT NULL ,gender VARCHAR(50) NOT NULL ,password VARCHAR(120) NOT NULL, course1 VARCHAR(50), course2 VARCHAR(50), course3 VARCHAR(50), course4 VARCHAR(50), course5 VARCHAR(50), course6 VARCHAR(50))")
	db.Exec("CREATE TABLE IF NOT EXISTS courses (id INTEGER PRIMARY KEY AUTOINCREMENT ,coursecode VARCHAR(50) NOT NULL ,coursename VARCHAR(50) NOT NULL)")
	defer db.Close()
	fmt.Println("http://localhost:5000/")
	http.HandleFunc("/staffRegister", staffRegister)
	http.HandleFunc("/studentRegister", studentRegister)
	http.HandleFunc("/courseRegister", courseRegister)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)
	http.HandleFunc("/get", getCourses)
	http.ListenAndServe(":5001", nil)

}
