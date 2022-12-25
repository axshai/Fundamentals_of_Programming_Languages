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

func parseClass(p syntaxParser) {
	p.writeBlockTag("class", false)
	p.writeToken(p.getNextToken()) //class
	p.writeToken(p.getNextToken()) // class name
	p.writeToken(p.getNextToken()) // {
	ParseClassVarDec(&p)
	ParseSubRoutineDec(&p)
	p.writeToken(p.getNextToken()) // }
	p.writeBlockTag("class", true)
}

func ParseClassVarDec(p *syntaxParser) {
	_, token := p.lookahead(1)
	for token == "static" || token == "field" {
		p.writeBlockTag("classVarDec", false)
		p.writeToken(p.getNextToken()) // field || static
		p.writeToken(p.getNextToken()) //<keyword> int </keyword> || <identifier> className </identifier>
		p.writeToken(p.getNextToken()) //<identifier> x </identifier>
		_, token = p.lookahead(1)
		for token == "," {
			p.writeToken(p.getNextToken()) // <symbol> , </symbol>
			p.writeToken(p.getNextToken()) //<identifier> y </identifier>
			_, token = p.lookahead(1)
		}
		p.writeToken(p.getNextToken()) //<identifier> ; </identifier>
		p.writeBlockTag("classVarDec", true)
		_, token = p.lookahead(1)
	}

}

func ParseSubRoutineDec(p *syntaxParser) {
	_, token := p.lookahead(1)
	for token == "constructor" || token == "method" || token == "function" {
		p.writeBlockTag("subroutineDec", false)
		p.writeToken(p.getNextToken()) // method || constructor || function
		p.writeToken(p.getNextToken()) //<keyword> type </keyword>
		p.writeToken(p.getNextToken()) //<identifier> funcNmae </identifier>
		p.writeToken(p.getNextToken()) //<symbol> ( </symbol>
		ParseParameterList(p)
		p.writeToken(p.getNextToken()) //<symbol> ) </symbol>
		ParseSubRoutineBody(p)
		p.writeBlockTag("subroutineDec", true)
		_, token = p.lookahead(1)
	}
}

func ParseParameterList(p *syntaxParser) {
	_, token := p.lookahead(1)
	p.writeBlockTag("parameterList", false)
	for token == "int" || token == "void" || token == "char" || token == "boolean" {
		p.writeToken(p.getNextToken()) //<keyword> type </keyword>
		p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
		_, token = p.lookahead(1)
		for token == "," {
			p.writeToken(p.getNextToken()) // <symbol> , </symbol>
			p.writeToken(p.getNextToken()) //<keyword> type </keyword>
			p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
			_, token = p.lookahead(1)
		}
		_, token = p.lookahead(1)
	}
	p.writeBlockTag("parameterList", true)

}
func ParseSubRoutineBody(p *syntaxParser) {
	p.writeBlockTag("subroutineBody", false)
	p.writeToken(p.getNextToken()) // <symbol> {</symbol>
	ParsevarDec(p)
	ParseStatments(p)
	p.writeToken(p.getNextToken()) // <symbol> }</symbol>
	p.writeBlockTag("subroutineBody", true)
}

func ParsevarDec(p *syntaxParser) {
	_, token := p.lookahead(1)
	for token == "var" {
		p.writeBlockTag("varDec", false)
		p.writeToken(p.getNextToken()) //<keyword> var </keyword>
		p.writeToken(p.getNextToken()) //<keyword> type </keyword>
		p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
		_, token = p.lookahead(1)
		for token == "," {
			p.writeToken(p.getNextToken()) // <symbol> , </symbol>
			p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
			_, token = p.lookahead(1)
		}
		p.writeToken(p.getNextToken()) //<keyword> ; </keyword>
		p.writeBlockTag("varDec", true)
		_, token = p.lookahead(1)
	}
}

func ParseStatments(p *syntaxParser) {
	var handler func(*syntaxParser)
	_, token := p.lookahead(1)
	exists := false
	p.writeBlockTag("statements", false)
	handler, exists = statmentHandlersMap[token]
	for exists {
		handler(p)
		_, token = p.lookahead(1)
		handler, exists = statmentHandlersMap[token]
	}
	p.writeBlockTag("statements", true)
}

func letStatment(p *syntaxParser) {
	p.writeBlockTag("letStatement", false)
	p.writeToken(p.getNextToken()) //let
	p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
	_, token := p.lookahead(1)
	if token == "[" {
		p.writeToken(p.getNextToken())
		ParseExpression(p)
		p.writeToken(p.getNextToken())
		_, token = p.lookahead(1)
	}
	p.writeToken(p.getNextToken()) // =
	ParseExpression(p)
	p.writeToken(p.getNextToken()) //;
	p.writeBlockTag("letStatement", true)
}

func ifStatment(p *syntaxParser) {
	p.writeBlockTag("ifStatement", false)
	p.writeToken(p.getNextToken()) // if
	p.writeToken(p.getNextToken()) //(
	ParseExpression(p)
	p.writeToken(p.getNextToken()) //)
	p.writeToken(p.getNextToken()) //{
	ParseStatments(p)
	p.writeToken(p.getNextToken()) //}
	_, token := p.lookahead(1)
	if token == "else" {
		p.writeToken(p.getNextToken()) //else
		p.writeToken(p.getNextToken()) //{
		ParseStatments(p)
		p.writeToken(p.getNextToken()) //}
	}
	p.writeBlockTag("ifStatement", true)
}

func whileStatment(p *syntaxParser) {
	p.writeBlockTag("whileStatement", false)
	p.writeToken(p.getNextToken()) // while
	p.writeToken(p.getNextToken()) //(
	ParseExpression(p)
	p.writeToken(p.getNextToken()) //)
	p.writeToken(p.getNextToken()) //{
	ParseStatments(p)
	p.writeToken(p.getNextToken()) //}
	p.writeBlockTag("whileStatement", true)

}
func doStatment(p *syntaxParser) {
	p.writeBlockTag("doStatement", false)
	p.writeToken(p.getNextToken()) // do
	ParseSubRoutineCall(p)
	p.writeToken(p.getNextToken()) // ;
	p.writeBlockTag("doStatement", true)

}
func returnStatment(p *syntaxParser) {
	p.writeBlockTag("returnStatement", false)
	p.writeToken(p.getNextToken()) // return
	_, token := p.lookahead(1)
	if token != ";" {
		ParseExpression(p)
		_, token = p.lookahead(1)
	}
	p.writeToken(p.getNextToken()) // ;
	p.writeBlockTag("returnStatement", true)

}

func ParseExpression(p *syntaxParser) {
	p.writeBlockTag("expression", false)
	ParseTerm(p)
	_, token := p.lookahead(1)
	for strings.Contains("+-*/|=", token) || token == "&amp;" || token == "&gt;" || token == "&lt;" {
		p.writeToken(p.getNextToken())
		ParseTerm(p)
		_, token = p.lookahead(1)
	}
	p.writeBlockTag("expression", true)
}

func ParseTerm(p *syntaxParser) {
	p.writeBlockTag("term", false)
	tType, token := p.lookahead(1)
	switch tType {
	case tokenTypeMap[integerConstant], tokenTypeMap[stringConstant], tokenTypeMap[keyword]:
		p.writeToken(p.getNextToken())
	case tokenTypeMap[symbol]:
		if token == "-" || token == "~" {
			p.writeToken(p.getNextToken())
			ParseTerm(p)
		} else { // (
			p.writeToken(p.getNextToken())
			ParseExpression(p)
			p.writeToken(p.getNextToken()) // )
		}
	case tokenTypeMap[identifier]:
		_, token1 := p.lookahead(2)
		if token1 == "(" || token1 == "." {
			ParseSubRoutineCall(p)
		} else if token1 == "[" {
			p.writeToken(p.getNextToken())
			p.writeToken(p.getNextToken())
			ParseExpression(p)
			p.writeToken(p.getNextToken())
		} else {
			p.writeToken(p.getNextToken())
		}
	}
	p.writeBlockTag("term", true)
}

func ParseSubRoutineCall(p *syntaxParser) {
	//p.writeBlockTag("subRoutineCall", false)
	p.writeToken(p.getNextToken())
	_, token := p.lookahead(1)
	p.writeToken(p.getNextToken())
	if token == "(" {
		ParseExpressionList(p)
		p.writeToken(p.getNextToken()) // )
	} else { // .
		p.writeToken(p.getNextToken()) // routineName
		p.writeToken(p.getNextToken()) // (
		ParseExpressionList(p)
		p.writeToken(p.getNextToken()) // )
	}
	//p.writeBlockTag("subRoutineCall", true)
}
func ParseExpressionList(p *syntaxParser) {
	p.writeBlockTag("expressionList", false)
	_, token := p.lookahead(1)
	if token != ")" {
		ParseExpression(p)
		_, token := p.lookahead(1)
		for token == "," {
			p.writeToken(p.getNextToken()) // <symbol> , </symbol>
			ParseExpression(p)
			_, token = p.lookahead(1)
		}
	}
	p.writeBlockTag("expressionList", true)
}

var statmentHandlersMap = map[string]func(*syntaxParser){}
