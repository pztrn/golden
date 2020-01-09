package ui

import (
	"net/http"
	"github.com/gorilla/mux"
	"path/filepath"
	"html/template"
	"log"
)

type RemoveAction struct {
	Action
}

func (self *RemoveAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("views", "layout.tmpl")
	fp := filepath.Join("views", "remove.tmpl")
	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		panic(err)
	}
	//
	vars := mux.Vars(r)
	echoTag := vars["echoname"]
	log.Printf("echoTag = %v", echoTag)

	//
	areaManager := self.Site.app.GetAreaManager()
	area, err1 := areaManager.GetAreaByName(echoTag)
	if (err1 != nil) {
		panic(err1)
	}
	log.Printf("area = %v", area)

	//
	msgHash := vars["msgid"]
	msg, err2 := self.Site.app.MessageBaseReader.GetMessageByHash(echoTag, msgHash)
	if (err2 != nil) {
		panic(err2)
	}

	//
	outParams := make(map[string]interface{})
	outParams["Areas"] = areaManager.GetAreas()
	outParams["Area"] = area
	outParams["Msg"] = msg
	tmpl.ExecuteTemplate(w, "layout", outParams)

}