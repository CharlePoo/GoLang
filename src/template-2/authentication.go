package main

import (
	"log"
	"net/http"
)

func pageLogout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "login")

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var uDetails UserDetails
	uDetails.ID = 0
	session.Values["userDetails"] = uDetails
	session.Save(req, w)
	http.Redirect(w, req, "/auth", http.StatusSeeOther)
}

func getSessionUserInfo(w http.ResponseWriter, req *http.Request) *UserDetails {
	session, err := store.Get(req, "login")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	val := session.Values["userDetails"]
	//var details = &UserDetails{}
	details, _ := val.(*UserDetails)
	return details
}

func checkIfAuth(w http.ResponseWriter, req *http.Request) {
	//http://www.gorillatoolkit.org/pkg/sessions

	session, err := store.Get(req, "login")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	val := session.Values["userDetails"]
	//var details = &UserDetails{}
	details, ok := val.(*UserDetails)
	//log.Println(details)
	if !ok {
		if details == nil {
			http.Redirect(w, req, "/auth", http.StatusSeeOther)
		} else if details.ID <= 0 {
			http.Redirect(w, req, "/auth", http.StatusSeeOther)
		}
		//http.Redirect(w, req, "/auth", http.StatusSeeOther)
	} else {
		if details.ID <= 0 {
			http.Redirect(w, req, "/auth", http.StatusSeeOther)
		}
	}
}

func authentication(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "auth.htmlgo", nil)

	if err != nil {
		log.Println(err)
	}
}

func getUserByEmail(email string) UserDetails {
	db := openDB()
	rows, err := db.Query("select ID,FirstName,LastName, BirthDate, CreatedDate from user where Email=?", email)
	if err != nil {
		log.Fatal(err)
	}

	var uDetails UserDetails

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uDetails.ID, &uDetails.FirstName, &uDetails.LastName, &uDetails.BirthDate, &uDetails.CreatedDate)
		uDetails.Email = email
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	return uDetails
}
