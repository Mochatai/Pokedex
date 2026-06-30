package helpFunc

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mochatai/pokedex/helpFunc/pokecache"
)

type configCommand struct {
	Url      string
	Next     string
	Previous string
	cache    *pokecache.Cache
	pokedex  map[string]pokemonInfo
}

type areaPokeRes struct {
	Name string `json:"name"`
}

type pokemonArea struct {
	Pokemon_encounters []pokedexNameMap `json:"pokemon_encounters"`
}

type pokedexNameMap struct {
	Pokemon pokedexName `json:"pokemon"`
}

type pokedexName struct {
	Name string `json:"name"`
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type pokemonInfo struct {
	Base_experience int            `json:"base_experience"`
	Height          int            `json:"height"`
	Weight          int            `json:"weight"`
	Stats           []pokemonStats `json:"stats"`
	Types           []pokemonTypes `json:"types"`
}

type pokemonStats struct {
	Base_stat int          `json:"base_stat"`
	Stat      pokemonState `json:"Stat"`
}

type pokemonState struct {
	Name string `json:"name"`
}

type pokemonTypes struct {
	Type pokemonType `json:"type"`
}

type pokemonType struct {
	Name string `json:"name"`
}

var con configCommand
var exploreName []string
var catchName []string
var inspectName []string
var commands map[string]cliCommand

func init() {
	con = configCommand{
		Url:      "https://pokeapi.co/api/v2/location-area/",
		Previous: "1",
		Next:     "1",
		cache:    pokecache.NewCache(time.Second * 10),
		pokedex:  make(map[string]pokemonInfo, 0),
	}

	exploreName = make([]string, 0)
	catchName = make([]string, 0)
	inspectName = make([]string, 0)

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "explore pokemons in the area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "to try to catch pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "to show pokemon info by providing a name",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "to show pokemon you have",
			callback:    commandPokedex,
		},
	}

}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage\n\n")

	for _, val := range commands {
		fmt.Printf("%v: %v\n", val.name, val.description)
	}

	return nil

}

func commandMap() error {

	counter := 0
	url := con.Url
	areasL := make([]string, 0)

	jsonAreaName := areaPokeRes{}

	client := &http.Client{}
	for i, _ := strconv.Atoi(con.Next); counter < 20; i++ {

		if val, ok := con.cache.Get(url + strconv.Itoa(i)); ok {
			areasL = append(areasL, string(val))
			counter++
			continue
		}

		resp, err := client.Get(url + strconv.Itoa(i))
		if err != nil {
			log.Fatalf("Failed function command map: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			break
		}

		err = json.NewDecoder(resp.Body).Decode(&jsonAreaName)
		if err != nil {
			os.Exit(0)
		}
		con.cache.Add(url+strconv.Itoa(i), []byte(jsonAreaName.Name))
		resp.Body.Close()
		areasL = append(areasL, jsonAreaName.Name)
		counter++
	}

	for _, i := range areasL {
		fmt.Println(i)
	}
	tempNext, _ := strconv.Atoi(con.Next)
	con.Next = strconv.Itoa(tempNext + counter)

	if res, _ := strconv.Atoi(con.Next); res >= 41 {
		temp := (res - 1) / 20
		con.Previous = strconv.Itoa(((temp - 2) * 20) + 1)
	}

	return nil
}

func commandMapb() error {

	tempRes, _ := strconv.Atoi(con.Next)

	if tempRes >= 41 {
		con.Next = con.Previous
		commandMap()
		return nil
	}

	fmt.Println("you're on the first page")
	return nil
}

func CleanInput(text string) []string {
	var str []string

	if len(text) == 0 {
		return str
	}

	lowerText := strings.ToLower(text)
	trimLowerText := strings.TrimSpace(lowerText)

	str = strings.Fields(trimLowerText)

	return str
}

func setEploreName(names []string) {
	exploreName = names
}

func setCatchName(names []string) {
	catchName = names
}

func setInspectName(names []string) {
	inspectName = names
}
func commandExplore() error {

	if len(exploreName) < 2 {
		return nil
	}

	areaName := exploreName[1]

	url := con.Url

	areaPokesNames := make([]string, 0)
	jsonAreaPokedx := pokemonArea{}

	client := http.Client{}
	fmt.Printf("Exploring %v...\nFound Pokemon:\n", areaName)
	if val, ok := con.cache.Get(url + areaName); ok {

		byteDataNames := string(val)

		for i := 0; i < len(byteDataNames); i++ {
			fmt.Printf("- %v\n", byteDataNames[i])
		}
		return nil
	}

	resp, err := client.Get(url + areaName)

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&jsonAreaPokedx)
	if err != nil {
		os.Exit(0)
	}

	resp.Body.Close()

	for i := 0; i < len(jsonAreaPokedx.Pokemon_encounters); i++ {
		areaPokesNames = append(areaPokesNames, jsonAreaPokedx.Pokemon_encounters[i].Pokemon.Name)
	}

	tempDataNames := strings.Join(areaPokesNames, " ")
	con.cache.Add(url+areaName, []byte(tempDataNames))

	for i := 0; i < len(areaPokesNames); i++ {
		fmt.Printf("- %v\n", areaPokesNames[i])
	}

	return nil

}

func commandCatch() error {

	if len(catchName) < 2 {
		return nil
	}
	pokeName := catchName[1]
	url := "https://pokeapi.co/api/v2/pokemon/"
	pokeInfo := pokemonInfo{}

	client := http.Client{}

	resp, err := client.Get(url + pokeName)

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		fmt.Print("Pokemon not found")
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&pokeInfo)
	if err != nil {
		os.Exit(0)
	}

	randomInt := rand.IntN(pokeInfo.Base_experience)

	fmt.Printf("Throwing a Pokeball at %v...\n", pokeName)

	if randomInt >= int(pokeInfo.Base_experience/2) {
		fmt.Printf("%v was caught!\n", pokeName)
		con.pokedex[pokeName] = pokeInfo
		return nil
	}

	fmt.Printf("%v escaped!\n", pokeName)

	return nil
}

func commandInspect() error {

	pokename := inspectName[1]
	if _, ok := con.pokedex[pokename]; !ok {
		fmt.Print("you have not caught that pokemon\n")
		return nil
	}

	fmt.Printf("Name: %v\n", pokename)
	fmt.Printf("Height: %v\n", con.pokedex[pokename].Height)
	fmt.Printf("Weight: %v\n", con.pokedex[pokename].Weight)
	fmt.Println("Stats: ")

	for i := 0; i < len(con.pokedex[pokename].Stats); i++ {
		fmt.Printf("   -%v: %v\n", con.pokedex[pokename].Stats[i].Stat.Name, con.pokedex[pokename].Stats[i].Base_stat)
	}

	fmt.Println("Types: ")

	for i := 0; i < len(con.pokedex[pokename].Types); i++ {
		fmt.Printf("   - %v\n", con.pokedex[pokename].Types[i].Type.Name)
	}

	return nil
}

func commandPokedex() error {

	fmt.Print("Your pokedex:\n")
	for key, _ := range con.pokedex {
		fmt.Printf("   - %v\n", key)
	}

	return nil
}
