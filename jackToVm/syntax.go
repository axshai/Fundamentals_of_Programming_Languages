package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Toknizer struct {
	file         *os.File
	tokensString string
}

func newToknizer(fileName string, jackFile string) Toknizer {
	f, _ := os.OpenFile(fileName+".xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	tokensString, _ := os.ReadFile(jackFile)
	f.WriteString("<tokens>\n")
	return Toknizer{file: f, tokensString: string(tokensString)}
}

func (t Toknizer) writeToken(tokentType string, token string) {
	strToken := fmt.Sprintf("<%s> %s </%s>", tokentType, token, tokentType)
	t.file.WriteString(strToken + "\n")
}

func (t Toknizer) closeToknizer() {
	t.file.WriteString("</tokens>\n")
	t.file.Close()
}

func (t *Toknizer) nextToken() (token, string) {
	for key, element := range tokenRegexMap {
		r, _ := regexp.Compile(element)
		if r.MatchString(t.tokensString) {
			token := r.FindString(t.tokensString)
			t.tokensString = t.tokensString[strings.Index(t.tokensString, token)+len(token):]
			return key, token
		}
	}
	return err, "error"
}

func (t Toknizer) isThereMoreTokens() bool {
	return t.tokensString != ""
}

func Index(s, token string) {
	panic("unimplemented")
}

type token int

const (
	keyword token = iota
	symbol
	integerConstant
	identifier
	stringConstant
	comment
	err
)

// map to check the token type type by regex
var tokenRegexMap = map[token]string{
	keyword:         `\b(?:class|method|function|constructor|int|boolean|char|void|var|static|field|let|do|if|else|while|return|true|false|null|this)\b`,
	symbol:          "{|}|[|]|(|)|.|,|;|+|-|*|/|&|<|>|=|~|[|]",
	integerConstant: "[0-9]+",
	identifier:      "^[a-zA-Z_][a-zA-Z_0-9]*",
	stringConstant:  "\"[^\n^\"]*\"",
	comment:         `//.*\n"|/*.**/`,
}

var tokenNameMap = map[token]string{
	keyword:         "keyword",
	symbol:          "symbol",
	integerConstant: "integerConstant",
	identifier:      "identifier",
	stringConstant:  "stringConstant",
}
