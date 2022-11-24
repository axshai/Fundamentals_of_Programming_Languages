package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	pathFinder(os.Args[1])
}

func pathFinder(rootPath string) {
	files, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if file.IsDir() {
			pathFinder(filepath.Join(rootPath, file.Name()))
		} else {
			r, _ := regexp.Compile(".*[.]vm")
			if r.MatchString(file.Name()) {
				currentFile = file.Name()
				createAsmFile(filepath.Join(rootPath, file.Name()))
			}
		}
	}
}

func createAsmFile(inputFilePath string) {
	outputFile := strings.Split(inputFilePath, ".")[0] + ".asm"
	if _, err := os.Stat(outputFile); err == nil {
		os.Remove(outputFile)
	}
	d1 := []byte(vmToAsmTraslator(inputFilePath))
	f, _ := os.Create(outputFile)
	f.Write(d1)
	f.Close()
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

	return hackCode
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

	switch segmant {
	case "constant":
		resString += "@" + args[2] + "\n"
		resString += "D = A" + "\n"
	case "pointer":
		resString += "@" + segmentsNameMap[segmant+args[2]] + "\n"
		resString += "A = M" + "\n"
		resString += "D = M" + "\n"
	case "static":
		staticLabel := fmt.Sprintf("%s%s%s", currentFile, ".", args[2])
		resString += "@" + staticLabel + "\n"
		resString += "D = M" + "\n"
	case "temp":
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += advanceABy(offset)
		resString += "D = M" + "\n"
	default:
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += "A = M" + "\n"
		resString += advanceABy(offset)
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
		resString += "@" + fmt.Sprintf("%s%s%s", currentFile, ".", args[2]) + "\n"
	case "temp":
		offset, _ := strconv.Atoi(args[2])
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += advanceABy(offset)
	default:
		offset, _ := strconv.Atoi(args[2])
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += "A = M" + "\n"
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

func advanceABy(offset int) string {
	resStr := ""
	for i := 0; i < offset; i++ {
		resStr += "A = A + 1" + "\n"
	}
	return resStr
}
