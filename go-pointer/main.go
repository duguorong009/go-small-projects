package main

import "fmt"

func main() {
	var answer int = 32

	var answerPtr *int = &answer

	fmt.Println(answerPtr)

	*answerPtr = 99

	fmt.Println(answer)
}
