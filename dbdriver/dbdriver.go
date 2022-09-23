package dbdriver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	once sync.Once
	db   *sql.DB
)

func init() {
	once.Do(initialiseDBconn)
}

func initialiseDBconn() {

}

// Todo data structure for a task with a description of what to do
type Todo struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Complete bool   `json:"complete"`
}

func dbConn() (db *sql.DB) {
	dbIpaddr := os.Getenv("DBIPADDRESS")
	dbPort := os.Getenv("DBPORT")
	dbDriver := os.Getenv("DBDRIVER")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASSWORD")
	dbName := "todo"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbIpaddr+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("SQL connection opened successfully")
	return db

}

func DatabaseGet(userId string) []Todo {

	var todos []Todo
	db := dbConn()

	addTable, err := db.Prepare("CREATE TABLE IF NOT EXISTS `" + userId + "` (ID varchar(255), Message varchar(255), Complete boolean)")
	if err != nil {
		panic(err.Error())
	}
	addTable.Exec()

	log.Println("Table " + userId + " created successfully")

	rows, err := db.Query("SELECT * FROM `" + userId + "`")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	db.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Message, &todo.Complete)
		todos = append(todos, todo)
	}
	return todos

}

func DatabaseAdd(userId string, ID string, Message string, Complete bool) {

	db := dbConn()

	r, err := db.Prepare("INSERT INTO `" + userId + "` (ID, Message, Complete) VALUES(?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	r.Exec(ID, Message, Complete)
	db.Close()
	log.Println("New item created successfully")
}

func DatabaseComplete(userId string, ID string) {

	db := dbConn()

	r, err := db.Prepare("UPDATE `" + userId + "` SET Complete=true WHERE ID=?")
	if err != nil {
		panic(err.Error())
	}
	r.Exec(ID)
	db.Close()
	log.Println("Item" + ID + "from " + userId + " set to completed status")
}

func DatabaseDelete(userId string, ID string) {

	db := dbConn()

	r, err := db.Prepare("DELETE FROM `" + userId + "` WHERE ID=?")
	if err != nil {
		panic(err.Error())
	}
	r.Exec(ID)
	db.Close()
	log.Println("Item " + ID + " from " + userId + " removed successfully")
}
