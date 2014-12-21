package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

//マシン名(hostname)を格納

//性能情報を格納

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}

	ret := make(map[string]interface{})
	//for linux
	ret["created_at"] = time.Now()
	ret["hostname"] = hostname

	keys := make([]string, 0, 100)
	values := make([]string, 0, 100)
	cmd := exec.Command("df", "-m")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Start()

	lines := 0
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		lines++
		s := strings.Fields(scanner.Text())
		switch lines {
		case 1:
			keys = s
		case 2:
			values = s
		}
		//keys = append(keys, scanner.Text())
		//fmt.Println(len(strings.Split(scanner.Text(), ",")))
	}

	for i, _ := range values {
		ret[keys[i]] = values[i]
	}

	cmd.Wait()

	j, err := json.MarshalIndent(ret, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(j))

}
