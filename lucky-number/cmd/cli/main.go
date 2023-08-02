package main

import (
	"fmt"

	// Import the random package.
	"lucky-number.dgr.net/internal/random"

	"github.com/fatih/color"
)

func main() {
	green := color.New(color.FgGreen)
	green.Printf("your lucky number is %d!\n", random.Number())

	// Call the random.Number() function to get the random number. Notice that
	// we use the package name as the accessor, just like we do for the standard
	// library packages.
	fmt.Printf("Your lucky number is %d!\n", random.Number())
}
