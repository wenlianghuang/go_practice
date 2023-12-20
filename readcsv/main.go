package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fd, error := os.Open("./data.csv")
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()

	fileReader := csv.NewReader(fd)
	fileReader.Comma = ','
	for {
		record, err := fileReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		a := record[0]
		b := record[1]
		c := record[2]
		fmt.Println(a, b, c)
	}
	/*
		records, error := fileReader.ReadAll()
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println(records[1])
	*/
}
