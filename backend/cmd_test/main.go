package main

import (
	"log"
	"os"

	"berita-acara-api/handlers"
	"berita-acara-api/models"
)

func main() {
	payload := models.BeritaAcara{
		Header: models.Header{
			Nama:        "Dimas Thaqif Attaulah",
			NppBni:      "901809",
			Departement: "WDL",
			Kelompok:    "BDD",
		},
		Rows: []models.Row{
			{Tanggal: "14 September 2026", Hari: "Senin", JamDatang: "15:19", JamPulang: "15:19", Keterangan: "tes tes"},
			{Tanggal: "15 September 2026", Hari: "Selasa", JamDatang: "15:20", JamPulang: "15:20", Keterangan: "testes"},
		},
		Signers: models.Signers{
			Saksi:             "KIPLI",
			HormatSaya:        "Dimas Thaqif Attaulah",
			MenyetujuiNama:    "Daniel Harry Simbolon",
			MenyetujuiJabatan: "Team Lead",
		},
	}

	pdfBytes, err := handlers.BuildPDF(payload)
	if err != nil {
		log.Fatalf("BuildPDF failed: %v", err)
	}

	if err := os.WriteFile(`C:\Users\901809\Documents\BA\test_output.pdf`, pdfBytes, 0644); err != nil {
		log.Fatalf("WriteFile failed: %v", err)
	}

	log.Println("Generated test_output.pdf")
}
