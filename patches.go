package main

import (
	"sort"
	"sync"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pkg/errors"
)

type PatchType int

const (
	first PatchType = iota

	// Allotment
	Potato
	Onion
	Cabbage
	Tomato
	Sweetcorn
	Strawberry
	Watermelon
	SnapeGrass

	// Flower
	Flower
	Marigold
	Rosemary
	Nasturtium
	Woad
	Limpwurt
	WhiteLily

	// Herb
	Herb
	Guam
	Marrentill
	Tarromin
	Harralander
	Gout
	Ranarr
	Toadflax
	Irit
	Avantoe
	Kwuarm
	Snapdragon
	Cadantine
	Lantadyme
	DwarfWeed
	Torstol

	// Hop
	Barley
	Hammerstone
	Asgarnian
	Jute
	Yanillian
	Krandorian
	Wildblood

	// Bush
	Redberries
	Cadavaberries
	Dwellberries
	Jangerberries
	Whiteberries
	PoisonIvy

	// Tree
	Oak
	Willow
	Maple
	Yew
	Magic

	// Fruit Tree
	FruitTree
	Apple
	Banana
	Orange
	Curry
	Pineapple
	Papaya
	Palm
	Dragonfruit

	// Special
	GiantSeaweed
	Grapes
	Mushroom
	Belladonna
	Hespori

	// Anima
	Anima
	Kronos
	Iasor
	Attas

	// Special Tree
	Teak
	Mahogany
	Calquat
	Crystal
	Spirit
	Celastrus
	Redwood

	// Cactus
	Cactus
	PotatoCactus

	DidYouMean
	Undefined
)

func (p PatchType) Names() []string {
	switch p {
	case Flower:
		return []string{"flowers", "flower"}
	case Marigold:
		return []string{"marigold"}
	case Rosemary:
		return []string{"rosemary", "rose"}
	case Nasturtium:
		return []string{"nasturtium", "nast"}
	case Woad:
		return []string{"woad"}
	case Limpwurt:
		return []string{"limpwurt", "limp", "lw"}
	case WhiteLily:
		return []string{"whitelily", "lily", "wl"}
	case Herb:
		return []string{"herbs", "herb"}
	case Guam:
		return []string{"guam"}
	case Marrentill:
		return []string{"marrentill", "marr", "mar"}
	case Tarromin:
		return []string{"tarromin", "tarr", "tar"}
	case Harralander:
		return []string{"harralander", "harra", "harr", "har"}
	case Gout:
		return []string{"gout"}
	case Ranarr:
		return []string{"ranarr", "ran"}
	case Toadflax:
		return []string{"toadflax", "toad"}
	case Irit:
		return []string{"irit"}
	case Avantoe:
		return []string{"avantoe", "ava"}
	case Kwuarm:
		return []string{"kwuarm"}
	case Snapdragon:
		return []string{"snapdragon", "snap"}
	case Cadantine:
		return []string{"cadantine", "cad"}
	case Lantadyme:
		return []string{"lantadyme", "lant"}
	case DwarfWeed:
		return []string{"dwarfweed", "dwarf"}
	case Torstol:
		return []string{"torstol", "tors", "torst"}
	case Barley:
		return []string{"barley"}
	case Hammerstone:
		return []string{"hammerstone", "hammer"}
	case Asgarnian:
		return []string{"asgarnian", "ash"}
	case Jute:
		return []string{"jute"}
	case Yanillian:
		return []string{"yanillian", "yan"}
	case Krandorian:
		return []string{"krandorian"}
	case Wildblood:
		return []string{"wildblood"}
	case FruitTree:
		return []string{"fruittree", "fruit", "ft"}
	case Apple:
		return []string{"apple"}
	case Banana:
		return []string{"banana", "nana"}
	case Orange:
		return []string{"orange"}
	case Curry:
		return []string{"curry"}
	case Pineapple:
		return []string{"pineapple", "pine"}
	case Papaya:
		return []string{"papaya", "pap"}
	case Palm:
		return []string{"palm"}
	case Dragonfruit:
		return []string{"dragonfruit", "dfruit", "dragon", "df"}
	case Teak:
		return []string{"teak"}
	case Mahogany:
		return []string{"mahogany", "mahog", "mah"}
	case Calquat:
		return []string{"calquat", "calq"}
	case Crystal:
		return []string{"crystal"}
	case Spirit:
		return []string{"spirit"}
	case Celastrus:
		return []string{"celastrus", "celast", "cel"}
	case Redwood:
		return []string{"redwood", "rw"}
	case GiantSeaweed:
		return []string{"giantseaweed", "seaweed", "sw"}
	case Grapes:
		return []string{"grapes", "grape"}
	case Mushroom:
		return []string{"mushroom", "shroom", "mush"}
	case Belladonna:
		return []string{"belladonna", "bella"}
	case Hespori:
		return []string{"hespori", "hesp"}
	case Potato:
		return []string{"potato", "potatoes", "potatos", "tater"}
	case Onion:
		return []string{"onion", "onions"}
	case Cabbage:
		return []string{"cabbage", "cabbages", "cabige"}
	case Tomato:
		return []string{"tomato", "tomatoes", "tomatos"}
	case Sweetcorn:
		return []string{"sweetcorn", "corn", "sc"}
	case Strawberry:
		return []string{"strawberry", "strawberries", "strawb", "sb"}
	case Watermelon:
		return []string{"watermelon", "watermelons", "melon", "melons"}
	case SnapeGrass:
		return []string{"snapegrass", "snape"}
	case Redberries:
		return []string{"redberries", "redberry", "rb"}
	case Cadavaberries:
		return []string{"cadavaberries", "cadavaberry", "cadavab", "cadava"}
	case Dwellberries:
		return []string{"dwellberries", "dwellberry", "dwell"}
	case Jangerberries:
		return []string{"jangerberries", "jangerberry", "janger", "jang", "jb"}
	case Whiteberries:
		return []string{"whiteberries", "whiteberry", "white", "wb"}
	case PoisonIvy:
		return []string{"poisonivy", "ivy"}
	case Oak:
		return []string{"oak"}
	case Willow:
		return []string{"willow"}
	case Maple:
		return []string{"maple"}
	case Yew:
		return []string{"yew"}
	case Magic:
		return []string{"magic", "mage"}
	case Cactus:
		return []string{"cactus", "cact"}
	case PotatoCactus:
		return []string{"potatocactus", "potcactus", "potcact", "pc"}
	case Anima:
		return []string{"anima"}
	case Kronos:
		return []string{"kronos"}
	case Iasor:
		return []string{"iasor"}
	case Attas:
		return []string{"attas"}
	default:
		return []string{}
	}
}

func (p PatchType) Name() string {
	names := p.Names()
	if len(names) > 0 {
		return names[0]
	}
	return "undefined"
}

func (p PatchType) getTickTime(offset, ticks int64) time.Time {
	tickRate := int64(p.TickRate().Minutes())
	calcOffset := (offset % tickRate * 60)
	unixNow := time.Now().Unix() + calcOffset

	currentTick := (unixNow - (unixNow % (tickRate * 60)))
	goalTick := currentTick + (ticks * tickRate * 60)

	return time.Unix(goalTick-calcOffset, 0)
}

func (p PatchType) Stages() uint {
	switch p {
	case Hespori:
		return 3
	case Barley, Hammerstone, GiantSeaweed, Belladonna:
		return 4
	case Asgarnian, Jute, Flower, Marigold, Rosemary, Nasturtium, Woad, Limpwurt, WhiteLily, Herb, Guam, Marrentill, Tarromin, Harralander, Gout, Ranarr, Toadflax, Irit, Avantoe, Kwuarm, Snapdragon, Cadantine, Lantadyme, DwarfWeed, Torstol, Potato, Onion, Cabbage, Tomato, Oak:
		return 5
	case Crystal, Celastrus, Redberries, Yanillian, Mushroom:
		return 6
	case Teak, Grapes, FruitTree, Apple, Banana, Orange, Curry, Pineapple, Papaya, Palm, Dragonfruit, Sweetcorn, Strawberry, Cadavaberries, Willow, Krandorian:
		return 7
	case Calquat, Mahogany, SnapeGrass, Dwellberries, Cactus, PotatoCactus, Wildblood, Kronos, Iasor, Attas:
		return 8
	case Watermelon, Jangerberries, Whiteberries, PoisonIvy, Maple:
		return 9
	case Redwood, Yew:
		return 11
	case Spirit:
		return 12
	case Magic:
		return 13
	default:
		return 0
	}
}

func (p PatchType) TickRate() time.Duration {
	switch p {
	case Grapes, Flower, Marigold, Rosemary, Nasturtium, Woad, Limpwurt, WhiteLily:
		return 5 * time.Minute
	case Potato, Onion, Cabbage, Tomato, Sweetcorn, Strawberry, Watermelon, SnapeGrass, PotatoCactus, Barley, Hammerstone, Asgarnian, Jute, Yanillian, Krandorian, Wildblood, GiantSeaweed:
		return 10 * time.Minute
	case Herb, Guam, Marrentill, Tarromin, Harralander, Gout, Ranarr, Toadflax, Irit, Avantoe, Kwuarm, Snapdragon, Cadantine, Lantadyme, DwarfWeed, Torstol, Redberries, Cadavaberries, Dwellberries, Jangerberries, Whiteberries, PoisonIvy:
		return 20 * time.Minute
	case Oak, Willow, Maple, Yew, Magic, Mushroom:
		return 40 * time.Minute
	case Cactus, Belladonna, Crystal:
		return 80 * time.Minute
	case Calquat, FruitTree, Apple, Banana, Orange, Curry, Pineapple, Papaya, Palm, Dragonfruit, Celastrus:
		return 160 * time.Minute
	case Spirit:
		return 320 * time.Minute
	case Redwood, Hespori, Kronos, Iasor, Attas, Teak, Mahogany:
		return 640 * time.Minute
	default:
		return -1 * time.Minute
	}
}

var (
	allPatchesCache    map[string]PatchType
	allPatchNamesCache []string
	cacheInitialized   bool
	cacheMutex         sync.Mutex
)

func initializePatchCache() {
	allPatchesCache = make(map[string]PatchType)
	allPatchNamesCache = nil

	for p := first; p < Undefined; p++ {
		names := p.Names()
		allPatchNamesCache = append(allPatchNamesCache, names...)
		for _, name := range names {
			allPatchesCache[name] = p
		}
	}

	sort.Strings(allPatchNamesCache)
	cacheInitialized = true
}

func AllPatchNames() []string {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if !cacheInitialized {
		initializePatchCache()
	}

	return allPatchNamesCache
}

func AllPatches() map[string]PatchType {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if !cacheInitialized {
		initializePatchCache()
	}

	return allPatchesCache
}

func FindPatchType(search string) (PatchType, error) {
	patch, exists := AllPatches()[search]
	if exists && patch != Undefined {
		return patch, nil
	}
	suggestions := fuzzy.RankFind(search, AllPatchNames())
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Distance < suggestions[j].Distance
	})
	switch {
	case len(suggestions) > 0 && suggestions[0].Distance < 10:
		return DidYouMean, errors.Errorf("Did you mean? %v", suggestions[0].Target)
	default:
		return Undefined, errors.Errorf("Not sure the patch :o")
	}
}
