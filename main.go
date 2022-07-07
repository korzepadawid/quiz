package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	FILE_FLAG_NAME    = "file"
	FILE_FLAG_DEFAULT = "./problems.csv"
	FILE_FLAG_USAGE   = "You can specify a path to CSV with questions so that the program can load new questions."

	TIME_FLAG_NAME    = "time"
	TIME_FLAG_DEFAULT = 3
	TIME_FLAG_USAGE   = "You can specify a path to CSV with questions so that the program can load new questions."
)

type Problem struct {
	Question string
	Answer   string
}

type Result struct {
	Score int
	Total int
}

func (p Problem) IsCorrect(answer string) bool {
	return strings.Trim(answer, " ") == p.Answer
}

func (r Result) PrintSummary() {
	fmt.Printf("User answered %d problems correctly. Scored %d of %d.\n", r.Score, r.Score, r.Total)
}

func (r *Result) IncreaseScore() {
	r.Score++
}

func main() {
	filename := flag.String(FILE_FLAG_NAME, FILE_FLAG_DEFAULT, FILE_FLAG_USAGE)
	time := flag.Int(TIME_FLAG_NAME, TIME_FLAG_DEFAULT, TIME_FLAG_USAGE)
	flag.Parse()

	file, err := os.Open(*filename)
	defer file.Close()

	if err != nil {
		panic("CSV with questions is essential to run the quiz. Try again or use -h if need a help.")
	}

	problems := readProblemsFromCsv(file)
	result := askQuestions(problems, *time)

	result.PrintSummary()
}

func readProblemsFromCsv(f *os.File) *[]Problem {
	reader := csv.NewReader(f)

	problems := make([]Problem, 0)

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		problem := Problem{Question: record[0], Answer: record[1]}
		problems = append(problems, problem)
	}

	return &problems
}

func askQuestions(problems *[]Problem, time int) Result {
	result := Result{Score: 0, Total: len(*problems)}

	for idx, problem := range *problems {

		inputChannel := make(chan string)
		timeExceedChannel := make(chan bool)

		go getUserInput(inputChannel, idx, problem)

		go timer(timeExceedChannel, time)

		select {
		case answer := <-inputChannel:
			if problem.IsCorrect(answer) {
				result.IncreaseScore()
			}
		case <-timeExceedChannel:
			fmt.Println("Time limit exceeded.")
			result.PrintSummary()
			os.Exit(0)
		}
	}

	return result
}

func getUserInput(channel chan<- string, problemID int, problem Problem) {
	fmt.Printf("%d. %s: ", problemID+1, problem.Question)
	var answer string
	fmt.Scanf("%s\n", &answer)
	channel <- answer
}

func timer(channel chan<- bool, timeLimit int) {
	time.Sleep(time.Duration(timeLimit) * time.Second)
	channel <- true
}
