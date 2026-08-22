package help

import (
	"log"
	"os"
	_ "modernc.org/sqlite"
	"github.com/jmoiron/sqlx"
)



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