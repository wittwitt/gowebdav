package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// AddToolsRouter tools router,
// pathPrefix /xxx
func AddToolsRouter(pathPrefix string, router *httprouter.Router) {
	router.GET(pathPrefix, Index)
	router.GET(pathPrefix+"/notepad", Notepad)
	router.ServeFiles(pathPrefix+"/assets/*filepath", http.Dir("static/tools/"))
}

// Index page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Xx page
func Xx(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "xx!\n")
}

// Hello router
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

//
var notepadTxt map[string][]string = make(map[string][]string)

// Notepad list
func Notepad(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qParams := r.URL.Query()
	if _, ok := qParams["clear"]; ok {
		notepadTxt = make(map[string][]string)
	}

	for k, v := range qParams {
		notepadTxt[k] = v
	}
	fmt.Fprintln(w, "?clear=ture to clear./r/n", notepadTxt)
}
