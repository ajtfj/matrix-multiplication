package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	MAX_THREADS       = 128
	BENCHMARK_SAMPLES = 1000
)

var (
	jobs chan int
	wg   *sync.WaitGroup
)

type Matrix struct {
	Data [][]int
	Lock *sync.Mutex
}

func (m *Matrix) Read(i int, j int) int {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	return m.Data[i][j]
}

func (m *Matrix) Write(i int, j int, value int) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Data[i][j] = value
}

func (m *Matrix) Inc(i int, j int, value int) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Data[i][j] += value
}

func (m *Matrix) GetDimensions() (i int, j int) {
	i = len(m.Data)

	if len(m.Data) > 0 {
		j = len(m.Data[0])
	}

	return
}

func (ma *Matrix) Print() {
	m, n := ma.GetDimensions()

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%d ", ma.Data[i][j])
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

func BenchmarkMatrixMulti(A *Matrix, B *Matrix, C *Matrix) {
	accum := 0
	for i := 1; i < BENCHMARK_SAMPLES; i++ {
		start := time.Now()
		MatrixMulti(A, B, C)
		accum += int(time.Since(start))
	}
	average := accum / BENCHMARK_SAMPLES
	fmt.Printf("Benchmark result (%d threads): %dns\n", MAX_THREADS, average)
}

func main() {
	A := &Matrix{
		Data: [][]int{
			{735, 342, 284, 173, 115},
			{591, 47, 728, 990, 782},
			{630, 662, 946, 123, 163},
			{812, 943, 812, 648, 470},
			{223, 573, 69, 541, 399},
			{113, 32, 770, 735, 399},
			{410, 159, 95, 290, 423},
			{48, 351, 97, 897, 995},
		},
		Lock: &sync.Mutex{},
	}
	B := &Matrix{
		Data: [][]int{
			{635, 189, 179},
			{665, 882, 328},
			{28, 791, 778},
			{414, 13, 882},
			{594, 437, 497},
		},
		Lock: &sync.Mutex{},
	}
	C := &Matrix{
		Data: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		Lock: &sync.Mutex{},
	}

	BenchmarkMatrixMulti(A, B, C)
}
