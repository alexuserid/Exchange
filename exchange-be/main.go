package main

import (
	"container/heap"
	"encoding/json"
	"log"
	"net/http"
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
	if err := newUser(r.Form["email"], r.Form["password"]); err != errNo {
		l := true
		code := http.StatusInternalServerError
		if errc, ok := err.(errorc); ok {
			l = errc.Log
			code = errc.Code
		}
		w.WriteHeader(code)
		if l {
			log.Printf("reg: newUser: %v", err)
		} else {
			json.NewEncoder(w).Encode(err)
		}
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("login: r.ParseForm: %v", err)
			return
		}
		sid, errc := newSid(r.Form["email"], r.Form["password"])
		if errc != errNo {
			w.WriteHeader(errc.Code)
			json.NewEncoder(w).Encode(errc.Text)
			if errc.Log {log.Printf("login: newSid: %v", errc.Text)}
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
	if err != errNo {
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(err.Text)
		if err.Log {log.Printf("dw: getUserInfo: %v", err)}
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
		err := dw(userInfo, p.Get("operation"), p.Get("currency"), p.Get("amount"))
		if err != errNo {
			w.WriteHeader(err.Code)
			json.NewEncoder(w).Encode(err.Text)
			if err.Log {log.Printf("dw: %v", err.Text)}
			return
		}
	}
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := getUserInfo(r)
	if err != errNo {
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(err.Text)
		if err.Log {log.Printf("dw: getUserInfo: %v", err.Text)}
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
			if err != errNo {
				w.WriteHeader(err.Code)
				json.NewEncoder(w).Encode(err.Text)
				if err.Log {log.Printf("limitOrder: %v", err.Text)}
			}
		case "market" :
			err := marketOrder(userInfo, p.Get("pair"), p.Get("amount"))
			if err != errNo {
				w.WriteHeader(err.Code)
				json.NewEncoder(w).Encode(err.Text)
				if err.Log {log.Printf("marketOrder: %v", err.Text)}
			}
		case "cancel":
			err := cancelOrder(userInfo, p.Get("pair"), p.Get("oid"))
			if err != errNo {
				w.WriteHeader(err.Code)
				json.NewEncoder(w).Encode(err.Text)
				if err.Log {log.Printf("cancelOrder: %v", err.Text)}
			}
		}
	}
}

func main() {
	heap.Init(&oq)
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
