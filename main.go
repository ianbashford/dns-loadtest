package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"inet.af/netaddr"
)

var hostnames [][]string

func readFileToArray(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	fmt.Printf("%d lines in %s\n", len(lines), path)
	return lines, scanner.Err()
}

func getHostnameList(input int) []string {
	fmt.Printf("returning file %d\n", input)
	return hostnames[input]
}

func generateLoad(threads int, lookups int) {
	var wg sync.WaitGroup
	wg.Add(threads)
	fmt.Println("Running for loop…")

	for i := 0; i < threads; i++ {
		go func(i int) {
			fmt.Printf("starting thread %d \n", i)
			hostnameFile := getHostnameList(i)
			defer wg.Done()
			var timings = []float64{}

			r := &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						Timeout: time.Millisecond * time.Duration(2000),
					}
					return d.DialContext(ctx, network, "192.168.3.10:5353")
				},
			}

			for j := 0; j < lookups; j++ {

				start := time.Now()
				addr, err := r.LookupHost(context.Background(), hostnameFile[j])

				if err != nil || len(addr) == 0 {
					//elapsed := time.Since(start)
					//fmt.Printf("%s in %d ms \n", err.Error(), elapsed.Milliseconds())
					//fmt.Println(len(addr))
					//fmt.Printf("error %s\n", hostnameFile[j])
				} else {
					elapsed := time.Since(start)
					timings = append(timings, float64(elapsed.Milliseconds()))
				}
			}
			var tot float64
			for k := 0; k < len(timings); k++ {
				tot += timings[k]
			}
			fmt.Printf("average lookup took %f ms \n", (tot / float64(len(timings))))

		}(i)
	}
	wg.Wait()
	fmt.Println("Finished for loop")
}

func main() {
	netaddr.IPv4(1, 2, 3, 4)
	var thr int = 8
	var loops int = 900
	var inputFiles = []string{"/home/ian/development/ip-test-data/ip00.txt", "/home/ian/development/ip-test-data/ip01.txt",
		"/home/ian/development/ip-test-data/ip02.txt", "/home/ian/development/ip-test-data/ip03.txt",
		"/home/ian/development/ip-test-data/ip04.txt", "/home/ian/development/ip-test-data/ip05.txt",
		"/home/ian/development/ip-test-data/ip06.txt", "/home/ian/development/ip-test-data/ip07.txt",
		"/home/ian/development/ip-test-data/ip08.txt", "/home/ian/development/ip-test-data/ip09.txt"}
	for _, s := range inputFiles {
		fmt.Printf("reading %s \n", s)
		a, _ := readFileToArray(s)
		hostnames = append(hostnames, a)
	}
	fmt.Printf("hostnames length %d \n", len(hostnames))
	generateLoad(thr, loops)
}
