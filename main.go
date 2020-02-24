package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	from := "...@gmail.com"
	pass := "..."
	to := "...@gmail.com"
	subject := "La Ley Lo"
	body := "Nabün len? Go ile yazdığım uygulama üzerinden gönderdim bu maili."
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("SMTP Error: %s", err)
		fmt.Fprintf(w, "SMTP Error, check logs")
		return
	}

	log.Print("Email Sent!")
	fmt.Fprintf(w, "Email sent baby!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))
}
