package main

import (
	"fmt"
)

type ManKind interface {
	Say(s string)
	SetMouth(m string)
	GetMouth() string
}
type Man struct {
	ManKind
	mouth string
}

func NewMan() ManKind {
	return &Man{mouth: "M0"}
}
func (this *Man) GetMouth() string {
	return this.mouth
}
func (this *Man) SetMouth(s string) {
	this.mouth = s
}
func (this *Man) Say(s string) {
	fmt.Printf("\n Speak with mouth[%s] : \"%s\"", this.GetMouth(), s)
}

type StrongMan struct {
	Man
}

func NewStrongMan() ManKind {
	sm := &StrongMan{}
	sm.SetMouth("M1")
	return sm
}

func main() {
	NewMan().Say("good luck!")
	&NewStrongMan().Say("good luck!")
}
