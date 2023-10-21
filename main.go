package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ensureHTTP(target string) string {
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		return "http://" + target
	}
	return target
}

func fetchLinks(target string, fileType string) ([]string, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Request failed with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	regexPattern := fmt.Sprintf(`(?i)src=["']([^"']+\.(?:%s)(?:\?[^\s"']+)?)["']`, fileType)
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindAllSubmatch(body, -1)

	var files []string
	for _, match := range matches {
		link := string(match[1])

		if !strings.HasPrefix(link, "http") {
			link = target + link
		}

		files = append(files, link)
	}

	return files, nil
}

func readTargets(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	linksFile := flag.String("l", "", "File containing list of target URLs")
	output := flag.String("o", "output.txt", "Output file name")
	downloadDir := flag.String("d", "", "Directory to download files to")
	flag.Parse()

	if *linksFile == "" {
		log.Fatal("Please specify a file containing target URLs.")
	}

	targets, err := readTargets(*linksFile)
	if err != nil {
		log.Fatalf("Could not read targets: %s", err)
	}

	fileTypes := []string{"js", "json", "env", "jsx"}
	var allLinks []string

	for _, target := range targets {
		target = ensureHTTP(target)
		for _, fileType := range fileTypes {
			links, err := fetchLinks(target, fileType)
			if err != nil {
				log.Printf("Error fetching links from %s: %s", target, err)
				continue
			}
			allLinks = append(allLinks, links...)
		}
	}

	if len(allLinks) == 0 {
		fmt.Println("No files found.")
		return
	}

	file, err := os.Create(*output)
	if err != nil {
		log.Fatalf("Could not create file: %s", err)
	}
	defer file.Close()

	for _, link := range allLinks {
		fmt.Println(link)
		file.WriteString(link + "\n")

		if *downloadDir != "" {
			parsedURL, err := url.Parse(link)
			if err != nil {
				log.Printf("Could not parse URL %s: %s", link, err)
				continue
			}

			filePath := filepath.Join(*downloadDir, filepath.Base(parsedURL.Path))
			err = downloadFile(filePath, link)
			if err != nil {
				log.Printf("Could not download file %s: %s", link, err)
				continue
			}
		}
	}

	fmt.Printf("Files saved to %s.\n", *output)
}
