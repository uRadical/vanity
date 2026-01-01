package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
	fmt.Printf("\tvanity <pkg> <repo> [output]\n")
	fmt.Printf("\tpkg \t - your desired package import path\n")
	fmt.Printf("\trepo \t - URL to repo\n")
	fmt.Printf("\toutput \t - optional path for where to place the generated index.html, current working dir will be used if omitted\n")
}

func validatePackage(pkg string) error {
	if !strings.Contains(pkg, ".") {
		return fmt.Errorf("package path should contain a domain (e.g., example.com/pkg)")
	}
	return nil
}

func validateRepo(repo string) error {
	u, err := url.Parse(repo)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
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

	return err
}

func output(file string) (*os.File, error) {
	return os.Create(file)
}

func main() {
	log.SetOutput(os.Stderr)
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	pkg, repo := os.Args[1], os.Args[2]

	if err := validatePackage(pkg); err != nil {
		log.Fatalf("Error: %v", err)
	}
	if err := validateRepo(repo); err != nil {
		log.Fatalf("Error: %v", err)
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

	if err := os.MkdirAll(filepath.Dir(f), 0755); err != nil {
		log.Fatalf("Error: failed to create output directory\n%v", err)
	}

	out, err := output(f)
	if err != nil {
		log.Fatalf("Error: failed to create output file\n%v", err)
	}
	defer out.Close()

	err = generateHTML(pkg, repo, out)
	if err != nil {
		log.Fatalf("Error: failed to generate HTML\n%v", err)
	}
}
