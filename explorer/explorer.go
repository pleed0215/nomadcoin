package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	bc "github.com/pleed0215/nomadcoin/blockchain"
)

var port string
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
	data := homeData{"Home", bc.BC().AllBlocks()}
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
			bc.BC().AddBlock(data)
			http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
	//tmpl.Execute(rw, data)
}

func Start(aPort int) {
	newServer := http.NewServeMux()
	port = fmt.Sprintf(":%d", aPort)
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir+ "partials/*.gohtml"))

	newServer.HandleFunc("/", home)
	newServer.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, newServer)) 
}