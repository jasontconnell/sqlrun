package process

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	_ "github.com/microsoft/go-mssqldb"
)

var sepreg *regexp.Regexp = regexp.MustCompile(`\W(?i)go\W`)

func RunAll(connstr string, paths []string) error {
	var err error
	for _, f := range paths {
		err = Run(connstr, f)
		if err != nil {
			return fmt.Errorf("executing %s. %w", f, err)
		}
	}

	return nil
}

func Run(connstr, path string) error {
	_, fn := filepath.Split(path)
	log.Println(fn)
	contents, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("couldn't open file %s. %w", path, err)
	}

	db, err := sql.Open("mssql", connstr)
	if err != nil {
		return fmt.Errorf("couldn't open connection to sql server %s. %w", connstr, err)
	}
	defer db.Close()

	sqls := sepreg.Split(string(contents), -1)
	for _, sql := range sqls {
		err := runsql(db, sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func runsql(db *sql.DB, contents string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("couldn't start transaction. %w", err)
	}
	_, err = db.Exec(contents)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("couldn't execute statement %s. %w", contents, err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("couldn't commit transaction. %w", err)
	}

	return nil
}
