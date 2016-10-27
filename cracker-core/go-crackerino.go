package main

import (
	"fmt"
	"os/exec"
	"sync"
	"os"
	"bufio"
	"unicode/utf8"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	//Bash Command
	cmdEnc  := "./encrypt"
	cmdDec  := "./decrypt"
	//Arguments to get passed to the command
	args := []string{"-s", "C330C9CBD01DFBA0", "E10C65124518DB05"}
	//Array to contain common CipherTexts
	found := make([]string, 6)

	f, err := os.Create("store.txt")
	check(err)
	defer f.Close()

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
					fmt.Printf("Can't get stdout:", err)
				}
				s:= string(out)
				f.WriteString(cmd.Args[2]+" "+s)
				f.Sync()
			}
		}(i, &wg)
	}
	//Generate Tasks
	for i := 0; i < 1000; i++ {
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
	scanner.Split(bufio.ScanWords)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		if(utf8.RuneCountInString(line) == 6){
		 currLine := line
			for scanner.Scan(){
				searchLine := scanner.Text()
				if(utf8.RuneCountInString(searchLine) == 6){
				 	if(currLine == searchLine){
						found = append(found, currLine)
					}
				}
			}
		}
	}
	fmt.Println(found)
	//Close File
	r.Close()

}
