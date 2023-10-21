# File Type Fetcher

## Overview
The File Type Fetcher is a Go script designed to extract and collect file links from specified target URLs. The script can be used to identify and list different types of files (e.g., .js, .json, .env, .jsx) present on web pages. The script can also download the identified files to a specified directory. This tool can be useful for bug bounters, pentesters, developers or web administrators to quickly identify and manage files on web servers and web sites.

## Features
- Extract file links from specified URLs.
- Ability to specify the target URLs via a file.
- Supports multiple file types.
- Option to download the found files to a specified directory.
- Writes the file links to a specified output file and prints them to the command line.

## Usage

### Prerequisites
- [Go](https://golang.org/dl/) installed on your machine.

### Installation
Clone this repository to your local machine:
```bash
git clone https://github.com/huseyinstif/FileTypeFetcher.git
cd FileLinkFetcher
```

### Running the Script
1. First, create a text file (e.g., `targets.txt`) containing the target URLs, one per line.
2. Then, run the script using the following command:
```bash
go run main.go -l targets.txt -o output.txt
```

This will extract the file links from the specified URLs and write them to output.txt.

### Downloading Files
To download the identified files to a specified directory, use the -d flag followed by the directory path:
```bash
go run main.go -l targets.txt -o output.txt -d downloads
```
This will download the files to the downloads directory.

### Flags
-l: Path to the file containing the list of target URLs (Required). <br>
-o: Output file name for saving the file links (Default: output.txt).<br>
-d: Directory to download files to. If not specified, files will not be downloaded.<br>

### Note
- Ensure that the target URLs in the input file start with `http://` or `https://`.
- The script does not create directories, so ensure that the specified download directory exists before running the script with the `-d` flag.

### Buy Me A Coffee
https://www.buymeacoffee.com/huseyintintas
