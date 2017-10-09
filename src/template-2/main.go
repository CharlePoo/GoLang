package main

//C:\GitHub\GoLang\src\hello
import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

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
	ID        int32
	FirstName string
	LastName  string
	Email     string
	Password  string
	Age       int32
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
	router.HandleFunc("/auth", login)
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
		"root:myPassw0rd@tcp(127.0.0.1:3306)/myFile")
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func apiRouters() {
	router.HandleFunc("/api/login", apiLogin)
	router.HandleFunc("/api/register", CreateUserEndPoint).Methods("POST")
	router.HandleFunc("/api/register", CreateUserEndPoint).Methods("POST")
}

func CreateUserEndPoint(w http.ResponseWriter, req *http.Request) {
	var uDetails UserDetails
	_ = json.NewDecoder(req.Body).Decode(&uDetails)
	uDetails.ID = 1
	uDetails.Age = 2
	fmt.Println(uDetails.FirstName)
	//people = append(people, person)
	//json.NewEncoder(w).Encode(people)

	hash, _ := HashPassword(uDetails.Password)

	db := openDB()
	insert, err := db.Query("INSERT INTO user(ID,FirstName,LastName,Email,Age) VALUES(?, ?, ?, ?, ?  )", 3, uDetails.FirstName, uDetails.LastName, uDetails.Email, hash, 0)
	defer insert.Close()

	if err != nil {
		//log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	defer db.Close()

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
	if !ok {
		if details == nil {
			http.Redirect(w, req, "/auth", http.StatusSeeOther)
		}
		//http.Redirect(w, req, "/auth", http.StatusSeeOther)
	}

	//if details != nil {
	//	http.Redirect(w, req, "/auth", http.StatusSeeOther)
	//}
	//session.Save(req, w)
}

func apiLogin(w http.ResponseWriter, req *http.Request) {
	//json.NewEncoder(w).Encode(people)
}

func login(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "login.htmlgo", nil)

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
