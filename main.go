package main

import (
	"fmt"
	"sync"
)

const (
	MAX_THREADS = 8
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
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.data[i][j]
}

func (m *Matrix) Write(i int, j int, value int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[i][j] = value
}

func (m *Matrix) Inc(i int, j int, value int) {
	m.lock.Lock()
	defer m.lock.Unlock()
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

func main() {
	A := NewMatrix([][]int{
		{735, 342, 284, 173, 115},
		{591, 47, 728, 990, 782},
		{630, 662, 946, 123, 163},
		{812, 943, 812, 648, 470},
		{223, 573, 69, 541, 399},
		{113, 32, 770, 735, 399},
		{410, 159, 95, 290, 423},
		{48, 351, 97, 897, 995},
	})
	B := NewMatrix([][]int{
		{635, 189, 179},
		{665, 882, 328},
		{28, 791, 778},
		{414, 13, 882},
		{594, 437, 497},
	})
	C := NewMatrix([][]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	})

	MatrixMulti(A, B, C)

	C.Print()
}
