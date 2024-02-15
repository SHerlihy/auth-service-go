package main

import (
	"os"
)

// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
// username:password@protocol(address)/dbname?param=value

//func ProdEnv(){
//    dbAddress := "terraform-20231115140424354000000001.chvww7wqh970.eu-west-2.rds.amazonaws.com"
//    os.Setenv("DSN", fmt.Sprintf("user:userpass@(%s)/users", dbAddress))
//    os.Setenv("SECRET_KEY", "secret")
//    os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
//}

func DevEnv(){
    os.Setenv("DBUSER", "authServiceGo")
    os.Setenv("DBPASS", "password")
    os.Setenv("DBADDR", "127.0.0.1")
    os.Setenv("DBPORT", "3306")
    os.Setenv("DBDATABASE", "users")
}
