package main

import "fmt"

type ManKind interface {
	Say(s string)
	GetMouth() string
}
type Man struct {
}

func NewMan() ManKind {
	return &Man{}
}
func (this *Man) GetMouth() string {
	return "M0"
}
func (this *Man) Say(s string) {
	fmt.Printf("\n Speak with mouth[%s] : \"%s\"", this.GetMouth(), s)
}

type StrongMan struct {
	Man
}

func NewStrongMan() ManKind {
	return &StrongMan{}
}
func (this *StrongMan) GetMouth() string {
	return "M1"
}
func main() {
	NewMan().Say("good luck!")
	NewStrongMan().Say("good luck!")
}
