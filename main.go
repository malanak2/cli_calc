package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/expr-lang/expr"
)

func validateExpression(expression string) bool {
	compPattern, err := regexp.Compile("\\d*.?\\d* *[+*/\\-] *\\d*.?\\d*") // Numbers[.]Numbers
	if err != nil {
		log.Fatal(err)
	}
	return compPattern.Match([]byte(expression))
}

func main() {
	for true {
		fmt.Print("Please enter a mathematical expression:")
		in := bufio.NewReader(os.Stdin)
		expres, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if !validateExpression(expres) {
			fmt.Print("Invalid expression. Please try again\n")
			continue
		}
		output, er := expr.Eval(expres, nil)

		if er != nil {
			log.Fatal("An unexpected error has occured. (", er, ")\n")
		}
		pattern, err := regexp.Compile(".*\\/0([^.]|$|\\.(0{4,}.*|0{1,4}([^0-9]|$))).*") // https://stackoverflow.com/a/41122334
		if er != nil {
			log.Fatal("An unexpected error has occured. (", er, ")\n")
		}
		out := pattern.Match([]byte(expres))
		if out {
			fmt.Print("Cannot divide by zero\n")
			continue
		}
		fmt.Print(output, "\n")
	}

}
