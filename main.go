package main

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"time"
)

func connectToOracle() *sql.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	server := viper.GetString("database.server")
	port := viper.GetString("database.port")
	service := viper.GetString("database.service")
	walletLocation := viper.GetString("database.walletLocation")
	connectionString := "oracle://" + username + ":" + password + "@" + server + ":" + port + "/" + service
	if walletLocation != "" {
		connectionString += "?TRACE FILE=trace.log&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(walletLocation)
	}
	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	return db
}

func main() {

	viper.SetConfigName("oraconfig")
	viper.SetConfigType("ini")
	viper.AddConfigPath("/etc/oratest")
	viper.AddConfigPath("$HOME/.oratest")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	var (
		d time.Time
	)
	db := connectToOracle()
	rows, err := db.Query("SELECT SYSTIMESTAMP AS D FROM DUAL")
	if err != nil {
		log.Panicln("error in query: ", err)
		return
	}
	defer rows.Close()
	for rows.Next() {

		rows.Scan(&d)

		fmt.Printf("d=%v\n", d.Month())
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()

}
