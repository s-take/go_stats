package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const OUTPUT_FILE = "counters.txt"
const DST = ""

func main() {
	content := []byte("\\Processor Information(*)\\% Processor Time")
	ioutil.WriteFile(OUTPUT_FILE, content, os.ModePerm)

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
	//cmd := exec.Command("df", "-m")
	cmd := exec.Command("cmd", "/C typeperf -sc 1 -cf counters.txt")
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
		s := strings.Split(scanner.Text(), "\",")
		switch lines {
		case 2:
			keys = s
		case 3:
			values = s
		}
	}

	re, _ := regexp.Compile("[a-zA-Z0-9]\\\\(.*)")

	for i, _ := range values {
		if i == 0 {
		} else {
			keys[i] = strings.Trim(keys[i], "\"")
			sb := re.FindStringSubmatch(keys[i])
			keys[i] = sb[1]
			values[i] = strings.Trim(values[i], "\"")
			ret[keys[i]] = values[i]
		}
	}

	cmd.Wait()

	j, err := json.MarshalIndent(ret, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(j))

	//httpリクエスト
	values := url.Values{}
	values.Set("json", string(j))

	resp, err := http.PostForm(DST, values)
	if err != nil {
		fmt.Println(err)
		return
	}
}
