package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

var threads = 0;
const maxthreads = 30;
var layers [][]string;
var currentLayer = 0;
var disallowed []string;
const useragent = "decsentCrawler";

func isIllegal(url string) (isDisallowed bool){
	for i := 0; i < len(disallowed); i++ {
			if strings.Contains(disallowed[i], url) {
			return false;
		}
	}
	return true;
}

func visit(url string){
	if isIllegal(url){
		fmt.Printf("[visit] Not allowed to visit %s\n", url);
		return;
	}
	result, err := http.Get(url);
	if err != nil{
		fmt.Printf("[visit] !!! Failed getting %s !!!\n", url);
		return;
	}
	tkn := html.NewTokenizer(result.Body);
	for{
		tt := tkn.Next();

		if tt == html.ErrorToken {
			if tkn.Err() == io.EOF {
                return
            }
			fmt.Printf("[visit] !!! Failed parsing %s !!!\n", url);
            return
		}
	}
}

func parseRobots(robotsFile string, source string){
	lines := strings.Split(robotsFile, "\n");
	fmt.Printf("[parseRobots] parsing robots for %s\n", source);
	applies := true;
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "User-agent:"){
			ua := strings.ReplaceAll(lines[i], "User-agent:", "");
			ua = strings.ReplaceAll(ua, " ", "");

			if ua == "*"{
				applies = true;
			}else{
				applies = strings.Contains(ua, useragent);
			}
			fmt.Printf("[parseRobots] Rules for %s\n", ua);
			if !applies{
				fmt.Println("[parseRobots] Doesnt apply.");
			}			
		} else if strings.HasPrefix(lines[i], "Allow:") && applies{
			url := strings.ReplaceAll(lines[i], "Allow:", "");
			url = strings.ReplaceAll(url, " ", "");
			url = source + url;
			
			if currentLayer+1 >= 0 && currentLayer+1 < len(layers) {
				layers[currentLayer+1] = append(layers[currentLayer+1], url)
			} else {
				fmt.Printf("[parseRobots] !!! currentLayer+1 is out of range (currentLayer+1=%d, len(layers)=%d) !!!\n[parseRobots] Automatically creating new layer...\n", currentLayer+1, len(layers))
				var empti []string;
				empti = append(empti, url)
				layers = append(layers, empti)
			}
			fmt.Printf("[parseRobots] + %s\n", url);
		}else if strings.HasPrefix(lines[i], "Disallow:") && applies{
			url := strings.ReplaceAll(lines[i], "Disallow:", "");
			url = strings.ReplaceAll(url, " ", "");
			url = source + url;
			fmt.Printf("[parseRobots] - %s\n", url);
			disallowed = append(disallowed, url);
		}
	}
}

func runLayer(layer []string){
	fmt.Printf("[runLayer] Preparing layer %d\n", currentLayer);
	for i := 0; i < len(layer); i++{
		fmt.Printf("[runLayer] Preparing url %s\n", layer[i]);
		// Parse the url
		parsed, err := url.Parse(layer[i]);
		if err != nil {
			fmt.Printf("[runLayer] Failed to parse %s\n", layer[i])
			return;
		}
		
		// Fetch the robots.txt
		robots, err := http.Get(parsed.String() + "/robots.txt");
		var hasRobotsFile = true;
		var sb string;
		if err != nil{
			hasRobotsFile = false;
		}else{
			body, err := ioutil.ReadAll(robots.Body);
			if err != nil {
				log.Fatalln(err);
			}
			sb = string(body);
			if strings.HasPrefix(sb, "<!DOCTYPE html>"){
				hasRobotsFile = false;
			}
		}

		if hasRobotsFile{
			parseRobots(sb, layer[i]);
		}else{
			fmt.Printf("[runLayer] (%s) No robots.txt found\n", layer[i]);
		}
	}
}

func main(){
	fmt.Println("Initialising crawler...");
	lines, err := readLines("links.txt");
	if err != nil{
		log.Fatalln(err);
	}
	fmt.Println("Started crawling...");
	var layer []string;
	for i := 0; i < len(lines); i++{
		fmt.Println(lines[i]);
		layer = append(layer, lines[i]);
	}
	layers = append(layers, layer)
	runLayer(layers[0])
}