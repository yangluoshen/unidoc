/*
 * Basic merging of PDF files.
 * Simply loads all pages for each file and writes to the output file.
 * See pdf_merge_advanced.go for a more advanced version which handles merging document forms (acro forms) also.
 *
 * Run as: go run pdf_merge.go output.pdf input1.pdf input2.pdf input3.pdf ...
 */

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	//unicommon "github.com/yangluoshen/unidoc/common"
	pdf "github.com/yangluoshen/unidoc/pdf/model"
)

func init() {
	// Debug log level.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

var (
	filePool []os.FileInfo
	root = "/home/shenweimin/virtual_disk/pdf/"
)

func main() {
	/*
	if len(os.Args) < 4 {
		fmt.Printf("Requires at least 3 arguments: output_path and 2 input paths\n")
		fmt.Printf("Usage: go run pdf_merge.go output.pdf input1.pdf input2.pdf input3.pdf ...\n")
		os.Exit(0)
	}

	outputPath := ""
	inputPaths := []string{}

	// Sanity check the input arguments.
	for i, arg := range os.Args {
		if i == 0 {
			continue
		} else if i == 1 {
			outputPath = arg
			continue
		}

		inputPaths = append(inputPaths, arg)
	}

	err := mergePdf(inputPaths, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
	*/

	initFilePool()
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i ++ {
		num := rand.Intn(30) + 50

		inputs := make([]string, 0)
		var output = fmt.Sprintf("./merge_pdf/%d.pdf", time.Now().UnixNano())
		for i := 0; i < num; i++ {
			f := filePool[rand.Intn(len(filePool))]
			inputs = append(inputs, root+f.Name())
		}

		fmt.Printf("%d: inputs:%+v\n", i, inputs)
		fmt.Printf("%d: output:%s\n", i, output)

		if err := mergePdf(inputs, output); err != nil {
			fmt.Println("merge failed:%s", err)
			return
		}
		fmt.Println("verify:", isValidPdf(output))
	}

	//_ = unlinkInvalidPdf()
	//fmt.Println(isValidPdf(root+"/459.xls.pdf"))
	//fmt.Println("verify:", isValidPdf(output))
}


func initFilePool() {
	filePool, _ = ioutil.ReadDir(root)
	fmt.Println("len:", len(filePool))
}

func mergePdf(inputPaths []string, outputPath string) error {
	pdfWriter := pdf.NewPdfWriter()

	for _, inputPath := range inputPaths {
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
			auth, err := pdfReader.Decrypt([]byte(""))
			if err != nil {
				return err
			}
			if !auth {
				return errors.New("Cannot merge encrypted, password protected document")
			}
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}
	}

	fWrite, err := os.Create(outputPath)
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


func ReaderToSeeker(r io.ReadCloser) io.ReadSeeker {
	b, _ := ioutil.ReadAll(r)
	return bytes.NewReader(b)
}

func unlinkInvalidPdf() error {
	files, _ := ioutil.ReadDir(root)
	for _, file := range files{
		name := path.Join(root, file.Name())
		f, err := os.Stat(name)
		if err != nil {
			continue
		}
		if f.Size() < 5856 {
			_ = os.Remove(name)
			fmt.Printf("remove:%s\n", name)
		}
	}

	return nil
}

func isValidPdf(path string) bool{
	f, err := os.Open(path)
	if err != nil {
		return true
	}

	_, err = pdf.NewPdfReader(f)
	if err != nil {
		return false
	}

	return true
}
