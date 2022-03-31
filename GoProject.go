// "A1" Arnas Juravicius (18257305) and Oisin McNamara(18237398)

package main

import (
	"fmt"
	"time"
	"sync"
	"runtime"
	"math/rand"
    "os"
)

type Matrix [][]int

func printMat(inM Matrix) {
    for _, i := range inM {
        for _, j := range i {
            fmt.Print(" ", j)
        }
        fmt.Println()
    }
}

func rowCount(inM Matrix) int {
    return(len(inM))
}

func colCount(inM Matrix) int {
    return(len(inM[0]))
}

func goCompare(inA Matrix, inB Matrix) int {
    var i, j int
    m := rowCount(inA)
    n := colCount(inA)
    errorCount := 0
 
    for i = 0; i < m; i++ {
       for j = 0; j < n; j++ {
          if (inA[i][j] != inB[i][j]) {
             errorCount++
          }
       }
    }
    return errorCount
}

func newMatrix(r, c int) Matrix {
    a := make([]int, c*r)
    m := make([][]int, r)
    lo, hi := 0, c
    for i := range m {
        m[i] = a[lo:hi:hi]
        lo, hi = hi, hi+c
    }
    return m
}

func transpose(a Matrix) Matrix {
    newArr := make([][]int, len(a[0]))
    for i := 0; i < len(a); i++ {
        for j := 0; j < len(a[0]); j++ {
            newArr[j] = append(newArr[j], a[i][j])
        }
    }
    return newArr
}

func makeArray(rows int, cols int) Matrix{
    newArr := make([][]int ,rows)
	for i := 0; i<rows;i++ {
		newArr[i] = make([]int, cols)
	}
	generate(newArr)
	return newArr
}

func generate(randMatrix Matrix) {
    rand.Seed(time.Now().UnixNano())
    for i, innerArray := range randMatrix {
        for j := range innerArray {
            randMatrix[i][j] = rand.Intn(100)
        }
    }
}

func doCalcSequential(inA Matrix, inB Matrix) Matrix {
    var i, j int
    m := rowCount(inA) 
    p := rowCount(inB)
    q := colCount(inB)
    k := 0	
    total := 0 

    nM := newMatrix(m, q) // create new matrix


    start := time.Now()
    for i = 0; i < m; i++ {
       for j = 0; j < q; j++ {
          for k = 0; k < p; k++ {
             total = total + inA[i][k]*inB[k][j]
          }
          nM[i][j] = total
          total = 0
       }
    }
    elapsed := time.Since(start)
    fmt.Printf("Time taken to calculate %s ",elapsed)
    return nM
}

func doCalc(inA Matrix, inB Matrix) Matrix {
    m := rowCount(inA)     // number of rows the first matrix
    q := colCount(inB)    // number of columns the second matrix

    nM := newMatrix(m, q)
    start := time.Now()
    var wg sync.WaitGroup
    for i := 0; i < m; i++ {
		wg.Add(1)
		go algo(&inA, &inB, &nM, i, &wg)
    }
	wg.Wait()
    elapsed := time.Since(start)
    fmt.Printf("Time taken to calculate %s ",elapsed)
    return nM
}

func algo(inA, inB, nM *Matrix, i int, wg *sync.WaitGroup) {
	total := 0
	defer wg.Done()
	for j := 0; j < colCount(*inB); j++ {
		for k := 0; k < rowCount(*inB); k++ {
			total = total + (*inA)[i][k]*(*inB)[k][j]
		}
		(*nM)[i][j] = total
		total = 0
	}
}

func doCalc2(inA Matrix, inB Matrix) Matrix {
    start := time.Now()
    inB = transpose(inB)
    m := rowCount(inA)     
    q := rowCount(inB)

    nM := newMatrix(m, q) // create new matrix
    var wg sync.WaitGroup
    for i := 0; i < m; i++ {
        for x := 0; x < q; x++ {
            wg.Add(1)
            go algo2(&inA[i], &inB[x], i, &nM, &wg, x)
        }
    }
	wg.Wait()
    elapsed := time.Since(start)
    fmt.Printf("Time taken to calculate %s ",elapsed)
    return nM
}

func algo2(inARows, inBRows *[]int, i int, nM *Matrix, wg *sync.WaitGroup, x int) {
	total := 0
	defer wg.Done()
    for y := 0; y < len(*inBRows); y++ {
		total = total + (*inARows)[y] * (*inBRows)[y]
	}
    (*nM)[i][x] = total
}

func doCalc3(inA Matrix, inB Matrix) Matrix {
    start := time.Now()
    inB = transpose(inB)

    m := rowCount(inA)    
    q := rowCount(inB)

    nM := newMatrix(m, q) // create new matrix

    var wg sync.WaitGroup
    
    for i := 0; i < q; i++ {
		wg.Add(1)
		go algo3(&inA, &inB, &nM, i, &wg)
    }
	wg.Wait()
    elapsed := time.Since(start)
    fmt.Printf("Time taken to calculate %s ",elapsed)
    return nM
}

func algo3(inA, inB, nM *Matrix, i int, wg *sync.WaitGroup) {
	total := 0
	defer wg.Done()
	for j := 0; j < rowCount(*inA); j++ {
		for k := 0; k < colCount(*inA); k++ {
			total = total + (*inB)[i][k]*(*inA)[j][k]
		}
		(*nM)[j][i] = total
		total = 0
	}
}

func main() {
    // Algorithms running of one CPU
    fmt.Println("Max CPUS: ", runtime.NumCPU())
	runtime.GOMAXPROCS(1)
    fmt.Println("Max CPUS: ", runtime.GOMAXPROCS(0))
	a := makeArray(1500, 3000)
	b := makeArray(3000, 2000)

    colsOfA := colCount(a)
    rowsOfB := rowCount(b)

    if colsOfA != rowsOfB {
        fmt.Println("Matrix must be m*n and n*l, matrix are not correct size")
		os.Exit(0)
	}
    
    // Sequential Algorithm
    c := doCalcSequential(a,b)
    fmt.Println()
    fmt.Println()

    // First algorithm rows X matrix
    d := doCalc(a,b)
    fmt.Println()
    same := goCompare(c,d)
    if same == 0 {
        fmt.Println("Matrix d is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()

    // Second algorithm rows X columns
    e := doCalc2(a, b)
    fmt.Println()
    same = goCompare(c,e)
    if same == 0 {
        fmt.Println("Matrix e is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()
    
    // Third algorithm columns X matrix
    f := doCalc3(a, b)
    fmt.Println()
    same = goCompare(c,f)
    if same == 0 {
        fmt.Println("Matrix f is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()


    // Algorithms running of 2 CPUs
	runtime.GOMAXPROCS(2)
    fmt.Println("Max CPUS: ", runtime.GOMAXPROCS(0))
    // First algorithm rows X matrix
    d = doCalc(a,b)
    fmt.Println()
    same = goCompare(c,d)
    if same == 0 {
        fmt.Println("Matrix d is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()

    // Second algorithm rows X columns
    e = doCalc2(a, b)
    fmt.Println()
    same = goCompare(c,e)
    if same == 0 {
        fmt.Println("Matrix e is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()
    
    // Third algorithm columns X matrix
    f = doCalc3(a, b)
    fmt.Println()
    same = goCompare(c,f)
    if same == 0 {
        fmt.Println("Matrix f is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()


    // Algorithms running of max CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())
    fmt.Println("Max CPUS: ", runtime.GOMAXPROCS(0))
    // First algorithm rows X matrix
    d = doCalc(a,b)
    fmt.Println()
    same = goCompare(c,d)
    if same == 0 {
        fmt.Println("Matrix d is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()

    // Second algorithm rows X columns
    e = doCalc2(a, b)
    fmt.Println()
    same = goCompare(c,e)
    if same == 0 {
        fmt.Println("Matrix e is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
    fmt.Println()
    
    // Third algorithm columns X matrix
    f = doCalc3(a, b)
    fmt.Println()
    same = goCompare(c,f)
    if same == 0 {
        fmt.Println("Matrix f is the same as sequential algorithm")
    } else {
        fmt.Println("Matrices are not the same")
    }
}