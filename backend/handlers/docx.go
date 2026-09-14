package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"

	"berita-acara-api/models"
)

// buildDOCX constructs a valid .docx (Office Open XML) file from scratch
func buildDOCX(payload models.BeritaAcara) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	files := map[string]string{
		"_rels/.rels":                  relsXML(),
		"[Content_Types].xml":          contentTypesXML(),
		"word/_rels/document.xml.rels": documentRelsXML(),
		"word/document.xml":            buildDocumentXML(payload),
		"word/styles.xml":              stylesXML(),
		"word/settings.xml":            settingsXML(),
	}

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write([]byte(content)); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func relsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
}

func contentTypesXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
  <Override PartName="/word/settings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/>
</Types>`
}

func documentRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings" Target="settings.xml"/>
</Relationships>`
}

func settingsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:settings xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:defaultTabStop w:val="720"/>
</w:settings>`
}

func stylesXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
          xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml">
  <w:style w:type="paragraph" w:default="1" w:styleId="Normal">
    <w:name w:val="Normal"/>
    <w:rPr>
      <w:rFonts w:ascii="Arial" w:hAnsi="Arial"/>
      <w:sz w:val="20"/>
    </w:rPr>
  </w:style>
</w:styles>`
}

func esc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// twip converts mm to twentieths of a point (twips). 1mm ≈ 56.69 twips
func twip(mm float64) int {
	return int(mm * 56.69)
}

func para(content, justify string, spaceAfter int) string {
	jc := ""
	if justify != "" {
		jc = fmt.Sprintf(`<w:jc w:val="%s"/>`, justify)
	}
	sa := ""
	if spaceAfter > 0 {
		sa = fmt.Sprintf(`<w:spacing w:after="%d"/>`, spaceAfter)
	}
	return fmt.Sprintf(`<w:p><w:pPr>%s%s</w:pPr>%s</w:p>`, jc, sa, content)
}

func run(text, bold, underline string, sz int) string {
	rpr := "<w:rPr>"
	if bold == "1" {
		rpr += "<w:b/>"
	}
	if underline != "" {
		rpr += fmt.Sprintf(`<w:u w:val="%s"/>`, underline)
	}
	if sz > 0 {
		rpr += fmt.Sprintf(`<w:sz w:val="%d"/><w:szCs w:val="%d"/>`, sz, sz)
	}
	rpr += `<w:rFonts w:ascii="Arial" w:hAnsi="Arial"/>`
	rpr += "</w:rPr>"
	return fmt.Sprintf(`<w:r>%s<w:t xml:space="preserve">%s</w:t></w:r>`, rpr, esc(text))
}

func kvRow(label, value string) string {
	// label cell ~32mm, colon cell ~4mm, value cell rest (total 160mm)
	labelTwips := twip(32)
	colonTwips := twip(4)
	valueTwips := twip(124)
	return fmt.Sprintf(`<w:tr>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr><w:p><w:r><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="20"/></w:rPr><w:t>%s</w:t></w:r></w:p></w:tc>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr><w:p><w:r><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="20"/></w:rPr><w:t>:</w:t></w:r></w:p></w:tc>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr><w:p><w:r><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="20"/></w:rPr><w:t>%s</w:t></w:r></w:p></w:tc>
</w:tr>`, labelTwips, esc(label), colonTwips, valueTwips, esc(value))
}

func tableCell(text, bold, align string, width int, borders bool) string {
	bdr := ""
	if borders {
		bdr = `<w:tcBorders>
  <w:top w:val="single" w:sz="4" w:color="000000"/>
  <w:left w:val="single" w:sz="4" w:color="000000"/>
  <w:bottom w:val="single" w:sz="4" w:color="000000"/>
  <w:right w:val="single" w:sz="4" w:color="000000"/>
</w:tcBorders>`
	}
	boldAttr := ""
	if bold == "1" {
		boldAttr = "<w:b/>"
	}
	return fmt.Sprintf(`<w:tc>
  <w:tcPr><w:tcW w:w="%d" w:type="dxa"/>%s<w:vAlign w:val="center"/></w:tcPr>
  <w:p>
    <w:pPr><w:jc w:val="%s"/><w:spacing w:before="60" w:after="60"/></w:pPr>
    <w:r><w:rPr>%s<w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">%s</w:t></w:r>
  </w:p>
</w:tc>`, width, bdr, align, boldAttr, esc(text))
}

func buildDocumentXML(payload models.BeritaAcara) string {
	var sb strings.Builder

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"
            xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
            xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"
            xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<w:body>
<w:sectPr>
  <w:pgSz w:w="11906" w:h="16838"/>
  <w:pgMar w:top="1417" w:right="1417" w:bottom="1134" w:left="1417" w:header="709" w:footer="709" w:gutter="0"/>
</w:sectPr>
`)

	// ── PT. BNI Header ──────────────────────────────────────────────────────
	sb.WriteString(para(run("PT. BANK NEGARA INDONESIA (Persero) Tbk", "1", "", 22), "left", 40))
	sb.WriteString(para(run("DIVISI RETAIL DIGITAL DELIVERY", "1", "", 22), "left", 160))

	// ── Title ────────────────────────────────────────────────────────────────
	sb.WriteString(para(run("BERITA ACARA", "1", "single", 26), "center", 160))

	// ── Intro ─────────────────────────────────────────────────────────────────
	sb.WriteString(para(run("Yang bertandatangan di bawah ini menerangkan bahwa:", "", "", 20), "left", 60))

	// ── KV Table (borderless) ────────────────────────────────────────────────
	kvTableW := twip(160)
	sb.WriteString(fmt.Sprintf(`<w:tbl>
<w:tblPr>
  <w:tblW w:w="%d" w:type="dxa"/>
  <w:tblBorders>
    <w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/><w:insideH w:val="none"/><w:insideV w:val="none"/>
  </w:tblBorders>
  <w:tblCellMar><w:top w:w="0" w:type="dxa"/><w:bottom w:w="0" w:type="dxa"/></w:tblCellMar>
</w:tblPr>`, kvTableW))

	sb.WriteString(kvRow("Nama", payload.Header.Nama))
	sb.WriteString(kvRow("NPP BNI", payload.Header.NppBni))
	sb.WriteString(kvRow("Departement", payload.Header.Departement))
	sb.WriteString(kvRow("Kelompok", payload.Header.Kelompok))
	sb.WriteString(`</w:tbl>`)
	sb.WriteString(para("", "left", 120))

	// Attendance Table Column widths in twips for 160mm total
	cw := []int{twip(28), twip(16), twip(23), twip(23), twip(70)}
	totalW := 0
	for _, w := range cw {
		totalW += w
	}

	sb.WriteString(fmt.Sprintf(`<w:tbl>
<w:tblPr>
  <w:tblW w:w="%d" w:type="dxa"/>
  <w:tblBorders>
    <w:top w:val="single" w:sz="4" w:color="000000"/>
    <w:left w:val="single" w:sz="4" w:color="000000"/>
    <w:bottom w:val="single" w:sz="4" w:color="000000"/>
    <w:right w:val="single" w:sz="4" w:color="000000"/>
    <w:insideH w:val="single" w:sz="4" w:color="000000"/>
    <w:insideV w:val="single" w:sz="4" w:color="000000"/>
  </w:tblBorders>
</w:tblPr>`, totalW))

	// Header row
	headers := []string{"TANGGAL", "HARI", "JAM DATANG", "JAM PULANG", "KETERANGAN"}
	sb.WriteString("<w:tr>")
	for i, h := range headers {
		sb.WriteString(fmt.Sprintf(`<w:tc>
  <w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="F0F0F0"/><w:vAlign w:val="center"/></w:tcPr>
  <w:p>
    <w:pPr><w:jc w:val="center"/><w:spacing w:before="60" w:after="60"/></w:pPr>
    <w:r><w:rPr><w:b/><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="18"/></w:rPr><w:t>%s</w:t></w:r>
  </w:p>
</w:tc>`, cw[i], esc(h)))
	}
	sb.WriteString("</w:tr>")

	// Data rows
	for _, row := range payload.Rows {
		vals := []string{row.Tanggal, row.Hari, row.JamDatang, row.JamPulang, row.Keterangan}
		aligns := []string{"center", "center", "center", "center", "left"}
		sb.WriteString("<w:tr>")
		for i, v := range vals {
			sb.WriteString(tableCell(v, "", aligns[i], cw[i], true))
		}
		sb.WriteString("</w:tr>")
	}

	sb.WriteString(`</w:tbl>`)
	sb.WriteString(para("", "left", 360))

	// ── Signer block ─────────────────────────────────────────────────────────
	signerTableW := totalW
	colW := signerTableW / 2

	// Saksi | Hormat Saya
	sb.WriteString(fmt.Sprintf(`<w:tbl>
<w:tblPr>
  <w:tblW w:w="%d" w:type="dxa"/>
  <w:tblBorders>
    <w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/>
    <w:insideH w:val="none"/><w:insideV w:val="none"/>
  </w:tblBorders>
</w:tblPr>`, signerTableW))

	// Labels row (Centered in each column)
	sb.WriteString(fmt.Sprintf(`<w:tr>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="20"/></w:rPr><w:t>Saksi</w:t></w:r></w:p></w:tc>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="20"/></w:rPr><w:t>Hormat Saya,</w:t></w:r></w:p></w:tc>
</w:tr>`, colW, colW))

	// Signature space rows (3 empty rows)
	for i := 0; i < 3; i++ {
		sb.WriteString(fmt.Sprintf(`<w:tr>
  <w:trPr><w:trHeight w:val="450"/></w:trPr>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr><w:p/></w:tc>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr><w:p/></w:tc>
</w:tr>`, colW, colW))
	}

	// Names row (Bold, Centered)
	sb.WriteString(fmt.Sprintf(`<w:tr>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="20"/></w:rPr><w:t>%s</w:t></w:r></w:p></w:tc>
  <w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/><w:tcBorders><w:top w:val="none"/><w:left w:val="none"/><w:bottom w:val="none"/><w:right w:val="none"/></w:tcBorders></w:tcPr>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/><w:rFonts w:ascii="Arial" w:hAnsi="Arial"/><w:sz w:val="20"/></w:rPr><w:t>%s</w:t></w:r></w:p></w:tc>
</w:tr>`, colW, esc(payload.Signers.Saksi), colW, esc(payload.Signers.HormatSaya)))

	sb.WriteString(`</w:tbl>`)
	sb.WriteString(para("", "left", 200))

	// Menyetujui section (centered)
	sb.WriteString(para(run("Menyetujui", "", "", 20), "center", 40))

	// Signature space
	for i := 0; i < 3; i++ {
		sb.WriteString(para("", "center", 0))
	}

	sb.WriteString(para(run(payload.Signers.MenyetujuiNama, "1", "", 20), "center", 40))
	sb.WriteString(para(run(payload.Signers.MenyetujuiJabatan, "", "", 20), "center", 40))

	sb.WriteString(`</w:body>
</w:document>`)

	return sb.String()
}
