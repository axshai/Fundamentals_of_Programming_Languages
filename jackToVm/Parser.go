package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type syntaxParser struct {
	outFile      *os.File
	indentations string
	tokensFile   *os.File
	fileScanner  *bufio.Scanner
}

func newParser(fileName string, tokensFileName string) syntaxParser {
	f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	tokensFile, _ := os.Open(tokensFileName)
	fileScanner := bufio.NewScanner(tokensFile)
	fileScanner.Split(bufio.ScanLines)
	return syntaxParser{outFile: f, tokensFile: tokensFile, indentations: "", fileScanner: fileScanner}
}

func (p *syntaxParser) getNextToken() (string, string) {
	typeR, _ := regexp.Compile("</.*>")
	tokenR, _ := regexp.Compile(">.*<")
	if p.fileScanner.Scan() {
		if tokenR.MatchString(p.fileScanner.Text()) {
			tokenStr := tokenR.FindString(p.fileScanner.Text())
			typeRStr := typeR.FindString(p.fileScanner.Text())
			return typeRStr[2 : len(typeRStr)-1], tokenStr[1 : len(tokenStr)-1]
		} else {
			return p.getNextToken()
		}
	} else {
		return "", ""
	}
}

func (p *syntaxParser) increasIndentation() {
	p.indentations += "  "
}

func (p *syntaxParser) decreasIndentation() {
	p.indentations = p.indentations[:len(p.indentations)-2]
}

// how can we now the tokenType?
func (p syntaxParser) writeToken(tokentType string, token string) {
	strToken := fmt.Sprintf("%s<%s> %s </%s>", p.indentations, tokentType, token, tokentType)
	p.outFile.WriteString(strToken + "\n")
}

func (p *syntaxParser) writeBlockTag(blockName string, closeTag bool) {
	tag := ""
	if !closeTag {
		tag = fmt.Sprintf("%s<%s>", p.indentations, blockName)
		p.increasIndentation()
	} else {
		p.decreasIndentation()
		tag = fmt.Sprintf("%s</%s>", p.indentations, blockName)
	}
	p.outFile.WriteString(tag + "\n")
}

func (p syntaxParser) closeToknizer() {
	p.outFile.Close()
	p.tokensFile.Close()
}
