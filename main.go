package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"inet.af/netaddr"
)

var hostnames []string

func shuffleArray(originalArray []string) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(originalArray); n > 0; n-- {
		randIndex := r.Intn(n)
		originalArray[n-1], originalArray[randIndex] = originalArray[randIndex], originalArray[n-1]
	}
	return originalArray
}

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
	return shuffleArray(lines), scanner.Err()
}

func getHostnameList(threads int, lookups int, input int) []string {
	startpoint := lookups * input
	//fmt.Printf("startpoint %d endpoint %d\n", startpoint+1, startpoint+lookups)
	return hostnames[startpoint+1 : startpoint+lookups]
}

func getResolver() *net.Resolver {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(2000),
			}
			return d.DialContext(ctx, network, "192.168.3.10:5353")
		},
	}
	return r
}

func doLookup(resolver *net.Resolver, hostname string) (float64, error) {
	start := time.Now()
	addr, err := resolver.LookupHost(context.Background(), hostname)
	if err != nil || len(addr) == 0 {
		return 0, err
	} else {
		return float64(time.Since(start).Milliseconds()), err
	}
}

func generateLoad(threads int, lookups int) {
	var wg sync.WaitGroup
	wg.Add(threads)
	fmt.Println("Running for loop…")

	for i := 0; i < threads; i++ {
		go func(i int) {
			fmt.Printf("starting thread %d \n", i)
			hostnameFile := getHostnameList(threads, lookups, i)
			defer wg.Done()
			var timings = []float64{}
			r := getResolver()

			time.Sleep(3 * time.Second)
			for j := 0; j < len(hostnameFile); j++ {
				newTiming, err := doLookup(r, hostnameFile[j])
				if err != nil {
				} else {
					timings = append(timings, newTiming)
				}
			}
			time.Sleep(3 * time.Second)
			var tot float64
			for k := 0; k < len(timings); k++ {
				tot += timings[k]
			}
			fmt.Printf("average lookup took %f ms  for %d lookups\n", (tot / float64(len(timings))), len(timings))

		}(i)
	}
	wg.Wait()
	fmt.Println("Finished for loop")
}

func main() {
	netaddr.IPv4(1, 2, 3, 4)
	var thr int = 8
	var loops int = 11
	var inputFiles = []string{"/home/ian/development/dns-loadtest/testdata/clean_test_data.txt"}
	for _, s := range inputFiles {
		fmt.Printf("reading %s \n", s)
		a, _ := readFileToArray(s)
		hostnames = a
	}
	fmt.Printf("hostnames length %d \n", len(hostnames))
	generateLoad(thr, loops)
}
