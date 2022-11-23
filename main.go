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
	if segmant == "constant" {
		resString += "@" + args[2] + "\n"
		resString += "D = A" + "\n"
	} else {
		if segmant == "pointer" {
			if args[2] == "0" {
				segmant = "this"
			} else {
				segmant = "that"
				offset = 0
			}
			resString += "@" + segmentsNameMap[segmant] + "\n"
			resString += "A = M" + "\n"
		} else if segmant == "static" {
			segmant = fmt.Sprintf("%s%s%s", os.Args[1], ".", args[2])
			resString += "@" + segmant + "\n"
			offset = 0

		} else {
			resString += "@" + segmentsNameMap[segmant] + "\n"
			if segmant != "temp" {
				resString += "A = M" + "\n"
			}
		}
		for i := 0; i < offset; i++ {
			resString += "A = A + 1" + "\n"
		}
		resString += "D = M" + "\n"
	}
	resString += "@SP" + "\n"
	resString += "A = M" + "\n"
	resString += "M = D" + "\n"
	resString += "@SP" + "\n"
	resString += "M = M + 1" + "\n"
	return resString
}

func popHandler(args []string) string {
	segmant := args[1]
	offset, _ := strconv.Atoi(args[2])
	resString := "@SP" + "\n"
	resString += "A = M - 1" + "\n"
	resString += "D = M" + "\n"
	if segmant == "pointer" {
		if args[2] == "0" {
			segmant = "this"
		} else {
			segmant = "that"
			offset = 0
		}
		resString += "@" + segmentsNameMap[segmant] + "\n"
	} else if segmant == "static" {
		segmant = fmt.Sprintf("%s%s%s", os.Args[1], ".", args[2])
		resString += "@" + segmant + "\n"
		offset = 0
	} else {
		resString += "@" + segmentsNameMap[segmant] + "\n"
	}
	for i := 0; i < offset; i++ {
		resString += "A = A + 1" + "\n"
	}
	resString += "M = D" + "\n"
	resString += "@SP" + "\n"
	resString += "M = M - 1" + "\n"
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
	trueLabel := "LABEL_T_" + strconv.Itoa(labelCounter)
	endLabel := "LABEL_E_" + strconv.Itoa(labelCounter)
	labelCounter++

	action := args[0]
	resString := "@sp" + "\n"
	resString += "A = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "D = M - D" + "\n"
	resString += "@" + trueLabel + "\n"
	if action == "eq" {
		resString += "D;JEQ" + "\n"
	} else if action == "gt" {
		resString += "D;JGT" + "\n"
	} else if action == "lt" {
		resString += "D;JLT" + "\n"
	}
	resString += "D = 0" + "\n"
	resString += "@" + endLabel + "\n"
	resString += "0;JMP" + "\n"
	resString += fmt.Sprintf("(%s)", trueLabel) + "\n"
	resString += "D = -1" + "\n"
	resString += fmt.Sprintf("(%s)", endLabel) + "\n"
	resString += "@sp" + "\n"
	resString += "A = M - 1" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "M = D" + "\n"
	resString += "@sp" + "\n"
	resString += "M = M - 1" + "\n"

	return resString
}
