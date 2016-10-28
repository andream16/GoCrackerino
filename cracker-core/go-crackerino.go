package main

import (
	"fmt"
	"os/exec"
	"sync"
	"os"
	"bufio"
	"runtime"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func main() {

	//Bash Command
	cmdEnc  := "./encrypt"
	cmdDec  := "./decrypt"
	//Arguments to get passed to the command
	args := []string{"-s", "C330C9CBD01DFBA0", "E10C65124518DB05"}
	//Array to contain common CipherTexts
	duplicates := make([]string, 0)

	f, err := os.Create("store.txt")
	check(err)
	defer f.Close()

	//Common Channel for the goroutines
	tasks := make(chan *exec.Cmd, 64)

	//Spawning 4 goroutines
	var wg sync.WaitGroup
	cores := MaxParallelism()
	for i := 0; i < cores; i++ {
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
					fmt.Printf("Can't get stdout:", err)
				}
				s:= string(out)
				f.WriteString(cmd.Args[2]+" "+s)
				f.Sync()
			}
		}(i, &wg)
	}
	//Generate Tasks
	for i := 0; i < 100000; i++ {
		key := string(fmt.Sprintf("%06x", i))
		tasks <- exec.Command(cmdEnc, args[0], key, args[1])
		tasks <- exec.Command(cmdDec, args[0], key, args[2])
	}
	close(tasks)
	// wait for the workers to finish
	wg.Wait()
	f.Close()
	fmt.Println("Done Writing Keys")

	r, _ := os.Open("store.txt")
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		text := line[7:]
		if !contains(duplicates, text) {
			duplicates = append(duplicates, text)
		} else {
				t, _:= os.Open("store.txt")
				dupScan := bufio.NewScanner(t)
			        //currLine := dupScan.Text()
				for dupScan.Scan() {
				    currLine   := dupScan.Text()
				    currCipher := currLine[7:23]
				    if( currCipher == text ){
					fmt.Println(currLine)
				    }
				}
				t.Close()
			}
	}
	//Close File
	r.Close()
}
