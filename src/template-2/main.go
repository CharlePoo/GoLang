package main

//C:\GitHub\GoLang\src\hello
import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	//"encoding/json"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type course struct {
	Number, Name, Units string
}

type semester struct {
	Term    string
	Courses []course
}

type year struct {
	Fall, Spring, Summer semester
}

type UserDetails struct {
	ID          int32
	FirstName   string
	LastName    string
	Email       string
	Password    string
	BirthDate   time.Time
	CreatedDate time.Time
}

type M map[string]interface{}

var tpl *template.Template
var store = sessions.NewCookieStore([]byte("myStorage"))

var router = mux.NewRouter()

func init() {

	tpl = template.Must(template.ParseGlob("templates/*"))

}

func main() {

	//This will initiallize session
	gob.Register(&UserDetails{})
	gob.Register(&M{})

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", foo)
	router.HandleFunc("/auth", authentication)
	router.HandleFunc("/logout", pageLogout)
	apiRouters()

	http.ListenAndServe(":8080", router)

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func openDB() *sql.DB {

	db, err := sql.Open("mysql",
		"root:myPassw0rd@tcp(127.0.0.1:3306)/myFile?parseTime=true")
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func apiRouters() {
	//router.HandleFunc("/api/login", apiLogin)
	router.HandleFunc("/api/register", CreateUserEndPoint).Methods("POST")
	router.HandleFunc("/api/login", LoginEndPoint).Methods("POST")
}

func CreateUserEndPoint(w http.ResponseWriter, req *http.Request) {
	var uDetails UserDetails
	_ = json.NewDecoder(req.Body).Decode(&uDetails)

	//people = append(people, person)
	//json.NewEncoder(w).Encode(people)
	temp := getUserByEmail(uDetails.Email)

	if temp.FirstName == "" {
		passwordHash, _ := HashPassword(uDetails.Password)
		fmt.Println(passwordHash)
		db := openDB()
		insert, err := db.Query("INSERT INTO user(FirstName,LastName,Email,BirthDate,Password,CreatedDate) VALUES(?, ?, ?, ?, ?, NOW() )", uDetails.FirstName, uDetails.LastName, uDetails.Email, uDetails.BirthDate, string(passwordHash))
		defer insert.Close()
		defer db.Close()
		if err != nil {
			//log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uDetails = getUserByEmail(uDetails.Email)
		session, err := store.Get(req, "login")
		session.Values["userDetails"] = uDetails
		session.Save(req, w)
	}

	json.NewEncoder(w).Encode(uDetails)
}

func LoginEndPoint(w http.ResponseWriter, req *http.Request) {
	var uDetails UserDetails
	_ = json.NewDecoder(req.Body).Decode(&uDetails)
	//uDetails = getUserByEmail(uDetails.Email)
	_tempPassword := uDetails.Password

	db := openDB()
	rows, err := db.Query("select ID,FirstName,LastName, BirthDate, CreatedDate, Password from user where Email=?", uDetails.Email)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uDetails.ID, &uDetails.FirstName, &uDetails.LastName, &uDetails.BirthDate, &uDetails.CreatedDate, &uDetails.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if !CheckPasswordHash(_tempPassword, uDetails.Password) {
		fmt.Println("wrong password")
		return
	}

	session, err := store.Get(req, "login")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["userDetails"] = uDetails
	session.Save(req, w)
	json.NewEncoder(w).Encode(uDetails)
}

func LogoutEndPoint(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "login")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var uDetails UserDetails
	session.Values["userDetails"] = uDetails
	session.Save(req, w)
}

func pageLogout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "login")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var uDetails UserDetails
	uDetails.ID = 0
	session.Values["userDetails"] = uDetails
	session.Save(req, w)
	http.Redirect(w, req, "/auth", http.StatusSeeOther)
}

func getUserByEmail(email string) UserDetails {
	db := openDB()
	rows, err := db.Query("select ID,FirstName,LastName, BirthDate, CreatedDate from user where Email=?", email)
	if err != nil {
		log.Fatal(err)
	}

	var uDetails UserDetails

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uDetails.ID, &uDetails.FirstName, &uDetails.LastName, &uDetails.BirthDate, &uDetails.CreatedDate)
		uDetails.Email = email
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	return uDetails
}

func checkIfAuth(w http.ResponseWriter, req *http.Request) {
	//http://www.gorillatoolkit.org/pkg/sessions

	session, err := store.Get(req, "login")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	val := session.Values["userDetails"]
	//var details = &UserDetails{}
	details, ok := val.(*UserDetails)
	log.Println(ok)
	if !ok {
		if details == nil {
			http.Redirect(w, req, "/auth", http.StatusSeeOther)
		} else if details.ID <= 0 {
			http.Redirect(w, req, "/auth", http.StatusSeeOther)
		}
		//http.Redirect(w, req, "/auth", http.StatusSeeOther)
	} else {
		if details.ID <= 0 {
			http.Redirect(w, req, "/auth", http.StatusSeeOther)
		}
	}
}

func apiLogin(w http.ResponseWriter, req *http.Request) {
	//json.NewEncoder(w).Encode(people)
}

func authentication(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "auth.htmlgo", nil)

	if err != nil {
		log.Println(err)
	}
}

func foo(w http.ResponseWriter, req *http.Request) {

	checkIfAuth(w, req)
	/*y := year{
		Fall: semester{
			Term: "Fall",
			Courses: []course{
				{"Fall 1", "Hello 1", "unit 1"},
				{"Fall 2", "Hello 2", "unit 2"},
				{"Fall 3", "Hello 3", "unit 3"},
			},
		},
		Spring: semester{
			Term: "Spring",
			Courses: []course{
				{"Spring 1", "Hello 1", "unit 1"},
				{"Spring 2", "Hello 2", "unit 2"},
				{"Spring 3", "Hello 3", "unit 3"},
			},
		},
	}*/

	err := tpl.ExecuteTemplate(w, "index.htmlgo", nil)

	if err != nil {
		log.Println(err)
	}

}
