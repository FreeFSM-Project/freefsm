package handlers

import (
	"bytes"
	"io"
	"net/http"
)

func writePDFResponse(w http.ResponseWriter, filename string, generate func(io.Writer) error) {
	writePDFResponseWithDisposition(w, filename, false, generate)
}

func writePDFResponseWithDisposition(w http.ResponseWriter, filename string, download bool, generate func(io.Writer) error) {
	var buf bytes.Buffer
	if err := generate(&buf); err != nil {
		http.Error(w, "PDF generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	disposition := "inline"
	if download {
		disposition = "attachment"
	}
	w.Header().Set("Content-Disposition", disposition+`; filename="`+filename+`"`)
	_, _ = buf.WriteTo(w)
}

func generatePDFBytes(generate func(io.Writer) error) ([]byte, error) {
	var buf bytes.Buffer
	if err := generate(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
