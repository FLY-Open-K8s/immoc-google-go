package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// 把十进制转成二进制输出。
func dec2bin(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		num := n % 2
		result = strconv.Itoa(num) + result
	}
	return result
}

func printFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func main() {
	fmt.Println(
		dec2bin(5),
		// 最后保留，是因为后面的)换行了
		dec2bin(13),
	)
	fmt.Println()
	fmt.Println("本源码内容如下")
	printFile("05.循环.go")
}
