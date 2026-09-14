'use client';

import { useState } from 'react';
import { Header, Row, Signers, BeritaAcaraPayload } from '@/types';
import AttendanceTable from './AttendanceTable';
import PreviewModal from './PreviewModal';
import { FileText, FileDown, Eye, Loader2, Building2 } from 'lucide-react';

const MONTHS = [
  'Januari','Februari','Maret','April','Mei','Juni',
  'Juli','Agustus','September','Oktober','November','Desember',
];

function formatTanggal(dateStr: string): string {
  if (!dateStr) return '';
  const [year, month, day] = dateStr.split('-').map(Number);
  return `${day} ${MONTHS[month - 1]} ${year}`;
}

function buildPayload(header: Header, rows: Row[], signers: Signers): BeritaAcaraPayload {
  return {
    header,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    rows: rows.map(({ id, ...rest }) => ({
      ...rest,
      tanggal: formatTanggal(rest.tanggal),
    })),
    signers,
  };
}

async function downloadFile(url: string, payload: BeritaAcaraPayload, filename: string) {
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(err.error ?? 'Server error');
  }

  const blob = await res.blob();
  const objectUrl = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = objectUrl;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(objectUrl);
}

const API = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

const inputClass =
  'w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition';

const labelClass = 'block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1';

interface FormField {
  key: keyof Header;
  label: string;
  placeholder: string;
}

const HEADER_FIELDS: FormField[] = [
  { key: 'nama',        label: 'Nama',        placeholder: 'Dimas Thaqif Attaulah' },
  { key: 'npp_bni',    label: 'NPP BNI',     placeholder: '901809' },
  { key: 'departement',label: 'Departement', placeholder: 'Wholesale Digital Delivery' },
  { key: 'kelompok',   label: 'Kelompok',    placeholder: 'BDD' },
];

export default function BeritaAcaraForm() {
  const [header, setHeader] = useState<Header>({
    nama: '', npp_bni: '', departement: '', kelompok: '',
  });

  const [rows, setRows] = useState<Row[]>([]);

  const [signers, setSigners] = useState<Signers>({
    saksi: '',
    hormat_saya: '',
    menyetujui_nama: '',
    menyetujui_jabatan: '',
  });

  const [loadingDocx, setLoadingDocx] = useState(false);
  const [loadingPdf, setLoadingPdf]   = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [error, setError]             = useState<string | null>(null);

  const updateHeader = (key: keyof Header, value: string) => {
    setHeader((prev) => {
      const next = { ...prev, [key]: value };
      // Keep hormat_saya in sync with nama if user hasn't manually changed it
      if (key === 'nama' && signers.hormat_saya === prev.nama) {
        setSigners((s) => ({ ...s, hormat_saya: value }));
      }
      return next;
    });
  };

  const handleExport = async (type: 'docx' | 'pdf') => {
    setError(null);
    const payload = buildPayload(header, rows, signers);
    const setSelf = type === 'docx' ? setLoadingDocx : setLoadingPdf;
    const url     = `${API}/api/generate/${type}`;
    const ext     = type === 'docx' ? 'docx' : 'pdf';
    const safeName = (header.nama || 'BeritaAcara').replace(/\s+/g, '_');
    const filename = `BeritaAcara_${safeName}.${ext}`;

    setSelf(true);
    try {
      await downloadFile(url, payload, filename);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Terjadi kesalahan. Pastikan backend berjalan.');
    } finally {
      setSelf(false);
    }
  };

  const payload = buildPayload(header, rows, signers);

  return (
    <>
      {showPreview && (
        <PreviewModal
          payload={{ ...payload, rows }}
          onClose={() => setShowPreview(false)}
        />
      )}

      <div className="space-y-6">

        {/* ── Company Letterhead Badge ─────────────────────────────────── */}
        <div className="flex items-start gap-4 p-5 rounded-xl bg-gradient-to-r from-blue-700 to-blue-900 text-white shadow-lg">
          <div className="mt-0.5 p-2 rounded-lg bg-white/20">
            <Building2 size={22} />
          </div>
          <div>
            <p className="font-bold text-base leading-tight">PT. BANK NEGARA INDONESIA (Persero) Tbk</p>
            <p className="text-blue-200 text-sm mt-0.5">DIVISI RETAIL DIGITAL DELIVERY</p>
          </div>
        </div>

        {/* ── Header Fields ────────────────────────────────────────────── */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
          <h2 className="text-sm font-bold text-gray-700 uppercase tracking-wider mb-4 flex items-center gap-2">
            <span className="inline-block w-1 h-4 rounded-full bg-blue-600"></span>
            Identitas
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {HEADER_FIELDS.map(({ key, label, placeholder }) => (
              <div key={key}>
                <label className={labelClass}>{label}</label>
                <input
                  type="text"
                  value={header[key]}
                  onChange={(e) => updateHeader(key, e.target.value)}
                  placeholder={placeholder}
                  className={inputClass}
                />
              </div>
            ))}
          </div>
        </div>

        {/* ── Attendance Table ─────────────────────────────────────────── */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
          <h2 className="text-sm font-bold text-gray-700 uppercase tracking-wider mb-4 flex items-center gap-2">
            <span className="inline-block w-1 h-4 rounded-full bg-blue-600"></span>
            Data Kehadiran
          </h2>
          <AttendanceTable rows={rows} onChange={setRows} />
        </div>

        {/* ── Signer Section ───────────────────────────────────────────── */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
          <h2 className="text-sm font-bold text-gray-700 uppercase tracking-wider mb-4 flex items-center gap-2">
            <span className="inline-block w-1 h-4 rounded-full bg-blue-600"></span>
            Penandatangan
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className={labelClass}>Saksi</label>
              <input
                type="text"
                value={signers.saksi}
                onChange={(e) => setSigners((s) => ({ ...s, saksi: e.target.value }))}
                placeholder="Nama saksi"
                className={inputClass}
              />
            </div>
            <div>
              <label className={labelClass}>Hormat Saya</label>
              <input
                type="text"
                value={signers.hormat_saya || header.nama}
                onChange={(e) => setSigners((s) => ({ ...s, hormat_saya: e.target.value }))}
                placeholder={header.nama || 'Nama pembuat'}
                className={inputClass}
              />
            </div>
            <div>
              <label className={labelClass}>Menyetujui — Nama</label>
              <input
                type="text"
                value={signers.menyetujui_nama}
                onChange={(e) => setSigners((s) => ({ ...s, menyetujui_nama: e.target.value }))}
                placeholder="Daniel Harry Hasudungan Simbolon"
                className={inputClass}
              />
            </div>
            <div>
              <label className={labelClass}>Menyetujui — Jabatan</label>
              <input
                type="text"
                value={signers.menyetujui_jabatan}
                onChange={(e) => setSigners((s) => ({ ...s, menyetujui_jabatan: e.target.value }))}
                placeholder="Team Lead"
                className={inputClass}
              />
            </div>
          </div>
        </div>

        {/* ── Error Banner ─────────────────────────────────────────────── */}
        {error && (
          <div className="flex items-start gap-3 p-4 rounded-xl bg-red-50 border border-red-200 text-red-700 text-sm">
            <span className="text-red-500 mt-0.5">⚠</span>
            <span>{error}</span>
            <button onClick={() => setError(null)} className="ml-auto text-red-400 hover:text-red-600">✕</button>
          </div>
        )}

        {/* ── Action Buttons ───────────────────────────────────────────── */}
        <div className="flex flex-wrap gap-3">
          <button
            onClick={() => handleExport('docx')}
            disabled={loadingDocx}
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-blue-600 hover:bg-blue-700 disabled:opacity-60 disabled:cursor-not-allowed text-white text-sm font-semibold shadow-md hover:shadow-lg active:scale-95 transition-all"
          >
            {loadingDocx
              ? <Loader2 size={16} className="animate-spin" />
              : <FileText size={16} />}
            {loadingDocx ? 'Generating...' : 'Export Word (.docx)'}
          </button>

          <button
            onClick={() => handleExport('pdf')}
            disabled={loadingPdf}
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-red-600 hover:bg-red-700 disabled:opacity-60 disabled:cursor-not-allowed text-white text-sm font-semibold shadow-md hover:shadow-lg active:scale-95 transition-all"
          >
            {loadingPdf
              ? <Loader2 size={16} className="animate-spin" />
              : <FileDown size={16} />}
            {loadingPdf ? 'Generating...' : 'Export PDF'}
          </button>

          <button
            onClick={() => setShowPreview(true)}
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-white border border-gray-200 hover:bg-gray-50 text-gray-700 text-sm font-semibold shadow-sm hover:shadow-md active:scale-95 transition-all"
          >
            <Eye size={16} />
            Preview
          </button>
        </div>
      </div>
    </>
  );
}
