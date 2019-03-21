
package main

import (
	"fmt"
	"os"

	"github.com/yangluoshen/unidoc/pdf/extractor"
	pdf "github.com/yangluoshen/unidoc/pdf/model"
	"reflect"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_extract_text.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	err := outputPdfText(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// outputPdfText prints out contents of PDF file to stdout.
func outputPdfText(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return err
		}

		pText, numChars, numMisses, err := ex.ExtractPageText()
		if err != nil {
			return err
		}

		/*
		text, err := ex.ExtractText()
		if err != nil {
			return err
		}
		*/

		fmt.Println("------------------------------")
		fmt.Printf("Page %d:\n", pageNum)
		fmt.Printf("\"%s\"\n", pText.ToText())
		fmt.Printf("numChars: %d\n", numChars)
		fmt.Printf("NumMisses: %d\n", numMisses)
		printText(pText)
		fmt.Println("------------------------------")
	}

	return nil
}

func printText(t *extractor.PageText) {
	v := reflect.ValueOf(*t)

	marks := v.FieldByName("marks")
	if marks.Len() <= 0 {
		return
	}

	for i:=0; i< marks.Len(); i++ {
		m := marks.Index(i)
		fmt.Println("mark:", m)

		text := m.FieldByName("text").String()
		fmt.Println("text:", text)

		orient := m.FieldByName("orient").Int()
		fmt.Println("orient:", orient)

		orientedStartX := m.FieldByName("orientedStart").FieldByName("X")
		orientedStartY := m.FieldByName("orientedStart").FieldByName("Y")
		fmt.Println("orientedStart:", orientedStartX, orientedStartY)

		fmt.Println("height:", m.FieldByName("height").Float())
		fmt.Println("spaceWidth:", m.FieldByName("spaceWidth").Float())
	}


}
