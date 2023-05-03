package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func addToDb(db *sql.DB, tableName string, fw ApiRequestClass) {
	fmt.Println("addToDb called")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Could not open the database: ", err)
		panic(err)
	} else {
		fmt.Println("Opened: ", os.Getenv("DATABASE_URL"))
	}
	defer db.Close()

	_, err = db.Query(`INSERT INTO public.weathertable (date, temp, feelslike, tempmin, tempmax, pressure, humidity, clouds)  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		time.Now(),
		fw.Main.Temp,
		fw.Main.FeelsLike,
		fw.Main.TempMin,
		fw.Main.TempMax,
		fw.Main.Pressure,
		fw.Main.Humidity,
		fw.Clouds.All)

	if err != nil {
		fmt.Println("Could not add row dating " + time.Now().Format("2006-01-02 15:04:05"))
		panic(err)
	}
}

func createDb(dbName, tableName string) *sql.DB {
	fmt.Println("createDb called")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Could not open the database: ", err)
		panic(err)
	} else {
		fmt.Println("Opened: ", os.Getenv("DATABASE_URL"))
	}
	defer db.Close()

	fmt.Println("Database creation: ", dbName)
	_, err = db.Exec(`CREATE DATABASE ` + dbName)
	if err != nil {
		fmt.Println("Database ", dbName, " already exists")
	} else {
		fmt.Println("Database created: ", dbName)
	}

	fmt.Println("Table creation")
	_, err = db.Exec("CREATE TABLE " + tableName +
		"( date timestamp, temp real, feelsLike real, tempMin real, tempMax real," +
		"pressure int, humidity int, clouds int )")
	if err != nil {
		fmt.Println("Table " + tableName + " already exists")
	} else {
		fmt.Println("Table" + tableName + "created")
	}
	fmt.Println("createDb done;")
	return db
}
