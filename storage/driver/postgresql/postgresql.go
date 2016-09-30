// Package postgresql provides a wrapper around the pq package.
package postgresql

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	info      Info
	infoMutex sync.RWMutex
)

// Info holds the details for the connection.
type Info struct {
	Username        string
	Password        string
	Database        string
	Hostname        string
	Port            int
	Parameter       string
	MigrationFolder string
	Extension       string
}

// SetConfig stores the config.
func SetConfig(i Info) {
	infoMutex.Lock()
	info = i
	infoMutex.Unlock()
}

// Config returns the config.
func Config() Info {
	infoMutex.RLock()
	defer infoMutex.RUnlock()
	return info

}

// ResetConfig removes the config.
func ResetConfig() {
	infoMutex.Lock()
	info = Info{}
	infoMutex.Unlock()
}

// *****************************************************************************
// Database Handling
// *****************************************************************************

// Connect to the database.
func Connect(specificDatabase bool) error {
	var err error

	// Connect to database and ping
	if SQL, err = sqlx.Connect("postgres", dsn(specificDatabase)); err != nil {
		return err

	}

	return err
}

// Disconnect the database connection.
func Disconnect() error {
	return SQL.Close()
}

// Create a new database.
func Create() error {
	// Create the database
	_, err := SQL.Exec(fmt.Sprintf(`CREATE DATABASE %v;`, Config().Database))

	return err
}

// Drop a database.
func Drop() error {
	// Drop the database
	_, err := SQL.Exec(fmt.Sprintf(`DROP DATABASE %v;`, Config().Database))

	return err
}

// *****************************************************************************
// Database Specific
// *****************************************************************************

var (
	// SQL wrapper
	SQL *sqlx.DB
)

// DSN returns the Data Source Name.
func dsn(includeDatabase bool) string {
	// Set defaults
	ci := setDefaults()

	// Build parameters
	param := ci.Parameter

	// If parameter is specified, add a question mark
	// Don't add one if a question mark is already there
	if len(ci.Parameter) > 0 && !strings.HasPrefix(ci.Parameter, "?") {
		param = "?" + ci.Parameter

	}

	// Example: postgres://pqgotest:password@localhost/pqgotest
	s := fmt.Sprintf("postgres://%v:%v@%v:%d/%v", ci.Username, ci.Password, ci.Hostname, ci.Port, param)

	if includeDatabase {
		s = fmt.Sprintf("postgres://%v:%v@%v:%d/%v%v", ci.Username, ci.Password, ci.Hostname, ci.Port, ci.Database, param)
	}

	return s
}

// setDefaults gets the default connection information.
func setDefaults() Info {
	ci := Config()

	return ci
}
