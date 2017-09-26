package main
//C:\GitHub\GoLang\src\hello
import(
	"log"
	"net/http"
	"text/template"
)

type course struct{
	Number, Name, Units string
}

type semester struct {
	Term string
	Courses []course
}

type year struct {
	Fall, Spring, Summer semester
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", foo)
	
	http.ListenAndServe(":8080", mux)
}

func foo(w http.ResponseWriter, req *http.Request) {
	
	y := year{
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
	}


	err := tpl.ExecuteTemplate(w, "index.htmlgo",y)
	if err != nil {
		log.Println(err)
	}

}
