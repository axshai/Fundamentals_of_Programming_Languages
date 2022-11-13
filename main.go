package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type Command int

const (
	cArithmetic Command = iota
	cPush
	cPop
	cLabel
	cGoto
	cIf
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
		parsLine(fileScanner.Text())
	}
	readFile.Close()
}

func parsLine(line string) {
	spaces := "^(\t|\\s)*"
	match, _ := regexp.MatchString(spaces+"//", line)
	fmt.Println(match)

}
