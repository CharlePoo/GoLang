package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func apiRouters() {
	//router.HandleFunc("/api/login", apiLogin)
	router.HandleFunc("/api/register", CreateUserEndPoint).Methods("POST")
	router.HandleFunc("/api/login", LoginEndPoint).Methods("POST")
	router.HandleFunc("/api/initialize", InitializeEndPoint).Methods("GET")
	router.HandleFunc("/api/getItemsByUserIdAndParentId", GetItemsByUserIdAndItemParentEndPoint).Methods("POST")
	router.HandleFunc("/api/getSubFiles", GetSubFilesEndPoint).Methods("POST")

}

func CreateUserEndPoint(w http.ResponseWriter, req *http.Request) {
	var uDetails *UserDetails
	_ = json.NewDecoder(req.Body).Decode(&uDetails)

	//people = append(people, person)
	//json.NewEncoder(w).Encode(people)
	temp := getUserByEmail(uDetails.Email)

	if temp.FirstName == "" {
		passwordHash, _ := HashPassword(uDetails.Password)

		db := openDB()
		insert, err := db.Query("INSERT INTO user(FirstName,LastName,Email,BirthDate,Password,CreatedDate) VALUES(?, ?, ?, ?, ?, NOW() )", uDetails.FirstName, uDetails.LastName, uDetails.Email, uDetails.BirthDate, string(passwordHash))
		defer insert.Close()
		defer db.Close()
		if err != nil {
			//log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_uDetails := getUserByEmail(uDetails.Email)
		uDetails.ID = _uDetails.ID
		session, _ := store.Get(req, "login")
		session.Values["userDetails"] = uDetails
		session.Save(req, w)

		CreateFolder("", uDetails, 0, 0) //this will create base folder for newly created user
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

	uDetails.Password = ""
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

func GetItemsByUserIdAndItemParentEndPoint(w http.ResponseWriter, req *http.Request) {
	var itemInfo ItemInfo
	_ = json.NewDecoder(req.Body).Decode(&itemInfo)

	itemsInfo := getItemInfo(w, req, itemInfo.Id)

	//log.Println(initializePage)
	json.NewEncoder(w).Encode(itemsInfo)
}

func GetSubFilesEndPoint(w http.ResponseWriter, req *http.Request) {
	var itemInfo ItemInfo
	_ = json.NewDecoder(req.Body).Decode(&itemInfo)

	//itemsInfo := getItemInfo(w, req, itemInfo.Id)
	uDetail := getSessionUserInfo(w, req)
	itemsInfo := ListFile(itemInfo.Path, uDetail)

	//log.Println(initializePage)
	json.NewEncoder(w).Encode(itemsInfo)
}

func getItemInfo(w http.ResponseWriter, req *http.Request, parentId int) []ItemInfo {
	uDetail := getSessionUserInfo(w, req)
	fmt.Println(uDetail)
	db := openDB()
	rows, err := db.Query("select Id,Name, ParentId, MaxUploadCount, MaxUploadSize, Public, ColumnIndex, RowIndex from files_and_folders where UserId=? and ParentId=?", 18, parentId)
	if err != nil {
		log.Fatal(err)
	}

	var itemArray []ItemInfo

	defer rows.Close()
	for rows.Next() {
		var itemInfo ItemInfo
		err := rows.Scan(&itemInfo.Id, &itemInfo.Name, &itemInfo.ParentId, &itemInfo.MaxUploadCount, &itemInfo.MaxUploadSize, &itemInfo.Public, &itemInfo.ColumnIndex, &itemInfo.RowIndex)
		itemArray = append(itemArray, itemInfo)
		fmt.Println(itemInfo)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	return itemArray
}

func InitializeEndPoint(w http.ResponseWriter, req *http.Request) {

	//uDetailList := getItemInfo(w, req, 0)
	initializePage := AppInfo{
		Desktop: DesktopSettings{
			ColumnSize: 10,
			RowSize:    10,
		},
	}

	ud := getSessionUserInfo(w, req)
	//log.Println(ud)

	initializePage.Items = ListFile("", ud)

	/*initializePage := AppInfo{
		Desktop: DesktopSettings{
			ColumnSize: 10,
			RowSize:    10,
		},
		Items: []ItemInfo{
			{
				MaxUploadCount: 2,
				MaxUploadSize:  2,
				Name:           "Folder 1",
				Public:         true,
				Id:             1,
				ParentId:       0,
				ColumnIndex:    0,
				RowIndex:       0,
			},
			{
				MaxUploadCount: 2,
				MaxUploadSize:  2,
				Name:           "Folder 2",
				Public:         true,
				Id:             2,
				ParentId:       0,
				ColumnIndex:    1,
				RowIndex:       1,
			},
		},
	}*/

	session, error := store.Get(req, "login")
	if error != nil {
		log.Fatalln(error)
	}

	val := session.Values["userDetails"]
	if uDetails, ok := val.(*UserDetails); ok {
		fmt.Println(uDetails)
		initializePage.User = *uDetails
	}

	//log.Println(initializePage)
	json.NewEncoder(w).Encode(initializePage)
}
