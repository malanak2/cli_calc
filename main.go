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
	compPattern, err := regexp.Compile("\\d*.*\\d* *[+*/\\-] *\\d*.*\\d*")
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
		fmt.Print(output, "\n")
	}

}
