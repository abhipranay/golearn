package main

import (
	"fmt"
	"log"
	"runtime"
)

type a interface {
	Hey()
}

func r() {
	if err := recover(); err != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		log.Printf("sps_service_panic:\n%s\n", buf)
	}
}

func die() {
	var b a
	b = nil
	b.Hey()
}

func main() {
	maxDiscount := 12600000
	os := 18000000
	as := 22000000

	discountPercentage := float64(maxDiscount) / float64(os)
	fmt.Println(discountPercentage)
	fmt.Println(as)
	// Rebate a percentage of the carrier shipping fee
	rebateAmount := int64(discountPercentage * 100000 * float64(as))/100000
	fmt.Println(rebateAmount)
}
