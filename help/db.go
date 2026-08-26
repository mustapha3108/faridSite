package help

import (
	"log"
	"os"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
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
		//fcmd
		`CREATE TABLE IF NOT EXISTS users (
			userId 		INTEGER PRIMARY KEY AUTOINCREMENT,
			userName 	TEXT NOT NULL UNIQUE,
			password 	TEXT NOT NULL,
			access 		TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS categories(
			categoryId   INTEGER PRIMARY KEY AUTOINCREMENT,
			categoryName TEXT NOT NULL unique,
			imagePath    TEXT NOT NULL
		)
		`,
		`CREATE TABLE IF NOT EXISTS projects (
			projectId 	INTEGER PRIMARY KEY AUTOINCREMENT,
			userId 		INTEGER NOT NULL REFERENCES users(userid) ON DELETE  SET NULL,
			categoryId  INTEGER NOT NULL REFERENCES categories(categoryId) ON DELETE SET NULL,
			projectName TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL,
			imagePaths 	TEXT NOT NULL,
			mImagePath  TEXT NOT NULL
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
			memberTitle 		TEXT NOT NULL,
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
			message 	TEXT NOT NULL,
		)`,//gotta add some columns here
		`CREATE TABLE IF NOT EXISTS jobApplications(
			apId		INTEGER PRIMARY KEY AUTOINCREMENT,
			firstName 	TEXT NOT NULL,
			lastName    TEXT NOT NULL,
			email 		TEXT NOT NULL,
			object 		TEXT NOT NULL,
			message 	TEXT NOT NULL,
			software    TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS partenairs(
			partenairId   INTEGER PRIMARY KEY AUTOINCREMENT,
			partenairName TEXT NOT NULL,
			imagePath     TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS messNot(
			messNotId     INTEGER PRIMARY KEY AUTOINCREMENT,
			messageId     INTEGER NOT NULL REFERENCES messages(messageId) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS canNot(
			canNotId     INTEGER PRIMARY KEY AUTOINCREMENT,
			candidatId   INTEGER NOT NULL REFERENCES jobApplication(apId) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS software(
			softwareId   INTEGER PRIMARY KEY AUTOINCREMENT,
			softwareName TEXT NOT NULL UNIQUE,
			required     INTEGER NOT NULL DEFAULT 0
		)`,
	}

	for _, v := range tables {
		if _, err := db.Exec(v); err != nil {
			log.Fatal("error creating:", v, "	", err)
		}
	}

	fPassword, err := Encode("farid5169")
	if err!=nil{
		log.Fatal("error creating: 	", err)
	}

	_,err = db.Exec(`INSERT OR IGNORE INTO users VALUES(?,?,?,?)`, 1, "farid", fPassword, "fcmd")
	if err!=nil {
		log.Fatal("error creating: 	", err)
	}

	devPassword, err := Encode("crow3108")
	if err!=nil{
		log.Fatal("error creating: 	", err)
	}

	_,err = db.Exec(`INSERT OR IGNORE INTO users VALUES(?,?,?,?)`, 1, "dev", devPassword, "sss")
	if err!=nil {
		log.Fatal("error creating: 	", err)
	}
	

	return db
}


//delete database
func DeleteTable(db *sqlx.DB) {
	tables := []string{"users", "categories", "projects", "ratings", "members", "contact", "messages", "jobApplications", "partenairs", "messNot", "canNot"}
	for _,v := range(tables) {
		if _, err:= db.Exec(`DROP TABLE ?`, v); err!=nil {
			log.Fatal("error creating: 	", err)
		}
	}
}