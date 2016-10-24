package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)


func main() {
	var (
		cmdOut []byte
		err    error
	)

	fmt.Println(strconv.FormatInt(10, 16))

	cmd  := "./encrypt"
	args := []string{"-s", "000000", "FFC7C9E5694ABFF7"}
	if cmdOut, err = exec.Command(cmd, args...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running "+cmd+" "+args[0]+args[1]+args[2], err)
		os.Exit(1)
	}
	sha := string(cmdOut)
	fmt.Println(sha)
}
