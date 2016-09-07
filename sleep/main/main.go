package main

import "flag"
import "fmt"
import "os"
import "time"

var bell = flag.Bool("bell", false,
	"Whether to send an ASCII bell after we finish sleeping.")

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Expected duration")
		os.Exit(1)
		return
	}
	duration, err := time.ParseDuration(flag.Arg(0))
	if err != nil {
		fmt.Printf("Could not parse \"%v\" as a duration: %v\n", flag.Arg(0), err)
		os.Exit(2)
		return
	}
	var deadline = time.Now().Add(duration)
	var remaining = duration

	// Print updates every second.
	var ticker = time.NewTicker(time.Second)
	for time.Now().Before(deadline) {
		// Home the cursor, clear to end of line, and print the remaining duration.
		fmt.Printf("\033[9999D\033[KSleeping for %s", remaining.String())
		remaining = remaining - time.Second

		// Await a tick.
		<-ticker.C
	}

	// Clear the remaining indicator before returning.
	fmt.Print("\033[9999D\033[K")

	// Print the final output line.
	fmt.Printf("Slept for %s\n", duration.String())

	if *bell {
		fmt.Print("\007")
	}
}
