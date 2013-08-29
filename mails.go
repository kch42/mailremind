package main

import (
	"bytes"
	"fmt"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"path"
	"text/template"
	"time"
)

func loadMailTpl(tplroot, name string) *template.Template {
	tpl, err := template.ParseFiles(path.Join(tplroot, name+".tpl"))
	if err != nil {
		log.Fatalf("Could not load mailtemplate %s: %s", name, err)
	}
	return tpl
}

var mailActivationcode *template.Template

func initMails() {
	tplroot, err := conf.GetString("paths", "mailtpls")
	if err != nil {
		log.Fatalf("Could not get paths.mailtpls from config: %s", err)
	}

	mailActivationcode = loadMailTpl(tplroot, "activationcode")
}

type activationcodeData struct {
	URL string
}

func SendActivationcode(to, acCode string, uid model.DBID) bool {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "To: %s\n", to)
	fmt.Fprintf(buf, "From: %s\n", MailFrom)
	fmt.Fprintf(buf, "Subject: Activation code for your mailremind account\n")
	fmt.Fprintf(buf, "Date: %s\n", time.Now().Format(time.RFC822))

	fmt.Fprintln(buf, "")

	url := fmt.Sprintf("%s/activate/U=%s&Code=%s", baseurl, uid, acCode)
	if err := mailActivationcode.Execute(buf, activationcodeData{url}); err != nil {
		log.Printf("Error while executing mail template (activationcode): %s", err)
		return false
	}

	return Mail(to, MailFrom, buf.Bytes())
}