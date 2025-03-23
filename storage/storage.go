package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"user-actions-api/models"
)

var Users []models.User
var Actions []models.Action

func LoadData() error {
	usersFile, err := os.ReadFile("data/users.json")
	if err != nil {
		return fmt.Errorf("failed to read users file: %w", err)
	}
	if err := json.Unmarshal(usersFile, &Users); err != nil {
		return fmt.Errorf("failed to parse users JSON: %w", err)
	}

	actionsFile, err := os.ReadFile("data/actions.json")
	if err != nil {
		return fmt.Errorf("failed to read actions file: %w", err)
	}
	if err := json.Unmarshal(actionsFile, &Actions); err != nil {
		return fmt.Errorf("failed to parse actions JSON: %w", err)
	}

	log.Printf("Loaded %d users and %d actions", len(Users), len(Actions))
	return nil
}