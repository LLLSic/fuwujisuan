package service

import (
	"net/http"
	"html/template"

	"github.com/unrolled/render"
)

func table_inHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method=="GET"{
			t:= template.Must(template.ParseFiles("templates/table_in.html"))
			t.Execute(w,nil)
		}
	}
}