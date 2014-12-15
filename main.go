package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const OUTPUT_FILE = "counters.txt"
const URL = ""

func main() {
	content := []byte(
		"\\Memory\\Available Bytes\n" +
			"\\Processor(*)\\% Processor Time\n" +
			"\\LogicalDisk(*)\\% Free Space\n",
	)
	ioutil.WriteFile(OUTPUT_FILE, content, os.ModePerm)

	out, err := exec.Command("cmd", "/C typeperf -sc 1 -cf counters.txt").Output()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
}
