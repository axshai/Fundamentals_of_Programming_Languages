package main

import "fmt"

func main() {
	t := newToknizer("shai", "Main.jack")
	for t.isThereMoreTokens() {
		typ, token := t.nextToken()
		fmt.Println(typ, token)
		if typ != comment && typ != b {
			t.writeToken(tokenTypeMap[typ], translateToken(token))
		}
	}
	t.closeToknizer()

}
