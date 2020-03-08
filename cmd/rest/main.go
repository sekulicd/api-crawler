package main

import (
	"api-crawler/core/collegescorecard/collegedomain"
	"api-crawler/interface/rest"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"runtime"
	"time"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

func main() {
	enableMemoryStatistics(time.Second)
	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()
	db.AutoMigrate(&collegedomain.School{})

	service := rest.NewService(db)
	service.StartServer()
}

func enableMemoryStatistics(interval time.Duration) {

	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				printMemoryStatistics()
				printNumOfRoutines()
			}
		}
	}()
}

func toGigabytes(bytes uint64) float64 {
	return float64(bytes) / GIGABYTE
}

func printMemoryStatistics() {
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	bytesAllocated := m2.Alloc
	bytesTotalAllocated := m2.TotalAlloc
	bytesHeapAllocated := m2.HeapAlloc

	fmt.Printf("Allocated: %.3fGB, Total allocated: %.3fGB, Heap allocated: %.3fGB\n", toGigabytes(bytesAllocated), toGigabytes(bytesTotalAllocated), toGigabytes(bytesHeapAllocated))
}

func printNumOfRoutines() {
	fmt.Printf("Num of go routines: %v\n", runtime.NumGoroutine())
}
