package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"berita-acara-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
)

// GeneratePDF handles POST /api/generate/pdf
func GeneratePDF(c *gin.Context) {
	var payload models.BeritaAcara
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := BuildPDF(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}

	safeName := sanitizeFilename(payload.Header.Nama)
	filename := fmt.Sprintf("BeritaAcara_%s.pdf", safeName)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// BuildPDF generates PDF bytes from BeritaAcara model
func BuildPDF(payload models.BeritaAcara) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(25, 25, 25)
	pdf.AddPage()

	pageW, _ := pdf.GetPageSize()
	leftM, _, rightM, _ := pdf.GetMargins()
	contentW := pageW - leftM - rightM // 210 - 25 - 25 = 160mm

	// ── Header ──────────────────────────────────────────────────────────────
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(contentW, 5.5, "PT. BANK NEGARA INDONESIA (Persero) Tbk", "", 1, "L", false, 0, "")
	pdf.CellFormat(contentW, 5.5, "DIVISI RETAIL DIGITAL DELIVERY", "", 1, "L", false, 0, "")

	pdf.Ln(8)

	// ── Title ────────────────────────────────────────────────────────────────
	pdf.SetFont("Arial", "BU", 13)
	pdf.CellFormat(contentW, 7, "BERITA ACARA", "", 1, "C", false, 0, "")

	pdf.Ln(8)

	// ── Intro sentence ───────────────────────────────────────────────────────
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(contentW, 5, "Yang bertandatangan di bawah ini menerangkan bahwa:", "", "L", false)
	pdf.Ln(3)

	// ── Key-value block ──────────────────────────────────────────────────────
	kvLabelW := 32.0
	kvColonW := 4.0
	kvValueW := contentW - kvLabelW - kvColonW

	writeKV := func(label, value string) {
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(kvLabelW, 5.5, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(kvColonW, 5.5, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(kvValueW, 5.5, value, "", 1, "L", false, 0, "")
	}

	writeKV("Nama", payload.Header.Nama)
	writeKV("NPP BNI", payload.Header.NppBni)
	writeKV("Departement", payload.Header.Departement)
	writeKV("Kelompok", payload.Header.Kelompok)

	pdf.Ln(6)

	// ── Attendance Table ─────────────────────────────────────────────────────
	colWidths := []float64{28, 16, 23, 23, 70}
	headers := []string{"TANGGAL", "HARI", "JAM DATANG", "JAM PULANG", "KETERANGAN"}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetLineWidth(0.2)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)

	for _, row := range payload.Rows {
		lineH := 4.5
		keteranganLines := splitTextToLines(pdf, row.Keterangan, colWidths[4]-4.0)
		numLines := len(keteranganLines)
		rowH := float64(numLines)*lineH + 3.0
		if rowH < 7.5 {
			rowH = 7.5
		}

		startY := pdf.GetY()
		startX := leftM

		// Col 0-3: draw with border + centered text
		simpleData := []string{row.Tanggal, row.Hari, row.JamDatang, row.JamPulang}
		for i := 0; i < 4; i++ {
			pdf.SetXY(startX, startY)
			pdf.CellFormat(colWidths[i], rowH, simpleData[i], "1", 0, "C", false, 0, "")
			startX += colWidths[i]
		}

		// Col 4 (Keterangan): draw border rect manually, then text
		pdf.Rect(startX, startY, colWidths[4], rowH, "D") // "D" = draw border only

		textH := float64(numLines) * lineH
		textStartY := startY + (rowH-textH)/2
		for idx, line := range keteranganLines {
			pdf.SetXY(startX+2.0, textStartY+float64(idx)*lineH)
			pdf.CellFormat(colWidths[4]-4.0, lineH, line, "", 0, "L", false, 0, "")
		}

		// Move cursor to next row
		pdf.SetXY(leftM, startY+rowH)
	}

	pdf.Ln(18)

	// ── Signer Block ─────────────────────────────────────────────────────────
	colW := contentW / 2

	// Row 1: Labels (Centered in each column)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(colW, 5, "Saksi", "", 0, "C", false, 0, "")
	pdf.CellFormat(colW, 5, "Hormat Saya,", "", 1, "C", false, 0, "")

	// Signature space
	pdf.Ln(22)

	// Row 2: Names (Bold, Centered in each column)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(colW, 5, payload.Signers.Saksi, "", 0, "C", false, 0, "")
	pdf.CellFormat(colW, 5, payload.Signers.HormatSaya, "", 1, "C", false, 0, "")

	pdf.Ln(10)

	// Menyetujui — centered
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(contentW, 5, "Menyetujui", "", 1, "C", false, 0, "")

	pdf.Ln(22)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(contentW, 5, payload.Signers.MenyetujuiNama, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(contentW, 5, payload.Signers.MenyetujuiJabatan, "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateDOCX handles POST /api/generate/docx
func GenerateDOCX(c *gin.Context) {
	var payload models.BeritaAcara
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docx, err := buildDOCX(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate DOCX: " + err.Error()})
		return
	}

	safeName := sanitizeFilename(payload.Header.Nama)
	filename := fmt.Sprintf("BeritaAcara_%s.docx", safeName)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Length", fmt.Sprintf("%d", len(docx)))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", docx)
}

// splitTextToLines splits text to fit within a given width, respecting newlines
func splitTextToLines(pdf *fpdf.Fpdf, text string, width float64) []string {
	var lines []string
	rawLines := strings.Split(text, "\n")
	for _, rawLine := range rawLines {
		words := strings.Fields(rawLine)
		if len(words) == 0 {
			continue
		}
		current := ""
		for _, w := range words {
			test := w
			if current != "" {
				test = current + " " + w
			}
			if pdf.GetStringWidth(test) > width && current != "" {
				lines = append(lines, current)
				current = w
			} else {
				current = test
			}
		}
		if current != "" {
			lines = append(lines, current)
		}
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

// sanitizeFilename removes/replaces characters not safe in filenames
func sanitizeFilename(name string) string {
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}
