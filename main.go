package main

import (
	"fmt"
)

func printSystem(A [][]float32, b, c []float32, signs []string) {
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(A[i]); j++ {
			fmt.Printf("(%v)*x%v", A[i][j], j+1)
			if j == len(A[i])-1 {
				break
			}
			fmt.Print(" + ")
		}
		fmt.Print(" ", signs[i], " ", b[i], "\n")
	}
	println()
}

func preprocess(A [][]float32, b, c []float32, signs []string) {
	fmt.Print("\nПриводим систему к каноническому виду\nДля этого умножим на -1 те строки где b < 0\n",
		"А также добавим дополнительные переменные, чтобы избавиться от неравенств\n",
		"Исходная система ограничений:\n",
	)
	printSystem(A, b, c, signs)

	for i := 0; i < len(b); i++ {
		if b[i] < 0 {
			b[i] *= -1
			for j := 0; j < len(A[i]); j++ {
				A[i][j] *= -1
			}
			if signs[i] == ">=" {
				signs[i] = "<="
			} else {
				signs[i] = ">="
			}
		}
		A[i] = append(A[i], 1)
	}
	println("Система приведенная к каноническому виду")
	printSystem(A, b, c, signs)
}

func main() {
	A := [][]float32{{4, 6}, {2, -5}, {7, 5}, {3, -2}}
	b := []float32{20, -27, 63, 23}
	c := []float32{2, 1}
	signs := []string{">=", ">=", "<=", "<="}
	// printSystem(A, b, c, signs)
	preprocess(A, b, c, signs)
}
