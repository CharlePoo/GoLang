package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func copyFile(source, destination string) {
	from, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()
}

func CreateFolder(path string, uDetails *UserDetails, ColumnIndex, RowIndex int) bool {
	basePath := "./AllFolders/" + strconv.Itoa(uDetails.ID)
	err := os.MkdirAll(basePath+"/"+path, os.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	} else if path != "" {
		var isParentPathDesktop = false
		splitedPath := strings.Split(path, ",")
		for index, val := range splitedPath {
			if index >= 3 && val != "" {
				isParentPathDesktop = true
				break
			}
		}

		if isParentPathDesktop {

			os.Mkdir(basePath+"/~settings~", os.ModePerm)

			/*db := openDB()
			insert, err := db.Query("INSERT INTO folder_settings(UserId,Path,ColumnIndex,RowIndex) VALUES(?, ?, ?, ? )", uDetails.ID, path, ColumnIndex, RowIndex)
			defer insert.Close()
			defer db.Close()
			if err != nil {
				log.Println(err)
				return false
			}*/
		}

	}

	return true
	//os.Rename("C:/testGoFolder", "C:/testGoFolder2")
	//os.Rename("C:/testGoFolder/test.txt", "C:/testGoFolder/Folder 1/test.txt")
}

func CopyOrRenameFile(oldPath, newPath string) {
	os.Rename(oldPath, newPath)
}

func FolderListFile(path string, uDetails *UserDetails) ItemInfo {

	oldPath := path
	path = strings.Replace(path, "m3g", "/", -1)
	path = strings.Replace(path, "3sp3", " ", -1)

	basePath := "./AllFolders/" + strconv.Itoa(uDetails.ID) + "/" + path
	parentArr := strings.Split(path, "/")

	var parentPath string

	for index, value := range parentArr {
		if len(parentArr) == index {
			break
		}
		parentPath = parentPath + "/" + value
	}

	parentPath = convertPath(parentPath)

	files, err := ioutil.ReadDir(basePath)

	var folderInfo ItemInfo

	filepath.Walk(basePath, func(_ string, info os.FileInfo, err error) error {

		folderInfo.Name = info.Name()
		folderInfo.IdPath = oldPath
		//itemInfo.Path = path + f.Name()
		folderInfo.Path = oldPath
		folderInfo.ParentId = 0
		folderInfo.ParentPath = parentPath
		folderInfo.IsFolder = info.IsDir()
		return err
	})

	if err != nil {
		log.Println(err)
		return folderInfo
	}

	var settings []DesktopFolderIndex
	//var itemArray []ItemInfo

	//Get json settings for desktop child folder index position
	jsonByte, error := ioutil.ReadFile(basePath + "/~settings~")
	if error != nil {
		log.Println(error)
		//return nil
	} else {
		json.Unmarshal(jsonByte, &settings)
	}

	for _, f := range files {

		if f.Name() == "~settings~" {
			continue
		}

		var itemInfo ItemInfo
		//var tempPath = strings.Replace(path+"m3g"+f.Name(), " ", "3sp3", -1) //replace space to 3sp3
		//tempPath = strings.Replace(tempPath, "/", "m3g", -1)                 //replace slash to m3g

		tempPath := convertPath(path + "m3g" + f.Name())

		itemInfo.Name = f.Name()
		itemInfo.IdPath = tempPath
		//itemInfo.Path = path + f.Name()
		itemInfo.Path = tempPath
		itemInfo.ParentId = 0
		itemInfo.ParentPath = parentPath
		itemInfo.IsFolder = f.IsDir()

		for _, val := range settings {
			if val.Name == f.Name() {
				itemInfo.ColumnIndex = val.ColumnIndex
				itemInfo.RowIndex = val.RowIndex
				break
			}
		}

		//itemArray = append(itemArray, itemInfo)
		folderInfo.Items = append(folderInfo.Items, itemInfo)
	}

	return folderInfo
}

func ListFile(path string, uDetails *UserDetails) []ItemInfo {

	path = strings.Replace(path, "m3g", "/", -1)
	path = strings.Replace(path, "3sp3", " ", -1)

	basePath := "./AllFolders/" + strconv.Itoa(uDetails.ID) + "/" + path
	parentArr := strings.Split(path, "/")

	var parentPath string

	for index, value := range parentArr {
		if len(parentArr) == index {
			break
		}
		parentPath = parentPath + "/" + value
	}

	parentPath = convertPath(parentPath)

	log.Println(basePath)
	files, err := ioutil.ReadDir(basePath)

	log.Println(files)
	if err != nil {
		log.Println(err)
		return nil
	}

	var settings []DesktopFolderIndex
	var itemArray []ItemInfo

	//Get json settings for desktop child folder index position
	jsonByte, error := ioutil.ReadFile(basePath + "/~settings~")
	if error != nil {
		log.Println(error)
		//return nil
	} else {
		json.Unmarshal(jsonByte, &settings)
	}

	for _, f := range files {

		if f.Name() == "~settings~" {
			continue
		}

		var itemInfo ItemInfo
		//var tempPath = strings.Replace(path+"m3g"+f.Name(), " ", "3sp3", -1) //replace space to 3sp3
		//tempPath = strings.Replace(tempPath, "/", "m3g", -1)                 //replace slash to m3g

		tempPath := convertPath(path + "m3g" + f.Name())

		itemInfo.Name = f.Name()
		itemInfo.IdPath = tempPath
		//itemInfo.Path = path + f.Name()
		itemInfo.Path = tempPath
		itemInfo.ParentId = 0
		itemInfo.ParentPath = parentPath
		itemInfo.IsFolder = f.IsDir()

		for _, val := range settings {
			if val.Name == f.Name() {
				itemInfo.ColumnIndex = val.ColumnIndex
				itemInfo.RowIndex = val.RowIndex
				break
			}
		}

		//log.Println(f.Name())
		itemArray = append(itemArray, itemInfo)
	}

	return itemArray
}

func convertPath(path string) string {
	var tempPath = strings.Replace(path, " ", "3sp3", -1) //replace space to 3sp3
	return strings.Replace(tempPath, "/", "m3g", -1)      //replace slash to m3g
}

func decodePath(path string) string {
	var tempPath = strings.Replace(path, "3sp3", " ", -1) //replace space to 3sp3
	return strings.Replace(tempPath, "m3g", "/", -1)      //replace slash to m3g
}
