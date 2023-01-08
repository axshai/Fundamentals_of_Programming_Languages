package main

import (
	"strings"
)

var className string

func parseClass(p syntaxParser) {
	p.writeBlockTag("class", false)
	p.writeToken(p.getNextToken()) //class
	className = getSecondvalue(p.lookahead(1))
	p.writeToken(p.getNextToken()) // class name
	p.writeToken(p.getNextToken()) // {
	ParseClassVarDec(&p)
	//classScopeTable.printTable()
	ParseSubRoutineDec(&p)
	p.writeToken(p.getNextToken()) // }
	p.writeBlockTag("class", true)
}

func ParseClassVarDec(p *syntaxParser) {
	_, token := p.lookahead(1)
	for token == "static" || token == "field" {
		p.writeBlockTag("classVarDec", false)
		seg := token
		p.writeToken(p.getNextToken()) // field || static
		ttype := getSecondvalue(p.lookahead(1))
		p.writeToken(p.getNextToken()) //<keyword> int </keyword> || <identifier> className </identifier>
		classScopeTable.insert(getSecondvalue(p.lookahead(1)), ttype, seg)
		p.writeToken(p.getNextToken()) //<identifier> x </identifier>
		_, token = p.lookahead(1)
		for token == "," {
			p.writeToken(p.getNextToken()) // <symbol> , </symbol>
			classScopeTable.insert(getSecondvalue(p.lookahead(1)), ttype, seg)
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
		isMethod := token == "method"
		methodScopeTable = newMethodScopeTable(isMethod)
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
		//methodScopeTable.printTable()
	}
}

func ParseParameterList(p *syntaxParser) {
	_, token := p.lookahead(1)
	p.writeBlockTag("parameterList", false)
	for token == "int" || token == "char" || token == "boolean" {
		p.writeToken(p.getNextToken()) //<keyword> type </keyword>
		name := getSecondvalue(p.lookahead(1))
		methodScopeTable.insert(name, token, "argument")
		p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
		_, token = p.lookahead(1)
		for token == "," {
			p.writeToken(p.getNextToken()) // <symbol> , </symbol>
			ttype := getSecondvalue(p.lookahead(1))
			p.writeToken(p.getNextToken()) //<keyword> type </keyword>
			name = getSecondvalue(p.lookahead(1))
			methodScopeTable.insert(name, ttype, "argument")
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
		ttype := getSecondvalue(p.lookahead(1))
		p.writeToken(p.getNextToken()) //<keyword> type </keyword>
		name := getSecondvalue(p.lookahead(1))
		methodScopeTable.insert(name, ttype, "local")
		p.writeToken(p.getNextToken()) //<identifier> varName </identifier>
		_, token = p.lookahead(1)
		for token == "," {
			p.writeToken(p.getNextToken()) // <symbol> , </symbol>
			name = getSecondvalue(p.lookahead(1))
			methodScopeTable.insert(name, ttype, "local")
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
		p.writeToken(p.getNextToken()) // [
		ParseExpression(p)
		p.writeToken(p.getNextToken()) // ]
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
		p.writeToken(p.getNextToken()) // +-*/|=><&
		ParseTerm(p)
		vw.writeArithmeticCmd(token)
		_, token = p.lookahead(1)
	}
	p.writeBlockTag("expression", true)
}

func ParseTerm(p *syntaxParser) {
	p.writeBlockTag("term", false)
	tType, token := p.lookahead(1)
	switch tType {
	case tokenTypeMap[integerConstant], tokenTypeMap[stringConstant], tokenTypeMap[keyword]:
		vw.writeConstantsPushCmd(token, tType)
		p.writeToken(p.getNextToken()) // integerConstant | stringConstant | keywordConstant
	case tokenTypeMap[symbol]:
		if token == "-" || token == "~" {
			p.writeToken(p.getNextToken()) // unary op (-, ~)
			ParseTerm(p)
			vw.writeArithmeticCmd(token)
		} else { // (
			p.writeToken(p.getNextToken()) // (
			ParseExpression(p)
			p.writeToken(p.getNextToken()) // )
		}
	case tokenTypeMap[identifier]:
		_, token1 := p.lookahead(2)
		if token1 == "(" || token1 == "." { // funcName() | className.funcName()
			ParseSubRoutineCall(p)
		} else if token1 == "[" {
			_, tName := p.lookahead(1)
			vw.writePushCmd(methodScopeTable.search(tName).varSeg, methodScopeTable.search(tName).varIndex)
			p.writeToken(p.getNextToken()) // varName
			p.writeToken(p.getNextToken()) // [
			ParseExpression(p)
			vw.writeArithmeticCmd("+")
			vw.writePopCmd("pointer", 1)
			vw.writePushCmd("that", 0)
			p.writeToken(p.getNextToken()) // ]
		} else {
			_, tName := p.lookahead(1)
			p.writeToken(p.getNextToken()) //varName
			vw.writePushCmd(methodScopeTable.search(tName).varSeg, methodScopeTable.search(tName).varIndex)
		}
	}
	p.writeBlockTag("term", true)
}

func ParseSubRoutineCall(p *syntaxParser) {
	//p.writeBlockTag("subRoutineCall", false)
	p.writeToken(p.getNextToken()) // routineName | className
	_, token := p.lookahead(1)
	p.writeToken(p.getNextToken()) // ( | .
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

func getSecondvalue(values ...string) string {
	return values[1]
}
