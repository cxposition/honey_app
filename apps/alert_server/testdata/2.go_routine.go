package main

import (
    "fmt"
    "sync"
    "time"
)

type Result struct {
    ID    int
    Value int
    Err   error
}

func worker(id int, out chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    time.Sleep(time.Duration(100+id*10) * time.Millisecond)
    out <- Result{ID: id, Value: id * id}
}

func main() {
    const n = 10
    results := make(chan Result, n)
    var wg sync.WaitGroup
    wg.Add(n)
    for i := 1; i <= n; i++ {
        go worker(i, results, &wg)
    }
    wg.Wait()
    close(results)

    var sum int
    var list []Result
    for r := range results {
        sum += r.Value
        list = append(list, r)
    }

    fmt.Printf("count=%d sum=%d\n", len(list), sum)
    for _, r := range list {
        fmt.Printf("id=%d value=%d\n", r.ID, r.Value)
    }
}
