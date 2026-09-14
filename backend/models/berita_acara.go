package models

// Header contains the personal information fields
type Header struct {
	Nama        string `json:"nama"`
	NppBni      string `json:"npp_bni"`
	Departement string `json:"departement"`
	Kelompok    string `json:"kelompok"`
}

// Row represents one attendance record
type Row struct {
	Tanggal    string `json:"tanggal"`
	Hari       string `json:"hari"`
	JamDatang  string `json:"jam_datang"`
	JamPulang  string `json:"jam_pulang"`
	Keterangan string `json:"keterangan"`
}

// Signers contains the names for the signature block
type Signers struct {
	Saksi             string `json:"saksi"`
	HormatSaya        string `json:"hormat_saya"`
	MenyetujuiNama    string `json:"menyetujui_nama"`
	MenyetujuiJabatan string `json:"menyetujui_jabatan"`
}

// BeritaAcara is the top-level request payload
type BeritaAcara struct {
	Header  Header `json:"header"`
	Rows    []Row  `json:"rows"`
	Signers Signers `json:"signers"`
}
