package main

import (
	"fmt"
	"os/exec"
	"os"
	"runtime"
	"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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

	/** Bash Commands **/
	cmdEnc  := "./encrypt"
	cmdDec  := "./decrypt"
	awkCmd := "awk"
	akwArgs := []string{"x[$2]++ == 1 { print $2 }", "store.txt"}
	grepCmd := "grep"
	grepArg := "store.txt"
	//Arguments to get passed to the encryption/decryption
	args := []string{"-s", "C330C9CBD01DFBA0", "E10C65124518DB05"}

	f, err := os.Create("store.txt")
	check(err)
	defer f.Close()

	//Common Channel for the goroutines
	tasks := make(chan *exec.Cmd, 64)

	//Spawning 8 goroutines
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
			for cmd := range tasks {
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
	for i := 0; i < 16777215; i++ {
		key := string(fmt.Sprintf("%06x", i))
		tasks <- exec.Command(cmdEnc, args[0], key, args[1])
		tasks <- exec.Command(cmdDec, args[0], key, args[2])
	}
	close(tasks)
	// Wait for the workers to finish
	wg.Wait()
	f.Close()
	fmt.Println("Done Writing Keys")

	//Execute Awk Command to find same CipherTexts By second Column in the File
	var (
		akwOut []byte
		awkErr error
	)
	akwCommand := exec.Command(awkCmd, akwArgs[0], akwArgs[1])
	akwOut, awkErr = akwCommand.Output()
	if awkErr != nil {
		fmt.Printf("Can't get stdout:", awkErr)
	}
	awkS := string(akwOut[0:16])
	fmt.Println(awkS)

	//Executing Grep Command to find all the lines containing the Common CipherText
	var (
		grepOut []byte
		grepErr error
	)
	grepCommand := exec.Command(grepCmd, awkS, grepArg)
	grepOut, grepErr = grepCommand.Output()
	if grepErr != nil {
		fmt.Printf("Can't get stdout:", grepErr)
	}
	grepS := string(grepOut)
	fmt.Println(grepS)

}
