package util

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

var dump *os.File

func init() {
	var err error
	dump, err = os.OpenFile("app.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal("cannot open dump file: ", err)
	}
}

func DumpLog(value interface{}) {
	if dump != nil {
		spew.Fdump(dump, value)
	}
}
