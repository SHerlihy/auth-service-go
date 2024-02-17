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

	"github.com/google/uuid"
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

func testAdd(t *testing.T, userDets RequestAccessUser) ResponseUserAuth {
    jsonTest1, _ := json.Marshal(userDets)
    w := httptest.NewRecorder()
    bodyTest1 := bytes.NewReader(jsonTest1)
    req := httptest.NewRequest("POST", "http://localhost/access", bodyTest1)
    AccessUser(testDB, w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
    bodyStruct := ResponseUserAuth{}
    json.Unmarshal(body, &bodyStruct)

    return bodyStruct


}

func testRegSessions(t *testing.T, userDets RequestAccessUser) ResponseUserAuth {
    resUserAuth := testAdd(t, userDets)

    err := uuid.Validate(resUserAuth.CurrentSession)
    if err != nil {
        t.Fatal(err)
    }

    if resUserAuth.PreviousSession != "" {
        t.Fatal("Prev session present")
    }

    return resUserAuth
}

func testRegister(t *testing.T, userDets RequestAccessUser, expId int) ResponseUserAuth {
    resUserAuth := testRegSessions(t, userDets)

    expectedId := fmt.Sprint(expId)
    if resUserAuth.Id != expectedId {
        t.Fatal(fmt.Sprintf("Id expected: %s\nId recieved: %s", expectedId, resUserAuth.Id))
    }

    return resUserAuth
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

    for i := range 5 {
        addTest := RequestAccessUser{
            Email: fmt.Sprintf("test%x@mail.com", i),
            Password: fmt.Sprintf("test%x", i),
        }

        testRegister(t, addTest, i+1)
    }

//    t.Log("test end")
//    t.Log(resp.Status)
//	t.Log(resp.Header.Get("Content-Type"))
//	t.Log(string(body))
}
