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
const perlPkg = "https://perldoc.perl.org/"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

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

func perl_dependencyCheck(libs []string, client *http.Client) {
	for _, val := range libs {
		url := perlPkg + trimQuotes(val)
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
	libs := []string{}
	result := strings.Split(data, "(")
	for i := 1; i < len(result); i++ {
		lib := strings.Split(result[i], "\n")
		//adding lib to new arr
		for _, val := range lib {
			val = strings.TrimSpace(val)
			if len(val) == 0 || val == ")" || val == "require" || val == "replace" {
				continue
			}
			libs = append(libs, strings.TrimSpace(val))
		}
	}
	return libs
}

func perl_parseLib(data string) []string {
	libs := []string{}
	result := strings.Split(data, "=> {")
	result = strings.Split(result[4], "}")
	result = strings.Split(result[0], "\n")
	for _, val := range result {
		val = strings.TrimSpace(val)
		result := strings.Split(val, "=>")
		libs = append(libs, strings.TrimSpace(result[0]))
	}
	//fmt.Println(reflect.TypeOf(result[0]).Kind())
	return libs
}

func perl_getLibs(url string, client *http.Client) []string {
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
	return perl_parseLib(respBody)
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
		-l, --lang programming language
		-v, --verbose  Print verbose logs to stderr.
	`)
}

func main() {
	target := flag.String("u", "", "specify url")
	localFile := flag.String("f", "", "specify path to local module file")
	lang := flag.String("l", "", "specify programming language")
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
		if *lang == "perl" {
			libs := perl_parseLib(string(data))
			perl_dependencyCheck(libs, client)
		} else if *lang == "go" {
			libs := parseLib(string(data))
			dependencyCheck(libs, client)
		}
	} else {
		if *lang == "perl" {
			libs := perl_getLibs(*target, client)
			perl_dependencyCheck(libs, client)
		} else if *lang == "go" {
			libs := getLibs(*target, client)
			dependencyCheck(libs, client)
		}

	}
}
