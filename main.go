package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	MAX_THREADS     = 8
	INPUT_FILE_NAME = "input.txt"
)

var (
	jobs chan int
	wg   *sync.WaitGroup
)

type Matrix struct {
	data [][]int
	lock *sync.Mutex
}

func (m *Matrix) Read(i int, j int) int {
	return m.data[i][j]
}

func (m *Matrix) Write(i int, j int, value int) {
	m.data[i][j] = value
}

func (m *Matrix) Inc(i int, j int, value int) {
	m.data[i][j] += value
}

func (m *Matrix) GetDimensions() (i int, j int) {
	i = len(m.data)

	if len(m.data) > 0 {
		j = len(m.data[0])
	}

	return
}

func (ma *Matrix) Print() {
	m, n := ma.GetDimensions()

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%d ", ma.data[i][j])
		}
		fmt.Println()
	}

}

func MatrixMultiRow(A *Matrix, B *Matrix, C *Matrix, i int) {
	defer wg.Done()

	_, n := A.GetDimensions()
	_, p := B.GetDimensions()
	for j := 0; j < p; j++ {
		for k := 0; k < n; k++ {
			Aik := A.Read(i, k)
			Bkj := B.Read(k, j)
			C.Inc(i, j, Aik*Bkj)
		}
	}

	<-jobs
}

func MatrixMulti(A *Matrix, B *Matrix, C *Matrix) {
	jobs = make(chan int, MAX_THREADS)
	wg = &sync.WaitGroup{}
	m, _ := C.GetDimensions()
	for i := 0; i < m; i++ {
		jobs <- i
		wg.Add(1)
		go MatrixMultiRow(A, B, C, i)
	}

	wg.Wait()
}

func NewMatrix(m [][]int) *Matrix {
	return &Matrix{m, &sync.Mutex{}}
}

func ParseInputRow(inputRow string) ([]int, error) {
	cells := strings.Split(inputRow, " ")
	var row []int
	for _, cell := range cells {
		item, err := strconv.ParseInt(cell, 10, 0)
		if err != nil {
			return nil, err
		}
		row = append(row, int(item))
	}

	return row, nil
}

func ReadInputMatrix(scanner *bufio.Scanner) ([][]int, error) {
	var matrix [][]int
	for scanner.Scan() {
		inputRow := scanner.Text()
		if inputRow == "" {
			break
		}

		row, err := ParseInputRow(inputRow)
		if err != nil {
			return nil, err
		}
		matrix = append(matrix, row)
	}

	return matrix, nil
}

func ReadInput() ([][]int, [][]int, error) {
	file, err := os.Open(INPUT_FILE_NAME)
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	A, err := ReadInputMatrix(scanner)
	if err != nil {
		return nil, nil, err
	}

	B, err := ReadInputMatrix(scanner)
	if err != nil {
		return nil, nil, err
	}

	return A, B, nil
}

func NewResultMatrix(A *Matrix, B *Matrix) (*Matrix, error) {
	m, An := A.GetDimensions()
	Bn, p := B.GetDimensions()
	if An != Bn {
		return nil, fmt.Errorf("incompatible input matrices")
	}

	c := make([][]int, m)
	for i := range c {
		c[i] = make([]int, p)
	}
	C := NewMatrix(c)

	return C, nil
}

func main() {
	inputA, inputB, err := ReadInput()
	if err != nil {
		log.Panic(err)
	}

	A := NewMatrix(inputA)
	B := NewMatrix(inputB)

	C, err := NewResultMatrix(A, B)
	if err != nil {
		log.Panic(err)
	}

	MatrixMulti(A, B, C)

	C.Print()
}
