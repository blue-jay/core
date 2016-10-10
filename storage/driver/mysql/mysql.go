// Package mysql provides a wrapper around the sqlx package.
package mysql

import (
	"fmt"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

// *****************************************************************************
// Thread-Safe Configuration
// *****************************************************************************

var (
	info      Info
	infoMutex sync.RWMutex
)

// Info holds the details for the MySQL connection.
type Info struct {
	Username        string
	Password        string
	Database        string
	Charset         string
	Collation       string
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
func (c Info) Connect(specificDatabase bool) (*sqlx.DB, error) {
	var err error

	// Connect to database and ping
	if SQL, err = sqlx.Connect("mysql", c.dsn(specificDatabase)); err != nil {
		return SQL, err
	}

	return SQL, err
}

// Disconnect the database connection.
func Disconnect() error {
	return SQL.Close()
}

// Create a new database.
func (c Info) Create() error {
	// Set defaults
	ci := c.setDefaults()

	// Create the database
	_, err := SQL.Exec(fmt.Sprintf(`CREATE DATABASE %v
				DEFAULT CHARSET = %v
				COLLATE = %v
				;`, ci.Database,
		ci.Charset,
		ci.Collation))
	return err
}

// Drop a database.
func (c Info) Drop() error {
	// Drop the database
	_, err := SQL.Exec(fmt.Sprintf(`DROP DATABASE %v;`, c.Database))
	return err
}

// *****************************************************************************
// MySQL Specific
// *****************************************************************************

var (
	// SQL wrapper
	SQL *sqlx.DB
)

// DSN returns the Data Source Name.
func (c Info) dsn(includeDatabase bool) string {
	// Set defaults
	ci := c.setDefaults()

	// Build parameters
	param := ci.Parameter

	// If parameter is specified, add a question mark
	// Don't add one if a question mark is already there
	if len(ci.Parameter) > 0 && !strings.HasPrefix(ci.Parameter, "?") {
		param = "?" + ci.Parameter
	}

	// Add collation
	if !strings.Contains(param, "collation") {
		if len(param) > 0 {
			param += "&collation=" + ci.Collation
		} else {
			param = "?collation=" + ci.Collation
		}
	}

	// Add charset
	if !strings.Contains(param, "charset") {
		if len(param) > 0 {
			param += "&charset=" + ci.Charset
		} else {
			param = "?charset=" + ci.Charset
		}
	}

	// Example: root:password@tcp(localhost:3306)/test
	s := fmt.Sprintf("%v:%v@tcp(%v:%d)/%v", ci.Username, ci.Password, ci.Hostname, ci.Port, param)

	if includeDatabase {
		s = fmt.Sprintf("%v:%v@tcp(%v:%d)/%v%v", ci.Username, ci.Password, ci.Hostname, ci.Port, ci.Database, param)
	}

	return s
}

// setDefaults sets the charset and collation if they are not set.
func (c Info) setDefaults() Info {
	ci := c

	if len(ci.Charset) == 0 {
		ci.Charset = "utf8"
	}
	if len(ci.Collation) == 0 {
		ci.Collation = "utf8_unicode_ci"
	}

	return ci
}
