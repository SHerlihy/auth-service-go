package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

var testDB *sql.DB

func resetDB(t *testing.T){
    DevEnv()

    accessStr := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s",
        os.Getenv("DBUSER"),
        os.Getenv("DBPASS"),
        os.Getenv("DBADDR"),
        os.Getenv("DBPORT"),
        os.Getenv("DBDATABASE"),
        )

    fmt.Println("accessStr")
    fmt.Println(accessStr)

    var err error
    testDB, err = sql.Open("mysql", accessStr)
    if err != nil {
        t.Fatal(err)
    }

    defer testDB.Close()

    pingErr := testDB.Ping()
    if pingErr != nil {
        t.Fatal(pingErr)
    }
    fmt.Println("Connected!")

    dropStmt, _ := testDB.Prepare("DROP TABLE IF EXISTS user")
    createStmt, _ := testDB.Prepare("CREATE TABLE user (email VARCHAR(255) NOT NULL UNIQUE, password VARCHAR(255) NOT NULL, id int NOT NULL AUTO_INCREMENT UNIQUE, session VARCHAR(255) NOT NULL, PRIMARY KEY (`id`));")

    dropStmt.Exec()
    createStmt.Exec()

    dropStmt.Close()
    createStmt.Close()
    testDB.Close()

    if err != nil {
        t.Log(err)
    }
}

func TestAddUsers(t *testing.T){
    resetDB(t)
}
