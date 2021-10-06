package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var banner string

type Asciiart struct {
	AsciiText   []string
	AsciiString string
}

func ReadLine(fn string, n int) string { // Returns the requested line in the font file
	f, _ := os.Open(fn)
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lineNum := 0; lineNum < n; lineNum++ {
		line, _ = bf.ReadString('\n')
	}
	return line[:len(line)-1]
}

func asciiArt(text string) []string {
	textArr := strings.Split(text, "\\n")
	var result []string

	for _, v := range text {
		if int(v) > 126 || int(v) < 32 {
			log.Fatal("non ascii character")
		}
	}

	for i := 0; i < len(textArr); i++ {
		if textArr[i] == "" { // Adds new line if multiple \n in a row
			result = append(result, "\n")
			continue
		}

		for z := 1; z < 9; z++ {
			for y := 0; y < len(textArr[i]); y++ {
				line := (int(textArr[i][y])-32)*9 + z + 1 // Finds the right place in the font file
				result = append(result, ReadLine(banner, line))
			}
			result = append(result, "\n")
		}
	}
	return result
}

func asciihandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/ascii-art" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "tamplates/index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		inputText := r.FormValue("text")
		banner = r.Form.Get("textstyle") + ".txt"
		result := asciiArt(inputText)

		art := Asciiart{
			AsciiText:   result,
			AsciiString: strings.Join(result, ""),
		}

		parsedTemplate, _ := template.ParseFiles("tamplates/asciiArt.html")
		err := parsedTemplate.Execute(w, art)

		if err != nil {
			log.Println("Error executing template :", err)
			return
		}

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
