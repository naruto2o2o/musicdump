package main

import (
	"bytes"
	"fmt"
	"github.com/yoki123/ncmdump"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func processFile(name string, output string) {
	fp, err := os.Open(name)
	if err != nil {
		log.Println(err)
		return
	}
	defer fp.Close()

	if meta, err := ncmdump.DumpMeta(fp); err != nil {
		log.Fatal(err)
	} else {
		if data, err := ncmdump.Dump(fp); err != nil {
			log.Fatal(err)
		} else {
			t := strings.Replace(name, ".ncm", "."+meta.Format, -1)
			output = output + "/" + filepath.Base(t)
			fmt.Println(output)
			if err = ioutil.WriteFile(output, data, 0644); err != nil {
				log.Fatal(err)
			} else {
				if cover, err := ncmdump.DumpCover(fp); err != nil {
					log.Fatal(err)
				} else {
					// tag信息补全
					switch meta.Format {
					case "mp3":
						addMP3Tag(output, cover, &meta)
					case "flac":
						addFLACTag(output, cover, &meta)
					}
				}
			}
		}
	}
}

func main() {
	files := make([]string, 0)

	if 2 > len(os.Args) {
		fmt.Println("Usage : ncmdum inputPath outputPath")
		return
	}

	path := os.Args[1]
	oPath := os.Args[2]

	if path == "" || oPath == "" {
		fmt.Println("Usage : ncmdum inputPath outputPath")
		return
	}

	if info, err := os.Stat(oPath); err != nil {
		log.Fatalf("output Path %s does not exist.", info)
		return
	}

	if info, err := os.Stat(path); err != nil {
		log.Fatalf("Path %s does not exist.", info)
	} else if info.IsDir() {
		filelist, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatalf("Error while reading %s: %s", path, err.Error())
		}
		for _, f := range filelist {
			files = append(files, filepath.Join(path, "./", f.Name()))
		}
	} else {
		files = append(files, path)
	}

	for _, filename := range files {
		if filepath.Ext(filename) == ".ncm" {
			processFile(filename, oPath)
		} else if strings.Contains(filepath.Ext(filename), "qmc") {
			var outInfo bytes.Buffer
			var opathName string
			ext := filepath.Ext(filename)
			baseName := filepath.Base(filename)

			switch ext {
			case ".qmcflac":
				opathName = oPath + "/" + strings.Replace(baseName, ".qmcflac", "", 1) + ".flac"
			}
			cmd := exec.Command("./qmcdump", filename, opathName)

			fmt.Println(opathName)
			cmd.Stdout = &outInfo
			err := cmd.Run()

			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(outInfo.String())
		}
	}
}
