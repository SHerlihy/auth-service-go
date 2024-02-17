package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
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
    t.Log("test begin")
    resetDB(t)
    accessStr := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s",
        os.Getenv("DBUSER"),
        os.Getenv("DBPASS"),
        os.Getenv("DBADDR"),
        os.Getenv("DBPORT"),
        os.Getenv("DBDATABASE"),
        )

    var err error
    testDB, err = sql.Open("mysql", accessStr)
    if err != nil {
        t.Fatal(err)
    }
    defer testDB.Close()

    addTest1 := RequestAccessUser{
        Email: "test1@mail.com",
        Password: "test1",
    }
    jsonTest1, _ := json.Marshal(addTest1)
    w := httptest.NewRecorder()
    bodyTest1 := bytes.NewReader(jsonTest1)
    req := httptest.NewRequest("POST", "http://localhost/access", bodyTest1)
    AccessUser(testDB, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

    t.Log("test end")
    t.Log(resp.Status)
	t.Log(resp.Header.Get("Content-Type"))
	t.Log(string(body))
}
