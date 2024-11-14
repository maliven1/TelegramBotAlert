package logic

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"todo-orion-bot/entity"
	eventMessage "todo-orion-bot/internal/event-message"
	"todo-orion-bot/storage"
)

func Split(text string) (status string, date string, triggerDate string, err error) {
	newText := strings.Split(text, " ")
	if len(newText) == 4 {
		status = newText[0]
		date = newText[2]
		trigger, err := time.Parse("02.01.2006", date)
		if err != nil {
			fmt.Println("Не верный формат времени")
		}
		triggerDate = trigger.Add(time.Duration(-1 * time.Hour)).Format("02.01.2006")
	}
	return
}

func AddSheet(data []entity.Param, db *storage.Storage) ([]int64, error) {
	var resID = make([]int64, 0)
	var id int64
	for _, v := range data {
		i, err := db.Save(v)
		if err != nil {
			fmt.Errorf("storage.sqlite.Save Error %v", err)
		}
		id = i
	}
	resID = append(resID, id)
	return resID, nil
}

func UpdateSheet(data []entity.Param, db *storage.Storage, id []int64) ([]int64, error) {
	var resID = make([]int64, 0)
	for i, v := range data {
		_, err := db.Update(v, id[i])
		if err != nil {
			return nil, fmt.Errorf("storage.sqlite.Save Error %v", err)
		}
	}
	return resID, nil
}

func CheckEvent(db *storage.Storage, id []int64) (string, error) {
	count := 0
	Name := ""
	for _, v := range id {
		Data, date, err := eventMessage.CheckDate(db, v)
		if err != nil {
			fmt.Errorf("storage.sqlite.CheckDate Error %v", err)
		}
		Name = Data.Name
		if date.Day() == time.Now().Day() && Data.Count == "" {
			Data.Count = "1"
			err := db.UpdateCount(Data, v)
			if err != nil {
				fmt.Errorf("storage.sqlite.UpdateCount Error %v", err)
			}
			return CreateMessageNow(Data), nil
		}
		if Data.Date != "" && Data.Count != "2" {
			Data.Count = "2"
			err := db.UpdateCount(Data, v)
			if err != nil {
				fmt.Errorf("storage.sqlite.UpdateCount Error %v", err)
			}
			return CreateMessageAfter(Data), nil
		}

		if date.Before(time.Now()) {
			count = 1

		}

	}
	if count != 0 {
		return CreateMessageTime(Name), nil
	}
	return "", nil
}

func CreateMessageTime(Name string) string {
	message := Name + "," + "Давно не вносил новую задачу"

	return message
}

func CreateMessageNow(data entity.EventData) string {
	uMessage := getUrgentMessage()
	message := data.Name + "," + uMessage + " Задача: " + data.Task

	return message
}
func CreateMessageAfter(data entity.EventData) string {
	oMesasge := getOverdueMessage()
	message := data.Name + "," + oMesasge + " Задача: " + data.Task

	return message
}

func getUrgentMessage() string {
	messages := []string{
		"Срочно займись этой задачей!",
		"Не откладывай, пора браться за работу!",
		"Эта задача требует твоего внимания сейчас!",
		"Необходимо немедленно приступить к этой задаче!",
		"Успех ждет, начни работу сейчас!",
		"Давайте сделаем это прямо сейчас!",
		"Эта задача не ждет, действуй!",
		"Скорее всего, это дело требует твоей руки!",
		"Пора взяться за дело и не откладывать!",
		"Будь настойчив, начни выполнять задачу!",
	}

	// Генерация случайного индекса
	src := rand.NewSource(time.Now().UnixNano())
	index := src.Int63() % int64(len(messages))

	return messages[index]
}
func getOverdueMessage() string {
	messages := []string{
		"Эта задача просрочена! Срочно занимайся ею!",
		"Время вышло! Необходимо немедленно начать работу над этой задачей!",
		"Срок выполнения задачи истек! Пора брать ее в работу!",
		"Ты просрочил задачу! Работай над ней прямо сейчас!",
		"Эта работа ждет! Пора взяться за просроченную задачу!",
		"Не откладывай! У тебя есть просроченная задача!",
		"Срочно займись этой задачей, срок которой истек!",
		"Не позволяй этой просроченной задаче оставаться без внимания!",
		"Проблема: задача просрочена! Действуй немедленно!",
		"Эту задачу нельзя оставлять без решения! Начинай работать прямо сейчас!",
	}

	// Генерация случайного индекса
	src := rand.NewSource(time.Now().UnixNano())
	index := src.Int63() % int64(len(messages))

	return messages[index]
}
