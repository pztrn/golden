package ui

import (
	"fmt"
	"github.com/vit1251/golden/pkg/stat"
	"github.com/vit1251/golden/pkg/tosser"
	"log"
	"net/http"
)

type NetmailComposeCompleteAction struct {
	Action
}

func NewNetmailComposeCompleteAction() (*NetmailComposeCompleteAction) {
	nm := new(NetmailComposeCompleteAction)
	return nm
}

func (self *NetmailComposeCompleteAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var tosserManager *tosser.TosserManager
	var statManager *stat.StatManager
	self.Container.Invoke(func(tm *tosser.TosserManager, sm *stat.StatManager) {
		tosserManager = tm
		statManager = sm
	})

	/* Parse */
	err1 := r.ParseForm()
	if err1 != nil {
		panic(err1)
	}

	/* Create netmail message */
	to := r.PostForm.Get("to")
	to_addr := r.PostForm.Get("to_addr")
	subj := r.PostForm.Get("subject")
	body := r.PostForm.Get("body")
	log.Printf("Compose netmail: to = %s subj = %s body = %s", to, subj, body)

	//
	nm := tosserManager.NewNetmailMessage()
	nm.Subject = subj
	nm.Body = body
	nm.To = to
	nm.ToAddr = to_addr

	/* Delivery message */
	if err := tosserManager.WriteNetmailMessage(nm); err != nil {
		log.Printf("Fail on WriteNetmailMessage: err = %+v", err)
	}

	/* Register packet */
	if err := statManager.RegisterOutPacket(); err != nil {
		log.Printf("Fail on RegisterInPacket: err = %+v", err)
	}
	if err := statManager.RegisterOutMessage(); err != nil {
		log.Printf("Fail on RegisterOutMessage: err = %+v", err)
	}

	/* Redirect */
	newLocation := fmt.Sprintf("/netmail")
	http.Redirect(w, r, newLocation, 303)
}
