package main

import (
	"time"
)

type PatchType int

const (
	Undefined PatchType = iota
	Flower
	Herb
	FruitTree
	Celastrus
	Redwood

	// Allotment
	Potato
	Onion
	Cabbage
	Tomato
	Sweetcorn
	Strawberry
	Watermelon
	SnapeGrass

	// Bushes
	Redberries
	Cadavaberries
	Dwellberries
	Jangerberries
	Whiteberries
	PoisonIvy

	// Trees
	Oak
	Willow
	Maple
	Yew
	Magic

	// Cactus
	Cactus
	PotatoCactus
)

func (p PatchType) Stages() uint {
	switch p {
	case Flower, Herb, Potato, Onion, Cabbage, Tomato, Oak:
		return 5
	case Celastrus, Redberries:
		return 6
	case FruitTree, Sweetcorn, Strawberry, Cadavaberries, Willow:
		return 7
	case SnapeGrass, Dwellberries, Cactus, PotatoCactus:
		return 8
	case Watermelon, Jangerberries, Whiteberries, PoisonIvy, Maple:
		return 9
	case Redwood, Yew:
		return 11
	case Magic:
		return 13
	}
	return 0
}

func (p PatchType) TickRate() time.Duration {
	switch p {
	case Flower:
		return 5 * time.Minute
	case Potato, Onion, Cabbage, Tomato, Sweetcorn, Strawberry, Watermelon, SnapeGrass, PotatoCactus:
		return 10 * time.Minute
	case Herb, Redberries, Cadavaberries, Dwellberries, Jangerberries, Whiteberries, PoisonIvy:
		return 20 * time.Minute
	case Oak, Willow, Maple, Yew, Magic:
		return 40 * time.Minute
	case Cactus:
		return 80 * time.Minute
	case FruitTree, Celastrus:
		return 160 * time.Minute
	case Redwood:
		return 640 * time.Minute
	}
	return -1
}

func FindPatchType(s string) PatchType {
	switch s {
	case "flowers", "flower":
		return Flower
	case "herb", "herbs":
		return Herb
	case "fruittree", "fruit", "ft":
		return FruitTree
	case "celastrus", "cel":
		return Celastrus
	case "redwood", "red", "rw":
		return Redwood
	case "potato", "potatos":
		return Potato
	case "onion", "onions":
		return Onion
	case "cabbage", "cabbages", "cabige":
		return Cabbage
	case "tomato", "tomatos":
		return Tomato
	case "sweetcorn", "corn", "sc":
		return Sweetcorn
	case "strawberry", "strawberries", "strawb", "sb":
		return Strawberry
	case "watermelon", "watermelons", "melon", "melons":
		return Watermelon
	case "snapegrass", "snape":
		return SnapeGrass
	case "redberries", "redberry", "rb":
		return Redberries
	case "cadavaberries", "cadavaberry", "cadavab", "cadava", "cada":
		return Cadavaberries
	case "dwellberries", "dwellberry", "dwell":
		return Dwellberries
	case "jangerberries", "jangerberry", "janger", "jb":
		return Jangerberries
	case "whiteberries", "whiteberry", "white", "wb":
		return Whiteberries
	case "poisonivy", "ivy":
		return PoisonIvy
	case "oak":
		return Oak
	case "willow":
		return Willow
	case "maple":
		return Maple
	case "yew":
		return Yew
	case "magic", "mage":
		return Magic
	case "cactus", "cact":
		return Cactus
	case "potatocactus", "potcactus", "potcact", "pc":
		return PotatoCactus
	}
	return Undefined
}
