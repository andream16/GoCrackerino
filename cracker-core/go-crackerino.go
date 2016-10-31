package main

import (
	"fmt"
	"os/exec"
	"os"
	"runtime"
	"sync"
)

//Error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Gets Max number of goroutines to get Called
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
	//Arguments to get passed to the encryption/decryption, change plainText and cipherText with yours
	args := []string{"-s", "yourPlainTex", "yourCipherText"}

	//Create a new file f
	f, err := os.Create("store.txt")
	check(err)
	//Defer Close in make sure we have done writing
	defer f.Close()

	//Common Channel for the goroutines
	tasks := make(chan *exec.Cmd, 64)

	//Create a WaitGroup for the goroutines
	var wg sync.WaitGroup
	//Getting Max number of goroutines
	cores := MaxParallelism()
	//Spawning MAX_Number goroutines
	for i := 0; i < cores; i++ {
		//Add a new gorotuine to the Waitgroup
		wg.Add(1)
		//Call a new Goroutine
		go func(num int, w *sync.WaitGroup) {
			//Until done
			defer w.Done()
			//Initialize a byte container
			var (
				out []byte
				err error
			)
			//Execute a passed Command
			for cmd := range tasks {
				//Get its output
				out, err = cmd.Output()
				//If error
				if err != nil {
					fmt.Printf("Can't get stdout:", err)
				}
				//Cast output to string
				s:= string(out)
				//Write the result on the file: key cipherText
				f.WriteString(cmd.Args[2]+" "+s)
				//Sync R/W on f
				f.Sync()
			}
		}(i, &wg)
	}
	/** Generate as many tasks as we need, 24-bit: 16777216, 28-bit: 268435456, 32-bit: 4294967296
	  *  goroutines will communicate between each other and pick one different task per time
	  *  CPU and Core utilization is managed by Go. The more tasks we have the more goroutines will be used.
	  *  Go scales based on the amount of tasks to do.
	**/
	for i := 0; i < 268435456; i++ {
		//Replace ?? with your preferred padding. 24-bit: 06x, 28-bit: 07x, 32-bit: 08x
		key := string(fmt.Sprintf("%??x", i))
		//Execute Encrypt
		tasks <- exec.Command(cmdEnc, args[0], key, args[1])
		//Execute Decrypt
		tasks <- exec.Command(cmdDec, args[0], key, args[2])
	}
	//When Done, close tasks
	close(tasks)
	// Wait for the workers to finish
	wg.Wait()
	//When all the workers are done, close the file
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
