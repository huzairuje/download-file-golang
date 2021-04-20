package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"greetings": "Hi, Welcome to PDF generator",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func PdfGen(w http.ResponseWriter, r *http.Request) {

	//pdf generator
	pdfGen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfGen.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	linkUrlFromRequest := r.URL.Query().Get("link")
	fmt.Println("link url ", linkUrlFromRequest)
	page := wkhtmltopdf.NewPage(linkUrlFromRequest)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	pdfGen.AddPage(page)

	err = pdfGen.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfGen.WriteFile("./output.pdf")
	if err != nil {
		log.Fatal(err)
	}

	filename := "output.pdf"

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename="+filename+"")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

	//stream the body to the client without fully loading it into memory

	io.Copy(w, f)
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/pdf", PdfGen)
	err := http.ListenAndServe(":9000", nil)

	if err != nil {
		fmt.Println(err)
	}
}
