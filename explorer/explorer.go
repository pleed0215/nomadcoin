package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	bc "github.com/pleed0215/nomadcoin/blockchain"
)

const port string = ":4000"
const templateDir string = "explorer/templates/"
var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks []*bc.Block
}
type addData struct {
	PageTitle string
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", bc.GetBlockchain().AllBlock()}
	//tmpl.Execute(rw, data)
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	data := addData{"Add"}
	switch(r.Method) {
		case "GET":
			templates.ExecuteTemplate(rw, "add", data)	
		case "POST":
			r.ParseForm()
			data:=r.Form.Get("blockData")
			bc.GetBlockchain().AddBlock(data)
			http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
	//tmpl.Execute(rw, data)
}

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir+ "partials/*.gohtml"))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil)) 
}