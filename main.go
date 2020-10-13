package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func mainE() error {
	args := os.Args
	if len(args) < 2 || len(args) > 4 {
		return errors.New("usage: sgix <file.idb> [<data> [<dir>]]")
	}
	idbfile := args[1]
	ents, err := readIDB(idbfile)
	if err != nil {
		return err
	}
	if len(args) < 3 {
		return nil
	}
	datafile := args[2]
	var dest string
	if len(args) >= 4 {
		dest = args[3]
		if dest == "" {
			return errors.New("invalid destination directory")
		}
		fmt.Println("Extracting...")
	} else {
		fmt.Println("Verifying...")
	}
	return extract(ents, datafile, dest)
}

func main() {
	if err := mainE(); err != nil {
		logrus.Errorln("Error:", err)
		os.Exit(1)
	}
}
