package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func accessUser(reqPayload RequestAccessUser)*http.Response{
    reqJSON, _ := json.Marshal(reqPayload)
    w := httptest.NewRecorder()
    reqBody := bytes.NewReader(reqJSON)
    req := httptest.NewRequest("POST", "http://localhost/access", reqBody)
    AccessUser(testDB, w, req)

	resp := w.Result()
    return resp
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

func successfulregisterSessions(t *testing.T, resUserAuth ResponseUserAuth) ResponseUserAuth {
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
    resUserAuth := testAdd(t, userDets)
    resUserAuth = successfulregisterSessions(t, resUserAuth)

    expectedId := fmt.Sprint(expId)
    if resUserAuth.Id != expectedId {
        t.Fatal(fmt.Sprintf("Id expected: %s\nId recieved: %s", expectedId, resUserAuth.Id))
    }

    return resUserAuth
}

func successfulLoginSessions(t *testing.T, resUserAuth ResponseUserAuth) ResponseUserAuth {
    err := uuid.Validate(resUserAuth.CurrentSession)
    if err != nil {
        t.Fatal(err)
    }

    err = uuid.Validate(resUserAuth.CurrentSession)
    if err != nil {
        t.Fatal(err)
    }

    return resUserAuth
}

func testSuccessfulLogin(t *testing.T, userDets RequestAccessUser, expId int) ResponseUserAuth {
    resUserAuth := testAdd(t, userDets)
    resUserAuth = successfulLoginSessions(t, resUserAuth)

    expectedId := fmt.Sprint(expId)
    if resUserAuth.Id != expectedId {
        t.Fatal(fmt.Sprintf("Id expected: %s\nId recieved: %s", expectedId, resUserAuth.Id))
    }

    return resUserAuth
}

func testUnsuccessfulLogin(t *testing.T, userDets RequestAccessUser, expId int) {
    accessResp := accessUser(userDets)

    if accessResp.StatusCode != 400 {
        t.Fatal(fmt.Sprintf("Expected status: 400\nRecieved status: %x", accessResp.StatusCode))
    }
}

func changeEmail(reqPayload RequestChangeEmail) *http.Response {
    reqJSON, _ := json.Marshal(reqPayload)
    w := httptest.NewRecorder()
    reqBody := bytes.NewReader(reqJSON)
    req := httptest.NewRequest("PUT", "http://localhost/email", reqBody)
    ChangeEmail(testDB, w, req)

	resp := w.Result()
    return resp
}

type ChangeEmailTestPayload struct {
    RequestChangeEmail
    Password string
}

func successfulChangeEmail(t *testing.T, testPayload ChangeEmailTestPayload) {
    changeResp := changeEmail(testPayload.RequestChangeEmail)

    if changeResp.StatusCode != 202 {
        t.Fatal(fmt.Sprintf("Expected status: 202\nRecieved status: %x", changeResp.StatusCode))
    }

    loginPayload := RequestAccessUser{
        Email: testPayload.RequestChangeEmail.Email,
        Password: testPayload.Password,
    }

    loginResp := testAdd(t, loginPayload)
    successfulLoginSessions(t, loginResp)
}

func changePassword(reqPayload RequestChangePassword) *http.Response {
    reqJSON, _ := json.Marshal(reqPayload)
    w := httptest.NewRecorder()
    reqBody := bytes.NewReader(reqJSON)
    req := httptest.NewRequest("PUT", "http://localhost/password", reqBody)
    ChangePassword(testDB, w, req)

	resp := w.Result()
    return resp
}

type ChangePasswordTestPayload struct {
    RequestChangePassword
    Email string
}

func successfulChangePassword(t *testing.T, testPayload ChangePasswordTestPayload) {
    changeResp := changePassword(testPayload.RequestChangePassword)

    if changeResp.StatusCode != 202 {
        t.Fatal(fmt.Sprintf("Expected status: 202\nRecieved status: %x", changeResp.StatusCode))
    }

    loginPayload := RequestAccessUser{
        Email: testPayload.Email,
        Password: testPayload.RequestChangePassword.Password,
    }

    loginResp := testAdd(t, loginPayload)
    successfulLoginSessions(t, loginResp)
}

func TestAddUsers(t *testing.T){
    t.Log("TestAddUsers")
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
}

func TestAccessUsers(t *testing.T){
    t.Log("TestAccessUsers")
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

    for i := range 5 {
        addTest := RequestAccessUser{
            Email: fmt.Sprintf("test%x@mail.com", i),
            Password: fmt.Sprintf("test%x", i),
        }

        testSuccessfulLogin(t, addTest, i+1)
    }
}

func TestUnsuccessfulAccess(t *testing.T){
    t.Log("TestUnsuccessfulAccess")
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

    for i := range 5 {
        addTest := RequestAccessUser{
            Email: fmt.Sprintf("test%x@mail.com", i),
            Password: fmt.Sprintf("test%x", i+1),
        }

        testUnsuccessfulLogin(t, addTest, i+1)
    }
}

func TestEmailChange(t *testing.T){
    t.Log("TestEmailChange")
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

    changeEmailTestPayloads := []ChangeEmailTestPayload{}
    for i := range 5 {
        addTest := RequestAccessUser{
            Email: fmt.Sprintf("test%x@mail.com", i),
            Password: fmt.Sprintf("test%x", i),
        }

        resUserAuth := testRegister(t, addTest, i+1)
        changeEmailPaylod := RequestChangeEmail{
            RequestUserAuth: resUserAuth.RequestUserAuth,
            Email: fmt.Sprintf("testA%x@mail.com", i),
        }

        changeEmailTestPayload := ChangeEmailTestPayload{
            RequestChangeEmail: changeEmailPaylod,
            Password: addTest.Password,
        }

        changeEmailTestPayloads = append(changeEmailTestPayloads, changeEmailTestPayload)
    }

    for _, changeEmailTestPayload := range changeEmailTestPayloads {
        successfulChangeEmail(t, changeEmailTestPayload)
    }
}

func TestPasswordChange(t *testing.T){
    t.Log("TestPasswordChange")
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

    changePasswordTestPayloads := []ChangePasswordTestPayload{}
    for i := range 5 {
        addTest := RequestAccessUser{
            Email: fmt.Sprintf("test%x@mail.com", i),
            Password: fmt.Sprintf("test%x", i),
        }

        resUserAuth := testRegister(t, addTest, i+1)
        changePasswordPaylod := RequestChangePassword{
            RequestUserAuth: resUserAuth.RequestUserAuth,
            Password: fmt.Sprintf("testA%x", i),
        }

        changePasswordTestPayload := ChangePasswordTestPayload{
            RequestChangePassword: changePasswordPaylod,
            Email: addTest.Email,
        }

        changePasswordTestPayloads = append(changePasswordTestPayloads, changePasswordTestPayload)
    }

    for _, changePasswordTestPayload := range changePasswordTestPayloads {
        successfulChangePassword(t, changePasswordTestPayload)
    }
}

