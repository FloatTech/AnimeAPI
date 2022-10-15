package runoob

import "testing"

func TestRun(t *testing.T) {
	ro := NewRunOOB("b6365362a90ac2ac7098ba52c13e352b")
	r, err := ro.Run(`package main

	import "fmt"
	
	func main() {
	   fmt.Println("Hello, World!aaa")
	}`, "go", "")
	if err != nil {
		t.Fatal(err)
	}
	if r != "Hello, World!aaa\n" {
		t.Fail()
	}
}
