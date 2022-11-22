package main

import "fmt"

type Command int

const (
	cArithmetic Command = iota
	cComp
	cPush
	cPop
	cLabel
	cGoto
	cIf
	comment
	err
)

var spaces = "(\t|\\s)*"

var labelCounter = 1

var cmdRegexMap = map[Command]string{
	cArithmetic: fmt.Sprintf("^%s(add|sub|neg|and|or|not).*", spaces),
	cComp:       fmt.Sprintf("^%s(eq|gt|lt).*", spaces),
	cPush:       fmt.Sprintf("^%spush.*", spaces),
	cPop:        fmt.Sprintf("^%spop.*", spaces),
	cLabel:      fmt.Sprintf("^%slabel:.*", spaces),
	cGoto:       fmt.Sprintf("^%sgoto.*", spaces),
	cIf:         fmt.Sprintf("^%sif.*", spaces),
	comment:     fmt.Sprintf("^%s//", spaces),
}

var cmdHandlersMap = map[Command]func([]string) string{
	cArithmetic: arithmaticHandler,
	//cComp:       ,
	cPush: pushHandler,
	cPop:  popHandler,
	//cLabel:      ,
	//cGoto:       ,
	//cIf:         ,
	//comment:     ,
}
