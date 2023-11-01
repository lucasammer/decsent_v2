package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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

type Userinfo struct {
	username    string
	password    string
	passwordSet bool
}

type URL struct {
	Scheme      string
	Opaque      string    // encoded opaque data
	User        *Userinfo // username and password information
	Host        string    // host or host:port
	Path        string    // path (relative paths may omit leading slash)
	RawPath     string    // encoded path hint (see EscapedPath method)
	ForceQuery  bool      // append a query ('?') even if RawQuery is empty
	RawQuery    string    // encoded query values, without '?'
	Fragment    string    // fragment for references, without '#'
	RawFragment string    // encoded fragment hint (see EscapedFragment method)
}

func visit(url *URL){
	// Visit said url
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