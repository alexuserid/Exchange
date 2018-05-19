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
	if errc := newUser(r.Form["email"], r.Form["password"]); errc != nil {
		if errc.Err != nil {
			log.Printf("reg: newUser: %v", errc)
		}
		w.WriteHeader(errc.Code)
		json.NewEncoder(w).Encode(errc.Text)
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
		if errc != nil {
			w.WriteHeader(errc.Code)
			json.NewEncoder(w).Encode(errc.Text)
			if errc.Err != nil {
				log.Printf("login: newSid: %v", errc)
			}
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
	userInfo, errc := getUserInfo(r)
	if errc != nil {
		w.WriteHeader(errc.Code)
		json.NewEncoder(w).Encode(errc.Text)
		if errc.Err != nil {
			log.Printf("dw: getUserInfo: %v", errc)
		}
		return
	}
	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode(userInfo.money)
		if err != nil {
			log.Printf("dw: json.NewEmcoder(w).Encode(userInfo.wallet): %v", err)
			return
		}
	}
	if r.Method == "POST" {
		p := r.URL.Query()
		errc := depositAndWithdraw(userInfo, p.Get("operation"), p.Get("currency"), p.Get("amount"))
		if errc != nil {
			w.WriteHeader(errc.Code)
			json.NewEncoder(w).Encode(errc.Text)
			if errc.Err != nil {
				log.Printf("dw: %v", errc)
			}
			return
		}
	}
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, errc := getUserInfo(r)
	if errc != nil {
		w.WriteHeader(errc.Code)
		json.NewEncoder(w).Encode(errc.Text)
		if errc.Err != nil {
			log.Printf("dw: getUserInfo: %v", errc)
		}
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
			errc := limitOrder(userInfo, p.Get("pair"), p.Get("amount"), p.Get("price"))
			if errc != nil {
				w.WriteHeader(errc.Code)
				json.NewEncoder(w).Encode(errc.Text)
				if errc.Err != nil {
					log.Printf("limitOrder: %v", errc)
				}
			}
		case "market" :
			errc := marketOrder(userInfo, p.Get("pair"), p.Get("amount"))
			if errc != nil {
				w.WriteHeader(errc.Code)
				json.NewEncoder(w).Encode(errc.Text)
				if errc.Err != nil {
					log.Printf("marketOrder: %v", errc)
				}
			}
		case "cancel":
			errc := cancelOrder(userInfo, p.Get("pair"), p.Get("oid"))
			if errc != nil {
				w.WriteHeader(errc.Code)
				json.NewEncoder(w).Encode(errc.Text)
				if errc.Err != nil {
					log.Printf("cancelOrder: %v", errc)
				}
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
