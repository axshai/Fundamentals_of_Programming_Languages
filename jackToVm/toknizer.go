package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Toknizer - Responsible for parsing jack file to tokens
type Toknizer struct {
	file         *os.File
	tokensString string
}

// constructor to Toknizer get the jack file to tokinze and the
// file name to where to write the tokens xml file to
func newToknizer(fileName string, jackFile string) Toknizer {
	f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	tokens, _ := os.ReadFile(jackFile)
	// delete the spacess
	tokensString := strings.TrimSpace(string(tokens))
	//the root of the xml file
	f.WriteString("<tokens>\n")
	return Toknizer{file: f, tokensString: string(tokensString)}
}

// function to write token that get the token type and the token value
// and create token like this <type> value </type>
func (t Toknizer) writeToken(tokentType string, token string) {
	strToken := fmt.Sprintf("<%s> %s </%s>", tokentType, token, tokentType)
	t.file.WriteString(strToken + "\n")
}

// function to finish the toknizer
func (t Toknizer) closeToknizer() {
	//the close for the root of the xml
	t.file.WriteString("</tokens>\n")
	//close the file
	t.file.Close()
}

// function that return the type and value of the next token in the jack file
func (t *Toknizer) nextToken() (token, string) {
	// go over all the tokens type and search for match
	for _, key := range []token{multiComment, comment, stringConstant, keyword, symbol, identifier, integerConstant} {
		r, _ := regexp.Compile(tokenRegexMap[key])
		if r.MatchString(t.tokensString) {
			//if you found token delete it from the file and delete the spaces
			token := r.FindString(t.tokensString)
			t.tokensString = t.tokensString[strings.Index(t.tokensString, token)+len(token):]
			t.tokensString = strings.TrimSpace(t.tokensString)
			return key, token
		}
	}
	return err, "error"
}

// function to check if we finish to toknize the file
func (t Toknizer) isThereMoreTokens() bool {
	return t.tokensString != ""
}

// function to translate the token to xml version
func translateToken(token string) string {
	if token == "<" {
		token = "&lt;"
	}
	if token == ">" {
		token = "&gt;"
	}
	if token == "&" {
		token = "&amp;"
	}
	if match, _ := regexp.MatchString(`".*"`, token); match {
		token = strings.ReplaceAll(token, `"`, "")
	}
	return token
}

type token int

// all the tokens type
const (
	keyword token = iota
	symbol
	integerConstant
	identifier
	stringConstant
	comment
	multiComment
	err
)

// map to check the token type type by regex
var tokenRegexMap = map[token]string{
	keyword:         `^(\b(?:class|method|function|constructor|int|boolean|char|void|var|static|field|let|do|if|else|while|return|true|false|null|this)\b)`,
	symbol:          "^({|}|\\[|\\]|\\(|\\)|\\.|,|;|\\+|-|\\*|\\/|&|\\||<|>|=|~)",
	integerConstant: "^([0-9])+",
	identifier:      "^([a-zA-Z_][a-zA-Z_0-9]*)",
	stringConstant:  "^(\"[^\n^\"]*\")",
	comment:         `^(//.*\n)`,
	multiComment:    `(?s)^/\*.*?\*/`,
}

// map to get the string of the token for the xml file
var tokenTypeMap = map[token]string{
	keyword:         "keyword",
	symbol:          "symbol",
	integerConstant: "integerConstant",
	identifier:      "identifier",
	stringConstant:  "stringConstant",
}
