package main

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"
)

func printSystem(A [][]float64, b, c []float64, signs []string, minmax string) {
	var bufTarget_function string
	// Вывод целевой функции
	for i := 0; i < len(c); i++ {
		bufTarget_function += fmt.Sprintf("(%v)*x%v", c[i], i+1)
		if i == len(c)-1 {
			bufTarget_function += fmt.Sprintf(" → %v", minmax)
			break
		}
		bufTarget_function += "+"
	}
	fmt.Printf("Z = %v\n", bufTarget_function)
	// Вывод системы ограничений
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
}

func printMatrix(matrix [][]coefficientOfValue) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Printf("(%v)%v ", matrix[i][j].coefficient, matrix[i][j].value)
			if j == len(matrix[i])-1 {
				fmt.Print("\n")
			}
		}
	}
}

func printTableau(tableau [][]float64) {
	vars := []string{}
	for i := 0; i < len(tableau[1])-1; i++ {
		vars = append(vars, fmt.Sprintf("%8s\t", fmt.Sprintf("x%v", i+1)))
	}

	// Проверка переменной на базисность
	basicVars := []string{}
	for i := 0; i < len(tableau)-1; i++ {
		for j := 0; j < len(tableau[1])-1; j++ {
			// если коэффициент равен 1 то мы пробегаемся по другим строкам и смотрим все ли они нули, если да, тогда число базисное
			if tableau[i][j] == 1 {
				isBasic := true
				for k := 1; k < len(tableau)-1; k++ {
					if tableau[k][j] != 0 && tableau[i][j] != tableau[k][j] {
						isBasic = false
					}
				}
				if isBasic {
					basicVars = append(basicVars, fmt.Sprintf("x%v", j+1))
				}
			}
		}
	}

	//Отрисовка таблицы
	fmt.Println("Базис\t" + strings.Join(vars, "") + fmt.Sprintf("%8s\t", "B"))
	for i := 0; i < len(tableau); i++ {
		if i < len(tableau)-1 {
			fmt.Printf("%v\t", basicVars[i])
		}

		if i == len(tableau)-1 {
			fmt.Print("C\t")
		}
		for j := 0; j < len(tableau[i]); j++ {
			fmt.Printf("%8.5g\t", math.Floor(tableau[i][j]*100000)/100000)
		}
		fmt.Println()
	}
}

func printTableauWithQ(tableau [][]float64, arrQ []float64, indexOfQ []int) {
	vars := []string{}
	for i := 0; i < len(tableau[1])-1; i++ {
		vars = append(vars, fmt.Sprintf("%8s\t", fmt.Sprintf("x%v", i+1)))
	}

	// Проверка переменной на базисность
	basicVars := []string{}
	for i := 0; i < len(tableau)-1; i++ {
		for j := 0; j < len(tableau[1])-1; j++ {
			// если коэффициент равен 1 то мы пробегаемся по другим строкам и смотрим все ли они нули, если да, тогда число базисное
			if tableau[i][j] == 1 {
				isBasic := true
				for k := 1; k < len(tableau)-1; k++ {
					if tableau[k][j] != 0 && tableau[i][j] != tableau[k][j] {
						isBasic = false
					}
				}
				if isBasic {
					basicVars = append(basicVars, fmt.Sprintf("x%v", j+1))
				}
			}
		}
	}

	//Отрисовка таблицы
	fmt.Println("Базис\t" + strings.Join(vars, "") + fmt.Sprintf("%8s\t", "B") + fmt.Sprintf("%8s\t", "Q"))

	for i := 0; i < len(tableau); i++ {
		if i < len(tableau)-1 {
			fmt.Printf("%v\t", basicVars[i])
		}

		if i == len(tableau)-1 {
			fmt.Print("C\t")
		}
		for j := 0; j < len(tableau[i]); j++ {
			fmt.Printf("%8.5g\t", math.Floor(tableau[i][j]*100000)/100000)
		}
		if slices.Index(indexOfQ, i) == -1 && i != len(tableau)-1 {
			fmt.Printf("%8s\t\n", "-")
			continue
		}
		if i != len(tableau)-1 && slices.Index(indexOfQ, i) != -1 {
			fmt.Printf("%8.5g\t", math.Floor(arrQ[slices.Index(indexOfQ, i)]*100000)/100000)
		}
		fmt.Println()
	}
}

func preprocess(A [][]float64, b, c []float64, signs []string, minmax string) [][]float64 {
	fmt.Print("\nПриводим систему к каноническому виду\nДля этого умножим на -1 те строки где b знак неравенства >=\n",
		"А также добавим дополнительные переменные, чтобы избавиться от неравенств\n",
		"Исходная система ограничений:\n",
	)
	variablesCount := len(c) + 1
	maxlimLen := 0
	printSystem(A, b, c, signs, minmax)
	// Инициализуем матрицу для красивого вывода
	coefNval := make([][]coefficientOfValue, len(b))

	for i := 0; i < len(b); i++ {
		if len(A[i]) >= maxlimLen {
			maxlimLen = len(A[i])
		}
		coefNval[i] = []coefficientOfValue{}
		// Если b оказалось отрицательным числом
		if signs[i] == ">=" {
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
		for j := 0; j < len(A[i]); j++ {
			coefNval[i] = append(coefNval[i], coefficientOfValue{A[i][j], fmt.Sprintf("x%v", j+1)})
		}
		// Добавление дополнительных переменных и свободного значения
		coefNval[i] = append(coefNval[i],
			coefficientOfValue{1, fmt.Sprintf("x%v", variablesCount)},
			coefficientOfValue{b[i], fmt.Sprintf("B%v", i+1)},
		)
		signs[i] = "="
		variablesCount++
	}

	fmt.Println("Система приведенная к каноническому виду:")
	printMatrix(coefNval)

	fmt.Println("Начальная симплекс таблица:")
	// Инициализуем таблицу/матрицу для передачи фунции simplex на итерации
	simplexTableau := make([][]float64, len(b)+1)
	basisSwing := 0
	// Добавляем в таблицу коэффициенты и базисные переменные
	for i := 0; i < len(simplexTableau)-1; i++ {

		simplexTableau[i] = make([]float64, (maxlimLen + len(b) + 1))
		copy(simplexTableau[i], A[i])
		simplexTableau[i][len(A[i])+basisSwing] = 1
		basisSwing++
		simplexTableau[i][len(simplexTableau[i])-1] = b[i]
	}
	// Добавляем строку с коэффициентами целевой функции
	simplexTableau[len(simplexTableau)-1] = make([]float64, (maxlimLen + len(b) + 1))
	copy(simplexTableau[len(simplexTableau)-1], c)

	printTableau(simplexTableau)

	return simplexTableau
}

func computeDeltas(simplexTableau [][]float64) []float64 {
	bVars_index := []int{}
	// Снова проверка на базисность
	for i := 0; i < len(simplexTableau)-1; i++ {
		for j := 0; j < len(simplexTableau[1])-1; j++ {
			if simplexTableau[i][j] == 1 {
				isBasic := true
				for k := 1; k < len(simplexTableau)-1; k++ {
					if simplexTableau[k][j] != 0 && simplexTableau[i][j] != simplexTableau[k][j] {
						isBasic = false
					}
				}
				if isBasic {
					bVars_index = append(bVars_index, j)
				}
			}
		}
	}

	deltaArr := []float64{}
	for i := 0; i < len(simplexTableau[1]); i++ {
		var delta float64
		for j := 0; j < len(simplexTableau)-1; j++ {

			C := simplexTableau[len(simplexTableau)-1][bVars_index[j]]
			a := simplexTableau[j][i]
			delta += C * a
		}
		delta -= simplexTableau[len(simplexTableau)-1][i]
		deltaArr = append(deltaArr, delta)
	}
	return deltaArr
}

func computeAnswer(matrix [][]float64) {
	bVars_index := []int{}
	// Снова проверка на базисность
	for i := 0; i < len(matrix)-1; i++ {
		for j := 0; j < len(matrix[1])-1; j++ {
			if matrix[i][j] == 1 {
				isBasic := true
				for k := 1; k < len(matrix)-1; k++ {
					if matrix[k][j] != 0 && matrix[i][j] != matrix[k][j] {
						isBasic = false
					}
				}
				if isBasic {
					bVars_index = append(bVars_index, j)
				}
			}
		}
	}

	// Секретный соус для вывода текущего плана x
	fixMatrix := make([][]int, len(bVars_index))
	for i := range bVars_index {
		fixMatrix[i] = []int{}
		fixMatrix[i] = append(fixMatrix[i], bVars_index[i], slices.Index(bVars_index, bVars_index[i]))
	}

	//Если честно без понятие не имею как оно работает, но оно сортирует матрицу в порядке возрастания по первому столбцу
	sort.Slice(fixMatrix, func(i, j int) bool {
		return fixMatrix[i][0] < fixMatrix[j][0]

	})

	currentPlanX := make([]float64, len(matrix[1])-1)

	for i := range bVars_index {
		currentPlanX[fixMatrix[i][0]] = matrix[fixMatrix[i][1]][len(matrix[1])-1]
	}

	fmt.Printf("Текущий план X: %v\n", currentPlanX)

	var targetF float64
	for i := 0; i < len(matrix[1])-1; i++ {
		targetF += matrix[len(matrix)-1][i] * currentPlanX[i]
	}
	fmt.Printf("Целевая функция: ")
	sumforF := ""
	for j := 0; j < len(matrix[1])-1; j++ {
		sumforF += fmt.Sprintf("%v*%.5g", matrix[len(matrix)-1][j], currentPlanX[j])
		if j == len(matrix[1])-2 {
			sumforF += " = "
			break
		}
		sumforF += " + "
	}

	//Вывод ответа
	ansArray := []float64{}
	for i := range matrix[len(matrix)-1] {
		if matrix[len(matrix)-1][i] != 0 {
			ansArray = append(ansArray, currentPlanX[i])
		}
	}

	pretty_ansArray := ""
	for i, val := range ansArray {
		pretty_ansArray += fmt.Sprintf("x%v = %.5g, ", i+1, val)
	}

	fmt.Printf("%s%.5g\n", sumforF, targetF)
	fmt.Printf("Ответ: %sF = %.5g", pretty_ansArray, targetF)
}

func simplex(matrix [][]float64, maxOrmin string) {
	var indexOfMinB, indexOfMinA, indexOfMinD int
	var minB float64 = math.MaxFloat64
	var minA float64
	var solvRow []float64

	// Найти минимальное B
	for i := 0; i < len(matrix)-1; i++ {
		if matrix[i][len(matrix[i])-1] <= minB {
			minB = matrix[i][len(matrix[i])-1]
			indexOfMinB = i
		}
	}

	for minB < 0 {
		minB = math.MaxFloat64
		fmt.Println("В столбце B присутствуют отрицательные значения")
		// Найти минимальное по модулю A
		minA = slices.Min(matrix[indexOfMinB][:(len(matrix[1]) - 1)])
		indexOfMinA = slices.Index(matrix[indexOfMinB], minA)
		//делим строку indexOfMinB на indexOfMinA
		for i := range matrix[indexOfMinB] {
			matrix[indexOfMinB][i] /= minA
		}

		//Из остальных строк вычитаем нашу строку, умношенную на сответствующий элемент строки из которой вычитаем
		solvRow := matrix[indexOfMinB]
		for i := 0; i < len(matrix)-1; i++ {
			if i == indexOfMinB {
				continue
			}
			solveColumnIndexValue := matrix[i][indexOfMinA]
			for j := 0; j < len(matrix[i]); j++ {
				matrix[i][j] -= solvRow[j] * solveColumnIndexValue
			}
		}
		// Найти минимальное B
		for i := 0; i < len(matrix)-1; i++ {
			if matrix[i][len(matrix[i])-1] <= minB {
				minB = matrix[i][len(matrix[i])-1]
				indexOfMinB = i
			}
		}
		// Проверка на решаемость
		isPositiveVars := true
		for _, val := range matrix[indexOfMinB][:len(matrix[indexOfMinB])-1] {
			if val < 0 {
				isPositiveVars = false
			}
		}
		if isPositiveVars {
			fmt.Printf("В строке %v отсутстуют отрицательные значения. Решение задачи не существует", indexOfMinB+1)
			return
		}
	}

	fmt.Printf("Обновленная симлекс таблица, в качестве минимального B и A были взяты B[%v] = %.5g и A[x%v] = %.5g:\n",
		indexOfMinB+1,
		minB,
		indexOfMinA+1,
		minA,
	)
	printTableau(matrix)

	// Вычисление дельты для первой итерации и ее отрисовка
	fmt.Println("Обновляем таблицу, добавляя к нем строку с дельтами")
	printTableau(matrix)
	deltas := computeDeltas(matrix)

	minD := math.MaxFloat64
	// Найти разрешающий столбец (Минимальную дельту)
	for i, val := range deltas {
		if val < minD {
			minD = val
			indexOfMinD = i
		}
	}

	// Проверка на ограниченность функции
	isNegative := true
	for i := range matrix[:len(matrix)-1] {
		if matrix[i][indexOfMinD] >= 0 {
			isNegative = false
		}
	}

	if isNegative {
		fmt.Println("Все значения разрешающего столбца неположительны.\nФункция не ограничена. Оптимальное решение отсутствует.")
		return
	}

	fmt.Print("Δ\t")

	// Проверка на оптимальность
	isOpt := true
	switch maxOrmin {
	case "max":
		for _, val := range deltas {
			if val < 0 && slices.Index(deltas, val) != len(deltas)-1 {
				isOpt = false
			}
			fmt.Printf("%8.5g\t", math.Floor(val*100000)/100000)
		}
	case "min":
		for _, val := range deltas {
			if val > 0 && slices.Index(deltas, val) != len(deltas)-1 {
				isOpt = false
			}
			fmt.Printf("%8.5g\t", math.Floor(val*100000)/100000)
		}
	}

	fmt.Println()

	if isOpt {
		fmt.Println("План оптимален")
		computeAnswer(matrix)
		return
	}

	// Начинаем цикл итерации
	for !isOpt {
		arrQ := []float64{}
		indexOfQ := []int{}

		// Найти разрешающий столбец (Минимальную дельту)
		for i, val := range deltas {
			if val < minD {
				minD = val
				indexOfMinD = i
			}
		}

		// Вычисление Q
		for i := range matrix {
			if i == len(matrix)-1 {
				continue
			}
			B := matrix[i][len(matrix[1])-1]
			a := matrix[i][indexOfMinD]
			if (B/a) <= 0 || a == 0 {
				continue
			}
			indexOfQ = append(indexOfQ, i)
			arrQ = append(arrQ, B/a)
		}

		// Нахождение минимального Q
		minQ := math.MaxFloat64
		var indexOfMinQ int
		for i, Q := range arrQ {
			if Q < minQ {
				minQ = Q
				indexOfMinQ = indexOfQ[i]
			}
		}

		// Делим строку indexOfMinQ на indexOfMinD
		solv_Element := matrix[indexOfMinQ][indexOfMinD]
		for i := range matrix[indexOfMinQ] {
			matrix[indexOfMinQ][i] /= solv_Element
		}

		// Из остальных строк вычитаем нашу строку, умношенную на сответствующий элемент строки из которой вычитаем
		solvRow = matrix[indexOfMinQ]
		for i := 0; i < len(matrix)-1; i++ {
			if i == indexOfMinQ {
				continue
			}
			solveColumnIndexValue := matrix[i][indexOfMinD]
			for j := 0; j < len(matrix[i]); j++ {
				matrix[i][j] -= solvRow[j] * solveColumnIndexValue
			}
		}

		// Вычисление дельты для итерации и ее отрисовка
		fmt.Println("Обновляем таблицу, добавляя к ней строку с дельтами")
		printTableauWithQ(matrix, arrQ, indexOfQ)
		deltas = computeDeltas(matrix)

		// Найти разрешающий столбец (Минимальную дельту)
		minD = math.MaxFloat64
		for i, val := range deltas {
			if val < minD {
				minD = val
				indexOfMinD = i
			}
		}

		// Проверка на ограниченность функции
		isNegative := true
		for i := range matrix[:len(matrix)-1] {
			if matrix[i][indexOfMinD] >= 0 {
				isNegative = false
			}
		}

		if isNegative {
			fmt.Println("Все значения разрешаюцего столбца отрицательны.\nФункция не ограничена. Оптимальное решение отсутствует.")
			return
		}

		fmt.Print("Δ\t")
		for i := range deltas {
			fmt.Printf("%8.5g\t", math.Floor(deltas[i]*100000)/100000)
		}

		// Проверка на оптимальность
		isOpt = true
		switch maxOrmin {
		case "max":
			for _, val := range deltas {
				if val < 0 && slices.Index(deltas, val) != len(deltas)-1 {
					isOpt = false
				}
			}
		case "min":
			for _, val := range deltas {
				if val > 0 && slices.Index(deltas, val) != len(deltas)-1 {
					isOpt = false
				}
			}
		}

		fmt.Println()
		if isOpt {
			fmt.Println("План оптимален")
			computeAnswer(matrix)
			return
		}
		fmt.Println("План не оптимален, продолжаем итерацию")
	}
}

type coefficientOfValue struct {
	coefficient float64
	value       string
}

// Пример из методички
func example1() {
	limits := [][]float64{{4, 6}, {2, -5}, {7, 5}, {3, -2}}
	freeColumn := []float64{20, -27, 63, 23}
	targetFunction := []float64{2, 1}
	signs := []string{">=", ">=", "<=", "<="}
	minmax := "max"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)

}

// Пример а
func example2() {
	limits := [][]float64{{1, 2}, {1, 1}, {1, 3}}
	freeColumn := []float64{8, 7, 3}
	targetFunction := []float64{2, 5}
	signs := []string{"<=", "<=", ">="}
	minmax := "max"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)
}

// Пример б
func example3() {
	limits := [][]float64{{1, 3}, {2, 1}, {2, 7}}
	freeColumn := []float64{9, 6, 2}
	targetFunction := []float64{2, 7}
	signs := []string{"<=", "<=", ">="}
	minmax := "min"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)
}

// Пример в
func example4() {
	limits := [][]float64{{1, 2}, {3, 1}, {2, 3}}
	freeColumn := []float64{8, 6, 3}
	targetFunction := []float64{1, 3}
	signs := []string{">=", ">=", ">="}
	minmax := "max"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)
}

// Пример г
func example5() {
	limits := [][]float64{{1, 2}, {2, 1}, {1, 3}}
	freeColumn := []float64{7, 6, 18}
	targetFunction := []float64{9, -3}
	signs := []string{"<=", "<=", ">="}
	minmax := "min"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)
}

// Пример 1
func example6() {
	limits := [][]float64{{4, 6}, {2, -5}, {7, 5}, {3, -2}}
	freeColumn := []float64{20, -27, 63, 23}
	targetFunction := []float64{2, 1}
	signs := []string{">=", ">=", "<=", "<="}
	minmax := "max"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)
}

// Пример 2
func example7() {
	limits := [][]float64{{4, 6}, {2, -5}, {7, 5}, {3, -2}}
	freeColumn := []float64{20, -27, 63, 23}
	targetFunction := []float64{2, 1}
	signs := []string{">=", ">=", "<=", "<="}
	minmax := "min"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)
}

// Пример из видео
func example8() {
	limits := [][]float64{{1, 3, 5, 3}, {2, 6, 1, 0}, {2, 3, 2, 5}}
	freeColumn := []float64{40, 50, 30}
	targetFunction := []float64{7, 8, 6, 5}
	signs := []string{"<=", "<=", "<="}
	minmax := "max"
	simplex(preprocess(limits, freeColumn, targetFunction, signs, minmax), minmax)
}

func main() {
	example8()

}
