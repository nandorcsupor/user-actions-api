package types

type ReferralIndexResponse map[int]int
type NextActionProbability map[string]float64

type ActionCount struct {
	Count int `json:"count"`
}