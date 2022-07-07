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
	FILE_FLAG_DEFAULT_FILE = "./problems.csv"
	FILE_FLAG_USAGE        = "You can specify a path to CSV with questions so that the program can load new questions."
)

type Problem struct {
	Question string
	Answer   string
}

type Result struct {
	Score int
	Total int
}

func main() {
	filename := flag.String("file", FILE_FLAG_DEFAULT_FILE, FILE_FLAG_USAGE)
	flag.Parse()

	file, err := os.Open(*filename)
	defer file.Close()

	if err != nil {
		panic("CSV with questions is essential to run the quiz. Try again or use -h if need a help.")
	}

	problems := readProblemsFromCsv(file)
	result := askQuestions(problems)

	result.printSummary()
}

func (p Problem) isCorrect(answer string) bool {
	return strings.Trim(answer, " ") == p.Answer
}

func (r Result) printSummary() {
	fmt.Printf("User answered %d problems correctly. Scored %d of %d.\n", r.Score, r.Score, r.Total)
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

func askQuestions(problems *[]Problem) Result {
	result := Result{Score: 0, Total: len(*problems)}

	for idx, problem := range *problems {

		inputChannel := make(chan string)
		timeExceedChannel := make(chan bool)

		go getUserInput(inputChannel, idx, problem)

		go timer(timeExceedChannel)

		select {
		case answer := <-inputChannel:
			if problem.isCorrect(answer) {
				result.Score++
			}
		case <-timeExceedChannel:
			fmt.Println("Time limit exceeded.")
			result.printSummary()
			os.Exit(0)
		}
	}

	return result
}

func getUserInput(channel chan<- string, problemId int, problem Problem) {
	fmt.Printf("%d. %s: ", problemId+1, problem.Question)
	var answer string
	fmt.Scanf("%s\n", &answer)
	channel <- answer
}

func timer(channel chan<- bool) {
	time.Sleep(2 * time.Second)
	channel <- true
}
