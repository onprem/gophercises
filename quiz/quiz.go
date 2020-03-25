package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	csvPath string
	limit   int
	shuffle bool
)

func init() {
	flag.StringVar(&csvPath, "csv", "problems.csv", `a csv file in the format of 'questions,answers'`)
	flag.IntVar(&limit, "limit", 30, "time limit for the quiz in seconds")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle questions before asking")
}

func askQuestions(questions [][]string, correctAnswers *int, exit chan bool) {
	var ans string
	for i, row := range questions {
		fmt.Printf("Question #%d: %s = ", i+1, row[0])
		fmt.Scanln(&ans)

		if strings.TrimSpace(ans) == row[1] {
			*correctAnswers++
		}
	}
	exit <- true
}

func timer(start chan bool, exit chan bool) {
	fmt.Printf("Time limit for the quiz: %ds\nPress enter to start the quiz. ", limit)
	fmt.Scanln()
	start <- true

	time.Sleep(time.Duration(limit) * time.Second)
	exit <- true
}

func main() {
	flag.Parse()

	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	csvr := csv.NewReader(file)
	questions, err := csvr.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}

	var correctAnswers int
	exit := make(chan bool)
	start := make(chan bool)

	go timer(start, exit)

	<-start
	go askQuestions(questions, &correctAnswers, exit)

	<-exit
	fmt.Printf("\nYou scored %d out of %d.\n", correctAnswers, len(questions))
}
