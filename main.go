package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
    Email string
    Password string
    Session string
}

type AuthUser struct {
    Id string `json:"id"` 
    Session string `json:"session"`
}

type ChangeEmail struct {
    AuthUser
    Email string `json:"email"`
}

type ChangePass struct {
    AuthUser
    Password string `json:"password"`
}

func main() {
//    env := flag.String("e", "dev", "environment: dev|prod")
//    flag.Parse()
//
//    if *env == "dev" {
//        env_vars.DevEnv()
//    }
//
//    if *env == "prod" {
//        env_vars.ProdEnv()
//    }
        
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
    db, err = sql.Open("mysql", accessStr)
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

    routeHandler := http.NewServeMux()

    routeHandler.HandleFunc("POST /access", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("/access request")

        if req.URL.Path != "/access" {
			http.NotFound(w, req)
			return
		}

        var user User

        err := json.NewDecoder(req.Body).Decode(&user)
        hadErr := HandleErr(w, err, http.StatusBadRequest)
        if hadErr {
            return
        }

        var (
            dbId int
        )
        hasEmailErr := db.QueryRow("SELECT id FROM user WHERE email = ?", user.Email).Scan(&dbId)

        if dbId > 0 {
            Login(w, req, user)
            return
        }

        if hasEmailErr == sql.ErrNoRows {
            Register(w, req, user)
            return
        }

        //Scan always has err != nil
        HandleErr(w, hasEmailErr, http.StatusBadRequest)
    })

    routeHandler.HandleFunc("PUT /email", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("/change email")

        if req.URL.Path != "/email" {
			http.NotFound(w, req)
			return
		}

        var user ChangeEmail

        err := json.NewDecoder(req.Body).Decode(&user)
        hadErr := HandleErr(w, err, http.StatusBadRequest)
        if hadErr {
            return
        }
    
        assignEmailStmt, err := db.Prepare(fmt.Sprintf("UPDATE user SET email = '%s' WHERE id = ?", user.Email))
        defer assignEmailStmt.Close()
        hadErr = HandleErr(w, err, http.StatusBadRequest)
        if hadErr {
            return
        }
    
        ChangeUser(w, req, user.Id, user.Session, assignEmailStmt)
    })

    routeHandler.HandleFunc("PUT /password", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("/change pass")

        if req.URL.Path != "/password" {
			http.NotFound(w, req)
			return
		}

        var user ChangePass

        err := json.NewDecoder(req.Body).Decode(&user)
        hadErr := HandleErr(w, err, http.StatusBadRequest)
        if hadErr {
            return
        }
    
        assignPassStmt, err := db.Prepare(fmt.Sprintf("UPDATE user SET password = '%s' WHERE id = ?", user.Password))
        defer assignPassStmt.Close()
        hadErr = HandleErr(w, err, http.StatusBadRequest)
        if hadErr {
            return
        }
    
        ChangeUser(w, req, user.Id, user.Session, assignPassStmt)
    })

    routeHandler.HandleFunc("DELETE /", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
    })

    server := &http.Server{
        Addr: ":8080",
        Handler: routeHandler,
    }

    log.Fatal(server.ListenAndServe())
}

func HandleErr(w http.ResponseWriter, err error, status int) bool {
    if err != nil {
        log.Println(err.Error())
        http.Error(w, err.Error(), status)
        return true
    }

    return false
}
