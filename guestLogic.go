package main

import (
	"bufio"
	"fmt"

	"html/template"
	"log"
	"net/http"
	"os"
)

type Guestbook struct {
	SignatureCount int
	Signature      []string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	Request()
}

func Request() {
	http.HandleFunc("/guestbook/", viewHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	http.HandleFunc("/guestbook/new", nevHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	check(err)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	signature := r.FormValue("signature")
	file, err := os.OpenFile("guestfile.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0666))
	if err != nil {
		file.Close()
	}
	check(err)
	_, err = fmt.Fprintf(file, signature)
	check(err)
	http.Redirect(w, r, "/guestbook", http.StatusFound)

}
func viewHandler(w http.ResponseWriter, r *http.Request) {

	signature := getString("guestfile.txt")
	html, err := template.ParseFiles("view.html")
	check(err)
	guestbook := Guestbook{
		SignatureCount: len(signature),
		Signature:      signature,
	}
	err = html.Execute(w, guestbook)
	check(err)
}
func getString(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil

	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines

}
func nevHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("new.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)

}
