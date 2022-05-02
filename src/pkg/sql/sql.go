// package sql

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/denisenkom/go-mssqldb"
// )

// // Connection Parameters
// type ConnectionParam struct {
// 	connectionString string
// 	server           string
// 	port             int
// 	user             string
// 	password         string
// 	database         string
// }

// func Init(cp ConnectionParam) (*sql.DB, error) {
// 	connString := cp.connectionString
// 	// Build connection string if property connection string is not set
// 	if connString == "" {
// 		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
// 			cp.server, cp.user, cp.password, cp.port, cp.database)
// 	}

// 	// Create connection pool
// 	db, err := sql.Open("sqlserver", connString)
// 	if err != nil {
// 		log.Fatal("Error creating connection pool: ", err.Error())
// 		return nil, err
// 	}
// 	ctx := context.Background()
// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 		return nil, err
// 	}
// 	fmt.Printf("Connected!")
// 	return db, nil
// }
