package main

import "fmt"

type Command int

// all the command types
const (
	cArithmetic Command = iota
	cComp
	cPush
	cPop
	cLabel
	cGoto
	cIf
	cComment
	cErr
)

var currentFile string
var spaces = "(\t|\\s)*"

var labelCounter = 1

// map to check the command type by regex
var cmdRegexMap = map[Command]string{
	cArithmetic: fmt.Sprintf("^%s(add|sub|neg|and|or|not).*", spaces),
	cComp:       fmt.Sprintf("^%s(eq|gt|lt).*", spaces),
	cPush:       fmt.Sprintf("^%spush.*", spaces),
	cPop:        fmt.Sprintf("^%spop.*", spaces),
	cLabel:      fmt.Sprintf("^%slabel:.*", spaces),
	cGoto:       fmt.Sprintf("^%sgoto.*", spaces),
	cIf:         fmt.Sprintf("^%sif.*", spaces),
	cComment:    fmt.Sprintf("^%s//", spaces),
}

// map to Use the appropriate function for the command
var cmdHandlersMap = map[Command]func([]string) string{
	cArithmetic: arithmaticHandler,
	cComp:       compHandler,
	cPush:       pushHandler,
	cPop:        popHandler,
	//cLabel:      ,
	//cGoto:       ,
	//cIf:         ,
	//cComment:     ,
}

//map to translate the segments from VM to hack
var segmentsNameMap = map[string]string{
	"static":   "STATIC",
	"argument": "ARG",
	"local":    "LCL",
	"this":     "THIS",
	"that":     "THAT",
	"pointer0": "THIS",
	"pointer1": "THAT",
	"temp":     "5",
}
