// =============================================================
// Goroutines & Channels — Go's Concurrency Model
// Run: go run cmd/examples/02_concurrency.go
//
// This is THE thing that makes Go special vs PHP.
// PHP: one request = one process (or thread with Swoole)
// Go:  millions of goroutines on a single process
// =============================================================
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// -----------------------------------------------------------
// 1. BASIC GOROUTINE
//    Just add "go" before a function call. That's it.
// -----------------------------------------------------------
func basicGoroutine() {
	fmt.Println("=== BASIC GOROUTINE ===")

	go func() {
		fmt.Println("  hello from goroutine!")
	}()

	// Without this, main() might exit before goroutine runs
	time.Sleep(100 * time.Millisecond)
}

// -----------------------------------------------------------
// 2. WAITGROUP — wait for goroutines to finish
//    Like PHP's pcntl_wait or Promise.all
// -----------------------------------------------------------
func waitGroupExample() {
	fmt.Println("\n=== WAITGROUP ===")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // increment counter

		go func(id int) {
			defer wg.Done() // decrement counter when done

			dur := time.Duration(rand.Intn(500)) * time.Millisecond
			time.Sleep(dur)
			fmt.Printf("  worker %d finished (took %v)\n", id, dur)
		}(i) // ⚠️ IMPORTANT: pass i as arg, don't capture loop variable!
	}

	wg.Wait() // block until counter is 0
	fmt.Println("  all workers done!")
}

// -----------------------------------------------------------
// 3. CHANNELS — goroutines communicate by passing messages
//    Think of it as a typed pipe between goroutines
// -----------------------------------------------------------
func channelBasics() {
	fmt.Println("\n=== CHANNELS ===")

	// Unbuffered channel — sender blocks until receiver is ready
	ch := make(chan string)

	go func() {
		ch <- "hello from goroutine" // send
	}()

	msg := <-ch // receive (blocks until message arrives)
	fmt.Println(" ", msg)

	// Buffered channel — can hold N messages without blocking
	bufCh := make(chan int, 3)
	bufCh <- 1
	bufCh <- 2
	bufCh <- 3
	// bufCh <- 4  // this would block — buffer is full!

	fmt.Println("  buffered:", <-bufCh, <-bufCh, <-bufCh)
}

// -----------------------------------------------------------
// 4. CHANNEL PATTERNS: range + close
//    Producer sends values, consumer reads with range
// -----------------------------------------------------------
func channelRange() {
	fmt.Println("\n=== CHANNEL RANGE ===")

	ch := make(chan int)

	// Producer goroutine
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i * i // send squares
		}
		close(ch) // signal no more values
	}()

	// Consumer — range reads until channel is closed
	for val := range ch {
		fmt.Printf("  received: %d\n", val)
	}
}

// -----------------------------------------------------------
// 5. SELECT — listen on multiple channels simultaneously
//    Like PHP's stream_select but for goroutine channels
// -----------------------------------------------------------
func selectExample() {
	fmt.Println("\n=== SELECT ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "result from service A"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "result from service B"
	}()

	// Wait for whichever comes first
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println(" ", msg)
		case msg := <-ch2:
			fmt.Println(" ", msg)
		}
	}
}

// -----------------------------------------------------------
// 6. WORKER POOL — very common interview pattern!
//    N workers process jobs from a shared channel
// -----------------------------------------------------------
func workerPool() {
	fmt.Println("\n=== WORKER POOL ===")

	const numWorkers = 3
	const numJobs = 8

	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate work
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
				results <- fmt.Sprintf("worker %d processed job %d", id, job)
			}
		}(w)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // no more jobs

	// Wait for all workers, then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for r := range results {
		fmt.Println(" ", r)
	}
}

// -----------------------------------------------------------
// 7. MUTEX — protect shared state
//    When channels aren't the right fit
// -----------------------------------------------------------
func mutexExample() {
	fmt.Println("\n=== MUTEX ===")

	var (
		counter int
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++ // safe with mutex
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("  counter = %d (should be 1000)\n", counter)

	// RWMutex — allows multiple readers OR one writer
	// Use when reads are much more frequent than writes
	var rwmu sync.RWMutex
	data := map[string]int{"a": 1}

	// Many readers can read simultaneously
	rwmu.RLock()
	_ = data["a"]
	rwmu.RUnlock()

	// Only one writer at a time
	rwmu.Lock()
	data["b"] = 2
	rwmu.Unlock()
	fmt.Println("  RWMutex data:", data)
}

// -----------------------------------------------------------
// 8. COMMON TRAP: Loop Variable Capture
//    ⚠️ THIS WILL COME UP IN INTERVIEWS
// -----------------------------------------------------------
func loopTrap() {
	fmt.Println("\n=== LOOP VARIABLE TRAP ===")

	// ❌ WRONG — all goroutines see the SAME variable
	fmt.Println("  WRONG way (all print same value):")
	var wg sync.WaitGroup
	values := []int{1, 2, 3, 4}

	for _, v := range values {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("    %d", v) // likely prints "4 4 4 4"
		}()
	}
	wg.Wait()
	fmt.Println()

	// ✅ RIGHT — pass as argument (creates a copy)
	fmt.Println("  RIGHT way (pass as arg):")
	for _, v := range values {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Printf("    %d", val)
		}(v) // v is copied here
	}
	wg.Wait()
	fmt.Println()

	// ✅ ALSO RIGHT — shadow the variable
	fmt.Println("  ALSO RIGHT (shadow variable):")
	for _, v := range values {
		v := v // shadow — creates new variable for each iteration
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("    %d", v)
		}()
	}
	wg.Wait()
	fmt.Println()

	// NOTE: Go 1.22+ fixed this for `for` loops, but interviewers
	// still ask about it. Know both the problem and the fix.
}

func main() {
	basicGoroutine()
	waitGroupExample()
	channelBasics()
	channelRange()
	selectExample()
	workerPool()
	mutexExample()
	loopTrap()

	fmt.Println("\n✅ All concurrency examples done!")
}
