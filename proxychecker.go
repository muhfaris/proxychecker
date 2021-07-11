package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/muhfaris/request"
)

const timeout15 = 15

// ProxyResp is response of proxy check
type ProxyResp struct {
	Addr string
	Time float64
}

// ProxyChecker is check proxy ip
func ProxyChecker(addr string, c chan ProxyResp) error {
	start := time.Now()
	var timeout = time.Duration(timeout15 * time.Second)
	URLProxy := &url.URL{Host: addr}
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(URLProxy)},
		Timeout:   timeout,
	}

	resp, err := client.Get("https://www.reddit.com/")
	if err != nil {
		c <- ProxyResp{Addr: addr, Time: -1}
		return err
	}

	defer resp.Body.Close()
	diff := time.Now().UnixNano() - start.UnixNano()
	c <- ProxyResp{Addr: addr, Time: float64(diff) / 1e9}
	return nil
}

func checkProxyFromIP(addr string) {
	if addr == "" {
		fmt.Printf("error address is empty")
		return
	}

	respChan := make(chan ProxyResp)
	go ProxyChecker(addr, respChan)
	for {
		fmt.Printf("trying to check ip proxy: %s\n", addr)
		r := <-respChan
		if r.Time > 1e-9 {
			fmt.Printf("the proxy is active, with response time %v", r.Time)
			return
		}

		fmt.Printf("proxy %s is not active", r.Addr)
		return
	}
}

func checkProxyFromURL(url, outfile string) error {
	resp := request.Get(
		&request.Config{
			URL: url,
		})

	if resp.Error != nil {
		return resp.Error.Err
	}

	data := strings.Split(strings.TrimSuffix(string(resp.Body), "\n"), "\n")
	respChan := make(chan ProxyResp)
	for _, addr := range data {
		/*
			addrs := strings.SplitN(addr, string(' '), 2)
			ip, port := addrs[0], addrs[1]
		*/

		go ProxyChecker(addr, respChan)
	}

	var proxies []ProxyResp
	for range data {
		r := <-respChan
		if r.Time > 1e-9 {
			proxies = append(proxies, r)
		}
	}

	createfile(proxies, outfile)
	return nil
}

func createfile(data []ProxyResp, outfile string) {
	file, err := os.OpenFile(outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, r := range data {
		_, _ = datawriter.WriteString(fmt.Sprintf("%s - %v \n", r.Addr, r.Time))
	}

	datawriter.Flush()

	fmt.Printf("active proxy save to %s", outfile)
	file.Close()
}

func checkProxyFromFile(filename, outfile string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error read file %s, %v", filename, err)
		return
	}

	data := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	respChan := make(chan ProxyResp)
	for _, addr := range data {
		/*
			addrs := strings.SplitN(addr, string(' '), 2)
			ip, port := addrs[0], addrs[1]
		*/

		go ProxyChecker(addr, respChan)
	}

	var proxies []ProxyResp
	for range data {
		r := <-respChan
		if r.Time > 1e-9 {
			proxies = append(proxies, r)
		}
	}

	createfile(proxies, outfile)
	fmt.Printf("active proxy save to %s", outfile)
}
