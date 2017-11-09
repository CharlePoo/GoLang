package main

//C:\GitHub\GoLang\src\hello
import (
	"encoding/gob"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	//"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

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
