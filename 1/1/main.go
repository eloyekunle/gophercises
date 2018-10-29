package main

import (
	"bufio"
	"encoding/csv"
	"os"
)

type Problem struct {
	Question string
	Answer   string
}

func main() {
	problemsFile, _ := os.Open("../problems.csv")
	reader := csv.NewReader(bufio.NewReader(problemsFile))

	inputReader := bufio.NewReader(os.Stdin)
}
