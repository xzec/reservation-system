package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}
	dbstring, command := args[0], args[1]

	db, err := goose.OpenDBWithDriver("postgres", dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err = goose.RunContext(context.Background(), command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
