package frac

import "fmt"

// Frac is an immutable number
type Frac struct {
	num int
	den int
}

func F(num int, den int) Frac {
	return Frac{num, den}
}

func (f Frac) String() string {
	return fmt.Sprintf("%d/%d", f.Num(), f.Den())
}

func (f Frac) Float() float64 {
	return float64(f.num) / float64(f.den)
}

func (f Frac) Simplify() Frac {
	// what happens with negative numbers?
	k := gcd(f.Num(), f.Den())
	return Frac{f.Num() / k, f.Den() / k}
}

func (f Frac) Abs() Frac {
	num := f.num
	den := f.den
	if num < 0 {
		num *= -1
	}
	if den < 0 {
		den *= -1
	}
	return Frac{num, den}
}

func (f Frac) Num() int {
	return f.num
}

func (f Frac) Den() int {
	return f.den
}

func (a Frac) Multiply(b Frac) Frac {
	return Frac{
		num: a.Num() * b.Num(),
		den: a.Den() * b.Den(),
	}.Simplify()
}

func (a Frac) Add(b Frac) Frac {
	return Frac{
		num: a.Num()*b.Den() + b.Num()*a.Den(),
		den: a.Den() * b.Den(),
	}.Simplify()
}

func (a Frac) Minus(b Frac) Frac {
	return a.Add(b.Multiply(Frac{-1, 1}))
}

// I just realised that we don't need this function... == works on struct with
// comparable fields
func (a Frac) Equal(b Frac) bool {
	a = a.Simplify()
	b = b.Simplify()
	return a.Num() == b.Num() && a.Den() == b.Den()
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}
