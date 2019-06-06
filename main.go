package main

import (
	"flag"
	"github.com/jasontconnell/sqlrun/config"
	"github.com/jasontconnell/sqlrun/process"
	"os"
	"log"
)

func main() {
	configFile := flag.String("config", "config.json", "config file with connection string")
	dir := flag.String("dir", "", "the directory to recursively search")
	p := flag.String("p", "", "priority prefixes. like tbl,Save,Get,Delete for filenames")
	flag.Parse()

	if *dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	files, err := process.GetSqlFiles(*dir, *p, "sql")
	if err != nil {
		log.Fatal("Got error getting sql files, ", err)
	}

	cfg := config.LoadConfig(*configFile)
	err = process.RunAll(cfg.ConnectionString, files)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}
