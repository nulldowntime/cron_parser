package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func printHelp() {
	fmt.Println(`Usage:`)
	fmt.Println(`	parse_cron "only one argument that's a valid single line crontab format without special strings like @yearly"`)
}

func parseCronLine(l string) error {
	// Of course the regex could be a lot more sophisticated
	re := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(.+)$`)
	validCronLine := re.MatchString(l)
	if !validCronLine {
		return errors.New("Cron line does not comply with standard format")
	}

	subMatch := re.FindStringSubmatch(l)

	minuteStr := subMatch[1]
	hourStr := subMatch[2]
	dayOfMonthStr := subMatch[3]
	monthStr := subMatch[4]
	dayOfWeekStr := subMatch[5]
	commandStr := subMatch[6]

	minutesParsed, err := parseMinute(minuteStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hoursParsed, err := parseHours(hourStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dayOfMonthParsed, err := parseDayOfMonth(dayOfMonthStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	monthParsed, err := parseMonth(monthStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dayOfWeekParsed, err := parseDayOfWeek(dayOfWeekStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%-14s%s\n", "minute", joinInts(minutesParsed))
	fmt.Printf("%-14s%s\n", "hour", joinInts(hoursParsed))
	fmt.Printf("%-14s%s\n", "day of month", joinInts(dayOfMonthParsed))
	fmt.Printf("%-14s%s\n", "month", joinInts(monthParsed))
	fmt.Printf("%-14s%s\n", "day of week", joinInts(dayOfWeekParsed))
	fmt.Printf("%-14s%s\n", "command", commandStr)

	return nil
}

func joinInts(sl []int) string {
	var joined string
	var slStr []string

	for _, e := range sl {
		slStr = append(slStr, strconv.Itoa(e))
	}

	joined = strings.Join(slStr, " ")
	return joined
}

func genStarList(s string, start int, limit int) ([]int, error) {
	var numList []int

	interval := strings.TrimLeft(s, "*/")
	intervalNum, err := strconv.Atoi(interval)
	if err != nil {
		return numList, fmt.Errorf("interval: %s is not a number", interval)
	}
	if intervalNum == 0 {
		return numList, errors.New("interval cannot be 0")
	}

	for i := start; i <= limit; i++ {
		res := i % intervalNum
		if res == 0 {
			numList = append(numList, i)
		}
	}

	return numList, nil
}

func getFirst(s string) string {
	var f string
	sl := strings.Split(s, "")
	f = sl[0]
	return f
}

func getLast(s string) string {
	var l string
	sl := strings.Split(s, "")
	l = sl[len(sl)-1]
	return l
}

func genDashList(s string, start int, limit int) ([]int, error) {
	var numList []int

	first := getFirst(s)
	firstNum, err := strconv.Atoi(first)
	if err != nil {
		return numList, fmt.Errorf("First element: %s is not a number", first)
	}
	last := getLast(s)
	lastNum, err := strconv.Atoi(last)
	if err != nil {
		return numList, fmt.Errorf("Last element: %s is not a number", last)
	}

	if firstNum < start {
		return numList, fmt.Errorf("First element: %d lower than allowed: %d", firstNum, start)
	}

	if lastNum > limit {
		return numList, fmt.Errorf("Last element: %d greater than allowed: %d", lastNum, limit)
	}

	if firstNum > lastNum {
		return numList, fmt.Errorf("firt element: %d greater than last element: %d", firstNum, lastNum)
	}

	for i := firstNum; i <= lastNum; i++ {
		numList = append(numList, i)
	}

	return numList, nil
}

func genCommaList(s string, start int, limit int) ([]int, error) {
	var strList []string
	var numList []int
	strList = strings.Split(s, ",")

	for _, e := range strList {
		num, err := strconv.Atoi(e)
		if err != nil {
			return numList, fmt.Errorf("comma list element: %s is not a number", e)
		}
		if num < start {
			return numList, fmt.Errorf("List element: %d lower than allowed: %d", num, start)
		}
		if num > limit {
			return numList, fmt.Errorf("List element: %d greater than allowed: %d", num, limit)
		}
		numList = append(numList, num)
	}
	return numList, nil
}

func genFullList(s string, start int, limit int) []int {
	var numList []int

	for i := start; i <= limit; i++ {
		numList = append(numList, i)
	}

	return numList
}

func parseNumDef(s string, start int, limit int) ([]int, error) {

	var numList []int

	reStar := regexp.MustCompile(`^\*\/\d+$`)
	if reStar.MatchString(s) {
		numList, err := genStarList(s, start, limit)
		// just pass the error for now
		return numList, err
	}

	reDash := regexp.MustCompile(`^\d+-\d+$`)
	if reDash.MatchString(s) {
		numList, err := genDashList(s, start, limit)
		// just pass the error for now
		return numList, err
	}

	//reComma := regexp.MustCompile(`^(\d+\,)+\d$`)
	reComma := regexp.MustCompile(`^\d+(,\d+)+$`)
	if reComma.MatchString(s) {
		numList, err := genCommaList(s, start, limit)
		// just pass the error for now
		return numList, err
	}

	if s == "*" {
		numList = genFullList(s, start, limit)
		return numList, nil
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		// this, somewhat implicitly, handles all the rest of wrong/unhandled
		// formats
		return numList, fmt.Errorf("Invalid crontab format, %s is not a valid cron definition", s)
	}

	numList = []int{num}
	return numList, nil
}

func parseMinute(m string) ([]int, error) {

	var minutes []int
	start := 0
	limit := 59

	minutes, err := parseNumDef(m, start, limit)
	if err != nil {
		return minutes, err
	}

	return minutes, nil
}

func parseHours(h string) ([]int, error) {
	var hours []int
	start := 0
	limit := 23

	hours, err := parseNumDef(h, start, limit)
	if err != nil {
		return hours, err
	}

	return hours, nil
}

func parseDayOfMonth(dom string) ([]int, error) {
	var daysOfMonth []int
	start := 1
	limit := 31

	daysOfMonth, err := parseNumDef(dom, start, limit)
	if err != nil {
		return daysOfMonth, err
	}

	return daysOfMonth, nil
}

func parseMonth(m string) ([]int, error) {
	var months []int
	start := 1
	limit := 12

	months, err := parseNumDef(m, start, limit)
	if err != nil {
		return months, err
	}

	return months, nil
}

func parseDayOfWeek(dow string) ([]int, error) {
	var daysOfWeek []int
	start := 1
	limit := 7

	daysOfWeek, err := parseNumDef(dow, start, limit)
	if err != nil {
		return daysOfWeek, err
	}

	return daysOfWeek, nil
}

func main() {
	switch argCount := len(os.Args); {
	case argCount == 1:
		//programName := os.Args[0]
		fmt.Println("Argument missing")
		printHelp()
		os.Exit(0)
	case argCount == 2:
		cronLine := os.Args[1]
		//fmt.Println(cronLine)
		err := parseCronLine(cronLine)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case argCount > 2:
		fmt.Println("Too many arguments")
		printHelp()
		os.Exit(1)
	}
}
