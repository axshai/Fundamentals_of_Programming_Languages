package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// syntaxParser - Responsible for creating syntax tree (As xml file)
// from tokens file
type syntaxParser struct {
	outFile      *os.File
	indentations string
	tokensBuffer [][]string
	lineNumber   int
}

// syntaxParser "constructor"
func newParser(fileName string, tokensFileName string) syntaxParser {
	f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	newP := syntaxParser{outFile: f, indentations: "", lineNumber: 0, tokensBuffer: [][]string{}}
	newP.initTokensBuffer(tokensFileName)
	return newP
}

// The function initializes the buffer of syntaxParser with all the tokens and their type
func (p *syntaxParser) initTokensBuffer(tokensFileName string) {
	tokensFile, _ := os.Open(tokensFileName)
	fileScanner := bufio.NewScanner(tokensFile)
	fileScanner.Split(bufio.ScanLines)
	typeR, _ := regexp.Compile("</.*>")
	tokenR, _ := regexp.Compile(">.*<")
	// for each line in the tokens xml
	for fileScanner.Scan() {
		// if kine contains token - extract the token and its type
		if tokenR.MatchString(fileScanner.Text()) {
			tokenStr := tokenR.FindString(fileScanner.Text())
			typeRStr := typeR.FindString(fileScanner.Text())
			p.tokensBuffer = append(p.tokensBuffer, []string{typeRStr[2 : len(typeRStr)-1], tokenStr[1 : len(tokenStr)-1]})
		}
	}
	tokensFile.Close()
}

// The function return the next token and its type
func (p *syntaxParser) getNextToken() (string, string) {
	// if ther are not more tokens - return empty strings
	if p.lineNumber >= len(p.tokensBuffer) {
		return "", ""
	}
	tType, token := p.tokensBuffer[p.lineNumber][0], p.tokensBuffer[p.lineNumber][1]
	p.lineNumber++
	return tType, strings.TrimSpace(token)
}

func (p *syntaxParser) lookahead(steps int) (string, string) {
	// if ther are not more tokens - return empty strings
	if (p.lineNumber + steps - 1) >= len(p.tokensBuffer) {
		return "", ""
	}
	tType, token := p.tokensBuffer[p.lineNumber+steps-1][0], p.tokensBuffer[p.lineNumber+steps-1][1]
	return tType, strings.TrimSpace(token)
}

// function to manage indentations in the syntax xml file - increase the indentations
func (p *syntaxParser) increasIndentation() {
	p.indentations += "  "
}

// function to manage indentations in the syntax xml file - decrease the indentations
func (p *syntaxParser) decreasIndentation() {
	p.indentations = p.indentations[:len(p.indentations)-2]
}

// the function gets token type and the token value
// and write the token in the following format: <type> value </type>
func (p syntaxParser) writeToken(tokentType string, token string) {
	strToken := fmt.Sprintf("%s<%s> %s </%s>", p.indentations, tokentType, token, tokentType)
	p.outFile.WriteString(strToken + "\n")
}

// the function writes start/end (according to isCloseTag) block tag
// and changes the indents accordingly
func (p *syntaxParser) writeBlockTag(blockName string, isCloseTag bool) {
	tag := ""
	if !isCloseTag {
		tag = fmt.Sprintf("%s<%s>", p.indentations, blockName)
		p.increasIndentation()
	} else {
		p.decreasIndentation()
		tag = fmt.Sprintf("%s</%s>", p.indentations, blockName)
	}
	p.outFile.WriteString(tag + "\n")
}

// function to finish the toknizer
func (p syntaxParser) closeToknizer() {
	p.outFile.Close()
}
