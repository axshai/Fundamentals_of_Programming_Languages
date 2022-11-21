package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
		cmdType, cmd := parsLine(fileScanner.Text())
		fmt.Println(cmdType, cmd)
	}
	readFile.Close()
}

func parsLine(line string) (Command, string) {
	spaces := "(\t|\\s)*"
	cmdRegexMap := map[Command]string{
		cArithmetic: fmt.Sprintf("^%s(add|sub|neg).*", spaces),
		cComp:       fmt.Sprintf("^%s(eq|gt|lt).*", spaces),
		cLogic:      fmt.Sprintf("^%s(and|or|not).*", spaces),
		cPush:       fmt.Sprintf("^%spush.*", spaces),
		cPop:        fmt.Sprintf("^%spop.*", spaces),
		cLabel:      fmt.Sprintf("^%slabel:.*", spaces),
		cGoto:       fmt.Sprintf("^%sgoto.*", spaces),
		cIf:         fmt.Sprintf("^%sif.*", spaces),
		comment:     fmt.Sprintf("^%s//", spaces),
	}
	for key, element := range cmdRegexMap {
		r, _ := regexp.Compile(element)
		if r.MatchString(line) {
			return key, r.FindString(line)
		}
	}
	return err, ""
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
