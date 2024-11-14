package eventMessage

import (
	"time"
	"todo-orion-bot/entity"
	"todo-orion-bot/storage"
)

func CheckDate(db *storage.Storage, id int64) (entity.EventData, time.Time, error) {
	var resData entity.EventData
	var resDate time.Time
	Data, err := db.Get(id)
	if err != nil {
		return resData, resDate, err
	}
	if Data.Date != "" {
		resDate, err := time.Parse("02.01.2006", Data.Date)
		if err != nil {
			return resData, resDate, err
		}
		if !resDate.After(time.Now()) && Data.Status != "Выполнено" {
			resData.Task = Data.Task
			resData.Date = Data.Date
			resData.Name = Data.Name
			resData.Status = Data.Status
			resData.Count = Data.Count
			return resData, resDate, nil
		}
	}

	return resData, resDate, nil
}
