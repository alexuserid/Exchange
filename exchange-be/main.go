package main

import (
//	"encoding/json"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
}

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

func tradeHandler(w http.ResponseWriter, r *http.Request) {
}

func dwHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sid")
	if err != nil {
		log.Printf("dw: r.Cookie: %v", err)
		return
	}
	log.Println(stringToB32(cookie.Value))
	uid := mapSidUid[stringToB32(cookie.Value)]
	log.Println(uid) //why does it print an array of zeros?
//	userInfo := mapUidUser[mapSidUid[stringToB32(cv)]]

	if r.Method == "GET" {
//		err := json.NewEncoder(w).Encode(userInfo.wallet)
//		if err != nil {
//			log.Printf("dw: json.NewEmcoder(w).Encode(userInfo.wallet)")
//		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("de: r.ParseForm: %v", err)
		}
		//think about how to identify which currency will be increased/decreased by recieved amount
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
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
