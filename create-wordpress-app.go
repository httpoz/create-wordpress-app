package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	flag "github.com/ogier/pflag"
)

// flags
var (
	site string
)

func init() {
	flag.StringVarP(&site, "site", "s", "", "The website's domain name.")
}

func main() {
	flag.Parse()

	fmt.Println("Starting app...")

	checkFlags()

	setUpDir()

	downloadAndExtractWordpress()
}

func checkFlags() {
	if flag.NFlag() == 0 {
		fmt.Println("Missing arguments.")
		fmt.Println("Available options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func setUpDir() {
	if _, err := os.Stat(site); err == nil {
		fmt.Println(site + " already exists. I cannot overwrite directories.")
		os.Exit(1)
	} else {
		os.Mkdir(site, os.FileMode(0755))
		fmt.Println("Site folder created.")
	}
}

func downloadAndExtractWordpress() {
	fmt.Println("Downloading WordPress tar ball")
	// Create the file
	out, err := os.Create("latest.tar.gz")
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	resp, err := http.Get("https://wordpress.org/latest.tar.gz")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Open("latest.tar.gz")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Println("Unpacking tar ball...")
	Untar(f, site)

	fmt.Println("Preparing Site")
	os.Rename(site+"/wordpress", site+"/site")

	fmt.Println("Done. Enjoy.")
}
