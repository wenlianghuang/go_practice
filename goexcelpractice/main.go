package main

import (
	"excelpractice/excelpackage"
	"os"
)

func main() {
	basicexcelfile := excelpackage.Basci()
	excelfolder := "excelfolder/"
	err := os.Rename(basicexcelfile, excelfolder+basicexcelfile)
	if err != nil {
		panic(err)
	}

	excelpackage.Chartinexcel()
}
