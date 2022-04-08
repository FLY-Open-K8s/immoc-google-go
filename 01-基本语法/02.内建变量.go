package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

// 1. 欧拉公式（e的i*pi次幂+1=0）：
func euler() {
	// go语言的内建变量就支持复数
	// 3+4i表示复数
	c := 3 + 4i
	fmt.Printf("Abs(3+4i)=%f\n", cmplx.Abs(c))
	// 底数
	//e^i(Pi)+1= (0+1.2246467991473515e-16i)
	//实部为0，虚部为一个很小的浮点数，复数的实部和虚部都是浮点数。
	fmt.Println("e^i(Pi)+1=", cmplx.Pow(math.E, 1i*math.Pi)+1)
	// E的多少次方
	fmt.Printf("e^i(Pi)+1=%.3f\n", cmplx.Exp(1i*math.Pi)+1)
}

// 2. 强制类型转换
//go语言只有强制类型转换，没有隐式转换。
func triangle() {
	var a, b int = 3, 4
	var c int
	c = int(math.Sqrt(float64(a*a + b*b)))
	fmt.Println(c)
}

func main() {
	euler()
	triangle()
}
