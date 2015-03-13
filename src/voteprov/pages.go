package voteprov

import (
    "net/http"
	"html/template"
	//"log"
)


var templates = template.Must(template.ParseGlob("templates/*"))


func Home(rw http.ResponseWriter, r *http.Request) {
	dc := DefaultContext{}
	AuthContext(rw, r, &dc)
	//log.Println("Session data: ", session.Values)
    err := templates.ExecuteTemplate(rw, "home", &dc) // nil arg is context
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}
