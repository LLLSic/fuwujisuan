package service

import (
	"net/http"
	"html/template"

	"github.com/unrolled/render"
)

func table_outHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method=="POST"{
			req.ParseForm()
			name := req.Form["name"][0]
			pass := req.Form["pass"][0]
			t:= template.Must(template.ParseFiles("templates/table_out.html"))
			t.Execute(w,map[string]interface{}{"name": name,"pass": pass})
		}
	}
}