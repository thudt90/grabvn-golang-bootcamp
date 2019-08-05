/*
Write a simple REPL program (Read–Eval–Print–Loop).
Take a formula from a user then print out the result. The formula must be in this format:
<first number> <arithmetic: + - * / > <second number>
Example:
> 1 + 2
1 + 2 = 3
> 2 * 10
2 * 10 = 20
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	oki        = iota
	noOperator = iota
	wrongInput = iota
)

func parseString(text string) error {
	textSplit := strings.Split(text, " ")
	//fmt.Println(textSplit[0])
	leftVal, err := strconv.ParseFloat(textSplit[0], 10)
	// fmt.Println(leftVal)
	if err != nil {
		return err
	}
	rightVal, err := strconv.ParseFloat(textSplit[2], 10)
	// fmt.Println(rightVal)
	if err != nil {
		return err
	}

	switch textSplit[1] {
	case "+":
		fmt.Println(leftVal, textSplit[1], rightVal, "=", leftVal+rightVal)
	case "-":
		fmt.Println(leftVal, textSplit[1], rightVal, "= ", leftVal-rightVal)
	case "*":
		fmt.Println(leftVal, textSplit[1], rightVal, "=", leftVal*rightVal)
	case "/":
		if rightVal != 0 {
			fmt.Println(leftVal, textSplit[1], rightVal, "= ", leftVal/rightVal)
		} else {
			return errors.New("Couldn't divide by zero")
		}
	default:
		return errors.New("Invalid operator")
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">")

	for scanner.Scan() {
		text := scanner.Text()
		// fmt.Println(text)
		err := parseString(text)

		if err != nil {
			fmt.Println("ERROR")
			// continue
		}
		fmt.Print(">")
	}
}
