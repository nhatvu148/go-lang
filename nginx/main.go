package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	port := flag.Int("port", 52525, "Port number")
	nginxDir := flag.String("nginxDir", "C:/Users/nhatv/Work/TechnoStar/jmu-dt/bin/nginx", "Nginx directory")

	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)

	// go editNginx(*port, fmt.Sprintf("%s/conf/nginx.conf", *nginxDir), &wg)
	go func(port int, nginxDir string, wg *sync.WaitGroup) {
		editNginx(port, fmt.Sprintf("%s/conf/nginx.conf", nginxDir), wg)
	}(*port, *nginxDir, &wg)

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Editing config took %s", elapsed)
}

func editNginx(port int, confFile string, wg *sync.WaitGroup) {
	defer wg.Done()
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
		bytesWritten, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
		fmt.Printf("Bytes Written: %d\n", bytesWritten)
		fmt.Printf("Available: %d\n", writer.Available())
		fmt.Printf("Buffered : %d\n", writer.Buffered())
	}
	writer.Flush()
}
