package process

import (
	"database/sql"
	"log"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
	"io/ioutil"
	"regexp"
)

var sepreg *regexp.Regexp = regexp.MustCompile(`\W(?i)go\W`)

func RunAll(connstr string, paths []string) error {
	var err error
	for _, f := range paths {
		err = Run(connstr, f)
		if err != nil {
			return errors.Wrapf(err, "executing %s", f)
		}
	}

	return nil
}

func Run(connstr, path string) error {
	log.Println(path)
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return errors.Wrapf(err, "couldn't open file %s", path)
	}

	db, err := sql.Open("mssql", connstr)
	if err != nil {
		return errors.Wrap(err, "couldn't open connection to sql server")
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
		return errors.Wrap(err, "couldn't start transaction")
	}
	_, err = db.Exec(contents)
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "couldn't execute statement %s", contents)
	}

	err = tx.Commit()

	if err != nil {
		return errors.Wrap(err, "couldn't commit transaction")
	}

	return nil
}
