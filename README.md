# TODOS

add error cases for mysql specific errors
http://go-database-sql.org/errors.html
if driverErr, ok := err.(*mysql.MySQLError); ok { // Now the error number is accessible directly
    if driverErr.Number == 1045 {
    	// Handle the permission-denied error
    }
}

