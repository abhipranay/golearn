package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("/Users/abhipranay.chauhan/Documents/work/shopee/abhipranay-60dfdd149ae5-sheet.json")
	checkError(err)

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)
	spreadSheetId := "1xItPkaD_bdx0MGKcbcJjIvrcXHbUvhdePUr-qtinoPc"
	spreadSheet, err := service.FetchSpreadsheet(spreadSheetId)
	checkError(err)
	for i, s := range spreadSheet.Sheets {
		fmt.Println(fmt.Sprintf("%d, %d", i, s.Properties.ID))
		for _, row := range s.Rows {
			for _, cell := range row {
				fmt.Print(cell.Value, " ")
			}
			fmt.Println()
		}
	}
	//fmt.Println(spreadSheet.Sheets)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}