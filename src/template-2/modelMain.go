package main

import "time"

type DirectionType int

const (
	LeftToRigth DirectionType = iota
	UpToBottom
)

type AppInfo struct {
	Desktop DesktopSettings
	User    UserDetails
	Items   []ItemInfo
}

type UserDetails struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	BirthDate   time.Time
	CreatedDate time.Time
}

type DesktopSettings struct {
	//Direction  DirectionType
	FolderIndexes []DesktopFolderIndex
	ColumnSize    int16
	RowSize       int16
}

type DesktopFolderIndex struct {
	Name        string
	ColumnIndex int
	RowIndex    int
}

type ItemInfo struct {
	//Direction      DirectionType
	Id             int
	ParentId       int
	Public         bool
	IdPath         string
	Path           string
	Name           string
	Descritpion    string
	ColumnIndex    int
	RowIndex       int
	MaxUploadCount int
	MaxUploadSize  int
	IsFolder       bool
}

type ParentChild struct {
	Shadow   bool
	ParendId int32
	ChildId  int32
}
