package runoob

import "testing"

func TestRun(t *testing.T) {
	ro := NewRunOOB("066417defb80d038228de76ec581a50a")
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
