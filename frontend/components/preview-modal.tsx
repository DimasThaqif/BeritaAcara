'use client';

import { BeritaAcaraPayload, Row } from '@/types';
import { formatTanggal } from '@/utils/date-utils';
import { X } from 'lucide-react';

interface Props {
  readonly payload: {
    readonly header: BeritaAcaraPayload['header'];
    readonly signers: BeritaAcaraPayload['signers'];
    readonly rows: readonly Row[];
  };
  readonly onClose: () => void;
}

export default function PreviewModal({ payload, onClose }: Readonly<Props>) {
  const { header, rows, signers } = payload;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center bg-black/60 backdrop-blur-sm overflow-y-auto py-8 px-4">
      <div className="relative bg-white w-full max-w-3xl rounded-xl shadow-2xl">

        {/* Modal toolbar */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 bg-gray-50 rounded-t-xl">
          <h2 className="text-base font-semibold text-gray-700">Preview Berita Acara</h2>
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-gray-500 hover:bg-gray-200 hover:text-gray-700 transition-colors"
          >
            <X size={18} />
          </button>
        </div>

        {/* Document preview — styled like real document */}
        <div className="p-10 font-serif text-[13px] text-gray-900 leading-relaxed" style={{ fontFamily: 'Arial, sans-serif' }}>

          {/* BNI Header */}
          <p className="font-bold text-[13px]">PT. BANK NEGARA INDONESIA (Persero) Tbk</p>
          <p className="font-bold text-[13px] mb-6">DIVISI RETAIL DIGITAL DELIVERY</p>

          {/* Title */}
          <p className="text-center font-bold underline text-[15px] mb-6 tracking-wide">
            BERITA ACARA
          </p>

          {/* Intro */}
          <p className="mb-3">Yang bertandatangan di bawah ini menerangkan bahwa:</p>

          {/* KV Block */}
          <table className="mb-5" style={{ borderCollapse: 'collapse' }}>
            <tbody>
              {[
                ['Nama', header.nama],
                ['NPP BNI', header.npp_bni],
                ['Departement', header.departement],
                ['Kelompok', header.kelompok],
              ].map(([label, value]) => (
                <tr key={label}>
                  <td className="pr-2 py-0.5 w-32 align-top">{label}</td>
                  <td className="pr-2 py-0.5 w-4 align-top">:</td>
                  <td className="py-0.5 align-top">{value}</td>
                </tr>
              ))}
            </tbody>
          </table>

          {/* Attendance Table */}
          <table className="w-full text-[12px] mb-8" style={{ borderCollapse: 'collapse' }}>
            <thead>
              <tr className="bg-gray-200">
                {['TANGGAL','HARI','JAM DATANG','JAM PULANG','KETERANGAN'].map((h) => (
                  <th
                    key={h}
                    className="border border-gray-400 px-2 py-1.5 text-center font-bold"
                  >
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.length === 0 && (
                <tr>
                  <td colSpan={5} className="border border-gray-400 text-center py-4 text-gray-400 italic">
                    Tidak ada data
                  </td>
                </tr>
              )}
              {rows.map((row) => (
                <tr key={row.id}>
                  <td className="border border-gray-400 px-2 py-1.5 text-center">
                    {formatTanggal(row.tanggal)}
                  </td>
                  <td className="border border-gray-400 px-2 py-1.5 text-center">{row.hari}</td>
                  <td className="border border-gray-400 px-2 py-1.5 text-center">{row.jam_datang}</td>
                  <td className="border border-gray-400 px-2 py-1.5 text-center">{row.jam_pulang}</td>
                  <td className="border border-gray-400 px-2 py-1.5">{row.keterangan}</td>
                </tr>
              ))}
            </tbody>
          </table>

          {/* Signer Block */}
          <table className="w-full text-[12px] mb-6" style={{ borderCollapse: 'collapse' }}>
            <tbody>
              <tr>
                <td className="w-1/2 align-top pb-1 text-center">Saksi</td>
                <td className="w-1/2 align-top pb-1 text-center">Hormat Saya,</td>
              </tr>
              <tr>
                <td className="h-16"></td>
                <td className="h-16"></td>
              </tr>
              <tr>
                <td className="font-bold pt-1 text-center">{signers.saksi}</td>
                <td className="font-bold pt-1 text-center">{signers.hormat_saya}</td>
              </tr>
            </tbody>
          </table>

          <div className="text-center mt-6">
            <p className="mb-14">Menyetujui</p>
            <p className="font-bold">{signers.menyetujui_nama}</p>
            <p>{signers.menyetujui_jabatan}</p>
          </div>
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 rounded-b-xl text-right">
          <button
            onClick={onClose}
            className="px-5 py-2 rounded-lg bg-gray-200 text-gray-700 hover:bg-gray-300 transition-colors text-sm font-medium"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  );
}
