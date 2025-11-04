package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
)

type result_t struct {
	value        float64
	isCalculated bool
}
type Symbol int

const (
	Addition Symbol = iota
	Subtraction
	Multiplication
	Division
	None
)

var symbolName = map[Symbol]string{
	Addition:       "+",
	Subtraction:    "-",
	Multiplication: "*",
	Division:       "/",
	None:           "no_symbol",
}
var symbolTo = map[string]Symbol{
	"+":         Addition,
	"-":         Subtraction,
	"*":         Multiplication,
	"/":         Division,
	"no_symbol": None,
}

func (s Symbol) String() string {
	return symbolName[s]
}

type value_placeholder_t struct {
	value  float64
	symbol Symbol
}

type expression_t struct {
	valA        value_placeholder_t
	valB        value_placeholder_t
	symbol_used Symbol
}

func (exp expression_t) String() string {
	res := ""
	if exp.valA.symbol == None {
		res = res + strconv.FormatFloat(exp.valA.value, 'f', -1, 64)
	} else {
		res = res + symbolName[exp.valA.symbol]
	}
	res = res + symbolName[exp.symbol_used]
	if exp.valB.symbol == None {
		res = res + strconv.FormatFloat(exp.valB.value, 'f', -1, 64)
	} else {
		res = res + symbolName[exp.valB.symbol]
	}
	return res
}

func (exp expression_t) Init(valA string, valB string, char rune) (expression_t, error) {
	switch valA {
	case "+", "-", "*", "/":
		debug("valA", valA, " ", string(char), " ", valB, " is in case sign\n")
		sym := symbolTo[valA]
		exp.valA = value_placeholder_t{value: 0, symbol: sym}
	default:
		debug("valA", valA, " ", string(char), " ", valB, " is in case default\n")
		numA, err := strconv.ParseFloat(valA, 64)
		if err != nil {
			return expression_t{}, errors.New("Invalid number: " + valA)
		}
		exp.valA = value_placeholder_t{value: numA, symbol: None}
	}
	switch valB {
	case "+", "-", "*", "/":
		debug("valB", valA, " ", string(char), " ", valB, " is in case sign\n")
		sym := symbolTo[valB]
		exp.valB = value_placeholder_t{value: 0, symbol: sym}
	default:
		debug("valB", valA, " ", string(char), " ", valB, " is in case default\n")
		numB, err := strconv.ParseFloat(valB, 64)
		if err != nil {
			return expression_t{}, errors.New("Invalid number: " + valB)
		}
		exp.valB = value_placeholder_t{value: numB, symbol: None}
	}
	exp.symbol_used = symbolTo[string(char)]
	return exp, nil

}

var isDebug = false

func debug(a ...any) {
	if isDebug {
		fmt.Print(a...)
	}
}

// parse_expression function    Parse expression contained in a string into and [][]string, which contains grouped up single symbol expressions, resembling a tree
func parse_expression(expression string) ([]expression_t, error) {
	number := ""
	expression_separated := []string{}
	nextCanBeNumSymbol := true
	// Convert expression to an array of values (numbers, symbols(+-*/)). This cannot be done by just seperating by spaces, since they are not guaranteed
	for index, char := range expression {
		// Ignore spaces
		if char == ' ' {
			continue
		}
		debug(index, ": ", string(char), "NumSymbol: ", nextCanBeNumSymbol, "\n")
		if strings.Contains("/*-+", string(char)) && nextCanBeNumSymbol {
			number = number + string(char)
			nextCanBeNumSymbol = false
			continue
		} else if strings.Contains("/*-+", string(char)) {
			nextCanBeNumSymbol = true
			expression_separated = append(expression_separated, number)
			number = ""
			// Dont want to add newline or the windows bs to our expression, now do we
			if !strings.Contains("\n\r", string(char)) {
				expression_separated = append(expression_separated, string(char))
			}
			continue
		}
		// Including decimal nunmbers
		if char >= '0' && char <= '9' || char == '.' {
			number = number + string(char)
			nextCanBeNumSymbol = false
			continue
		}
		// Next char kinda has to be a symbol in list, so if not, invalid expression
		if !strings.Contains("/*-+\n\r", string(char)) {
			return []expression_t{}, errors.New("Invalid character: " + string(char))
		}
		// Add last numer
		expression_separated = append(expression_separated, number)
		number = ""
		// Dont want to add newline or the windows bs to our expression, now do we
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
	groups := []expression_t{}
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
				index_next = len(expression_separated) - 2
			}
			if index == 1 {
				if len(expression_separated) < 4 {
					// last_num symbol next_num
					exp, err := expression_t{}.Init(expression_separated[index-1], expression_separated[index+1], rune(i[0]))
					if err != nil {
						return []expression_t{}, err
					}
					groups = append(groups, exp)
				} else {
					// last_symbol symbol next_symbol where next_symbol is ideally next + or -, if there is none then it is the next multiplication / divison

					exp, err := expression_t{}.Init(expression_separated[index-1], expression_separated[index_next], rune(i[0]))
					if err != nil {
						return []expression_t{}, err
					}
					groups = append(groups, exp)
					continue
				}
				continue
			}
			// if index is second to last
			if len(expression_separated) <= index+2 {
				// 6 + 8 + 5 for example
				if strings.Contains("/*", expression_separated[index-2]) {
					// last_symbol symbol next_num

					exp, err := expression_t{}.Init(expression_separated[index-2], expression_separated[index+1], rune(i[0]))
					if err != nil {
						return []expression_t{}, err
					}
					groups = append(groups, exp)
				} else {
					// last_num symbol next_num

					exp, err := expression_t{}.Init(expression_separated[index-1], expression_separated[index+1], rune(i[0]))
					if err != nil {
						return []expression_t{}, err
					}
					groups = append(groups, exp)

				}
				continue
			}
			if strings.Contains("/*", expression_separated[index-2]) {
				// last_symbol symbol next_symbol

				exp, err := expression_t{}.Init(expression_separated[index-2], expression_separated[index_next], rune(i[0]))
				if err != nil {
					return []expression_t{}, err
				}
				groups = append(groups, exp)
				continue
			}
			// last_symbol symbol next_symbol where next_symbol is ideally next + or -, if there is none then it is the next multiplication / divison

			exp, err := expression_t{}.Init(expression_separated[index-1], expression_separated[index_next], rune(i[0]))
			if err != nil {
				return []expression_t{}, err
			}
			groups = append(groups, exp)
			continue

		case "*", "/":
			if index == 1 {
				// last_num symbol next_num

				exp, err := expression_t{}.Init(expression_separated[index-1], expression_separated[index+1], rune(i[0]))
				if err != nil {
					return []expression_t{}, err
				}
				groups = append(groups, exp)
				continue
			}
			if strings.Contains("/*", expression_separated[index-2]) {
				// last_symbol symbol next_num

				exp, err := expression_t{}.Init(expression_separated[index-2], expression_separated[index+1], rune(i[0]))
				if err != nil {
					return []expression_t{}, err
				}
				groups = append(groups, exp)

				continue
			}
			// last_num symbol next_num

			exp, err := expression_t{}.Init(expression_separated[index-1], expression_separated[index+1], rune(i[0]))
			if err != nil {
				return []expression_t{}, err
			}
			groups = append(groups, exp)

		}
	}
	for _, i := range groups {
		debug(i.String(), ", ")
	}
	debug("\n")
	return groups, nil
}
func (exp expression_t) calculate_group() (float64, error) {
	switch exp.symbol_used {
	case Addition:
		debug("Calculating result: ", exp.valA.value, "+", exp.valB.value, "\n")
		return exp.valA.value + exp.valB.value, nil
	case Subtraction:
		debug("Calculating result: ", exp.valA.value, "-", exp.valB.value, "\n")
		return exp.valA.value - exp.valB.value, nil
	case Multiplication:
		debug("Calculating result: ", exp.valA.value, "*", exp.valB.value, "\n")
		return exp.valA.value * exp.valB.value, nil
	case Division:
		debug("Calculating result: ", exp.valA.value, "/", exp.valB.value, "\n")
		if exp.valB.value == 0 {
			// Crashes the program, what can you do :)
			return 0, errors.New("cannot divide by zero. (" + strconv.FormatFloat(exp.valA.value, 'f', -1, 64) + symbolName[exp.symbol_used] + strconv.FormatFloat(exp.valB.value, 'f', -1, 64) + ")")
		}
		return exp.valA.value / exp.valB.value, nil
	}
	return 0, errors.New("An unexpected error has occured in calculate_group function. (" + exp.valA.symbol.String() + symbolName[exp.symbol_used] + exp.valB.symbol.String() + ")")
}

// calculate_expression function    Calculates the result of an expression from parsed expression
func calculate_expression(parsed_expression []expression_t) (float64, error) {
	results := []result_t{}
	debug("Result init\n")
	for range parsed_expression {
		results = append(results, result_t{isCalculated: false})
	}
	lastRes := result_t{}
	// Until uncalculated results remain or until broken from
	for slices.Contains(results, result_t{}) {
		for index, expression := range parsed_expression {
			// No need to recalculate
			if results[index].isCalculated {
				continue
			}
			debug("result ", index, " is empty\n")
			// See if numbers are in both places
			if expression.valA.symbol == None && expression.valB.symbol == None {
				result, err := expression.calculate_group()
				if err != nil {
					return 0, err
				}
				results[index] = result_t{value: result, isCalculated: true}

				lastRes = results[index]
			} else {
				a, b := result_t{isCalculated: false}, result_t{isCalculated: false}
				if expression.valA.symbol != None {
					// Find the closest occurance of the symbol to the left
					for i := index; i >= 0; i-- {
						if parsed_expression[i].symbol_used == expression.valA.symbol {
							// Only if it is calculated
							if results[i].isCalculated {
								debug("Found expression a with result ", results[i].value, "\n")
								a = results[i]
								break
							}
						}
					}
				} else {
					a = result_t{value: expression.valA.value, isCalculated: true}
				}
				if expression.valB.symbol != None {
					// Find the closest occurance of the symbol to the right
					for i := index; i < len(parsed_expression); i++ {
						if parsed_expression[i].symbol_used == expression.valB.symbol {
							// Only if it is calculated
							if results[i].isCalculated {
								debug("Found expression b with result ", results[i].value, "\n")
								b = results[i]
								break
							}
						}
					}
				} else {
					b = result_t{value: expression.valB.value, isCalculated: true}
				}
				if (a.isCalculated) && (b.isCalculated) {
					group := expression_t{value_placeholder_t{a.value, None}, value_placeholder_t{b.value, None}, expression.symbol_used}
					result, err := group.calculate_group()
					if err != nil {
						return 0, err
					}
					results[index] = result_t{result, true}
					lastRes = results[index]
				}
			}
		}
	}
	return lastRes.value, nil
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

		parsed_expression, err := parse_expression(expres)
		if err != nil {
			fmt.Print(err.Error(), "\n")
			continue
		}

		result, err := calculate_expression(parsed_expression)
		if err != nil {
			fmt.Print(err.Error(), "\n")
			continue
		}
		fmt.Print("Result: ", result, "\n")
		if isDebug {

			output, er := expr.Eval(expres, nil) // Eval math
			if er != nil {
				log.Fatal("An unexpected error has occured. (", er, ")\n")
			}
			debug("Verification output from lib:", output, "\n")
		}
	}
}
