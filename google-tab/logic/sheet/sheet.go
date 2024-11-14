package sheet

import (
	"context"
	"fmt"
	"log"
	"os"

	entity "todo-orion-bot/entity"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func SheetSearch(sheetID string, rageSheet string, name string) ([]entity.Param, []int64) {
	var info entity.Param
	var data []entity.Param
	// Replace 'path/to/your/credentials.json' with the actual path to your downloaded JSON key file.
	var id []int64
	credentialsFile := os.Getenv("JSON_KEY_PATH")

	// Replace 'your-spreadsheet-id' with the ID of the spreadsheet you want to read.
	spreadsheetID := sheetID
	// Replace 'Sheet1!A1:B10' with the range of cells you want to read.
	readRange := rageSheet
	// Initialize Google Sheets API
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		log.Fatalf("Unable to initialize Sheets API: %v", err)
	}
	// Read data from the specified range
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Print the values from the response
	if len(resp.Values) > 0 {
		fmt.Println("Data from sheet:")
		for i, row := range resp.Values {
			for i, cell := range row {
				str, ok := cell.(string)
				if ok {
					if i == 0 {
						info.Task = str
						info.Name = name
					} else if i == 1 {
						info.Status = str
					} else if i == 3 {
						info.Date = str
					}
				} else {
					fmt.Println("Failed")
				}
			}
			data = append(data, info)
			id = append(id, int64(i+1))
		}
	} else {
		fmt.Println("No data found.")
	}
	return data, id
}
