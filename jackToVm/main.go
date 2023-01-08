package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var currentFile string

func main() {
	// t := Table{}
	// t.insert("temp0", "int", "local")
	// t.insert("temp1", "int", "local")
	// t.insert("temp2", "int", "argument")
	// t.insert("temp3", "int", "argument")
	// t.insert("temp4", "int", "local")
	// t.insert("temp5", "int", "local")
	// fmt.Println(t["temp3"].varType)
	// fmt.Println(t["temp3"].varSeg)
	// fmt.Println(t["temp3"].varIndex)

	// fmt.Println(t["temp4"].varType)
	// fmt.Println(t["temp4"].varSeg)
	// fmt.Println(t["temp4"].varIndex)
	// fmt.Println(t.countType("local"))
	// fmt.Println(t.countType("argument"))
	// t.printTable()
	pathFinder(os.Args[1])
}

// The function recives a path to the project root folder (e.g. pro 10)
// and for each folder that contains jack files it triggers the compilation process to create token xml file
// in the same folder.
func pathFinder(rootPath string) {
	jackF, _ := regexp.Compile(".*[.]jack")
	files, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err)
	}
	// loop over the items in the dir
	for _, item := range files {
		// if it is a dir: make recursive call
		if item.IsDir() {
			pathFinder(filepath.Join(rootPath, item.Name()))
			// else - it is a jack files directory
		} else {
			// if it is a jack file - translate it to tokens xml
			if jackF.MatchString(item.Name()) {
				currentFile = item.Name()
				outPutFile := strings.Split(currentFile, ".")[0]
				outPutFile = filepath.Join(rootPath, outPutFile)
				currentFile = filepath.Join(rootPath, currentFile)
				jackToTokensTraslator(currentFile, outPutFile)
			}
		}
	}
}

// The function recives a path a to jack file, and path to the output file
// and translate the jack file into tokens inside xml file in the output file
func jackToTokensTraslator(inputFile string, outputFile string) {
	toknizerOutputFile := outputFile + "Tk.xml"
	parserOutputFile := outputFile + "k.xml"
	vmFileoutput := outputFile + ".vm"
	fmt.Println(toknizerOutputFile)
	// Clean from files from previous runs
	if _, err := os.Stat(toknizerOutputFile); err == nil {
		os.Remove(toknizerOutputFile)
	}
	if _, err := os.Stat(parserOutputFile); err == nil {
		os.Remove(parserOutputFile)
	}
	// create the tokenizer struct to translate the file
	t := newToknizer(toknizerOutputFile, inputFile)
	for t.isThereMoreTokens() {
		typ, token := t.nextToken()
		if typ != comment && typ != multiComment && typ != err {
			t.writeToken(tokenTypeMap[typ], translateToken(token))
		}
	}
	t.closeToknizer()
	initParsFuncs()
	classScopeTable = Table{}
	vw = newVmWriter(vmFileoutput)
	p := newParser(parserOutputFile, toknizerOutputFile)
	parseClass(p)
	p.closeToknizer()
}

func initParsFuncs() {
	statmentHandlersMap["let"] = letStatment
	statmentHandlersMap["if"] = ifStatment
	statmentHandlersMap["while"] = whileStatment
	statmentHandlersMap["do"] = doStatment
	statmentHandlersMap["return"] = returnStatment
}
