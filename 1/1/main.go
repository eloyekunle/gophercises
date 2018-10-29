package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	csvFilePtr := flag.String("csv", "../problems.csv", "A CSV file in the format of 'question,answer'")
	flag.Parse()

	csvFile, _ := os.Open(*csvFilePtr)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var totalProblems, correctAnswers int
	inputReader := bufio.NewReader(os.Stdin)

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		totalProblems++
		fmt.Printf("Problem #%d: %s = ", totalProblems, line[0])
		ans, _ := inputReader.ReadString('\n')
		ans = strings.TrimSuffix(ans, "\n")

		if ans == line[1] {
			correctAnswers++
		}
	}
	fmt.Printf("Answered %d correctly out of %d.\n", correctAnswers, totalProblems)
}
