package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/0xN0x/go-artifactsmmo"
)

var logger = log.New(os.Stdout, "", 2)

func usage() {
	fmt.Println("Usage: go run main.go <api-token> <character-name>")
}

func GatherOre(client *artifactsmmo.ArtifactsMMO, code string, quantity int) {
	// Move to the map
	switch code {
	case "copper":
		client.Move(2, 0)
	}

	// Mine the ore
	for i := 0; i < quantity; i++ {
		res, err := client.Gather()

		if err != nil {
			logger.Println("Error:", err)
			return
		}

		for _, item := range res.Details.Items {
			for _, inventoryItem := range res.Character.Inventory {
				if inventoryItem.Code == item.Code {
					logger.Printf("[Gathering] Got +%d %s (total: %d)\n", item.Quantity, item.Code, inventoryItem.Quantity)
				}
			}
		}

		time.Sleep(time.Duration(res.Cooldown.RemainingSeconds) * time.Second)
	}
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	client := artifactsmmo.NewClient(os.Args[1], os.Args[2])
	character, err := client.GetCharacterInfo(os.Args[2])
	if err != nil {
		logger.Println("Error:", err)
		return
	}

	logger.Printf("Welcome, %s (XP: %d/%d)!\n", character.Name, character.Xp, character.MaxXp)

	logger.Printf("[Task(%s)] %d/%d %s", character.TaskType, character.TaskProgress, character.TaskTotal, character.Task)

	// If task complete, move to task manager to get new task
	if character.TaskProgress == character.TaskTotal {
		client.Move(1, 2)
		client.CompleteTask()
	}

	// GatherOre(client, "copper", 50)
	// return

	// client.Move(0, 1)
	for {
		fight, err := client.Fight()
		if err != nil {
			logger.Println("Error:", err)

			if err.Error() == "character in cooldown" {
				time.Sleep(time.Duration(fight.Cooldown.RemainingSeconds) * time.Second)
				continue
			} else {
				return
			}
		}

		logger.Printf("[%s] +%d XP, +%d Gold\n", fight.Fight.Result, fight.Fight.Xp, fight.Fight.Gold)

		time.Sleep(time.Duration(fight.Cooldown.RemainingSeconds) * time.Second)
	}
}
