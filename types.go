package main

import "fmt"

type Command int

const (
	cArithmetic Command = iota
	cComp
	cLogic
	cPush
	cPop
	cLabel
	cGoto
	cIf
	comment
	err
)

var spaces = "(\t|\\s)*"

var cmdRegexMap = map[Command]string{
	cArithmetic: fmt.Sprintf("^%s(add|sub|neg).*", spaces),
	cComp:       fmt.Sprintf("^%s(eq|gt|lt).*", spaces),
	cLogic:      fmt.Sprintf("^%s(and|or|not).*", spaces),
	cPush:       fmt.Sprintf("^%spush.*", spaces),
	cPop:        fmt.Sprintf("^%spop.*", spaces),
	cLabel:      fmt.Sprintf("^%slabel:.*", spaces),
	cGoto:       fmt.Sprintf("^%sgoto.*", spaces),
	cIf:         fmt.Sprintf("^%sif.*", spaces),
	comment:     fmt.Sprintf("^%s//", spaces),
}

