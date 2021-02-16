package main

import "testing"

func BenchmarkCreateCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CreateCopy()
	}
}

func BenchmarkCreatePointerPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CreatePointer()
	}
}
