package excelpackage

import (
	"github.com/xuri/excelize/v2"
)

func Basci() string {
	f := excelize.NewFile()
	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet2", "A2", "Hello World")
	f.SetCellValue("Sheet1", "B2", 100)
	f.SetActiveSheet(index)

	excelsavename := "Book1.xlsx"
	if err := f.SaveAs(excelsavename); err != nil {
		panic(err)
	}

	//err := os.Rename(excelsavename, "excelfolder/")
	//if err != nil {
	//panic(err)
	//}

	return excelsavename
}
