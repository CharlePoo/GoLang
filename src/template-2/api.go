package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CreateUserEndPoint(w http.ResponseWriter, req *http.Request) {
	var uDetails UserDetails
	_ = json.NewDecoder(req.Body).Decode(&uDetails)

	//people = append(people, person)
	//json.NewEncoder(w).Encode(people)
	temp := getUserByEmail(uDetails.Email)

	if temp.FirstName == "" {
		passwordHash, _ := HashPassword(uDetails.Password)
		fmt.Println(passwordHash)
		db := openDB()
		insert, err := db.Query("INSERT INTO user(FirstName,LastName,Email,BirthDate,Password,CreatedDate) VALUES(?, ?, ?, ?, ?, NOW() )", uDetails.FirstName, uDetails.LastName, uDetails.Email, uDetails.BirthDate, string(passwordHash))
		defer insert.Close()
		defer db.Close()
		if err != nil {
			//log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uDetails = getUserByEmail(uDetails.Email)
		session, err := store.Get(req, "login")
		session.Values["userDetails"] = uDetails
		session.Save(req, w)
	}

	json.NewEncoder(w).Encode(uDetails)
}

func LoginEndPoint(w http.ResponseWriter, req *http.Request) {
	var uDetails UserDetails
	_ = json.NewDecoder(req.Body).Decode(&uDetails)
	//uDetails = getUserByEmail(uDetails.Email)
	_tempPassword := uDetails.Password

	db := openDB()
	rows, err := db.Query("select ID,FirstName,LastName, BirthDate, CreatedDate, Password from user where Email=?", uDetails.Email)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uDetails.ID, &uDetails.FirstName, &uDetails.LastName, &uDetails.BirthDate, &uDetails.CreatedDate, &uDetails.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if !CheckPasswordHash(_tempPassword, uDetails.Password) {
		fmt.Println("wrong password")
		return
	}

	session, err := store.Get(req, "login")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["userDetails"] = uDetails
	session.Save(req, w)
	json.NewEncoder(w).Encode(uDetails)
}

func LogoutEndPoint(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "login")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var uDetails UserDetails
	session.Values["userDetails"] = uDetails
	session.Save(req, w)
}
