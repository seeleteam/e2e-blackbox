package store

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/seeleteam/go-seele/database"

	"github.com/seeleteam/go-seele/database/leveldb"
)

// ...
const (
	DbName   = "Seele-blacke2e-test"
	CoverKey = "Seele-cover-test"
)

// DB ...
var db database.Database

func init() {
	var err error
	if db, err = prepareDB(DbName); err != nil {
		fmt.Println("create db err:", err)
		os.Exit(1)
	}
}

func prepareDB(dbName string) (database.Database, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(usr.HomeDir, dbName)
	return leveldb.NewLevelDB(dbPath)

}

// Save the e2e test result
func Save(date string, coverbyte []byte) {

	db.Put([]byte(date+CoverKey), coverbyte)
}

// Get the e2e test result
func Get(date string) (coverbyte []byte) {

	coverbyte, err := db.Get([]byte(date + CoverKey))
	if err != nil {
		fmt.Println("get cover result err:", err)
		return
	}

	return coverbyte
}
