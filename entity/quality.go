package entity

type Quality uint8

const (
	VeryBad Quality = iota + 1
	Bad
	Middle
	Good
	VeryGood
)
