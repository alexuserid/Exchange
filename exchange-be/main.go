package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func regHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("reg: r.ParseForm: %v", err)
			return
		}
		errf := newUser(r.Form["email"], r.Form["password"], w)
		if errf != nil {
			log.Printf("reg: newUser: %v", errf)
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("login: r.ParseForm: %v", err)
			return
		}
		sid, errf := newSid(r.Form["email"], r.Form["password"], w)
		if errf != nil {
			log.Printf("login: newSid: %v", errf)
			return
		}
		cookieLogin := http.Cookie{Name: "sid", Value: sid, Path: "/", MaxAge: 3600, HttpOnly: true}
		http.SetCookie(w, &cookieLogin)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "sid", MaxAge: 0})
}

func dwHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := getUserInfo(r)
	if err != nil {
		log.Printf("dw: getUserInfo: %v", err)
		return
	}
	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode(userInfo.money)
		if err != nil {
			log.Printf("dw: json.NewEmcoder(w).Encode(userInfo.wallet)")
			return
		}
	}
	if r.Method == "POST" {
		p := r.URL.Query()
		err := dw(w, userInfo, p.Get("operation"), p.Get("currency"), p.Get("amount"))
		if err != nil {
			log.Printf("dw: %v", err)
			return
		}
	}
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := getUserInfo(r)
	if err != nil {
		log.Printf("trade: getUserInfo: %v", err)
		return
	}
	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode(userInfo)
		if err != nil {
			log.Printf("trade: json.NewEncoder(w).Encode(userInfo)")
			return
		}
	}
	if r.Method == "POST" {
		p := r.URL.Query()
		switch p.Get("order") {
		case "limit" :
			err := limitOrder(userInfo, p.Get("pair"), p.Get("amount"), p.Get("price"))
			if err != nil {
				log.Printf("limitOrder: %v", err)
			}
		case "now" :
			err := nowOrder(userInfo, p.Get("pair"), p.Get("amount"))
			if err != nil {
				log.Printf("nowOrder: %v", err)
			}
		case "cancel":
			err := cancelOrder(userInfo, p.Get("pair"), p.Get("oid"))
			if err != nil {
				log.Printf("cancelOrder: %v", err)
			}
		}
	}
}

func main() {
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/dw", dwHandler)
	http.HandleFunc("/trade", tradeHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
