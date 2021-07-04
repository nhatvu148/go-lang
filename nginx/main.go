package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows"
)

const TH32CS_SNAPPROCESS = 0x00000002

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func main() {
	start := time.Now()
	port := flag.Int("port", 12345, "Port number")
	nginxDir := flag.String("nginxDir", "C:/Users/nhatv/Work/TechnoStar/jmu-dt/bin/nginx", "Nginx directory")

	flag.Parse()

	var wg sync.WaitGroup
	sig := make(chan struct{})
	wg.Add(2)

	go editNginx(*port, fmt.Sprintf("%s/conf/nginx.conf", *nginxDir), &wg, sig)
	go restartNginx(*nginxDir, &wg, sig)

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Editing config took %s", elapsed)
}

func restartNginx(nginxDir string, wg *sync.WaitGroup, sig chan struct{}) {
	defer wg.Done()
	<-sig

	findKillProcess()
	findKillProcess()
}

func findKillProcess() {
	procs, err := processes()
	if err != nil {
		log.Fatal(err)
	}
	nginx_dt := findProcessByName(procs, "nginx_dt.exe")
	if nginx_dt != nil {
		// found it
		KillProcess(nginx_dt.ProcessID)
	}
}

func editNginx(port int, confFile string, wg *sync.WaitGroup, sig chan struct{}) {
	defer wg.Done()
	defer close(sig)
	file, err := os.Create(confFile)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)

	lines := []string{
		"", "#user  nobody;", "worker_processes  1;", "", "#error_log  logs/error.log;", "#error_log  logs/error.log  notice;",
		"#error_log  logs/error.log  info;", "", "#pid        logs/nginx.pid;", "", "", "events {", "    worker_connections  1024;", "}", "", "",
		"http {", "    include       mime.types;", "    default_type  application/octet-stream;", "", "    #log_format  main  '$remote_addr - $remote_user [$time_local] \"$request\" '",
		"    #                  '$status $body_bytes_sent \"$http_referer\" '", "    #                  '\"$http_user_agent\" \"$http_x_forwarded_for\"';",
		"", "    #access_log  logs/access.log  main;", "", "    sendfile        on;", "    #tcp_nopush     on;", "", "    #keepalive_timeout  0;",
		"    keepalive_timeout  65;", "", "    #gzip  on;", "", "    upstream backend-server {", "        server localhost:4000;",
		"    }", "", "    server {", fmt.Sprintf("        listen       %d;", port), "        server_name  localhost;", "", "        location / {",
		"            root   ../../client/build;", "            index  index.html;", "", "            try_files $uri /index.html;",
		"        }", "", "        location /api/ {", "            proxy_pass http://backend-server;", "        }", "",
		"        error_page  404              /404.html;", "", "        # redirect server error pages to the static page /50x.html",
		"        #", "        error_page   500 502 503 504  /50x.html;", "        location = /50x.html {", "            root   html;",
		"        }", "    }", "", "    server {", "        listen       5005;", "        server_name  localhost;", "",
		"        location / {", "            root ../../cug_viewer/dist/example-cug-viewer;", "		    index index.html;", "",
		"            try_files $uri /index.html;", "        }", "    }", "}",
	}

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
		// fmt.Printf("Bytes Written: %d\n", bytesWritten)
		// fmt.Printf("Available: %d\n", writer.Available())
		// fmt.Printf("Buffered : %d\n", writer.Buffered())
	}
	writer.Flush()
}

func processes() ([]WindowsProcess, error) {
	handle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = windows.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func findProcessByName(processes []WindowsProcess, name string) *WindowsProcess {
	for _, p := range processes {
		if strings.ToLower(p.Exe) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func newWindowsProcess(e *windows.ProcessEntry32) WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func KillProcess(pid int) {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println(err)
	}
	for _, p := range processes {
		if p.Pid == int32(pid) {
			fmt.Println(p.Pid)
			p.Kill()
			// p.Terminate()
		}
	}
}
