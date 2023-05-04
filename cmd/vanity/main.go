package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

//go:embed template.gohtml
var redirectTemplate string

func usage() {
	fmt.Println("vanity")
	fmt.Printf("\nUsage:\n")
	fmt.Printf("\tFlags:\n")
	fmt.Printf("\t-h\t - display usage information\n")

	fmt.Printf("\n\tArguments:\n")
	fmt.Printf("\tvanity [pkg] [url]\n")
	fmt.Printf("\tpkg \t - your desired package import path\n")
	fmt.Printf("\trepo \t - URL to repo\n")
	fmt.Printf("\toutput \t - path to place the generated index.html\n")
}

func generateHTML(pkg, repo string, out io.Writer) error {
	t, err := template.New("redirectTemplate").Parse(redirectTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(out, struct {
		Package string
		Repo    string
		Name    string
	}{
		pkg,
		repo,
		path.Base(repo),
	})

	if err != nil {
		return err
	}

	return nil
}

func output(file string) io.Writer {
	f, _ := os.Create(file)
	return f
}

func main() {
	log.SetOutput(os.Stderr)
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 3 {
		usage()
		os.Exit(0)
	}

	f, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if len(os.Args) > 3 {
		f = strings.Replace(os.Args[3], "~", os.Getenv("HOME"), -1)
	}

	if !strings.HasSuffix(f, "/") {
		f = f + "/"
	}
	f = f + "index.html"

	err = generateHTML(os.Args[1], os.Args[2], output(f))
	if err != nil {
		log.Fatalf("Error: failed to generate HTML\n%v", err)
	}
}
