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

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pageW, _ := pdf.GetPageSize()
	leftM, _, rightM, _ := pdf.GetMargins()
	contentW := pageW - leftM - rightM

	// ── Header ──────────────────────────────────────────────────────────────
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(contentW, 6, "PT. BANK NEGARA INDONESIA (Persero) Tbk", "", 1, "L", false, 0, "")
	pdf.CellFormat(contentW, 6, "DIVISI RETAIL DIGITAL DELIVERY", "", 1, "L", false, 0, "")

	pdf.Ln(6)

	// ── Title ────────────────────────────────────────────────────────────────
	pdf.SetFont("Arial", "BU", 13)
	pdf.CellFormat(contentW, 8, "BERITA ACARA", "", 1, "C", false, 0, "")

	pdf.Ln(6)

	// ── Intro sentence ───────────────────────────────────────────────────────
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(contentW, 5, "Yang bertandatangan di bawah ini menerangkan bahwa:", "", "L", false)
	pdf.Ln(3)

	// ── Key-value block ──────────────────────────────────────────────────────
	kvLabelW := 35.0
	kvColonW := 5.0
	kvValueW := contentW - kvLabelW - kvColonW

	writeKV := func(label, value string) {
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(kvLabelW, 6, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(kvColonW, 6, ":", "", 0, "C", false, 0, "")
		pdf.CellFormat(kvValueW, 6, value, "", 1, "L", false, 0, "")
	}

	writeKV("Nama", payload.Header.Nama)
	writeKV("NPP BNI", payload.Header.NppBni)
	writeKV("Departement", payload.Header.Departement)
	writeKV("Kelompok", payload.Header.Kelompok)

	pdf.Ln(4)

	// ── Attendance Table ─────────────────────────────────────────────────────
	colWidths := []float64{32, 22, 24, 24, contentW - 32 - 22 - 24 - 24}
	headers := []string{"TANGGAL", "HARI", "JAM DATANG", "JAM PULANG", "KETERANGAN"}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(220, 220, 220)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)

	for _, row := range payload.Rows {
		cells := []string{row.Tanggal, row.Hari, row.JamDatang, row.JamPulang, row.Keterangan}

		// Calculate the max height needed for this row (for keterangan wrapping)
		lineH := 5.0
		keteranganLines := splitTextToLines(pdf, row.Keterangan, colWidths[4]-2)
		rowH := float64(len(keteranganLines)) * lineH
		if rowH < 10 {
			rowH = 10
		}

		startY := pdf.GetY()
		startX := leftM

		// Draw each cell
		for i := 0; i < 4; i++ {
			pdf.SetXY(startX, startY)
			pdf.CellFormat(colWidths[i], rowH, cells[i], "1", 0, "C", false, 0, "")
			startX += colWidths[i]
		}

		// Draw keterangan cell with multi-line text
		pdf.SetXY(startX, startY)
		pdf.MultiCell(colWidths[4], lineH, row.Keterangan, "1", "L", false)

		// Ensure next row starts below
		if pdf.GetY() < startY+rowH {
			pdf.SetY(startY + rowH)
		}
	}

	pdf.Ln(12)

	// ── Signer Block ─────────────────────────────────────────────────────────
	colW := contentW / 2

	// Row 1: Labels
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(colW, 6, "Saksi", "", 0, "L", false, 0, "")
	pdf.CellFormat(colW, 6, "Hormat Saya,", "", 1, "L", false, 0, "")

	// Signature space
	pdf.Ln(20)

	// Row 2: Names
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(colW, 6, payload.Signers.Saksi, "", 0, "L", false, 0, "")
	pdf.CellFormat(colW, 6, payload.Signers.HormatSaya, "", 1, "L", false, 0, "")

	pdf.Ln(10)

	// Menyetujui — centered
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(contentW, 6, "Menyetujui", "", 1, "C", false, 0, "")

	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(contentW, 6, payload.Signers.MenyetujuiNama, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(contentW, 6, payload.Signers.MenyetujuiJabatan, "", 1, "C", false, 0, "")

	// ── Output ───────────────────────────────────────────────────────────────
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}

	safeName := sanitizeFilename(payload.Header.Nama)
	filename := fmt.Sprintf("BeritaAcara_%s.pdf", safeName)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
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

// splitTextToLines splits text to fit within a given width
func splitTextToLines(pdf *fpdf.Fpdf, text string, width float64) []string {
	words := strings.Fields(text)
	var lines []string
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
