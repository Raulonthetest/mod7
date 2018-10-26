package main

import (
	"flag"
	"fmt"
	"mod7/oem"
	"mod7/tendigit"
	"time"
)

func main() {
	b := flag.Bool("b", false, "Generate both keys.")
	o := flag.Bool("o", false, "Generate an OEM key.")
	d := flag.Bool("d", false, "Generate a 10-digit key (aka CD Key).")
	r := flag.Int("r", 1, "Generate n keys.")
	t := flag.Bool("t", false, "Show how long the generation took.")
	flag.Parse()
	switch {
	case *r < 1:
		*r = 1
	}
	var started time.Time
	if *t {
		started = time.Now()
	}
	CDKeych := make(chan string)
	OEMKeych := make(chan string)
	for i := 0; i < *r; i++ {
		switch {
		case *d:
			go tendigit.Generate10digit(CDKeych)
			fmt.Println(<-CDKeych)
		case *o:
			go oem.GenerateOEM(OEMKeych)
			fmt.Println(<-OEMKeych)
		case *b:
			go oem.GenerateOEM(OEMKeych)
			go tendigit.Generate10digit(CDKeych)
			fmt.Println(<-CDKeych)
			fmt.Println(<-OEMKeych)
		default:
			fmt.Println("You must specify what you want to generate! Usage:")
			flag.PrintDefaults()
			return
		}
	}
	if *t {
		var ended time.Duration
		switch {
		case time.Since(started).Round(time.Second) > 1:
			ended = time.Since(started).Round(time.Millisecond)
		default:
			ended = time.Since(started).Round(time.Microsecond)
		}
		switch {
		case *r > 1:
			fmt.Printf("Took %s to generate %d keys.\n", ended, *r)
		case *r == 1:
			fmt.Printf("Took %s to generate %d key.\n", ended, *r)
		}
	}
}
