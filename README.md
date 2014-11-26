killroutines
============

A small go library to safely signal multiple goroutines at once (i.e. to tell them to terminate).

Not really production code (yet).  Written for [ircstatus](https://github.com/kd5pbo/ircstatus "ircstatus repo").

Contrived Example
-----------------
    package main
    
    import "github.com/kd5pbo/killroutines"
    import "log"
    import "time"
    
    /* Routine that will square a number, or kill it all on a negative number */
    func routine(i int, done *killroutines.K, in chan int, out chan int,
            dead chan int) {
            log.Printf("I am routine %v", i)
            defer log.Printf("Routine %v ending", i)
            defer func() { dead <- i }()
            for {
                    select {
                    /* Get a task */
                    case n := <-in:
                            /* Processing time */
                            time.Sleep(time.Second)
                            /* Signal all the workers to end if the task is bad */
                            if n < 0 {
                                    log.Printf("Got a negative number: %v", n)
                                    /* No fear of closing a closed channel */
                                    done.Signal()
                                    log.Printf("Closed channel")
                                    return
                            }
                            out <- n * n
                    case <-done.Chan(): /* Signal to give up */
                            return
                    }
            }
    }
    
    /* Spawn a bunch of routines, wait, signal their death in one fell swoop */
    func main() {
            /* Struct to notify goroutines to exit */
            done := killroutines.New()
            /* Channel to send in work */
            in := make(chan int)
            /* Channel to get result of goroutines' work */
            out := make(chan int)
            /* Channel to record dead goroutines */
            dead := make(chan int)
    
            /* Make 8 goroutines to kill off */
            for i := 0; i < 8; i++ {
                    log.Printf("Starting routine %v", i)
                    go routine(i, done, in, out, dead)
            }
    
            /* Give them a range of numbers */
            go func() {
                    for _, n := range []int{4, 5, 88, 4, 12, 900, 1, 8, 0, -3, -2,
                            8, 44} {
                            log.Printf("Sending %v for processing", n)
                            in <- n
                    }
            }()
    
            /* Wait for results of work, or done */
            alive := 8
            for alive > 0 {
                    select {
                    case o := <-out: /* Work output */
                            log.Printf("Got result: %v", o)
                    case r := <-dead:
                            log.Printf("Routine %v noted dead", r)
                            alive--
                    }
            }
    }
