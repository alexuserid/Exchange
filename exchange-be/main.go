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

func getUserInfo(r *http.Request) (user, error) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		return user{}, err
	}
	uid := mapSidUid[stringToB32(cookie.Value)]
	return mapUidUser[uid.id], nil
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := getUserInfo(r)
	if err != nil {
		log.Printf("trade: getUserInfo: %v", err)
	}
	log.Println(userInfo)
}

func dwHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := getUserInfo(r)
	if err != nil {
		log.Printf("dw: getUserInfo: %v", err)
	}

	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode(userInfo.wallet)
		if err != nil {
			log.Printf("dw: json.NewEmcoder(w).Encode(userInfo.wallet)")
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("de: r.ParseForm: %v", err)
		}
		if r.Form["deposiLt"] != nil {
//			deposit
		}
		if r.Form["withdraw"] != nil {
//			withdraw
		}
	}
}

func main() {
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/trade", tradeHandler)
	http.HandleFunc("/dw", dwHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
