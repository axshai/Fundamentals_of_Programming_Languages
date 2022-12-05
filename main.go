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
	outR, _ := regexp.Compile(".*[.]cmp")
	asmR, _ := regexp.Compile(".*[.]asm")
	vmR, _ := regexp.Compile(".*[.]vm")
	initVmR, _ := regexp.Compile(".*[.]vm")
	outPutFile := ""
	vmFiles := []string{}
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
			if outR.MatchString(item.Name()) {
				outPutFile = item.Name()[:strings.Index(item.Name(), ".cmp")]
				outPutFile = filepath.Join(rootPath, outPutFile)
			} else if asmR.MatchString(item.Name()) {
				os.Remove(filepath.Join(rootPath, item.Name()))
			} else if initVmR.MatchString(item.Name()) {
				vmFiles = append([]string{filepath.Join(rootPath, item.Name())}, vmFiles...)
				currentFile = item.Name()
			} else if vmR.MatchString(item.Name()) {
				vmFiles = append(vmFiles, filepath.Join(rootPath, item.Name()))
				currentFile = item.Name()
			}
			if len(outPutFile) > 0 {
				createAsmFile(vmFiles, outPutFile)
			}
		}
	}
}

// The function recives a path to the vm file - calls to the translator,
// and writes the the resulting hack code to asm file.
func createAsmFile(inputFilesPath []string, outputFileName string) {
	outputFile := outputFileName + ".asm"
	hackCode := []byte{}
	// Clean from files from previous runs
	if _, err := os.Stat(outputFile); err == nil {
		os.Remove(outputFile)
	}
	// translte (compile) the file.
	for _, file := range inputFilesPath {
		hackCode = append(hackCode, []byte(vmToAsmTraslator(file))...)
	}
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
	offset := 0
	if len(args) > 2 {
		offset, _ = strconv.Atoi(args[2])
	}
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
		staticLabel := fileNamePrefix(args[2])
		resString += "@" + staticLabel + "\n"
		resString += "D = M" + "\n"
	case "temp":
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += advanceABy(offset, "+")
		resString += "D = M" + "\n"
	case "LCL", "ARG", "THIS", "THAT":
		resString += "@" + args[1] + "\n"
		resString += "D = M" + "\n"
	default:
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += "A = M" + "\n"
		resString += advanceABy(offset, "+")
		resString += "D = M" + "\n"
	}
	// common code for all aegments
	resString += "@SP" + "\n"
	resString += "A = M" + "\n"
	resString += "M = D" + "\n"
	resString += movePointer("SP", "+")
	return resString
}

// pop handler - to translate vm push command
func popHandler(args []string) string {
	segmant := args[1]

	// common code for all aegments
	resString := topStackPeek("SP")

	// The translation according to the different segments
	switch segmant {
	case "pointer":
		resString += "@" + segmentsNameMap[segmant+args[2]] + "\n"
		resString += "A = M" + "\n"
	case "static":
		resString += "@" + fileNamePrefix(args[2]) + "\n"
	case "temp":
		offset, _ := strconv.Atoi(args[2])
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += advanceABy(offset, "+")
	case "LCL", "ARG", "THIS", "THAT":
		resString += "@" + segmant + "\n"
	default:
		offset, _ := strconv.Atoi(args[2])
		resString += "@" + segmentsNameMap[segmant] + "\n"
		resString += "A = M" + "\n"
		resString += advanceABy(offset, "+")
	}
	// common code for all aegments
	resString += "M = D" + "\n"
	resString += movePointer("SP", "-")
	return resString
}

// arithmatic handler - to translate vm arithmatic commands (neg, not sub, add, or, and)
func arithmaticHandler(args []string) string {
	action := args[0]

	// common code for all operations
	resString := topStackPeek("SP")

	// unaries operations translation (neg, not)
	if action == "neg" {
		return resString[:strings.Index(resString, "D")] + "M = -M" + "\n"
	} else if action == "not" {
		return resString[:strings.Index(resString, "D")] + "M = !M" + "\n"
	}
	// binaries operations translation (add, or, sub, and)
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
	resString += movePointer("SP", "-")
	return resString
}

// comp handler - to translate vm comparison commands (eq, gt, lt)
func compHandler(args []string) string {
	// the labels for the commands (with numbering for identification)
	trueLabel := "LABEL_T_" + strconv.Itoa(labelCounter)
	endLabel := "LABEL_E_" + strconv.Itoa(labelCounter)
	labelCounter++

	action := args[0]
	// common code to take the top 2 arguments in the stack to compare them
	resString := topStackPeek("SP")
	resString += "A = A - 1" + "\n"
	resString += "D = M - D" + "\n"
	// if the comparation is true jump to true label
	resString += "@" + trueLabel + "\n"
	// translation acoording to the different comperations
	if action == "eq" {
		resString += "D;JEQ" + "\n"
	} else if action == "gt" {
		resString += "D;JGT" + "\n"
	} else if action == "lt" {
		resString += "D;JLT" + "\n"
	}
	resString += "D = 0" + "\n"
	// if the comparation is false jump to end label
	resString += "@" + endLabel + "\n"
	resString += "0;JMP" + "\n"
	resString += fmt.Sprintf("(%s)", trueLabel) + "\n"
	resString += "D = -1" + "\n"
	resString += fmt.Sprintf("(%s)", endLabel) + "\n"
	// common code to update the stack and the sp
	resString += "@SP" + "\n"
	resString += "A = M - 1" + "\n"
	resString += "A = A - 1" + "\n"
	resString += "M = D" + "\n"
	resString += movePointer("SP", "-")
	return resString
}

func labelHandler(args []string) string {
	return fmt.Sprint("(", fileNamePrefix(args[1]), ")") + "\n"
}
func gotoHandler(args []string) string {
	resString := "@" + fileNamePrefix(args[1]) + "\n"
	resString += "0;JMP" + "\n"
	return resString
}

func ifGotoHndler(args []string) string {
	resString := topStackPeek("SP")
	resString += "@" + fileNamePrefix(args[1]) + "\n"
	resString += "D;JNE" + "\n"
	resString += movePointer("SP", "-")
	return resString
}

func callHandler(args []string) string {
	retLabel := fmt.Sprint(args[1], ".", "ReturnAddress", "_", strconv.Itoa(labelCounter))
	n, _ := strconv.Atoi(args[2])
	segs := []string{"LCL", "ARG", "THIS", "THAT"}
	labelCounter += 1

	resString := pushHandler([]string{"push", "constant", retLabel})
	for _, seg := range segs {
		resString += pushHandler([]string{"push", seg})
	}
	resString += "@SP" + "\n"
	resString += "A = M" + "\n"
	resString += advanceABy(n+5, "-")
	resString += "D = A" + "\n"
	resString += "@ARG" + "\n"
	resString += "M = D" + "\n"
	resString += "@SP" + "\n"
	resString += "D=M" + "\n"
	resString += "@LCL" + "\n"
	resString += "M = D" + "\n"
	resString += gotoHandler([]string{"goto", args[1]})
	resString += fmt.Sprint("(", retLabel, ")") + "\n"
	return resString

}

func functionHandler(args []string) string {
	k, _ := strconv.Atoi(args[2])
	resString := fmt.Sprint("(", fileNamePrefix(args[1]), ")") + "\n"
	for i := 0; i < k; i++ {
		resString += pushHandler([]string{"push", "constant", "0"})
	}

	return resString
}

func returnHandler(args []string) string {
	segs := []string{"THAT", "THIS", "ARG", "LCL"}
	// FRAME = LCL
	resString := "@LCL" + "\n"
	resString += "D = M" + "\n"
	// RAM[13] = (LOCAL - 5)
	resString += "@5" + "\n"
	resString += "A=D-A" + "\n"
	resString += "D=M" + "\n"
	resString += "@13" + "\n"
	resString += "M=D" + "\n"
	// *ARG = pop() - put ret value in its place
	resString += popHandler([]string{"pop", "argument", "0"})
	// SP = ARG+1
	resString += "@ARG" + "\n"
	resString += "D = M" + "\n"
	resString += "@SP" + "\n"
	resString += "M = D + 1" + "\n"
	// SEGMENTS = *(FRAM-i): i=1...5

	for _, seg := range segs {
		resString += restoreSegmants(seg)
	}

	// goto RET
	resString += "@13" + "\n"
	resString += "A=M" + "\n"
	resString += "0;JMP" + "\n"
	return resString
}

// ---------------------------------------------------------------------------------------
// helper function - given int n return ("A = A +/- 1") * n
func advanceABy(steps int, direction string) string {
	resStr := ""
	for i := 0; i < steps; i++ {
		resStr += fmt.Sprintf("A = A %s 1\n", direction)
	}
	return resStr
}

func fileNamePrefix(l string) string {
	return fmt.Sprint(currentFile, ".", l)
}

func topStackPeek(topPointer string) string {
	resString := "@" + topPointer + "\n"
	resString += "A = M - 1" + "\n"
	resString += "D = M" + "\n"
	return resString
}

func movePointer(pointer string, direction string) string {
	resString := "@" + pointer + "\n"
	resString += fmt.Sprintf("M = M %s 1\n", direction)
	return resString
}

func restoreSegmants(seg string) string {
	resString := topStackPeek("LCL")
	resString += "@" + seg + "\n"
	resString += "M = D" + "\n"
	resString += movePointer("LCL", "-")
	return resString
}
