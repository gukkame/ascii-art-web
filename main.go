package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func asciihandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/ascii-art" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "tamplates/index.html")
		return

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		var asciiarr []string
		intext := r.FormValue("text")
		asciiarr = append(asciiarr, intext)
		asciiarr = append(asciiarr, r.Form.Get("textstyle")+".txt")

		fmt.Println(asciiarr[0])
		fmt.Println(asciiarr[1])
		fmt.Println(len(asciiarr))

		var s []string
		if len(asciiarr) > 1 {

			incomingSentence := asciiarr[0]
			var arr []int
			var word string
			//Printing words so that \n marks the start of a new line

			input_slice := strings.Split(incomingSentence, "\\n")
			for _, word := range input_slice {

				var ch, position rune
				// var i int

				for _, ch = range word {
					position = (ch - 32) * 9
					arr = append(arr, int(position)) // array of NUMBERS of first lines of input charachters
				}

			}
			s = Printer(word, asciiarr[1], arr)

		} else {
			fmt.Println("Please enter more than 0 arguments")
		}

		AsciiString := strings.Join(s, "")
		fmt.Fprintf(w, AsciiString)
		//fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./tamplates")))
	http.HandleFunc("/ascii-art", asciihandler)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func Printer(input, fontName string, originalArray []int) (s []string) {
	for lineCounter := 1; lineCounter < 9; lineCounter++ { //loop through the heigth of the ascii characters
		for i := range originalArray {
			oneLineofOneLetter, _ := ReadLine(fontName, originalArray[i]+lineCounter)
			fmt.Print(oneLineofOneLetter) //just print a bunch of same "height" lines of different letters in a row
			s = append(s, oneLineofOneLetter)

		}
		s = append(s, "\n")
		fmt.Print(" ")
		fmt.Print("\n") //and start a new line again
	}
	return s

}

func ReadLine(fontName string, startingLineEach int) (line string, err error) {

	var lastLineOfScanner int

	FontFile, err := os.Open(fontName) //open file
	if err != nil {                    //check for errors
		log.Fatal(err)
	}
	defer FontFile.Close()                //close file
	scanner := bufio.NewScanner(FontFile) // make a scaner machine aka not doing anything yet
	for scanner.Scan() {                  //loop through lines until you reached the wished line aka start machine

		if lastLineOfScanner == startingLineEach {
			return scanner.Text(), scanner.Err() //scanner.Text() is where you store the result of scanner you can print it aswell if you make it equal to a value

		}
		lastLineOfScanner++
	}
	return line, io.EOF // it will return the line when the scanner mathces reading to the input line number aka gets to it
}
