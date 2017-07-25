package main

import (
	"fmt"
	"log"

	"github.com/LindsayBradford/go-dbf/godbf"
)

func main() {
	// NOT WORK
	// csKOI8R windows-1251 windows-1250
	// WORK - DOS 866 - Russian OEM
	dbfTable, err := godbf.NewFromFile("./101-20170701/NAMES.DBF", "866")
	if err != nil {
		log.Fatalf("Ошибка открытия %v", err)
	}

	fmt.Printf("%v\n", len(dbfTable.FieldNames()))

	for i := 0; i <= dbfTable.NumberOfRecords()-1; i++ {
		for y := 0; y <= len(dbfTable.FieldNames())-1; y++ {
			fmt.Printf("%v;", dbfTable.FieldValue(i, y))
		}
		fmt.Println()
	}
}
