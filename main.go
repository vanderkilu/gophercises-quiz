package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func importCSV(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}


func main() {
	filePath := flag.String("f", "problems.csv", "Path to csv file to use for the quiz")
	limit := flag.Int("t", 10, "Timer for the trivia")
	flag.Parse()

	csvFile, err := importCSV(*filePath)
	if err != nil {
		panic("CsvImportError: path csv file is incorrect")
	}

	//read the contents of the csv file
	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1
	data, err := csvReader.ReadAll()
	if err != nil {
		panic("CsvReaderError: unable to read the csv content")
	}

	fmt.Println("Welcome to the quiz app")
	fmt.Println("Type a single answer for each question")
	fmt.Println("====================================")
	fmt.Printf("Press any key to begin the test\n")
	fmt.Scanf("\n")

	scanner := bufio.NewScanner(os.Stdin)
	userInputChan := make(chan string)
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	totalCorrect := 0
	totalWrong := 0

	for index, row := range data {
		question, answer := row[0], row[1]
		fmtQuestion := fmt.Sprintf("%d. %s = ", index+1, question)
		fmt.Print(fmtQuestion)

		go func() {
			scanner.Scan()
			userInput := scanner.Text()
			userInputChan <- userInput
		}()

			select {
			case <-timer.C:
				fmt.Printf("\nTime up\n")
				fmt.Printf("Correct answers: %d / %d\n", totalCorrect, len(data))
				fmt.Printf("Wrong answers: %d / %d\n", totalWrong, len(data))
				return
			case userInput := <-userInputChan:
				if userInput == answer {
					totalCorrect += 1
				} else {
					totalWrong += 1
				}
			}
	}

	fmt.Printf("\nYou completed the quiz\n")
	fmt.Printf("Correct answers: %d / %d\n", totalCorrect, len(data))
	fmt.Printf("Wrong answers: %d / %d\n", totalWrong, len(data))

}
