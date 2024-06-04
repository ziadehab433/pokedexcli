package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")

		scanner.Scan()
		input := scanner.Text()

		err := evaluate(input)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func evaluate(input string) error {
	args := cleanInput(input)
	commands := getCommands()

	if !isValidCommand(args[0]) {
		return errors.New("command not found")
	}

	err := commands[args[0]].callback(args)
	if err != nil {
		return err
	}

	return nil
}

func cleanInput(input string) []string {
	lowered := strings.ToLower(input)
	return strings.Split(lowered, " ")
}
