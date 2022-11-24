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
	inputFile := os.Args[1]
	outputFile := strings.Split(inputFile, ".")[0] + ".asm"
	if _, err := os.Stat(outputFile); err == nil {
		os.Remove(outputFile)
	}
	d1 := []byte(vmToAsmTraslator(inputFile))
	os.WriteFile(outputFile, d1, 0x666)
}

func vmToAsmTraslator(filename string) string {
	readFile, err := os.Open(filename)
	var hackCode string
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		cmdType, cmd := parsLine(fileScanner.Text())
		args := strings.Split(cmd, " ")
		if (cmdType != cComment) && (cmdType != cErr) {
			hackCode += cmdHandlersMap[cmdType](args)
		}
	}
	readFile.Close()

	return initRam() + hackCode
}

func parsLine(line string) (Command, string) {
	for key, element := range cmdRegexMap {
		r, _ := regexp.Compile(element)
		if r.MatchString(line) {
			return key, line
		}
	}
	return cErr, ""
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
	resString := "@SP" + "\n"
	resString += "A = M - 1" + "\n"
	resString += "D = M" + "\n"
	switch segmant {
	case "pointer":
		resString += "@" + segmentsNameMap[segmant+args[2]] + "\n"
	case "static":
		resString += "@" + fmt.Sprintf("%s%s%s", os.Args[1], ".", args[2]) + "\n"
	default:
		offset, _ := strconv.Atoi(args[2])
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += advanceABy(offset)

	}
	resString += "M = D" + "\n"
	resString += "@SP" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

func arithmaticHandler(args []string) string {
	action := args[0]
	resString := "@SP" + "\n"
	resString += "A = M - 1" + "\n"
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
	resString += "@SP" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

func compHandler(args []string) string {
	trueLabel := "LABEL_T_" + strconv.Itoa(labelCounter)
	endLabel := "LABEL_E_" + strconv.Itoa(labelCounter)
	labelCounter++

	action := args[0]
	resString := "@SP" + "\n"
	resString += "A = M - 1" + "\n"
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
	resString += "@SP" + "\n"
	resString += "A = M - 1" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "M = D" + "\n"
	resString += "@SP" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

func initRam() string {
	res := ""
	res += "@256\n"
	res += "D = A\n"
	res += "@SP\n"
	res += "M = D\n"
	return res
}

func advanceABy(offset int) string {
	resStr := ""
	for i := 0; i < offset; i++ {
		resStr += "A = A + 1" + "\n"
	}
}
