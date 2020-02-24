package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type EmailType struct {
	Template string
	Subject  string
	Content  template.HTML
}

type ResponseType struct {
	Status string
	Message string
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func middleWare(handlerFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunction(w, r)
	}
}

func sendEmail(to string, subject string, body string) {
	from := os.Getenv("USERNAME")
	pass := os.Getenv("PASSWORD")
	msg :=
		"MIME-Version: 1.0;\nContent-Type: text/html;\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("SMTP Error: %s", err)
	}

	log.Print("Email Sent!")
}

func sendMailHandler(w http.ResponseWriter, r *http.Request) {
	var emailType EmailType

	// Get template file name
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Couldn't read request body.")
	}

	jsonError := json.Unmarshal(reqBody, &emailType)
	if jsonError != nil {
		log.Fatal("Json parse error")
	}

	t := parseTemplate("./templates/" + emailType.Template + ".html")

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, emailType); err != nil {
		log.Fatal("Couldn't get template html")
	}

	emailBody := tpl.String()
	go sendEmail("gotest@zahidefe.net", emailType.Subject, emailBody)
	responseType := &ResponseType{Status: "success", Message: "Email queued!"}
	response, err := json.Marshal(responseType)
	if err != nil {
		log.Fatal("Couldn't encode")
	}
	fmt.Fprintf(w, string(response))
}

func parseTemplate(templateFile string) *template.Template {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatal("Template file couldn't be parsed")
	}

	return tmpl
}

func main() {
	loadEnv()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", middleWare(sendMailHandler))
	log.Fatal(http.ListenAndServe(":8080", router))
}
