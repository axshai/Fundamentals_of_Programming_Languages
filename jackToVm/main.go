package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
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
				outPutFile := strings.Split(currentFile, ".")[0] + "Tk.xml"
				outPutFile = filepath.Join(rootPath, outPutFile)
				currentFile = filepath.Join(rootPath, currentFile)
				fmt.Println(outPutFile)
				jackToTokensTraslator(currentFile, outPutFile)
			}
		}
	}
}

// The function recives a path a to jack file, and path to the output file
// and translate the jack file into tokens inside xml file in the output file
func jackToTokensTraslator(inputFile string, outputFile string) {
	// Clean from files from previous runs
	if _, err := os.Stat(outputFile); err == nil {
		os.Remove(outputFile)
	}
	// create the tokenizer struct to translate the file
	t := newToknizer(outputFile, inputFile)
	for t.isThereMoreTokens() {
		typ, token := t.nextToken()
		if typ != comment && typ != multiComment && typ != err {
			t.writeToken(tokenTypeMap[typ], translateToken(token))
		}
	}
	t.closeToknizer()
}
