package database

import (
	"context"
	"database/sql"
	"fmt"
	"main/config"

	_ "github.com/microsoft/go-mssqldb"
)

type Database struct {
	connString string
	db         *sql.DB
}

func NewDatabase(config config.ConfigManager) Database {
	fmt.Println("ConnectDb New")
	connectionString := config.GetDatabaseConnectionString()
	return Database{
		connString: connectionString,
	}
}

func (d *Database) Connect() error {
	conn, err := sql.Open("sqlserver", d.connString)
	if err != nil {
		return err
	}

	d.db = conn
	return nil
}

func (d *Database) Disconnect() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *Database) Query(query string, args ...any) (*sql.Rows, error) {
	err := d.Connect()
	if err != nil {
		return nil, err
	}
	defer d.Disconnect()

	ctx := context.Background()
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *Database) QueryRow(query string, args ...any) (*sql.Row, error) {
	err := d.Connect()
	if err != nil {
		return nil, err
	}
	defer d.Disconnect()

	ctx := context.Background()
	row := d.db.QueryRowContext(ctx, query, args...)
	if row == nil {
		err = fmt.Errorf("QueryRowContext return nil")
	}
	return row, err
}

func (d *Database) Execute(query string, args ...any) error {
	err := d.Connect()
	if err != nil {
		return err
	}
	defer d.Disconnect()

	ctx := context.Background()
	_, err = d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// helper function to convert rows.Scan() to map[string]interface{}
func (d *Database) RowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
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

// OLD
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// OLD
func (d *Database) ExecuteStoredProcedure(procedure string, params map[string]interface{}) (sql.Result, error) {
	var args []interface{}

	for i, v := range params {
		args = append(args, sql.Named(i, v))
	}

	ctx := context.Background()
	result, err := d.db.ExecContext(ctx, procedure, args...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// OLD
func (d *Database) ExecuteStoredProcedureWithResult(procedure string, params map[string]interface{}) ([]map[string]interface{}, error) {
	var args []interface{}

	ctx := context.Background()

	for i, v := range params {
		args = append(args, sql.Named(i, v))
	}

	rows, err := d.db.QueryContext(ctx, procedure, args...)
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
		for i := range values {
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
