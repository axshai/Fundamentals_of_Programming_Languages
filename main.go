package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	inputPath := "try.vm"
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
		args := strings.Split(cmd, " ")
		fmt.Println(args)
		hackCode := cmdHandlersMap[cmdType](args)
		fmt.Println(hackCode) //
	}
	readFile.Close()
}

func parsLine(line string) (Command, string) {
	for key, element := range cmdRegexMap {
		r, _ := regexp.Compile(element)
		if r.MatchString(line) {
			return key, r.FindString(line)
		}
	}
	return err, ""
}

func pushHandler(args []string) string {
	segmant := args[1]
	offset, _ := strconv.Atoi(args[2])
	resString := ""
	if segmant != "constant" {
		resString += "@" + segmant + "\n"
		resString += "A = M" + "\n"
		for i := 0; i < offset; i++ {
			resString += "A = A + 1" + "\n"
		}
		resString += "D = M" + "\n"
	} else {
		resString += "@" + args[2] + "\n"
		resString += "D = A" + "\n"
	}
	resString += "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "M = D" + "\n"
	resString += "@sp" + "\n"
	resString += "M = M + 1"
	return resString
}

func popHandler(args []string) string {
	segmant := args[1]
	offset, _ := strconv.Atoi(args[2])
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

func arithmaticHandler(args []string) string {
	action := args[0]
	resString := "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "A = A - 1" + "\n"
	if action == "neg" {
		resString += "M = -M" + "\n"
		return resString
	} else if action == "not" {
		resString += "M = !M" + "\n"
		return resString
	}
	resString += "D = M" + "\n"
	resString += "A = A - 1" + "\n"
	if action == "sub" {
		resString += "M = M - D" + "\n"
	} else if action == "add" {
		resString += "M = D + M" + "\n"
	} else if action == "or" {
		resString += "M = D | M" + "\n"
	} else if action == "and" {
		resString += "M = D & M" + "\n"
	}
	resString += "@sp" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

func compHandler(args []string) string {
	trueLabel := fmt.Sprintf("(LABEL_T_%s)", strconv.Itoa(labelCounter))
	endLabel := fmt.Sprintf("(LABEL_E_%s)", strconv.Itoa(labelCounter))
	labelCounter++

	action := args[0]
	resString := "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M - D" + "\n"
	resString += "@" + "trueLabel" + "\n"
	if action == "eq" {
		resString += "D;JEQ" + "\n"
	} else if action == "gt" {
		resString += "D;JGT" + "\n"
	} else if action == "lt" {
		resString += "D;JLT" + "\n"
	}
	resString += "M = 0" + "\n"
	resString += "@" + "endLabel" + "\n"
	resString += "0;JMP" + "\n"
	resString += trueLabel + "\n"
	resString += "M = 1" + "\n"
	resString += endLabel + "\n"
	resString += "@sp" + "\n"
	resString += "M = M - 1" + "\n"

	return resString
}
