package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Programm for replacing old sequence with new for files in catalog")

	required := []string{"p", "o", "n"}

	pathPtr := flag.String("p", "", "the path to catalog")
	oldStringPtr := flag.String("o", "", "the old sequence")
	newStringPtr := flag.String("n", "", "the new sequence")
	flag.Parse()
	argsWithProg := os.Args[1:]
	fmt.Println(argsWithProg)

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s argument\n", req)
			fmt.Fprint(os.Stderr,
				"requered flags\n"+
					"-p the path to catalog\n"+
					"-o the old sequence\n"+
					"-n the new sequence\n")
			os.Exit(2)
		}
	}
	result := ReplaceSequenceInCatalog(*pathPtr, *oldStringPtr, *newStringPtr)

	http.HandleFunc("/docker", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Дратути<br>%s</h1>", result)
	})

	http.ListenAndServe(":8080", nil)
}

func ReplaceSequenceInCatalog(path string, oldSequence string, newSequence string) string {
	result := "Files that have been processed"
	fmt.Println("path:", path)
	fmt.Println("oldSequence:", oldSequence)
	fmt.Println("newSequence:", newSequence)

	err := filepath.Walk(path,
		func(innerPath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !fileInfo.IsDir() {
				fmt.Println("Work with file: ", innerPath)
				result += "<br>" + innerPath

				oldContent, err := os.ReadFile(innerPath)
				if err != nil {
					panic(err)
				}
				fmt.Println("Old content of file:", string(oldContent))

				newContent := strings.Replace(string(oldContent), oldSequence, newSequence, -1)
				fmt.Println("New content of file:", newContent)

				err = os.WriteFile(innerPath, []byte(newContent), 0)
				if err != nil {
					panic(err)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return result
}
