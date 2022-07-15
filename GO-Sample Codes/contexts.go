package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {
	wg := &sync.WaitGroup{}
	bg := context.Background()
	ctx, cancel := context.WithCancel(bg)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	wg.Add(1)
	go Dispatcher(ctx, 2, wg)
	for range c {
		fmt.Println("Ending. Received SIGINT")
		cancel()
		break
	}
	wg.Wait()
}

// Dispatcher accepts a context and the number of jobs to dispatch.
// Creates a new cancellable context for each job and starts the Do function
// on another Go Routine. When the dispatcher context dies, all jobs are also killed.
func Dispatcher(ctx context.Context, count int, wg *sync.WaitGroup) {
	defer wg.Done()
	cancelfuncs := make([]context.CancelFunc, 0, count)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Dispatch Terminated. Jobs Left:", count)
			for _, cancel := range cancelfuncs {
				cancel()
			}
			return
		default:
		}
		if count != 0 {
			jobctx, cancel := context.WithCancel(ctx)
			cancelfuncs = append(cancelfuncs, cancel)
			wg.Add(1)
			go Do(jobctx, wg, RandomString(10))
			time.Sleep(time.Second * 1)
			count--
		}
	}
}
func Do(ctx context.Context, wg *sync.WaitGroup, msg string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("TERMINATING JOB", msg)
			return
		default:
		}
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
	}
}

// RandomString generates a random of string of size n.
// Only Alpha-Numeric Sequences are generated by this function.
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}