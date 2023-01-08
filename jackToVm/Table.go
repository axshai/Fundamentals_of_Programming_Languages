package main

type Row struct {
	varType  string
	varSeg   string
	varIndex int
}

type Table map[string]Row

func (t *Table) insert(name string, ttype string, seg string) {
	if _, isTablaContainsVar := (*t)[name]; isTablaContainsVar {
		return
	}
	(*t)[name] = Row{varType: ttype, varSeg: seg, varIndex: t.countSeg(seg)}
}

// func (t Table) printTable() {
// 	columns := []string{"type:", "seg:", "index:"}
// 	fmt.Println("name:", "\t", columns)
// 	for name, row := range t {
// 		fmt.Println(name, "\t", row)
// 	}
// }

func (t Table) search(name string) Row {
	if _, isTablaContainsVar := t[name]; isTablaContainsVar {
		return t[name]
	}
	if classScopeTable[name].varSeg == "static" {
		return classScopeTable[name]
	} else {
		return Row{varType: classScopeTable[name].varType, varSeg: "this", varIndex: classScopeTable[name].varIndex}
	}
}

func (t Table) countSeg(seg string) int {
	counter := 0
	for _, row := range t {
		if row.varSeg == seg {
			counter++
		}
	}
	return counter
}

func newMethodScopeTable(methodType string) Table {
	methodScopeTable := Table{}
	if methodType == "method" {
		methodScopeTable.insert("this", className, "argument")
	}
	return methodScopeTable
}

var classScopeTable Table
var methodScopeTable Table
