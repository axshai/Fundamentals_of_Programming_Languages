package main

func main() {
	t := newToknizer("shai", "Main.jack")
	for t.isThereMoreTokens() {
		typ, token := t.nextToken()
		t.writeToken(tokenNameMap[typ], token)
	}
	t.closeToknizer()

}
