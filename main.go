package main

import (
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func mainE() error {
	args := os.Args
	if len(args) != 3 {
		return errors.New("usage: sgix <file> <dir>")
	}
	basename := args[1]
	dest := args[2]
	if strings.HasSuffix(basename, ".idb") {
		basename = basename[:len(basename)-4]
	} else if strings.HasSuffix(basename, ".sw") {
		basename = basename[:len(basename)-3]
	}
	ents, err := readIDB(basename + ".idb")
	if err != nil {
		return err
	}
	return extract(ents, basename+".sw", dest)
}

func main() {
	if err := mainE(); err != nil {
		logrus.Errorln("Error:", err)
		os.Exit(1)
	}
}
