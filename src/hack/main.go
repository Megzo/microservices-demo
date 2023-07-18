package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var ttydPortNumber int

func waitForPort(port int) {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			conn.Close()
			fmt.Printf("Port %d is now open.\n", port)
			return
		}
		time.Sleep(time.Second)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	proxyPort := ttydPortNumber
	ttydPortNumber++
	cmd := exec.Command("ttyd", "-o", "--base-path", os.Getenv("BASEPATH"), "-p", strconv.Itoa(proxyPort), "bash")
	log.Printf("Running ttyd and waiting for it to finish while proxying client to: " + "http://localhost:" + strconv.Itoa(proxyPort) + os.Getenv("BASEPATH") + "/")
	//http.Redirect(w, r, "http://"+hostnameWithoutPort(r.Host)+":"+strconv.Itoa(proxyPort)+"/", http.StatusSeeOther)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Check if the tty port is open and wait until it becomes available
	waitForPort(proxyPort)

	// use goroutine waiting, manage process
	// this is important, otherwise the process becomes in S mode
	go func() {
		err = cmd.Wait()
		if err == nil {
			log.Println("ttyd command finished with no error")
		} else {
			log.Printf("ttyd command finished with error: %v", err)
		}
	}()

	targetURL, err := url.Parse("http://localhost:" + strconv.Itoa(proxyPort) + os.Getenv("BASEPATH"))
	if err != nil {
		fmt.Println("Error parsing target URL:", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
	//proxy := httputil.NewSingleHostReverseProxy(targetURL)
	//return func(w http.ResponseWriter, r *http.Request) {
	//	proxy.ServeHTTP(w, r)
	//}
}

func main() {
	ttydPortNumber = 9001
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
