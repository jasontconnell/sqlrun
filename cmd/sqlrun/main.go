package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jasontconnell/sqlrun/config"
	"github.com/jasontconnell/sqlrun/process"
)

func main() {
	configFile := flag.String("config", "config.json", "config file with connection string")
	dir := flag.String("dir", "", "the directory to recursively search")
	p := flag.String("p", "", "priority prefixes. like tbl,Save,Get,Delete for filenames")
	ex := flag.String("exclude", "", "exclude prefixes. like tbl,Save,Get,Delete for filenames")
	filter := flag.String("filter", "", "filter to include filenames only containing the filter")
	dry := flag.Bool("dry", false, "dry run")
	flag.Parse()

	start := time.Now()
	if *dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	baseDir := *dir
	if !filepath.IsAbs(baseDir) {
		cwd, _ := os.Getwd()
		baseDir = filepath.Join(cwd, baseDir)
	}

	files, err := process.GetSqlFiles(baseDir, *p, *ex, *filter, "sql")
	if err != nil {
		log.Fatal("Got error getting sql files, ", err)
	}

	if *dry {
		log.Println("Dry run")
		for _, f := range files {
			_, filename := filepath.Split(f)
			log.Println(filename)
		}
		return
	}

	cfg := config.LoadConfig(*configFile)

	log.Println("Running", len(files), "sql files against", cfg.ConnectionString)
	err = process.RunAll(cfg.ConnectionString, files)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done", time.Since(start))
}
