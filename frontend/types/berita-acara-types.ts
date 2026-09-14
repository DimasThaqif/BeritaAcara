export interface Header {
  nama: string;
  npp_bni: string;
  departement: string;
  kelompok: string;
}

export interface Row {
  id: string; // client-side only, not sent to API
  tanggal: string;
  hari: string;
  jam_datang: string;
  jam_pulang: string;
  keterangan: string;
}

export interface Signers {
  saksi: string;
  hormat_saya: string;
  menyetujui_nama: string;
  menyetujui_jabatan: string;
}

export interface BeritaAcaraPayload {
  header: Header;
  rows: Omit<Row, 'id'>[];
  signers: Signers;
}
