/*
 * Basic PDF split example: Splitting by page range.
 *
 * Run as: go run pdf_split.go input.pdf output
 */

package main

import (
	"fmt"
	"os"
	pdf "github.com/yangluoshen/unidoc/pdf/model"
	"gopkg.in/h2non/bimg.v1"
	"bytes"
)

const (
	pdfTpl = "output/%s_%d.pdf"
	pngTpl = "output/%s_%d.png"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_split.go input.pdf <page_from> <page_to> output.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	err := split(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Complete")
}

func split(inputPath string) error {

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	// convert
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return err
		}
		err = pageToPdf(page, fmt.Sprintf(pdfTpl, f.Name()[:len(f.Name())-4], i))
		if err != nil {
			panic(err)
		}
		err = pageToPng(page, fmt.Sprintf(pngTpl, f.Name()[:len(f.Name())-4], i))
		if err != nil {
			panic(err)
		}
	}

	// statistics
	for i := 1; i <= numPages; i++ {
		pdff, err := os.Open(fmt.Sprintf(pdfTpl, f.Name()[:len(f.Name())-4], i))
		if err != nil {
			panic(err)
		}
		pngf, err := os.Open(fmt.Sprintf(pngTpl, f.Name()[:len(f.Name())-4], i))
		if err != nil {
			panic(err)
		}

		info1, _ := pdff.Stat()
		info2, _ := pngf.Stat()

		f, err := os.OpenFile("statistics.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("OpenFile failed:%s", err))
		}
		if _, err := f.Write([]byte(fmt.Sprintf("%.2fk, %.2fk\n", float64(info1.Size())/1024, float64(info2.Size())/1024))); err != nil {
			panic(fmt.Errorf("write failed:%s", err))
		}
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
	return nil
}

func pageToPdf(page *pdf.PdfPage, outFile string) error {
	pdfWriter := pdf.NewPdfWriter()
	err := pdfWriter.AddPage(page)
	if err != nil {
		return err
	}

	fWrite, err := os.Create(outFile)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}
	return nil
}

func pageToPng(page *pdf.PdfPage, outFile string) error {

	pdfWriter := pdf.NewPdfWriter()
	err := pdfWriter.AddPage(page)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = pdfWriter.Write(&buf)
	if err != nil {
		return err
	}

	//w, h := getPageWH(page)
	//fmt.Println("w, h:", w, h)

	o := &bimg.Options{
		//Width: int(w),
		//Height: int(h),
		Quality: 200,
		Type: bimg.PNG,
		Zoom: 2,
	}
	newImage, err := bimg.NewImage(buf.Bytes()).Process(*o)
	if err != nil {
		return err
	}

	bimg.Write(outFile, newImage)
	return nil
}

func getPageWH(page *pdf.PdfPage) (w, h float64) {
	bbox, _ := page.GetMediaBox()

	w = (*bbox).Urx - (*bbox).Llx
	h = (*bbox).Ury - (*bbox).Lly
	return
}
