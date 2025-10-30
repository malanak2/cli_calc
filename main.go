package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
)

type node struct {
	parent *node
	right  *node
	left   *node
	value  string
}

func validateExpression(expression string) bool {
	compPattern, err := regexp.Compile(`^\d+.?\d* *[+*/\-] *\d+.?\d*`) // {beginExpr} Numbers[.]Numbers [+-/*] Numbers[.]Numbers
	if err != nil {
		log.Fatal(err)
	}
	res1 := compPattern.Match([]byte(expression))
	res2, err := regexp.Match("[a-zA-Z]", []byte(expression)) // No letters
	if err != nil {
		log.Fatal("An unexpected error has occured. (", err, ")\n")
	}
	pattern, err := regexp.Compile(`.*\/0([^.]|$|\.(0{4,}.*|0{1,4}([^0-9]|$))).*`) // https://stackoverflow.com/a/41122334
	if err != nil {
		log.Fatal("An unexpected error has occured. (", err, ")\n")
	}
	res3 := pattern.Match([]byte(expression))
	return res1 && !res2 && !res3
}
func parse_expression(expression string) [][]string {
	number := ""
	nodes := []node{}
	expression_separated := []string{}
	var root = node{value: ""}
	nodes = append(nodes, root)
	// Convert expression to an array of values (numbers, symbols(+-*/)). This cannot be done by just seperating by spaces, since they are not guaranteed
	for _, char := range expression {
		// Ignore spaces
		if char == ' ' {
			continue
		}
		// Including decimal nunmbers
		if char >= '0' && char <= '9' || char == '.' {
			number = number + string(char)
			continue
		}
		// Next char kinda has to be a symbol in list, so if not, invalid expression
		if !strings.Contains("/*-+\n", string(char)) {
			log.Fatal("An unexpected error has occured. ( Invalid character found: ", string(char), " )\n")
			return [][]string{}
		}
		expression_separated = append(expression_separated, number)
		number = ""
		if string(char) != "\n" {
			expression_separated = append(expression_separated, string(char))
		}
	}
	// An expression really shouldn't end in a symbol, so expect a number to be left over
	if number != "" {
		expression_separated = append(expression_separated, number)
		number = ""
	}
	for _, i := range expression_separated {
		fmt.Print(i, ", ")
	}
	fmt.Print("\n")
	// group
	groups := [][]string{}
	for index, i := range expression_separated {
		if !strings.Contains("/*-+", i) {
			continue
		}
		switch i {
		case "+", "-":
			// fínd next addition or subtraction
			index_next := index + 2
			for index_next < len(expression_separated) && strings.Contains("/*", expression_separated[index_next]) {
				index_next += 2
			}
			if index_next >= len(expression_separated) { // || strings.Contains("/*", expression_separated[index_next]) {
				index_next = index + 2
			}
			if index == 1 {
				if len(expression_separated) < 4 {
					// last_num symbol next_num
					groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index+1]})
				} else {
					// last_symbol symbol next_symbol where next_symbol is ideally next + or -, if there is none then it is the next multiplication / divison
					groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index_next]})
					continue
				}
				continue
			}
			// if index is second to last
			if len(expression_separated) <= index+2 {
				// 6 + 8 + 5 for example
				if strings.Contains("/*", expression_separated[index-2]) {
					// last_symbol symbol next_num
					groups = append(groups, []string{expression_separated[index-2], i, expression_separated[index+1]})
				} else {
					// last_num symbol next_num
					groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index+1]})
				}
				continue
			}
			if strings.Contains("/*", expression_separated[index-2]) {
				// last_symbol symbol next_symbol
				groups = append(groups, []string{expression_separated[index-2], i, expression_separated[index_next]})
				continue
			}
			// last_symbol symbol next_symbol where next_symbol is ideally next + or -, if there is none then it is the next multiplication / divison
			groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index_next]})
			continue

		case "*", "/":
			if index == 1 {
				groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index+1]})
				continue
			}
			if strings.Contains("/*", expression_separated[index-2]) {
				groups = append(groups, []string{expression_separated[index-2], i, expression_separated[index+1]})
				continue
			}
			groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index+1]})
		}
	}
	for _, i := range groups {
		fmt.Print(i, ", ")
	}
	fmt.Print("\n")
	return groups
}
func calculate_expression(parsed_expression [][]string) string {
	results := []string{}
	fmt.Print("Result init\n")
	for range parsed_expression {
		results = append(results, "")
	}
	for slices.Contains(results, "") {
		for index, expression := range parsed_expression {
			if results[index] != "" {
				continue
			}
			fmt.Print("result ", index, " is empty (", expression[0], ", ", expression[2], ")\n")
			// TODO: Calculate expressions that are complex
			a, a_err := strconv.ParseFloat(expression[0], 64)
			b, b_err := strconv.ParseFloat(expression[2], 64)

			if a_err == nil && b_err == nil {
				fmt.Print("Calculating result ", index, "\n")
				switch expression[1] {
				case "+":
					results[index] = strconv.FormatFloat(a+b, 'f', -1, 64)
					continue
				case "-":
					results[index] = strconv.FormatFloat(a-b, 'f', -1, 64)
					continue
				case "*":
					results[index] = strconv.FormatFloat(a*b, 'f', -1, 64)
					continue
				case "/":
					results[index] = strconv.FormatFloat(a/b, 'f', -1, 64)
					continue
				}
			} else {
				if a_err != nil {
					// TODO: Get previous group of this symbol and check for value
				}
				if b_err != nil {
					// TODO: Get next group of this symbol and check for value
				}
				// TODO: If both values exist, calculate result
				fmt.Print("A: ", a_err.Error(), ", B: ", b_err.Error(), "\n")
			}
		}
	}
	return results[0]
}
func main() {
	//for {
	fmt.Print("Please enter a mathematical expression or q to quit:")
	in := bufio.NewReader(os.Stdin)
	expres, err := in.ReadString('\n')
	// expres := "6 + 8 * 5 + 7"
	if err != nil {
		log.Fatal(err)
	}
	if expres == "q" {
		return
	}
	if !validateExpression(expres) {
		fmt.Print("Invalid expression. Please try again\n")
		return
		// continue
	}

	output, er := expr.Eval(expres, nil) // Eval math
	if er != nil {
		log.Fatal("An unexpected error has occured. (", er, ")\n")
	}
	parsed_expression := parse_expression(expres)
	fmt.Print("Result: ", calculate_expression(parsed_expression), "\n")
	fmt.Print("Verification output from lib:", output, "\n")
	//}

}
