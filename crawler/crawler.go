package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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