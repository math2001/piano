package main

import "fmt"

type Frac struct {
	Num int
	Den int
}

func (f Frac) String() string {
	return fmt.Sprintf("%d/%d", f.Num, f.Den)
}

func (f Frac) Multiply(target Frac) Frac {
	num := target.Num * f.Num
	den := target.Den * f.Den

	k := gcd(num, den)
	return Frac{num / k, den / k}
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}
