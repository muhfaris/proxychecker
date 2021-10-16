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
	"golang.org/x/net/proxy"
)

const (
	timeout15 = 15

	// TestURL is target test url
	TestURL = "https://www.reddit.com"
)

// ProxyResp is response of proxy check
type ProxyResp struct {
	Addr string
	Time float64
}

// ProxyChecker is check proxy ip
func ProxyChecker(addr string, c chan ProxyResp, sock, target string) error {
	start := time.Now()
	var timeout = time.Duration(timeout15 * time.Second)

	var client = &http.Client{}
	switch sock {
	case "sock5":
		// create a socks5 dialer
		dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: dialer.Dial,
			},
		}

	case "sock4":

	default:
		URLProxy := &url.URL{Host: addr}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(URLProxy),
			},
			Timeout: timeout,
		}

	}

	if target == "" {
		target = TestURL
	}

	resp, err := client.Get(target)
	if err != nil {
		c <- ProxyResp{Addr: addr, Time: -1}
		return err
	}

	defer resp.Body.Close()
	diff := time.Now().UnixNano() - start.UnixNano()
	c <- ProxyResp{Addr: addr, Time: float64(diff) / 1e9}
	return nil
}

func checkProxyFromIP(addr, sock, target string) {
	if addr == "" {
		fmt.Printf("error address is empty")
		return
	}

	respChan := make(chan ProxyResp)
	go ProxyChecker(addr, respChan, sock, target)
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

func checkProxyFromURL(url, outfile, sock, target string) error {
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

		go ProxyChecker(addr, respChan, sock, target)
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

func checkProxyFromFile(filename, outfile, sock, target string) {
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

		go ProxyChecker(addr, respChan, sock, target)
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
