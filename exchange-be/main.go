package main

import (
	"log"
	"net/http"
)

func sidChecker(r *http.Request) bool {
	_, err := r.Cookie("sid")
	if err != nil {
		return true
	}
	return false
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if sidChecker(r) {
		http.Redirect(w, r, "/trade", http.StatusSeeOther)
	}
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	if sidChecker(r) {
		http.Redirect(w, r, "/trade", http.StatusSeeOther)
	}
	if r.Method == "POST" {
		err := r.ParseForm()
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
	if sidChecker(r) {
		http.Redirect(w, r, "/trade", http.StatusSeeOther)
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("login: r.ParseForm: %v", err)
		}
		sid, errf := newSid(r.Form["email"], r.Form["password"], w)
		if errf != nil {
			log.Printf("login: newSid: %v", errf)
		}
		cookieLogin := http.Cookie{Name: "sid", Value: sid, Path: "/", MaxAge: 3600, HttpOnly: true}
		http.SetCookie(w, &cookieLogin)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "sid", MaxAge: 0})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	if !sidChecker(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/trade", tradeHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
