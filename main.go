package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const goPkg = "https://pkg.go.dev/"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"

func dependencyCheck(libs []string, client *http.Client) {
	for _, val := range libs {
		url := goPkg + val
		pkgUrl := strings.Split(url, " ")
		resp, err := client.Get(pkgUrl[0])
		if err != nil {
			log.Fatalln(err)
		}
		status := resp.StatusCode
		if status == 200 {
			fmt.Println(url, string(colorGreen), status)
		} else {
			fmt.Println(url, string(colorRed), status)
		}

	}
}

func parseLib(data string) []string {
	result := strings.Split(data, "(")
	lib := strings.Split(result[1], "\n")
	//adding lib to new arr
	libs := []string{}
	for _, val := range lib {
		val = strings.TrimSpace(val)
		if len(val) == 0 || val == ")" {
			continue
		}
		libs = append(libs, strings.TrimSpace(val))
	}
	return libs
}

func getLibs(url string, client *http.Client) []string {
	//parse github raw string
	//ex: https://raw.githubusercontent.com/<project>
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	// read resp body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	respBody := string(body)

	return parseLib(respBody)
}

func usage() {
	fmt.Fprint(os.Stderr, `Usage: dependency-check [flag] [URL/localFile]
	DependencyConfusion is tool used for finding any library used by the project that might be vulnerable to dependency confusion attack.
	Project with following languages supported:
	- Golang
	- python
	- c/c++

	Flags:
		-u, --url  provide github raw url
		-f, --file path to local module file 
		-v, --verbose  Print verbose logs to stderr.
	`)
}

func main() {
	target := flag.String("u", "", "specify url")
	localFile := flag.String("f", "", "specify path to local module file")
	flag.Parse()
	if len(*target) == 0 && len(*localFile) == 0 {
		usage()
		os.Exit(1)
	}
	//os.Setenv("http_proxy", "http://127.0.0.1:8080")
	//proxy: https://iamninad.com/posts/burp-suite-for-web-app-testing-go-lang/
	//proxyURL, e := url.Parse(os.Getenv("http_proxy"))
	//if e != nil {
	//	panic(e)
	//}
	//client := &http.Client{
	//	//Timeout: 20 * time.Second,
	//	Transport: &http.Transport{
	//		Proxy:           http.ProxyURL(proxyURL),
	//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	},
	//}

	client := &http.Client{}
	if len(*target) == 0 {
		data, err := os.ReadFile(*localFile)
		if err != nil {
			log.Fatalln(err)
		}
		libs := parseLib(string(data))
		dependencyCheck(libs, client)
	} else {
		libs := getLibs(*target, client)
		dependencyCheck(libs, client)
	}
}
