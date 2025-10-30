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

var isDebug = false

func debug(a ...any) {
	if isDebug {
		fmt.Print(a...)
	}
}
func validateExpression(expression string) bool {
	containsExpressionPattern, err := regexp.Compile(`^\d+.?\d* *[+*/\-] *\d+.?\d*`) // {beginExpr} Numbers[.]Numbers [+-/*] Numbers[.]Numbers
	if err != nil {
		log.Fatal(err)
	}
	containsExpressionResult := containsExpressionPattern.Match([]byte(expression))
	containsLettersResult, err := regexp.Match("[a-zA-Z]", []byte(expression)) // No letters
	if err != nil {
		log.Fatal("An unexpected error has occured. (", err, ")\n")
	}
	pattern, err := regexp.Compile(`.*\/0([^.]|$|\.(0{4,}.*|0{1,4}([^0-9]|$))).*`) // https://stackoverflow.com/a/41122334
	if err != nil {
		log.Fatal("An unexpected error has occured. (", err, ")\n")
	}
	containsDivisionByZeroResult := pattern.Match([]byte(expression))
	return containsExpressionResult && !containsLettersResult && !containsDivisionByZeroResult
}

// parse_expression function    Parse expression contained in a string into and [][]string, which contains grouped up single symbol expressions, resembling a tree
func parse_expression(expression string) [][]string {
	number := ""
	expression_separated := []string{}
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
		if !strings.Contains("/*-+\n\r", string(char)) {
			log.Fatal("An unexpected error has occured. ( Invalid character found: ", string(char), " )\n")
			return [][]string{}
		}
		// Dont want to add newline or the windows bs to our expression, now do we
		expression_separated = append(expression_separated, number)
		number = ""
		if !strings.Contains("\n\r", string(char)) {
			expression_separated = append(expression_separated, string(char))
		}
	}
	// An expression really shouldn't end in a symbol, so expect a number to be left over
	// if number != "" {
	// 	expression_separated = append(expression_separated, number)
	// 	number = ""
	// }
	//
	// Windows is stupid, remove empty elements
	for _, i := range expression_separated {
		if i == "" {
			expression_separated = expression_separated[:len(expression_separated)-1]
		}
		debug(i, ", ")
	}
	debug("\n")
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
				// last_num symbol next_num
				groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index+1]})
				continue
			}
			if strings.Contains("/*", expression_separated[index-2]) {
				// last_symbol symbol next_num
				groups = append(groups, []string{expression_separated[index-2], i, expression_separated[index+1]})
				continue
			}
			// last_num symbol next_num
			groups = append(groups, []string{expression_separated[index-1], i, expression_separated[index+1]})
		}
	}
	for _, i := range groups {
		debug(i, ", ")
	}
	debug("\n")
	return groups
}
func calculate_group(a float64, b float64, char string) string {
	switch char {
	case "+":
		debug("Calculating result: ", a, "+", b, "\n")
		return strconv.FormatFloat(a+b, 'f', -1, 64)
	case "-":
		debug("Calculating result: ", a, "-", b, "\n")
		return strconv.FormatFloat(a-b, 'f', -1, 64)
	case "*":
		debug("Calculating result: ", a, "*", b, "\n")
		return strconv.FormatFloat(a*b, 'f', -1, 64)
	case "/":
		debug("Calculating result: ", a, "/", b, "\n")
		if b == 0 {
			// Crashes the program, what can you do :)
			log.Fatal("Cannot divide by zero")
		}
		return strconv.FormatFloat(a/b, 'f', -1, 64)
	}
	return ""
}

// calculate_expression function    Calculates the result of an expression from parsed expression
func calculate_expression(parsed_expression [][]string) string {
	results := []string{}
	debug("Result init\n")
	for range parsed_expression {
		results = append(results, "")
	}
	lastRes := ""
	// Until uncalculated results remain or until broken from
	for slices.Contains(results, "") {
		for index, expression := range parsed_expression {
			// No need to recalculate
			if results[index] != "" {
				continue
			}
			debug("result ", index, " is empty (", expression[0], ", ", expression[2], ")\n")
			// See if numbers are in both places
			a, a_err := strconv.ParseFloat(expression[0], 64)
			b, b_err := strconv.ParseFloat(expression[2], 64)

			if a_err == nil && b_err == nil {
				results[index] = calculate_group(a, b, expression[1])

				lastRes = results[index]
			} else {
				a, b := "", ""
				if a_err != nil {
					// Find the closest occurance of the symbol to the left
					for i := index; i >= 0; i-- {
						if parsed_expression[i][1] == expression[0] {
							// Only if it is calculated
							if results[i] != "" {
								debug("Found expression a with result ", results[i], "\n")
								a = results[i]
								break
							}
						}
					}
				} else {
					a = expression[0]
				}
				if b_err != nil {
					// Find the closest occurance of the symbol to the right
					for i := index; i < len(parsed_expression); i++ {
						if parsed_expression[i][1] == expression[2] {
							// Only if it is calculated
							if results[i] != "" {
								debug("Found expression b with result ", results[i], "\n")
								b = results[i]
								break
							}
						}
					}
				} else {
					b = expression[2]
				}
				if a != "" && b != "" {
					resA, _ := strconv.ParseFloat(a, 64)
					resB, _ := strconv.ParseFloat(b, 64)
					results[index] = calculate_group(resA, resB, expression[1])
					lastRes = results[index]
				}
			}
		}
	}
	return lastRes
}
func main() {
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		isDebug = true
		debug("Launching in debug mode...\n")
	}
	for {
		fmt.Print("Please enter a mathematical expression or q to quit:")
		in := bufio.NewReader(os.Stdin)
		expres, err := in.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(expres) == "q" {
			fmt.Print("Bye!")
			return
		}
		if !validateExpression(expres) {
			fmt.Print("Invalid expression. Please try again\n")
			continue
		}

		output, er := expr.Eval(expres, nil) // Eval math
		if er != nil {
			log.Fatal("An unexpected error has occured. (", er, ")\n")
		}
		parsed_expression := parse_expression(expres)
		fmt.Print("Result: ", calculate_expression(parsed_expression), "\n")
		debug("Verification output from lib:", output, "\n")
	}
}
