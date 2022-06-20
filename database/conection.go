package database

import (
	"database/sql"
	"fmt"
	"log"

	//Postgres Driver imported
	_ "github.com/lib/pq"
	config "github.com/spf13/viper"
)

func ConnectDB() *sql.DB {
	config.SetConfigFile(".env")
	config.ReadInConfig()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Get("host"), config.GetInt("port"), config.Get("user"), config.Get("password"), config.Get("database"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	
	log.Println("DB connected")

	autoMigrations(db)

	return db
}

func autoMigrations(db *sql.DB) {

	if _, err := db.Exec(Usuarios); err != nil {
		log.Fatal(err)
	}

	db.Exec(index_users)

	if _, err := db.Exec(Friends); err != nil {
		log.Fatal(err)
	}
	
}

