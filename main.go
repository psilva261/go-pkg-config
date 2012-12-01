/*
Copyright 2012 Philip Silva
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	ErrNotFound     = errors.New("pc file could not be found")
	PKG_CONFIG_PATH = os.Getenv("PKG_CONFIG_PATH")
)

func parseConfig(filename string) (cfg map[string]string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	cfg = make(map[string]string)
	vars := make(map[string]string)
	r := bufio.NewReader(f)
	prelude := true
	for {
		lb, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		l := string(lb)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		if strings.Contains(l, ":") {
			prelude = false
		}
		if prelude {
			sl := strings.Split(l, "=")
			if len(sl) < 2 {
				continue
			}
			name := sl[0]
			content := sl[1]
			from := strings.Index(content, "${")
			to := strings.Index(content, "}")
			if from >= 0 && to >= 0 {
				a := content[:from]
				if from+2 < to {
					b := vars[content[from+2:to]]
					c := content[to+1:]
					vars[name] = a + b + c
				}
			} else {
				vars[name] = content
			}
		} else {
			sl := strings.Split(l, ": ")
			if len(sl) < 2 {
				continue
			}
			name := sl[0]
			content := sl[1]
			from := strings.Index(content, "${")
			to := strings.Index(content, "}")
			if from >= 0 && to >= 0 {
				cfg[strings.ToLower(name)] = content[:from] + vars[content[from+2:to]] + content[to+1:]
			} else {
				cfg[strings.ToLower(name)] = content
			}
		}
	}
	return cfg, nil
}

func locatePC(name string) (fullPath string, err error) {
	configPaths := strings.Split(PKG_CONFIG_PATH, ":")
	for _, p := range configPaths {
		fullPath, err = find(p, name+".pc")
		if err == nil {
			return
		}
	}
	return "", ErrNotFound
}

func find(path string, filename string) (targetPath string, err error) {
	dir, err := os.Open(path)
	if err != nil {
		return
	}
	defer dir.Close()
	fis, err := dir.Readdir(-1)
	if err != nil {
		return
	}
	for _, fi := range fis {
		if fi.Name() == filename {
			return path + "/" + filename, nil
		} else if fi.IsDir() {
			targetPath, err = find(path+"/"+fi.Name(), filename)
			if err == nil {
				return
			}
		}
	}
	return "", ErrNotFound
}

func main() {
	var target string
	params := make([]string, 0, 3)
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--") {
			params = append(params, arg[2:])
		} else {
			target = arg
		}
	}
	if target == "" {
		log.Fatalf("no target specified")
	}

	fn, err := locatePC(target)
	if err != nil {
		log.Fatalf("locate pc: %v", err)
	}
	cfg, err := parseConfig(fn)
	if err != nil {
		log.Fatalf("parsing: %v", err)
	}
	for _, p := range params {
		fmt.Print(cfg[p] + " ")
	}
	fmt.Printf("\n")
}
