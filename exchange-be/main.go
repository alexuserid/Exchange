package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/starius/status"
)

func regHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("reg: r.ParseForm: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := register(r.Form.Get("email"), r.Form.Get("password")); err != nil {
		log.Printf("reg: register: %v", err)
		w.WriteHeader(status.Code(err))
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("login: r.ParseForm: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sid, err := login(r.Form.Get("email"), r.Form.Get("password"))
	if err != nil {
		log.Printf("login: login: %v", err)
		w.WriteHeader(status.Code(err))
		return
	}
	cookieLogin := http.Cookie{Name: "sid", Value: sid, Path: "/", MaxAge: 3600, HttpOnly: true}
	http.SetCookie(w, &cookieLogin)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "sid", MaxAge: 0})
}

func dwHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := getUserInfo(r)
	if err != nil {
		log.Printf("dw: getUserInfo: %v", err)
		w.WriteHeader(status.Code(err))
		return
	}
	if r.Method == "GET" {
		if err := json.NewEncoder(w).Encode(userInfo.money); err != nil {
			log.Printf("dw: json.NewEncoder(w).Encode(userInfo.wallet): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if r.Method == "POST" {
		p := r.URL.Query()
		if err := depositAndWithdraw(userInfo, p.Get("operation"), p.Get("currency"), p.Get("amount")); err != nil {
			log.Printf("dw: depositAndWithdraw: %v", err)
			w.WriteHeader(status.Code(err))
			return
		}
	}
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := getUserInfo(r)
	if err != nil {
		log.Println("trade: userInfo: %v", err)
		w.WriteHeader(status.Code(err))
		return
	}
	if r.Method == "GET" {
		if err := json.NewEncoder(w).Encode(userInfo); err != nil {
			log.Printf("trade: json.NewEncoder(w).Encode(userInfo)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if r.Method == "POST" {
		p := r.URL.Query()
		switch p.Get("order") {
		case "limit":
			if err := limitOrder(userInfo, p.Get("pair"), p.Get("amount"), p.Get("price")); err != nil {
				log.Printf("trade: limitOrder: %v", err)
				w.WriteHeader(status.Code(err))
			}
		case "market":
			if errc := marketOrder(userInfo, p.Get("pair"), p.Get("amount")); errc != nil {
				log.Printf("trade: markerOrder: %v", err)
				w.WriteHeader(status.Code(err))
			}
		case "cancel":
			if errc := cancelOrder(userInfo, p.Get("pair"), p.Get("oid")); errc != nil {
				log.Printf("trade: cancelOrder: %v", err)
				w.WriteHeader(status.Code(err))
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
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
