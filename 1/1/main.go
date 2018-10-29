package main

import (
	"bufio"
	"os"
	"path/filepath"
)

type Problem struct {
	Question string
	Answer   string
}

func main() {
	problemsFile, _ := filepath.Abs("../problems.csv")
	reader := bufio.NewReader(os.Stdin)
}
