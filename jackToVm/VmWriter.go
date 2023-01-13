package main

import (
	"fmt"
	"os"
	"strconv"
)

type MethodDetails struct {
	name      string
	localsNum int
	mthodType string
}

var vw VmWriter

// VmWriter - Responsible for transltaing jack file to VM
type VmWriter struct {
	file          *os.File
	labelCounter  int
	currentMethod MethodDetails
}

func newVmWriter(fileName string) VmWriter {
	f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return VmWriter{file: f, labelCounter: 3, currentMethod: MethodDetails{}}
}

func (v VmWriter) writeConstantsPushCmd(token string, constType string) {
	//tokenTypeMap[integerConstant], tokenTypeMap[stringConstant], tokenTypeMap[keyword]
	switch constType {
	case tokenTypeMap[integerConstant]:
		index, _ := strconv.Atoi(token)
		v.writePushCmd("constant", index)
	case tokenTypeMap[stringConstant]:
		v.writePushCmd("constant", len(token))
		v.writeCallCmd("String.new", 1)
		for _, c := range token {
			v.writePushCmd("constant", int(c))
			v.writeCallCmd("String.appendChar", 2)
		}
	case tokenTypeMap[keyword]:
		if token == "this" {
			v.writePushCmd("pointer", 0)
		} else {
			v.writePushCmd("constant", keywordConstMap[token])
		}
	}

}

func (v VmWriter) writeCallCmd(funcNmae string, numOfArgs int) {
	v.file.WriteString(fmt.Sprintf("call %s %d\n", funcNmae, numOfArgs))
}

func (v VmWriter) writeArithmeticCmd(jackOp string, isUnary bool) {
	if isUnary && jackOp == "-" {
		v.file.WriteString("neg" + "\n")
	} else {
		v.file.WriteString(opMap[jackOp] + "\n")
	}
}

func (v VmWriter) writePushCmd(seg string, index int) {
	v.file.WriteString(fmt.Sprintf("push %s %d\n", seg, index))
}
func (v VmWriter) writePopCmd(seg string, index int) {
	v.file.WriteString(fmt.Sprintf("pop %s %d\n", seg, index))
}

func (v VmWriter) writeGoTo(label string) {
	v.file.WriteString(fmt.Sprintf("goto %s\n", label))
}

func (v VmWriter) writeIfGoTo(label string) {
	v.file.WriteString(fmt.Sprintf("if-goto %s\n", label))
}

func (v VmWriter) writeLabel(label string) {
	v.file.WriteString(fmt.Sprintf("label %s\n", label))
}

func (v VmWriter) writeReturn() {
	// if vw.currentMethod.mthodType == "constructor" {
	// 	v.writePushCmd("pointer", 0)
	// }
	v.file.WriteString("return\n")
}

func (v VmWriter) writeFuncDec() {
	v.file.WriteString(fmt.Sprintf("function %s.%s %d\n", className, v.currentMethod.name, v.currentMethod.localsNum))
}

func (v *VmWriter) generateLabelSofix(labels ...string) []string {
	for i := range labels {
		labels[i] += strconv.Itoa(v.labelCounter)
	}
	v.labelCounter++
	return labels
}

func (v VmWriter) closeVmWriter() {
	v.file.Close()
}

var opMap = map[string]string{
	"+":     "add",
	"-":     "sub",
	"*":     "call Math.multiply 2",
	"/":     "call Math.divide 2",
	"&amp;": "and",
	"|":     "or",
	"&lt;":  "lt",
	"&gt;":  "gt",
	"=":     "eq",
	"~":     "not",
}

var keywordConstMap = map[string]int{
	"true":  -1,
	"false": 0,
	"null":  0,
	//"this":  "pointer 0",
}
