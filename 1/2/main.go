package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type problem struct {
	question, answer string
}

func readCSV(filename string) ([]problem, error) {
	csvFile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()

	out := []problem{}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		out = append(out, problem{line[0], line[1]})
	}

	return out, nil
}

func main() {
	csvFilePtr := flag.String("csv", "../problems.csv", "A CSV file in the format of 'question,answer'")
	limitPtr := flag.Int("limit", 2, "The time limit for the quiz in seconds.")
	flag.Parse()

	problems, _ := readCSV(*csvFilePtr)

	var correctAnswers int
	inputReader := bufio.NewReader(os.Stdin)

	done := make(chan bool)
	ticker := time.NewTicker(time.Second * time.Duration(*limitPtr))

	go func() {
		for i, q := range problems {
			fmt.Printf("Problem #%d: %s = ", i+1, q.question)

			ans, _ := inputReader.ReadString('\n')
			ans = strings.TrimSpace(ans)

			if ans == q.answer {
				correctAnswers++
			}
		}
		done <- true
	}()

	select {
	case <-done:
	case <-ticker.C:
		fmt.Printf("\nAnswered %d correctly out of %d.\n", correctAnswers, len(problems))
	}
}
