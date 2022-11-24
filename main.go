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

// The function recives a path to the project root folder (e.g. pro 7)
// and for each folder that contains vm file it triggers the compilation process to create asm file
// in the same folder.
func pathFinder(rootPath string) {
	files, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err)
	}
	// loop over the items in the dir
	for _, item := range files {
		// if it is a dir: make recursive call
		if item.IsDir() {
			pathFinder(filepath.Join(rootPath, item.Name()))
			// if it is a vm file - trigger the compilation process to create asm file
		} else {
			r, _ := regexp.Compile(".*[.]vm")
			if r.MatchString(item.Name()) {
				currentFile = item.Name()
				createAsmFile(filepath.Join(rootPath, item.Name()))
			}
		}
	}
}

// The function recives a path to the vm file - calls to the translator,
// and writes the the resulting hack code to asm file.
func createAsmFile(inputFilePath string) {
	outputFile := strings.Split(inputFilePath, ".")[0] + ".asm"
	// Clean from files from previous runs
	if _, err := os.Stat(outputFile); err == nil {
		os.Remove(outputFile)
	}
	// translte (compile) the file.
	hackCode := []byte(vmToAsmTraslator(inputFilePath))
	f, _ := os.Create(outputFile)
	f.Write(hackCode)
	f.Close()
}

// The function recives a path a to vm file, and returns string which is the translation
// of this file to hack code
func vmToAsmTraslator(filename string) string {
	readFile, err := os.Open(filename)
	var hackCode string
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	// loop over the lines in the vm file
	for fileScanner.Scan() {
		// parse the line - get the type of the cmd and its arguments
		cmdType, args := parsLine(fileScanner.Text())
		// if the line contains a valid command - translate it.
		if (cmdType != cComment) && (cmdType != cErr) {
			hackCode += cmdHandlersMap[cmdType](args)
		}
	}
	readFile.Close()

	return hackCode
}

// the function gets a single line from vm file and parses it - returns
// the type of the cmd and
func parsLine(line string) (Command, []string) {
	for key, element := range cmdRegexMap {
		r, _ := regexp.Compile(element)
		if r.MatchString(line) {
			return key, strings.Split(line, " ")
		}
	}
	return cErr, []string{}
}

// push handler - to translate vm push command
func pushHandler(args []string) string {
	segmant := args[1]
	offset, _ := strconv.Atoi(args[2])
	resString := ""

	// The translation according to the different segments
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
	// common code for all aegments
	resString += "@SP" + "\n"
	resString += "A = M" + "\n"
	resString += "M = D" + "\n"
	resString += "@SP" + "\n"
	resString += "M = M + 1" + "\n"
	return resString
}

// pop handler - to translate vm push command
func popHandler(args []string) string {
	segmant := args[1]

	// common code for all aegments
	resString := "@SP" + "\n"
	resString += "A = M - 1" + "\n"
	resString += "D = M" + "\n"

	// The translation according to the different segments
	switch segmant {
	case "pointer":
		resString += "@" + segmentsNameMap[segmant+args[2]] + "\n"
		resString += "A = M" + "\n"
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
	// common code for all aegments
	resString += "M = D" + "\n"
	resString += "@SP" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

// arithmatic handler - to translate vm arithmatic commands (neg, not sub, add, or, and)
func arithmaticHandler(args []string) string {
	action := args[0]

	// common code for all operations
	resString := "@SP" + "\n"
	resString += "A = M - 1" + "\n"

	// unaries operations translation (neg, not)
	if action == "neg" {
		resString += "M = -M" + "\n"
		return resString
	} else if action == "not" {
		resString += "M = !M" + "\n"
		return resString
	}
	// binaries operations translation (add, or, sub, and)
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
	// common code for binaries operations
	resString += "@SP" + "\n"
	resString += "M = M - 1" + "\n"
	return resString
}

// comp handler - to translate vm comparison commands (eq, gt, lt)
func compHandler(args []string) string {
	// the labels for the commands (with numbering for identification)
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

// helper function - given int n return ("A = A + 1") * n
// (calculate offset from segment base)
func advanceABy(offset int) string {
	resStr := ""
	for i := 0; i < offset; i++ {
		resStr += "A = A + 1" + "\n"
	}
	return resStr
}
