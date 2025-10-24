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
	compPattern, err := regexp.Compile("^\\d+.?\\d* *[+*/\\-] *\\d+.?\\d*") // {beginExpr} Numbers[.]Numbers [+-/*] Numbers[.]Numbers
	if err != nil {
		log.Fatal(err)
	}
	res1 := compPattern.Match([]byte(expression))
	res2, err := regexp.Match("[a-zA-Z]", []byte(expression))                        // No letters
	pattern, err := regexp.Compile(".*\\/0([^.]|$|\\.(0{4,}.*|0{1,4}([^0-9]|$))).*") // https://stackoverflow.com/a/41122334
	if err != nil {
		log.Fatal("An unexpected error has occured. (", err, ")\n")
	}
	res3 := pattern.Match([]byte(expression))

	return res1 && !res2 && !res3
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

		output, er := expr.Eval(expres, nil) // Eval math
		if er != nil {
			log.Fatal("An unexpected error has occured. (", er, ")\n")
		}
		fmt.Print(output, "\n")
	}

}
