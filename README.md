# Berita Acara Generator

Aplikasi web untuk membuat dokumen **Berita Acara** PT. Bank Negara Indonesia (Persero) Tbk — Divisi Retail Digital Delivery.

---

## Prasyarat

Pastikan sudah terinstal di komputer kamu:

| Tool | Versi Minimum | Download |
|------|--------------|---------|
| **Node.js** | v18+ | https://nodejs.org |
| **Go** | v1.21+ | https://go.dev/dl |
| **Git** | (opsional) | https://git-scm.com |

---

## Instalasi

### 1. Clone / Download Proyek

```powershell
# Clone (jika pakai git)
git clone <repo-url>
cd BA

# Atau langsung masuk ke folder proyek
cd C:\Users\901809\Documents\BA
```

---

### 2. Install Dependensi Backend (Go)

```powershell
cd backend
go mod tidy
```

> Perintah ini akan mengunduh semua library Go yang dibutuhkan:
> - `gin-gonic/gin` — HTTP router
> - `gin-contrib/cors` — CORS middleware
> - `go-pdf/fpdf` — PDF generator

---

### 3. Install Dependensi Frontend (Next.js)

```powershell
cd frontend
npm install
```

> Menginstal Next.js, React, Tailwind CSS, lucide-react, dan semua dependensi lainnya.

---

## Cara Menjalankan

Buka **dua terminal** secara bersamaan.

### Terminal 1 — Jalankan Backend (Go API)

```powershell
cd C:\Users\901809\Documents\BA\backend
go run .
```

Jika berhasil, akan muncul:
```
🚀 Berita Acara API listening on :8080
```

---

### Terminal 2 — Jalankan Frontend (Next.js)

```powershell
cd C:\Users\901809\Documents\BA\frontend
npm run dev
```

Jika berhasil, akan muncul:
```
▲ Next.js
- Local: http://localhost:3000
✓ Ready
```

---

### Buka di Browser

Setelah kedua server berjalan, buka:

```
http://localhost:3000
```

---

## Penggunaan

1. **Isi Identitas** — Nama, NPP BNI, Departement, Kelompok
2. **Tambah Data Kehadiran** — Klik tombol `+ Tambah Baris`, isi tanggal (HARI terisi otomatis), jam datang, jam pulang, dan keterangan
3. **Isi Penandatangan** — Saksi, Hormat Saya, Menyetujui (nama & jabatan)
4. **Export Dokumen**:
   - `Export Word (.docx)` — Unduh file Word
   - `Export PDF` — Unduh file PDF
   - `Preview` — Lihat tampilan dokumen sebelum diunduh

---

## Struktur Proyek

```
BA/
├── backend/                  ← Go REST API (port 8080)
│   ├── main.go
│   ├── go.mod
│   ├── handlers/
│   │   ├── generate.go       ← Handler PDF
│   │   └── docx.go           ← Handler DOCX
│   └── models/
│       └── berita_acara.go   ← Model data
│
└── frontend/                 ← Next.js App (port 3000)
    ├── app/
    │   ├── layout.tsx
    │   ├── page.tsx
    │   └── globals.css
    ├── components/
    │   ├── BeritaAcaraForm.tsx
    │   ├── AttendanceTable.tsx
    │   └── PreviewModal.tsx
    ├── types/
    │   └── index.ts
    └── .env.local            ← NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Konfigurasi

### Mengubah URL API Backend

Edit file `frontend/.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

Ganti dengan URL server Go kamu jika berjalan di host/port yang berbeda.

---

## Troubleshooting

| Masalah | Solusi |
|---------|--------|
| Port 8080 sudah dipakai | Ganti port di `backend/main.go` baris `r.Run(":8080")` |
| Port 3000 sudah dipakai | Jalankan `npm run dev -- -p 3001` |
| CORS error di browser | Pastikan backend sudah berjalan di port 8080 |
| `go mod tidy` gagal | Pastikan koneksi internet aktif dan Go terinstal dengan benar |
| `npm install` gagal | Hapus folder `node_modules` lalu jalankan ulang |
