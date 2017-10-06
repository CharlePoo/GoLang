package main

//C:\GitHub\GoLang\src\hello
import (
	"encoding/gob"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	//"encoding/json"
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
	UserName  string
	Email     string
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

	gob.Register(&UserDetails{})
	gob.Register(&M{})

	//mux = http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/", foo)
	router.HandleFunc("/auth", login)

	http.ListenAndServe(":8080", router)

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

func apiRouters() {
	router.HandleFunc("/api/login", apiLogin)
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
