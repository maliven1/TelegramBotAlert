package worker

import (
	"fmt"
	"time"
	"todo-orion-bot/google-tab/logic"
	"todo-orion-bot/google-tab/logic/sheet"
	"todo-orion-bot/storage"
)

func Created(db *storage.Storage, Sheet string, RangeSheet string, name string, ch chan string) {
	data, _ := sheet.SheetSearch(Sheet, RangeSheet, name)
	_, err := logic.AddSheet(data, db)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {

			data, id := sheet.SheetSearch(Sheet, RangeSheet, name)
			_, _ = logic.UpdateSheet(data, db, id)
			data, id = sheet.SheetSearch(Sheet, RangeSheet, name)
			_, err = logic.AddSheet(data, db)
			if err != nil {
				fmt.Println(err)
			}
			db.Delete(int64(len(id)))

			time.Sleep(3 * time.Second)
			messageEvent, err := logic.CheckEvent(db, id)
			if err != nil {
				fmt.Println(err)
			}
			if messageEvent != "" {
				ch <- messageEvent
			}
			time.Sleep(1 * time.Hour)
		}
	}()
}
