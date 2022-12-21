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
	tokensBuffer [][]string
	lineNumber   int
}

func newParser(fileName string, tokensFileName string) syntaxParser {
	f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	newP := syntaxParser{outFile: f, indentations: "", lineNumber: 0, tokensBuffer: [][]string{}}
	newP.initTokensBuffer(tokensFileName)
	return newP
}

func (p *syntaxParser) initTokensBuffer(tokensFileName string) {
	tokensFile, _ := os.Open(tokensFileName)
	fileScanner := bufio.NewScanner(tokensFile)
	fileScanner.Split(bufio.ScanLines)
	typeR, _ := regexp.Compile("</.*>")
	tokenR, _ := regexp.Compile(">.*<")
	for fileScanner.Scan() {
		if tokenR.MatchString(fileScanner.Text()) {
			tokenStr := tokenR.FindString(fileScanner.Text())
			typeRStr := typeR.FindString(fileScanner.Text())
			p.tokensBuffer = append(p.tokensBuffer, []string{typeRStr[2 : len(typeRStr)-1], tokenStr[1 : len(tokenStr)-1]})
		}
	}
	tokensFile.Close()
}

func (p *syntaxParser) getNextToken() (string, string) {
	if p.lineNumber >= len(p.tokensBuffer) {
		return "", ""
	}
	tType, token := p.tokensBuffer[p.lineNumber][0], p.tokensBuffer[p.lineNumber][1]
	p.lineNumber++
	return tType, token
}

func (p *syntaxParser) backToPrevToken() {
	p.lineNumber--
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
}

func parseClass(p syntaxParser) {
	p.writeBlockTag("class", false)
	p.writeToken(p.getNextToken()) //class
	p.writeToken(p.getNextToken()) // class name
	p.writeToken(p.getNextToken()) // {
	ParseClassVarDec(p)
	ParseSubRoutineDec(p)
	p.writeToken(p.getNextToken()) // }
	p.writeBlockTag("class", true)
}

func ParseClassVarDec(p syntaxParser) {
	tType, token := p.getNextToken()
	flag := false
	for token == "static" || token == "field" {
		flag = true
		p.writeBlockTag("classVarDec", false)
		p.writeToken(tType, token)     // field || static
		p.writeToken(p.getNextToken()) //<keyword> int </keyword> || <identifier> className </identifier>
		p.writeToken(p.getNextToken()) //<identifier> x </identifier>
		tType, token = p.getNextToken()
		for token == "," {
			p.writeToken(tType, token)     // <symbol> , </symbol>
			p.writeToken(p.getNextToken()) //<identifier> y </identifier>
			tType, token = p.getNextToken()
		}
		p.backToPrevToken()
		p.writeToken(p.getNextToken()) //<identifier> ; </identifier>
		tType, token = p.getNextToken()
	}
	p.backToPrevToken()
	if flag {
		p.writeBlockTag("classVarDec", true)
	}
}

func ParseSubRoutineDec(p syntaxParser) {
	tType, token := p.getNextToken()
	flag := false
	for token == "constructor" || token == "method" || token == "function" {
		flag = true
		p.writeBlockTag("subroutineDec", false)
		p.writeToken(tType, token)     // method || constructor || function
		p.writeToken(p.getNextToken()) //<keyword> type </keyword>
		p.writeToken(p.getNextToken()) //<identifier> funcNmae </identifier>
		p.writeToken(p.getNextToken()) //<symbol> ( </symbol>
		ParseParameterList(p)
		p.writeToken(p.getNextToken())  //<symbol> ) </symbol>
		tType, token = p.getNextToken() // method || constructor || function
	}
	p.backToPrevToken()
	if flag {
		p.writeBlockTag("subroutineDec", true)
	}
}

func ParseParameterList(p syntaxParser) {
	tType, token := p.getNextToken()
	p.writeBlockTag("parameterList", false)
	for token == "int" || token == "void" || token == "char" || token == "boolean" {
		p.writeToken(tType, token)     //<keyword> type </keyword>
		p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
		tType, token = p.getNextToken()
		for token == "," {
			p.writeToken(tType, token)     // <symbol> , </symbol>
			p.writeToken(p.getNextToken()) //<keyword> type </keyword>
			p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
			tType, token = p.getNextToken()
		}
		p.backToPrevToken()
		tType, token = p.getNextToken()
	}
	p.backToPrevToken()
	p.writeBlockTag("parameterList", true)

}
func ParseSubRoutineBody(p syntaxParser) {
	p.writeBlockTag("subroutineBody", false)
	p.writeToken(p.getNextToken()) // <symbol> {</symbol>
	ParsevarDec(p)
	ParseStatments(p)
	p.writeToken(p.getNextToken()) // <symbol> }</symbol>
	p.writeBlockTag("subroutineBody", true)

}

func ParsevarDec(p syntaxParser) {
	tType, token := p.getNextToken()
	flag := false
	if token == "var" {
		p.writeBlockTag("varDec", false)
		flag = true
	}
	for token == "var" {
		p.writeToken(tType, token)     //<keyword> var </keyword>
		p.writeToken(p.getNextToken()) //<keyword> type </keyword>
		p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
		tType, token = p.getNextToken()
		for token == "," {
			p.writeToken(tType, token)     // <symbol> , </symbol>
			p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
			tType, token = p.getNextToken()
		}
		p.backToPrevToken()
		p.writeToken(p.getNextToken()) //<keyword> ; </keyword>
		tType, token = p.getNextToken()
	}
	p.backToPrevToken()
	if flag {
		p.writeBlockTag("varDec", true)
	}
}

func ParseStatments(p syntaxParser) {
	var handler func(syntaxParser, string, string)
	tType, token := p.getNextToken()
	exists := false
	flag := false
	if handler, exists = statmentHandlersMap[token]; exists {
		p.writeBlockTag("statements", false)
		flag = true
	}
	for exists {
		handler(p, tType, token)
		tType, token = p.getNextToken()
	}
	p.backToPrevToken()
	if flag {
		p.writeBlockTag("statements", true)
	}

}

func letStatment(p syntaxParser, tType string, token string) {
	p.writeBlockTag("letStatement", false)
	p.writeToken(tType, token) //<identifier> varName </identifier>
	tType, token = p.getNextToken()
	if token == "[" {
		p.writeToken(tType, token)
		ParseExpression(p)
		p.writeToken(p.getNextToken())
		tType, token = p.getNextToken()
	}
	p.writeToken(tType, token) // =
	ParseExpression(p)
	p.writeToken(p.getNextToken()) //;

}
func ifStatment(p syntaxParser, tType string, token string)     {}
func whileStatment(p syntaxParser, tType string, token string)  {}
func doStatment(p syntaxParser, tType string, token string)     {}
func returnStatment(p syntaxParser, tType string, token string) {}

func ParseExpression(p syntaxParser) {

}

var statmentHandlersMap = map[string]func(syntaxParser, string, string){
	"let":    letStatment,
	"if":     ifStatment,
	"while":  whileStatment,
	"do":     doStatment,
	"return": returnStatment,
}
