package main

//C:\GitHub\GoLang\src\hello
import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
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

var tpl *template.Template
var store = sessions.NewCookieStore([]byte("storage-login"))

func init() {

	tpl = template.Must(template.ParseGlob("templates/*"))

}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", foo)

	http.ListenAndServe(":8080", mux)
}

func foo(w http.ResponseWriter, req *http.Request) {

	session, _ := store.Get(req, "storage-login")
	session.Values["testSession"] = "test session"
	session.Save(req, w)

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
