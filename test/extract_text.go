
package main

import (
	"fmt"
	"os"

	"github.com/yangluoshen/unidoc/pdf/extractor"
	pdf "github.com/yangluoshen/unidoc/pdf/model"
	"strconv"
	"github.com/yangluoshen/unidoc/pdf/creator"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run pdf_insert_text.go input.pdf <page> <xpos> <ypos> \"text\" output.pdf\n")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := outputPdfText(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// outputPdfText prints out contents of PDF file to stdout.
func outputPdfText(inputPath, out string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		panic(err)
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	c := creator.New()

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

		pText, _, _, err := ex.ExtractPageText()
		if err != nil {
			return err
		}

		c.NewPage()

		for _, tm := range pText.Marks {
			p := c.NewParagraph(tm.Text)
			p.SetPos(tm.OrientedStart.X, c.Context().PageHeight - tm.OrientedStart.Y)
			_ = c.Draw(p)
		}
	}
	err = c.WriteToFile(out)

	return nil
}

func printText(t *extractor.PageText) {
	for _, m := range t.Marks {
		fmt.Printf("%v,", m.Text)
		fmt.Printf("%v,", m.Orient)
		fmt.Printf("%v,", m.OrientedStart)
		fmt.Printf("%v,", m.OrientedEnd)
		fmt.Printf("%v,", m.Height)
		fmt.Printf("%v\n", round(m.SpaceWidth))
	}
}

func round(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	return value
}
