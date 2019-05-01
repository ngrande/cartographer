package main

import (
	"io"
	"io/ioutil"
	"bufio"
	"net/http"
	"log"
	"os"
	"os/user"
	"fmt"
	"strings"
	"path/filepath"
	"flag"

	"github.com/ngrande/cartographer/convert"
)

type handler_data struct {
		content string
		exec func(http.ResponseWriter, *http.Request, string)
}

type custom_handler struct {
	mux map[string]handler_data
}

func (h * custom_handler) errorHandler(writer http.ResponseWriter, req *http.Request, status int) {
	writer.WriteHeader(status)
	fmt.Fprintf(writer, "Error: %d", status)
}

func (h *custom_handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if handler, ok := h.mux[req.URL.String()]; ok {
		handler.exec(writer, req, handler.content)
		return
	}

	h.errorHandler(writer, req, http.StatusNotFound)
}

func path_handler(writer http.ResponseWriter, req *http.Request, content string) {
	log.Printf("Serving request for: %s", req.URL.String())
	writer.WriteHeader(http.StatusOK)
	io.WriteString(writer, content)
}

func resolve_path(path string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed getting current user: %v", err)
	}

	if path == "~" {
		return usr.HomeDir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(usr.HomeDir, path[2:])
	} else {
		return path
	}
}

func map_dir(dir string, level string) map[string]handler_data {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed reading dir '%s': %v", dir, err)
	}

	mux := make(map[string]handler_data)
	for _, f := range files {

		if f.IsDir() {
			next := map_dir(filepath.Join(dir, f.Name()), filepath.Join(level, f.Name()))
			for k, v := range next {
				mux[k] = v
			}
			continue
		}

		fpath := filepath.Join(dir, f.Name())
		url := filepath.Join(level, f.Name())

		fdata := ""
		if strings.HasSuffix(f.Name(), "md") {
			converted, err := convert.MarkdownToHTML(fpath)
			if err != nil {
				log.Fatalf("Failed converting markdown to html for file '%s': %v", fpath, err)
			}

			fdata = converted
		} else {
			file, err := os.Open(fpath)
			if err != nil {
				log.Fatalf("Failed reading file '%s': %v", fpath, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var buf strings.Builder
			for scanner.Scan() {
				buf.WriteString(scanner.Text() + "\n")
			}

			fdata = buf.String()
		}

		mux[url] = handler_data{ content: fdata, exec: path_handler }
		if strings.HasPrefix(f.Name(), "index") {
			if _, ok := mux[level]; ok {
				log.Fatalf("Multiple index files detected for level: %s", level)
			}
			mux[level] = handler_data{ content: fdata, exec: path_handler }
		}

	}

	return mux
}

var logdirFlag = flag.String("logdir", "~/log/", "Directory where logfile will be written to")
var addrFlag = flag.String("addr", "0.0.0.0:8080", "Address to listen to(ip:port)")
var dirFlag  = flag.String("dir", "~/cartographer", "root directory which to serve")

func main() {
	flag.Parse()

	logdir := resolve_path(*logdirFlag)
	if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create logfile directory!")
	}

	logfile := filepath.Join(logdir, "cartographer.log")
	file, err := os.OpenFile(logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetPrefix("Server: ")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	dir := resolve_path(*dirFlag)
	addr := *addrFlag

	log.Printf("Address: %s", addr)
	log.Printf("Directory: %s", dir)

	mux := map_dir(dir, "/")
	if _, ok := mux["/"]; !ok {
		log.Fatalf("Failed to get the index file")
	}

	server := http.Server{ Addr: addr, Handler: &custom_handler{mux: mux}, }

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen and serve like a good boi: %v", err)
	}
}
