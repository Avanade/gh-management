package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

type DB struct {
	*sql.DB
}

type ConnectionParam struct {
	ConnectionString string
	Server           string
	Port             int
	User             string
	Password         string
	Database         string
}

func Init(cp ConnectionParam) (*DB, error) {
	connString := cp.ConnectionString
	// Build connection string if property connection string is not set
	if connString == "" {
		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
			cp.Server, cp.User, cp.Password, cp.Port, cp.Database)
	}

	// Create connection pool
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		return nil, err
	}

	// Check connection
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) ExecuteStoredProcedure(procedure string, params map[string]interface{}) (sql.Result, error) {
	var args []interface{}

	for i, v := range params {
		args = append(args, sql.Named(i, v))
	}

	ctx := context.Background()
	result, err := db.ExecContext(ctx, procedure, args...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DB) ExecuteStoredProcedureWithResult(procedure string, params map[string]interface{}) ([]map[string]interface{}, error) {
	var args []interface{}

	ctx := context.Background()

	for i, v := range params {
		args = append(args, sql.Named(i, v))
	}

	rows, err := db.QueryContext(ctx, procedure, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i, _ := range values {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		result := make(map[string]interface{})
		for i, val := range values {
			result[columns[i]] = val
		}
		results = append(results, result)
	}

	return results, nil
}
