package main

import (
	"log"
	"html/template"
	"net/http"
)

var (
	tIndex, tReg, tLogin, tTrade *template.Template
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tIndex.Execute(w, nil)
	}
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tReg.Execute(w, nil)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tLogin.Execute(w, nil)
	}
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tTrade.Execute(w, nil)
	}
}


func templParser(file string) (*template.Template) {
	tmp, err := template.ParseFiles(file)
	if err != nil {
		log.Printf("template.ParseFiles(%s): %v", file, err)
	}
	return tmp
}

func init() {
	tIndex = templParser("html/index.html")
	tReg = templParser("html/reg.html")
	tLogin = templParser("html/login.html")
	tTrade = templParser("html/trade.html")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/trade", tradeHandler)

	http.ListenAndServe(":8080", nil)
}
