package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/ziadehab433/pokedexcli/internal/pokeAPI"
)

type cliCommand struct {
	name     string
	desc     string
	callback func([]string) error
}

type config struct {
	next string
	prev string
}

var conf config = config{
	next: "https://pokeapi.co/api/v2/location-area/",
	prev: "",
}

var pokemons map[string]pokeAPI.Pokemon = make(map[string]pokeAPI.Pokemon)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:     "help",
			desc:     "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name:     "exit",
			desc:     "Exit the pokedex",
			callback: commandExit,
		},
		"map": {
			name:     "map",
			desc:     "Displays the names of 20 location areas in the Pokemon world",
			callback: commandMap,
		},
		"mapb": {
			name:     "mapb",
			desc:     "Displays the names of previous 20 location areas in the Pokemon world",
			callback: commandMapb,
		},
		"explore": {
			name:     "explore",
			desc:     "Let's you explore a certain location",
			callback: commandExplore,
		},
		"catch": {
			name:     "catch",
			desc:     "Let's your catch a pokemon so you can add it to your pokedex",
			callback: commandCatch,
		},
		"inspect": {
			name:     "inspect",
			desc:     "Displays information about an already caught Pokemon",
			callback: commandInspect,
		},
		"pokedex": {
			name:     "pokedex",
			desc:     "Displays the pokemons you've catched so far",
			callback: commandPokedex,
		},
	}
}

func isValidCommand(str string) bool {
	commands := getCommands()
	for name := range commands {
		if name == str {
			return true
		}
	}

	return false
}

func commandHelp(args []string) error {
	if len(args) > 1 {
		return errors.New("the help command can only be used without any subcommands")
	}

	fmt.Print(
		"\n",
		"Welcome to the Pokedex! \n",
		"Usage: \n",
		"\n",
	)

	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%s: %s \n", command.name, command.desc)
	}

	fmt.Println()
	return nil
}

func commandExit(args []string) error {
	os.Exit(0)
	return nil
}

func commandMap(args []string) error {
	if len(args) > 1 {
		return errors.New("the map command can only be used without any subcommands")
	}

	pokeData, err := pokeAPI.GetLocations(conf.next)
	if err != nil {
		return err
	}

	for _, l := range pokeData.Results {
		fmt.Println(l.Name)
	}

	if pokeData.Previous == nil {
		updateConfig(pokeData.Next, "")
	} else {
		updateConfig(pokeData.Next, pokeData.Previous.(string))
	}

	return nil
}

func commandMapb(args []string) error {
	if len(args) > 1 {
		return errors.New("the mapb command can only be used without any subcommands")
	}

	if conf.prev == "" {
		return errors.New("cannot go back, already on the first page")
	}

	pokeData, err := pokeAPI.GetLocations(conf.prev)
	if err != nil {
		return err
	}

	for _, l := range pokeData.Results {
		fmt.Println(l.Name)
	}

	if pokeData.Previous == nil {
		updateConfig(pokeData.Next, "")
	} else {
		updateConfig(pokeData.Next, pokeData.Previous.(string))
	}
	fmt.Println(conf)

	return nil
}

func commandExplore(args []string) error {
	if len(args) != 2 {
		return errors.New("too many or too few commands")
	}

	sub := args[1]

	url := "https://pokeapi.co/api/v2/location-area/" + sub
	locationData, err := pokeAPI.ExploreLocation(url)
	if err != nil {
		return err
	}

	pokemons := locationData.PokemonEncounters

	fmt.Println("Exploring", sub, "...")
	fmt.Println("Found Pokemon:")
	for _, res := range pokemons {
		fmt.Println("   -", res.Pokemon.Name)
	}

	return nil
}

func commandCatch(args []string) error {
	if len(args) != 2 {
		return errors.New("too many or too few commands")
	}

	name := args[1]

	if isCaught(name) {
		fmt.Printf("You already caught %s", name)
		return nil
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + name
	pokemonData, err := pokeAPI.GetPokemon(url)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a PokeBall at %s... \n", name)

	baseExp := pokemonData.BaseExperience
	r := rand.IntN(200)
	if r < baseExp {
		fmt.Printf("%s escaped! \n", name)
		return nil
	}

	fmt.Printf("%s was caught! \n", name)
	pokemons[name] = pokemonData
	return nil
}

func commandInspect(args []string) error {
	if len(args) != 2 {
		return errors.New("too many or too few commands")
	}

	name := args[1]

	if !isCaught(name) {
		fmt.Printf("You haven't caught that pokemon yet... \n")
		return nil
	}

	pokemon := pokemons[name]
	printPokemonDetails(pokemon)
	return nil
}

func commandPokedex(args []string) error {
	if len(args) < 1 {
		return errors.New("too many arguments to function")
	}

	if len(pokemons) == 0 {
		fmt.Println("You haven't caught any pokemons yet...")
		return nil
	}

	fmt.Println("Your Pokedex: ")
	for n := range pokemons {
		fmt.Println("  -", n)
	}
	return nil
}

func printPokemonDetails(pokemon pokeAPI.Pokemon) {
	fmt.Printf("Name: %s \n", pokemon.Name)
	fmt.Printf("Height: %d \n", pokemon.Height)
	fmt.Printf("Weight: %d \n", pokemon.Weight)

	fmt.Println("Stats: ")
	for _, v := range pokemon.Stats {
		fmt.Printf("  - %s: %d \n", v.Stat.Name, v.BaseStat)
	}

	fmt.Println("Types: ")
	for _, v := range pokemon.Types {
		fmt.Printf("  - %s\n", v.Type.Name)
	}
}

func isCaught(name string) bool {
	for n := range pokemons {
		if n == name {
			return true
		}
	}
	return false
}

func updateConfig(next, prev string) {
	conf.next = next
	conf.prev = prev
}
