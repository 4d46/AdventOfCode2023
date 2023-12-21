package main

import (
	"fmt"
	"maps"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	LT = iota
	GT
)

type ruleSet struct {
	workflows map[string]workflow
}

type workflow struct {
	rules []rule
}

type rule struct {
	end      bool
	operand  string
	operator int
	value    int
	target   string
}

type part struct {
	x, m, a, s int
}

type partLimits struct {
	min   map[string]int
	max   map[string]int
	valid bool
}

const example1Str = `px{a<2006:qkq,m>2090:A,rfg}
pv{a>1716:R,A}
lnx{m>1548:A,A}
rfg{s<537:gd,x>2440:R,A}
qs{s>3448:A,lnx}
qkq{x<1416:A,crn}
crn{x>2662:A,R}
in{s<1351:px,qqz}
qqz{s>2770:qs,m<1801:hdj,R}
gd{a>3333:R,R}
hdj{m>838:A,pv}

{x=787,m=2655,a=1222,s=2876}
{x=1679,m=44,a=2067,s=496}
{x=2036,m=264,a=79,s=2244}
{x=2461,m=1339,a=466,s=291}
{x=2127,m=1623,a=2188,s=1013}`

func main() {
	fmt.Printf("Advent of Code 2023 - Day %2d\n", 19)

	// Load input
	// input := example1Str
	input := loadFileContents("workflowAndParts.txt")

	// Parse input
	ruleSet, parts := parseInput(input)

	// fmt.Println(rules)
	// fmt.Println(parts)

	// Loop over parts processing each one
	total := 0
	for _, part := range parts {
		total += processPart(part, &ruleSet, "in")
	}

	fmt.Printf("Part 1 Total: %d\n", total)

	// Part 2
	// Loop over rules calculating what sets of part values will pass
	// acceptableParts := generaliseRules(ruleSet, "crn")
	// totalParts := walkRuleSet(ruleSet, "crn", makeMaxPartLimits(), 0,true)
	// totalParts := walkRuleSet(ruleSet, "qkq", makeMaxPartLimits(), 0, true)
	totalParts := walkRuleSet(ruleSet, "in", makeMaxPartLimits(), 0, false)

	fmt.Printf("Part 2 Total parts: %d\n", totalParts)
}

// Walk through ruleSet like a tree, amending partlimits as we go.  When we hit the end of a branch
// calculate the number of valid combinations of part values and return that back up
func walkRuleSet(rules ruleSet, wflowName string, branchLimits partLimits, stack int, debug bool) int {
	var totalParts int

	// Confirm that the branch limits are valid and return 0 if not because nothing will come of this branch
	if !branchLimits.valid {
		return 0
	}
	if debug {
		fmt.Printf("nm:%s%s\n", strings.Repeat("\t", stack), wflowName)
	}

	// Check in case we are already at the end
	switch wflowName {
	case "A":
		// This is an accept rule, calculate the number of valid combinations of part values
		// fmt.Printf("Accept rule: %s\n", ruleName)
		// spew.Dump(branchLoopLimits)
		result := calculateValidPartCombinations(branchLimits)
		if debug {
			fmt.Print("A |")
			printPartLimits(branchLimits, stack)
			fmt.Printf("%sAccept rule: %s, result: %d\n", strings.Repeat("\t", stack), wflowName, result)
		}
		return result
	case "R":
		// This is a reject rule, nothing from this path is valid
		if debug {
			fmt.Print("R |")
			printPartLimits(branchLimits, stack)
			fmt.Printf("%sReject rule: %s\n", strings.Repeat("\t", stack), wflowName)
		}
		return 0
	}

	// Fetch workflow
	wflow := rules.workflows[wflowName]

	// Check workflow has rules
	if len(wflow.rules) == 0 {
		panic("Workflow has no rules")
	}
	if debug {
		fmt.Print("in|")
		printPartLimits(branchLimits, stack)
		fmt.Println()
	}

	limits := branchLimits
	// Loop over rules in workflow
	for _, rule := range wflow.rules {
		if debug {
			printRule(rule, stack)
			fmt.Print("br|")
			printPartLimits(limits, stack)
		}
		// Check if this is an end rule
		if rule.end {
			// This is an end rule
			// Check if this is an accept rule
			switch rule.target {
			case "A":
				// This is an accept rule
				totalParts += calculateValidPartCombinations(limits)
			case "R":
				// This is a reject rule, nothing from this path is valid, ignore
				// return []partLimits{}
			default:
				// Process next named rule and add any valid combinations to the total
				totalParts += walkRuleSet(rules, rule.target, limits, stack+1, debug)
			}
		} else {
			// This is a normal rule
			// Calculate the new split limits
			trueLimits, falseLimits := splitPartLimits(limits, rule)
			if debug {
				fmt.Print("aT|")
				printPartLimits(trueLimits, stack)
				fmt.Print("aF|")
				printPartLimits(falseLimits, stack)
			}
			// Start a new branch for the true limits
			totalParts += walkRuleSet(rules, rule.target, trueLimits, stack+1, debug)
			// Continue this branch with the false limits
			limits = falseLimits
		}
	}
	return totalParts
}

// Print rule
func printRule(rule rule, indent int) {
	fmt.Print(strings.Repeat("\t", indent))
	if rule.end {
		fmt.Printf("End:%s\n", rule.target)
	} else {
		fmt.Printf("%s%s%d:%s\n", rule.operand, decodeOperator(rule.operator), rule.value, rule.target)
	}
}

// Print the part limits
func printPartLimits(limits partLimits, indent int) {
	fmt.Print(strings.Repeat("\t", indent))
	// Loop over each operand
	for _, operand := range []string{"x", "m", "a", "s"} {
		fmt.Printf("%s: %4d➜%4d ", operand, limits.min[operand], limits.max[operand])
	}
	fmt.Println()
}

// Calculate the number of valid combinations of part values
func calculateValidPartCombinations(limits partLimits) int {
	var total int

	// Check if the limits are valid, return 0 if not
	if !limits.valid {
		return 0
	}

	total = limits.max["x"] - limits.min["x"] + 1

	// Calculate the number of valid combinations
	for _, operand := range []string{"m", "a", "s"} {
		total *= limits.max[operand] - limits.min[operand] + 1
	}

	return total
}

// func walkRuleSet(rules ruleSet, ruleName string, branchLimits []partLimits) int {
// 	var resultingLimits []partLimits
// 	// Fetch workflow
// 	wflow := rules.workflows[ruleName]
// 	fmt.Printf("ruleName: %s\n", ruleName)
// 	var branchLoopLimits []partLimits
// 	branchLoopLimits = append(branchLoopLimits, branchLimits...)
// 	// Check in case we are already at the end
// 	switch ruleName {
// 	case "A":
// 		// This is an accept rule
// 		// fmt.Printf("Accept rule: %s\n", ruleName)
// 		// spew.Dump(branchLoopLimits)
// 		return branchLimits
// 	case "R":
// 		// This is a reject rule, nothing from this path is valid
// 		return []partLimits{}
// 	}

// 	// Loop over rules in workflow
// 	for _, rule := range wflow.rules {
// 		// Check if this is an end rule
// 		if rule.end {
// 			// This is an end rule
// 			// Check if this is an accept rule
// 			switch rule.target {
// 			case "A":
// 				// This is an accept rule
// 				fmt.Printf("Accept rule: %s\n", ruleName)
// 				// spew.Dump(branchLoopLimits)
// 				return branchLimits
// 			case "R":
// 				// This is a reject rule, nothing from this path is valid, ignore
// 				// return []partLimits{}
// 			default:
// 				// Process next named rule
// 				return walkRuleSet(rules, rule.target, branchLimits)
// 			}
// 		} else {

// 			// This is a normal rule
// 			// Apply rule to part limits to get new limits and call target with result
// 			branchLoopLimits = applyRuleToPartLimits(rule, branchLoopLimits)
// 			// spew.Dump(branchLoopLimits)
// 			resultingLimits = append(resultingLimits, walkRuleSet(rules, rule.target, branchLoopLimits)...)
// 			// spew.Dump(resultingLimits)
// 		}
// 	}
// 	return resultingLimits
// }

// Generalise the rules to find the acceptable ranges of part values
// func generaliseRules(rules ruleSet, ruleName string) map[string]partLimits {
// 	// Initialise map of acceptable part values
// 	acceptableParts := make(map[string]partLimits)

// 	// Loop over rules
// 	for ruleName, wflow := range rules.workflows {
// 		// Loop over rules in workflow
// 		for _, rule := range wflow.rules {
// 			// Check if this is an end rule
// 			if rule.end {
// 				// This is an end rule
// 				// Check if this is an accept rule
// 				if rule.target == "A" {
// 					// This is an accept rule
// 					// Get the part limits for this rule
// 					limits := getPartLimits(rule)
// 					// fmt.Println(limits)

// 					// Add the limits to the acceptable parts map
// 					acceptableParts[ruleName] = limits
// 				}
// 			}
// 		}
// 	}

//		return acceptableParts
//	}
// func generaliseWorkflow(rules ruleSet, workflowName string) []partLimits {
// 	// Initialise map of acceptable part values
// 	var acceptableParts []partLimits

// 	var partLimits []partLimits
// 	partLimits = append(partLimits, makeMaxPartLimits())

// 	// Loop over rules in workflow
// 	for _, rule := range rules.workflows[workflowName].rules {
//                 if rule.end {
//                         // This is an end rule
//                         // Check if this is an accept rule
//                         switch rule.target {
//                         case "A":
//                                 // This is an accept rule, start with everything and reduce limits as we go
//                                  []partLimits{makeMaxPartLimits()}
//                         case "R":
//                                 // This is a reject rule, ignore
//                         default:
//                                 // Process next named rule
//                                 return generaliseWorkflow(rules, rule.target)
//                         }
//                 }
// 		acceptableParts = applyRuleToPartLimits(rule, partLimits)
// 	}

// 	return acceptableParts
// }

// func generaliseRule(rules ruleSet, workflowName string) []partLimits {
// 	// Initialise map of acceptable part values
// 	var acceptableParts []partLimits

//         // If our worflow is Accept or Reject then return the max limits or empty respectively
//         if workflowName == "A" {
//                 return []partLimits{makeMaxPartLimits()}
//         } else if workflowName == "R" {
//                 return []partLimits{}
//         }

// 	// Fetch workflow
// 	wflow := rules.workflows[workflowName]

// 	// Loop over rules in workflow
// 	for _, rule := range wflow.rules {
// 	// 	// Check if this is an end rule
//                 if rule.end {
//                         // This is an end rule
//                         // Check if this is an accept rule
//                         switch rule.target {
//                         case "A":
//                                 // This is an accept rule
//                                 return []partLimits{makeMaxPartLimits()}
//                         case "R":
//                                 // This is a reject rule
//                                 return []partLimits{}
//                         default:
//                                 // Process next named rule
//                                 return generaliseRule(rules, rule.target)

//                                 // Add the limits to the acceptable parts map
//                                 acceptableParts = append(acceptableParts, limits)
//                         }
//                 }
// 	}
// 	return acceptableParts
// }

func makeMaxPartLimits() partLimits {
	var limits partLimits
	limits.min = make(map[string]int)
	limits.max = make(map[string]int)
	for _, key := range []string{"x", "m", "a", "s"} {
		limits.min[key] = 1
		limits.max[key] = 4000
	}
	limits.valid = true
	return limits
}

// Split the part limits into 2, based on a rule
func splitPartLimits(limit partLimits, rule rule) (partLimits, partLimits) {
	var trueLimits, falseLimits partLimits

	// Copy limits
	trueLimits.max = maps.Clone(limit.max)
	trueLimits.min = maps.Clone(limit.min)
	falseLimits.max = maps.Clone(limit.max)
	falseLimits.min = maps.Clone(limit.min)

	operand := rule.operand

	// Check if the operand matches
	switch rule.operator {
	case LT:
		trueLimits.max[operand] = rule.value - 1
		falseLimits.min[operand] = rule.value
	case GT:
		trueLimits.min[operand] = rule.value + 1
		falseLimits.max[operand] = rule.value
	default:
		panic("Invalid operator")
	}
	validatePartLimits(&trueLimits)
	validatePartLimits(&falseLimits)

	return trueLimits, falseLimits
}

func validatePartLimits(limits *partLimits) {
	for _, key := range []string{"x", "m", "a", "s"} {
		if limits.min[key] > limits.max[key] {
			limits.valid = false
			return
		}
	}
	limits.valid = true
}

// func applyRuleToPartLimits(rule rule, limits []partLimits) []partLimits {
// 	var newLimits []partLimits

// 	// spew.Dump(limits)
// 	// spew.Dump(rule)
// 	// if rule.end {
// 	// 	// This is an end rule
// 	// 	// Check if this is an accept rule
// 	// 	switch rule.target {
// 	// 	case "A":
// 	// 		// This is an accept rule, start with everything and reduce limits as we go
// 	// 		return []partLimits{makeMaxPartLimits()}
// 	// 	case "R":
// 	// 		// This is a reject rule
// 	// 		return []partLimits{}
// 	// 	default:
// 	// 		// Process next named rule
// 	// 		return generaliseWorkflow(rules, rule.target)
// 	// 	}
// 	// }

// 	// Loop over limits
// 	for _, limit := range limits {
// 		// Check if the operand matches
// 		switch rule.operand {
// 		case "x":
// 			// Check if the x operand matches
// 			switch rule.operator {
// 			case LT:
// 				limit.xMax = rule.value - 1
// 			case GT:
// 				limit.xMin = rule.value + 1
// 			default:
// 				panic("Invalid operator")
// 			}
// 		case "m":
// 			// Check if the m operand matches
// 			switch rule.operator {
// 			case LT:
// 				limit.mMax = rule.value - 1
// 			case GT:
// 				limit.mMin = rule.value + 1
// 			default:
// 				panic("Invalid operator")
// 			}
// 		case "a":
// 			// Check if the a operand matches
// 			switch rule.operator {
// 			case LT:
// 				limit.aMax = rule.value - 1
// 			case GT:
// 				limit.aMin = rule.value + 1
// 			default:
// 				panic("Invalid operator")
// 			}
// 		case "s":
// 			// Check if the s operand matches
// 			switch rule.operator {
// 			case LT:
// 				limit.sMax = rule.value - 1
// 			case GT:
// 				limit.sMin = rule.value + 1
// 			default:
// 				panic("Invalid operator")
// 			}
// 		default:
// 			panic("Invalid operand")
// 		}
// 		newLimits = append(newLimits, limit)
// 	}

// 	return newLimits
// }

// Process a part
func processPart(p part, rules *ruleSet, ruleName string) int {
	fmt.Printf("Rule: %s\n", ruleName)

	// Check in case the next target is accept or reject
	switch ruleName {
	case "A":
		// This is an accept rule
		return p.x + p.m + p.a + p.s
	case "R":
		// This is a reject rule
		return 0
	}

	// Get workflow for this part
	wflow := rules.workflows[ruleName]

	// Loop over rules
	for _, rule := range wflow.rules {
		// Check if this is an end rule
		if rule.end {
			// This is an end rule
			switch rule.target {
			case "A":
				// This is an accept rule
				return p.x + p.m + p.a + p.s
			case "R":
				// This is a reject rule
				return 0
			default:
				// Process next named rule
				return processPart(p, rules, rule.target)
			}
		} else {
			// This is a normal rule
			// Check if the part matches the rule
			if checkRule(p, rule) {
				// This part matches the rule
				// Process the next rule
				return processPart(p, rules, rule.target)
			}
		}
	}

	// If we get here then we didn't match any rules, it should be impossible to get here
	panic("No matching rule")
}

// Check if a part matches a rule
func checkRule(p part, rule rule) bool {
	// Check if the operand matches
	switch rule.operand {
	case "x":
		// Check if the x operand matches
		switch rule.operator {
		case LT:
			return p.x < rule.value
		case GT:
			return p.x > rule.value
		default:
			panic("Invalid operator")
		}
	case "m":
		// Check if the m operand matches
		switch rule.operator {
		case LT:
			return p.m < rule.value
		case GT:
			return p.m > rule.value
		default:
			panic("Invalid operator")
		}
	case "a":
		// Check if the a operand matches
		switch rule.operator {
		case LT:
			return p.a < rule.value
		case GT:
			return p.a > rule.value
		default:
			panic("Invalid operator")
		}
	case "s":
		// Check if the s operand matches
		switch rule.operator {
		case LT:
			return p.s < rule.value
		case GT:
			return p.s > rule.value
		default:
			panic("Invalid operator")
		}
	default:
		panic("Invalid operand")
	}
}

// Parse the input to create a map of rules and set of part data
func parseInput(input string) (ruleSet, []part) {
	var rules ruleSet
	var parts []part

	// Initialise rules map
	rules.workflows = make(map[string]workflow)

	// Split input into lines
	lines := strings.Split(input, "\n")

	// Loop over lines
	for _, line := range lines {
		// Check if line is blank
		if len(line) == 0 {
			continue
		} else if line[0] == '{' {
			// This is a part
			parts = append(parts, parsePart(line))
		} else {
			// This is a workflow
			wflowName, wflow := parseWorkflow(line)
			rules.workflows[wflowName] = wflow
		}
	}
	return rules, parts
}

// Parse a workflow from a string
func parseWorkflow(line string) (string, workflow) {
	var wflow workflow

	// Split line into parts
	parts := strings.Split(line, "{")

	// Get workflow name
	wflowName := parts[0]

	// Get workflow rules
	wflow.rules = parseRules(strings.TrimRight(parts[1], "}"))

	return wflowName, wflow
}

// Parse a set of rules from a string
func parseRules(line string) []rule {
	var rules []rule

	// Split line into rules
	ruleStrs := strings.Split(line, ",")

	// Loop over rules
	for _, ruleStr := range ruleStrs {
		rules = append(rules, parseRule(ruleStr))
	}

	return rules
}

// Parse a rule from a string
func parseRule(line string) rule {
	var rule rule

	// Check if this is an end state
	if !strings.ContainsRune(line, ':') {
		rule.end = true
		rule.target = line
		return rule
	}

	// Otherwise it's a normal rule
	// fmt.Println("Ruleline:" + line)

	// Split line into rule and target parts
	parts := strings.Split(line, ":")

	// Get target
	rule.target = parts[1]

	// Use regex to split rule into operand, operator and value
	r, _ := regexp.Compile("([xmas])([<>])([0-9]+)")
	matches := r.FindStringSubmatch(parts[0])
	if len(matches) != 4 {
		panic("Invalid rule")
	}

	// Get operand
	rule.operand = matches[1]

	// Get operator
	rule.operator = encodeOperator(matches[2])

	// Get value
	var err error
	rule.value, err = strconv.Atoi(matches[3])
	if err != nil {
		panic(err)
	}

	return rule
}

func encodeOperator(operator string) int {
	switch operator {
	case "<":
		return LT
	case ">":
		return GT
	default:
		panic("Invalid operator")
	}
}

func decodeOperator(operator int) string {
	switch operator {
	case LT:
		return "<"
	case GT:
		return ">"
	default:
		panic("Invalid operator")
	}
}

// Parse a part from a string
func parsePart(line string) part {
	var part part

	// Split line into part details
	parts := strings.Split(strings.Trim(line, "{}"), ",")

	// Loop over part details
	for _, partStr := range parts {
		// Split part into key and value
		kv := strings.Split(partStr, "=")
		// spew.Dump(kv)

		// Get key
		key := kv[0]

		// Get value
		value, err := strconv.Atoi(kv[1])
		if err != nil {
			panic(err)
		}

		// Store value in part
		switch key {
		case "x":
			part.x = value
		case "m":
			part.m = value
		case "a":
			part.a = value
		case "s":
			part.s = value
		default:
			panic("Invalid key")
		}
	}

	return part
}

// Load file contents into a string and return it
func loadFileContents(filename string) string {
	// Read contents of file into a string
	fileBytes, err := os.ReadFile(filename) // just pass the file name
	if err != nil {
		panic(err)
	}

	return string(fileBytes) // convert content to a 'string'
}
