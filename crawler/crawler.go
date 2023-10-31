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

func visit(url string){
	// Visit said url
}

func runLayer(layer []string){
	fmt.Printf("[runLayer] Preparing layer %d\n", currentLayer);
	for i := 0; i < len(layer); i++{
		fmt.Printf("[runLayer] Preparing url %s\n", layer[i]);
		// Parse the url
		parsed, err := url.Parse(layer[i]);
		if err != nil {
			fmt.Println("[runLayer] Failed to parse %s", layer[i])
			return;
		}
		
		// Fetch the robots.txt
		robots, err := http.Get(parsed.String() + "/robots.txt");
		var hasRobotsFile = true;
		if err != nil{
			hasRobotsFile = false;
			fmt.Println("[runLayer] (%s) No robots.txt found", layer[i]);
		}else{
			body, err := ioutil.ReadAll(robots.Body);
			if err != nil {
				log.Fatalln(err);
			}
			sb := string(body);
			if strings.HasPrefix(sb, "<!DOCTYPE html>"){
				hasRobotsFile = false;
			}
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
}