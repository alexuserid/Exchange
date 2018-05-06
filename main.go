package main

import (
	"log"
	"net/http"
	"text/template"
)

func templateParseAndExecute(file string, w http.ResponseWriter) {
	t, err := template.ParseFiles(file)
	if err != nil {
		log.Printf("template.ParseFiles: %v", err)
	}
	t.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templateParseAndExecute("html/index.html", w)
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateParseAndExecute("html/reg.html", w)
	}
	if r.Method == "POST" {
		err := r.ParseForm
		if err != nil {
			log.Printf("reg: r.ParseForm: %v", err)
		}
		errf := newUser(r.Form["email"], r.Form["password"], w)
		if errf != nil {
			log.Printf("reg: newUser: %v", errf)
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateParseAndExecute("html/login.html", w)
	}
	if r.Method == "POST" {
		err := r.ParseForm
		if err != nil {
			log.Printf("login: r.ParseForm: %v", err)
		}
		sid, errf := newSid(r.Form["email"], r.Form["password"], w)
		if errf != nil {
			log.Printf("login: newSid: %v", errf)
		}

		cookieLogin := http.Cookie{Name: "sessionid", Value: sid, Path: "/", MaxAge: 86400}
		http.SetCookie(w, &cookieLogin)
	}
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Trade page will be here soon"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/trade", tradeHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
