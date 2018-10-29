package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	problemsFile, _ := os.Open("../problems.csv")
	reader := csv.NewReader(bufio.NewReader(problemsFile))
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
		fmt.Println(line[0])
		ans, _ := inputReader.ReadString('\n')
		ans = strings.TrimSuffix(ans, "\n")

		if ans == line[1] {
			correctAnswers++
		}
	}
	fmt.Println(correctAnswers, totalProblems)
}
