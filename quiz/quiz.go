package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

var csvPath string

func init() {
	flag.StringVar(&csvPath, "csv", "problems.csv", `a csv file in the format of 'questions,answers'`)
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

	var (
		ans            string
		correctAnswers int
	)
	for i, row := range questions {
		fmt.Printf("Question #%d: %s = ", i, row[0])
		fmt.Scanln(&ans)

		if ans == row[1] {
			correctAnswers++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correctAnswers, len(questions))
}
