package main

type SomeStruct struct {
	A, B, C int
	D, E, F string
	G, H, I bool
}

//go:noinline
func CreateCopy() SomeStruct {
	return SomeStruct{
		A: 123, B: 456, C: 789,
		D: "ABC", E: "DEF", F: "HIJ",
		G: true, H: true, I: true,
	}
}
//go:noinline
func CreatePointer() *SomeStruct {
	return &SomeStruct{
		A: 123, B: 456, C: 789,
		D: "ABC", E: "DEF", F: "HIJ",
		G: true, H: true, I: true,
	}
}
