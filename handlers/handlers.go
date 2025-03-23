package handlers

import (
	"math"
	"sort"

	"github.com/gofiber/fiber/v2"

	"user-actions-api/models"
	"user-actions-api/storage"
	"user-actions-api/types"
)

// GetUserByID returns a user by ID
func GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	for _, user := range storage.Users {
		if user.ID == id {
			return c.JSON(user)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "User not found",
	})
}

// GetUserActionCount returns the total number of actions for a user
func GetUserActionCount(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	count := 0
	for _, action := range storage.Actions {
		if action.UserID == id {
			count++
		}
	}

	return c.JSON(types.ActionCount{Count: count})
}

// GetNextActionBreakdown returns a breakdown of next actions after a specific action type
func GetNextActionBreakdown(c *fiber.Ctx) error {
	actionType := c.Params("type")
	if actionType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Action type is required",
		})
	}

	actionIndices := make(map[int]int, len(storage.Actions))
	for i, action := range storage.Actions {
		actionIndices[action.ID] = i
	}

	sortedActions := make([]models.Action, len(storage.Actions))
	copy(sortedActions, storage.Actions)
	sort.Slice(sortedActions, func(i, j int) bool {
		return sortedActions[i].CreatedAt.Before(sortedActions[j].CreatedAt)
	})

	userActions := make(map[int][]models.Action)
	for _, action := range sortedActions {
		userActions[action.UserID] = append(userActions[action.UserID], action)
	}

	nextActions := make(map[string]int)
	totalNextActions := 0

	for _, userActionList := range userActions {
		for i := 0; i < len(userActionList)-1; i++ {
			if userActionList[i].Type == actionType {
				nextActionType := userActionList[i+1].Type
				nextActions[nextActionType]++
				totalNextActions++
			}
		}
	}

	probabilities := make(types.NextActionProbability)
	if totalNextActions > 0 {
		for actionType, count := range nextActions {
			probability := float64(count) / float64(totalNextActions)
			probabilities[actionType] = math.Round(probability*100) / 100
		}
	}

	return c.JSON(probabilities)
}

// GetReferralIndices returns the referral index for all users
func GetReferralIndices(c *fiber.Ctx) error {
	referralGraph := make(map[int][]int)
	
	for _, action := range storage.Actions {
		if action.Type == "REFER_USER" && action.TargetUser != 0 {
			referralGraph[action.UserID] = append(referralGraph[action.UserID], action.TargetUser)
		}
	}

	referralIndices := make(types.ReferralIndexResponse)
	
	for _, user := range storage.Users {
		referralIndices[user.ID] = 0
	}
	
	visited := make(map[int]bool)
	
	var countReferrals func(userID int) int
	countReferrals = func(userID int) int {
		if count, exists := referralIndices[userID]; exists && count > 0 {
			return count
		}
		
		visited[userID] = true
		
		count := len(referralGraph[userID])
		
		for _, referredUserID := range referralGraph[userID] {
			if !visited[referredUserID] {
				count += countReferrals(referredUserID)
			}
		}
		
		referralIndices[userID] = count
		
		visited[userID] = false
		
		return count
	}
	
	for _, user := range storage.Users {
		if !visited[user.ID] {
			countReferrals(user.ID)
		}
	}
	
	return c.JSON(referralIndices)
}