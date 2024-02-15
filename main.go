package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var db *sql.DB

type User struct {
    Email string
    Password string
    Session string
}

type AccessSuccess struct {
    Id string `json:"id"` 
    Session string `json:"session"`
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

    fmt.Println("Hello, World!")

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

    // Get a database handle.
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
    routeHandler.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
    })

    routeHandler.HandleFunc("POST /access", func(w http.ResponseWriter, req *http.Request) {
        fmt.Println("/access request")

        if req.URL.Path != "/access" {
			http.NotFound(w, req)
			return
		}

        var user User

        err := json.NewDecoder(req.Body).Decode(&user)
        hadErr := handleErr(w, err, http.StatusBadRequest)
        if hadErr {
            return
        }

        var (
            dbId int
        )
        hasEmailErr := db.QueryRow("SELECT id FROM user WHERE email = ?", user.Email).Scan(&dbId)

        if dbId > 0 {
            login(w, req, user)
            return
        }

        if hasEmailErr == sql.ErrNoRows {
            register(w, req, user)
            return
        }

        log.Println("has email err")
        hadErr = handleErr(w, hasEmailErr, http.StatusBadRequest)
        if hadErr {
            return
        }
        //add error cases for mysql specific errors
        //http://go-database-sql.org/errors.html
        //if driverErr, ok := err.(*mysql.MySQLError); ok { // Now the error number is accessible directly
	    //    if driverErr.Number == 1045 {
	    //    	// Handle the permission-denied error
	    //    }
        //}

    })

    routeHandler.HandleFunc("POST /isAuth", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
    })

    routeHandler.HandleFunc("PUT /", func(w http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
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

func login(w http.ResponseWriter, req *http.Request, user User) {
    var (
        dbPass string
    )

    userRow := db.QueryRow("SELECT password FROM user WHERE email = ?", user.Email)

    err := userRow.Scan(&dbPass)
    hadErr := handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    if dbPass != user.Password {
        log.Fatal("password doesn't match")
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    assignSessionStmt, err := db.Prepare("UPDATE user SET session = ? WHERE email = ?")
    defer assignSessionStmt.Close()
    hadErr = handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    sessionBytes, err := uuid.NewUUID()
    hadErr = handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    sessionId := sessionBytes.String() 
    _, err = assignSessionStmt.Exec(sessionId, user.Email)
    hadErr = handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    addedRow := db.QueryRow("SELECT id, session FROM user WHERE email = ?", user.Email)
    
    var (
        addId int
        addSession string
    )

    err = addedRow.Scan(&addId, &addSession)
    hadErr = handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    retVals := AccessSuccess{
        string(addId),
        sessionId,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(retVals)
}

func register(w http.ResponseWriter, req *http.Request, user User) {
    sessionBytes, err := uuid.NewUUID()
    hadErr := handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    sessionId := sessionBytes.String() 

    // adding session here in case of error so don't have to reverse change
    createUserStmt, err := db.Prepare("INSERT INTO user (email, password, session) VALUES (?, ?, ?)")
    defer createUserStmt.Close()
    hadErr = handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    _, err = createUserStmt.Exec(user.Email, user.Password, sessionId)
    hadErr = handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    addedRow := db.QueryRow("SELECT id, session FROM user WHERE email = ?", user.Email)
    
    var (
        addId int
        addSession string
    )

    err = addedRow.Scan(&addId, &addSession)
    hadErr = handleErr(w, err, http.StatusBadRequest)
    if hadErr {
        return
    }

    retVals := AccessSuccess{
        string(addId),
        sessionId,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(retVals)
}

func handleErr(w http.ResponseWriter, err error, status int) bool {
    if err != nil {
        log.Println(err.Error())
        http.Error(w, err.Error(), status)
        return true
    }

    return false
}
