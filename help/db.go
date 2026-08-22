package help

import (
	"log"
	"mime/multipart"
	"os"
	_ "modernc.org/sqlite"
	"github.com/jmoiron/sqlx"
)

// structs
type User struct {
	UserId   int    `db:"userId"`
	UserName string `db:"userName" form:"UserName" validate:"required"`
	Password string `db:"password" form:"Password" validate:"required"`
	Access   int    `db:"access"   form:"Access"   validate:"number"`
}
 
type Project struct {
	ProjectId   int                     `db:"projectId"`
	UserId      int                     `db:"userId"`
	ProjectName string                  `db:"projectName" form:"ProjectName" validate:"required"`
	Description string                  `db:"description" form:"Description" validate:"required"`
	ImagePaths  string                  `db:"imagePaths"`
	Images      []*multipart.FileHeader `form:"Images" db:"-"`
}
 
type Rating struct {
	RatingId int    `db:"ratingId"`
	Name     string `db:"name" form:"Name" validate:"required"`
	Comment  string `db:"comment" form:"Comment"`
	Rating int `db:"rating" form:"Rating" validate:"required,gte=0,lte=50"`
}
 
type Member struct {
	MemberId          int                   `db:"memberId"`
	MemberName        string                `db:"memberName"        form:"MemberName" validate:"required"`
	MemberTitle       string                `db:"memberTitle"       form:"MemberTitle"`
	MemberDescription string                `db:"memberDescription" form:"MemberDescription"`
	MemberImagePath   string                `db:"memberImagePath"`
	MemberImage       *multipart.FileHeader `form:"MemberImage" db:"-" validate:"required"`
}
 
type Contact struct {
	Address  string `db:"address"  form:"Address"`
	Baladya  string `db:"baladya"  form:"Baladya"`
	Wilaya   string `db:"wilaya"   form:"Wilaya"`
	Email    string `db:"email"    form:"Email" validate:"omitempty,email"`
	Number   string `db:"number"   form:"Number"`
	Location string `db:"location" form:"Location"`
}
 
type Message struct {
	MessageId int    `db:"messageId"`
	FirstName string `db:"firstName" form:"FirstName" validate:"required"`
	LastName  string `db:"lastName"  form:"LastName" validate:"required"`
	Email     string `db:"email"     form:"Email" validate:"required,email"`
	Object    string `db:"object"    form:"Object" validate:"required"`
	Message   string `db:"message"   form:"Message" validate:"required"`
}
 
type JobApplication struct {
	ApId int    `db:"apId"`
	FirstName     string `db:"firstName" form:"FirstName" validate:"required"`
	LastName      string `db:"lastName"  form:"LastName" validate:"required"`
	Email         string `db:"email"     form:"Email" validate:"required,email"`
	Object        string `db:"object"    form:"Object" validate:"required"`
	Message       string `db:"message"   form:"Message" validate:"required"`
}

// database
func BasicDb() *sqlx.DB {
	dsn := os.Getenv("db_url")

	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		log.Fatal("database connection error:", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatal("error enabling foreign keys:", err)
	}

	tables := []string{
		//1=all, 2=create, 3=delete, 4=modify, 5=create+delete, 6=create+modify, 7=delete+modify
		`CREATE TABLE IF NOT EXISTS users (
			userId 		INTEGER PRIMARY KEY AUTOINCREMENT,
			userName 	TEXT NOT NULL UNIQUE,
			password 	TEXT NOT NULL,
			access 		INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			projectId 	INTEGER PRIMARY KEY AUTOINCREMENT,
			userId 		INTEGER NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
			projectName TEXT NOT NULL,
			description TEXT NOT NULL,
			imagePaths 	TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS ratings(
			ratingId 	INTEGER PRIMARY KEY AUTOINCREMENT,
			name 		TEXT NOT NULL,
			comment 	TEXT,
			rating 		INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 50)
		);`,
		`CREATE TABLE IF NOT EXISTS members(
			memberId 			INTEGER PRIMARY KEY AUTOINCREMENT,
			memberName 			TEXT NOT NULL,
			memberTitle 		TEXT,
			memberDescription   TEXT,
			memberImagePath 	TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS contact(
			address  TEXT,
			baladya  TEXT,
			wilaya   TEXT,
			email    TEXT,
			number   TEXT,
			location TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS messages(
			messageId 	INTEGER PRIMARY KEY AUTOINCREMENT,
			firstName 	TEXT NOT NULL,
			lastName    TEXT NOT NULL,
			email 		TEXT NOT NULL,
			object 		TEXT NOT NULL,
			message 	TEXT NOT NULL
		)`,//gotta add some columns here
		`CREATE TABLE IF NOT EXISTS jobApplication(
			apId		INTEGER PRIMARY KEY AUTOINCREMENT,
			firstName 	TEXT NOT NULL,
			lastName    TEXT NOT NULL,
			email 		TEXT NOT NULL,
			object 		TEXT NOT NULL,
			message 	TEXT NOT NULL
		)`,
	}

	for _, v := range tables {
		if _, err := db.Exec(v); err != nil {
			log.Fatal("error creating:", v, "	", err)
		}
	}

	return db
}