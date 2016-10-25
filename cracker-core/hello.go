package main

import (
	"fmt"
	"os/exec"
	"sync"
)


func main() {

	//Bash Command
	cmd  := "./encrypt"
	//Arguments to get passed to the command
	args := []string{"-s", "000000", "FFC7C9E5694ABFF7"}

	//Common Channel for the goroutines
	tasks := make(chan *exec.Cmd, 64)

	//Spawning 4 goroutines
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(num int, w *sync.WaitGroup) {
			defer w.Done()
			var (
				out []byte
				err error
			)
			for cmd := range tasks { // this will exit the loop when the channel closes
				out, err = cmd.Output()
				if err != nil {
					fmt.Printf("can't get stdout:", err)
				}
				fmt.Printf("%s", string(out))
			}
		}(i, &wg)
	}
	//Generate Tasks
	for i := 0; i < 100000; i++ {
		tasks <- exec.Command(cmd, args...)
	}
	close(tasks)

	// wait for the workers to finish
	wg.Wait()

	fmt.Println("Done")
}
