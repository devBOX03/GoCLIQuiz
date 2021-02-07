package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "csv file in the format of 'question, answer'")
	timeLimit := flag.Int("timelimit", 30, "time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Not able to open file: %s \n", *csvFilename)
		os.Exit(1)
	}
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Printf("Not able to read file %s \n", *csvFilename)
		os.Exit(1)
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correctAnswerCount := 0

quizloop:
	for index, problem := range problems {
		fmt.Printf("%d) %s = \n", index+1, problem.question)
		answerChanel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s \n", &answer)
			answerChanel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break quizloop
		case answer := <-answerChanel:
			if answer == problem.answer {
				correctAnswerCount++
			}
		}		
	}
	fmt.Printf("You got %d correct answer out of %d", correctAnswerCount, len(problems))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for index, line := range lines {
		problems[index] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

type problem struct {
	question string
	answer   string
}
