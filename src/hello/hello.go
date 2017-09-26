package main

//C:\GitHub\GoLang\src\hello
import (
	"log"
		"net/http"
		"fmt"
		"html/template"
	)

	//type hotdog int

	//func (h hotdog) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//	fmt.Fprint(w, "Hellow World! First Go web app! 22")
	//}

	var tpl *template.Template



	func init() {
		tpl = template.Must(template.ParseGlob("templates/*"))
	}

	func main(){
		//var x hotdog
		//http.ListenAndServe(":8080", x)
		
		mux := http.NewServeMux()
		mux.HandleFunc("/", foo)
		mux.HandleFunc("/dog", bar)
		mux.HandleFunc("/cat", barred)
		
		http.ListenAndServe(":8080", mux)

	}

	func foo(w http.ResponseWriter, req *http.Request) {
		//fmt.Fprint(w, "Hellow World! First Go web app! From foo!")
		err := tpl.ExecuteTemplate(w, "index.htmlgo",4)
		if err != nil {
			log.Println(err)
		}

	}

	func bar(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hellow World! First Go web app! Dog bar!")
	}

	func barred(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hellow World! First Go web app! Cat barred!")
	}