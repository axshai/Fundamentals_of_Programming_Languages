package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Command int

const (
	cArithmetic Command = iota
	cComp
	cLogic
	cPush
	cPop
	cLabel
	cGoto
	cIf
	comment
	err
)

func main() {
	inputPath := "BasicTest.vm"
	readFromFile(inputPath)
}

func readFromFile(filename string) {
	readFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		cmdType := parsLine(fileScanner.Text())
	}
	readFile.Close()
}

func parsLine(line string) Command {
	spaces := "^(\t|\\s)*"
	cmdRegexMap := map[Command]string{
		cLogic:      "^and *|^or *|^not *",
		cComp:       "^eq *|^gt *|^lt *",
		cArithmetic: "^add *|^sub *|^neg",
		cPush:       "^push *",
		cPop:        "^pop *",
		cLabel:      "^label:",
		cGoto:       "^goto *",
		cIf:         "^if *",
		comment:     "^//*",
	}
	t1 := strings.ReplaceAll(line, "", "")
	fmt.Println(t1)
	t := strings.Split(t1, " ")
	fmt.Println(t[0])
	for key, element := range cmdRegexMap {
		match, _ := regexp.MatchString(spaces+element, line)
		if match {
			return key
		}
	}
	return err
}

func pushHandler(segmant string, offset int) string {
	resString := ""
	if segmant != "constant" {
		resString += "@" + segmant + "\n"
		resString += "A = M" + "\n"
		for i := 0; i < offset; i++ {
			resString += "A = A + 1" + "\n"
		}
		resString += "D = M" + "\n"
	} else {
		resString += "@" + strconv.Itoa(offset) + "\n"
		resString += "D = A" + "\n"
	}
	resString += "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "M = D" + "\n"
	resString += "@sp" + "\n"
	resString += "M = M + 1"
	return resString
}

func popHandler(segmant string, offset int) string {
	resString := "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M" + "\n"
	resString += "@" + segmant + "\n"
	resString += "A = M" + "\n"
	for i := 0; i < offset; i++ {
		resString += "A = A + 1" + "\n"
	}
	resString += "M = D" + "\n"
	resString += "@sp" + "\n"
	resString += "M = M - 1"
	return resString
}

func plusMinusHandler(sign string) string {
	resString := "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M" + "\n"
	resString += "A = A - 1" + "\n"
	if sign == "sub" {
		resString += "M = D - M" + "\n"
	} else {
		resString += "M = D + M" + "\n"
	}
	resString += "@sp" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

func negHandler(sign string) string {
	resString := "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M" + "\n"

	resString += "M = -M" + "\n"

	resString += "@sp" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

func compHandler(comp string) string {
	resString := "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M" + "\n"
	resString += "M = -M" + "\n"
	resString += "@sp" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}
