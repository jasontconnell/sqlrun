package main

import (
	"flag"
	"fmt"
	"github.com/jasontconnell/sqlrun/config"
	"github.com/jasontconnell/sqlrun/process"
	"os"
)

func main() {
	dir := flag.String("dir", "", "the directory to recursively search")
	p := flag.String("p", "", "priority prefixes. like tbl,Save,Get,Delete for filenames")
	flag.Parse()

	if *dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	files, err := process.GetSqlFiles(*dir, *p, "sql")

	if err != nil {
		fmt.Println("Got error getting sql files", err)
	}

	cfg := config.LoadConfig("config.json")
	err = process.RunAll(cfg.ConnectionString, files)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("done")
}
