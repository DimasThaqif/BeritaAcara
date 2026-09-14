'use client';

import { Row } from '@/types';
import { Trash2 } from 'lucide-react';

const HARI_MAP: Record<number, string> = {
  0: 'Minggu',
  1: 'Senin',
  2: 'Selasa',
  3: 'Rabu',
  4: 'Kamis',
  5: 'Jumat',
  6: 'Sabtu',
};

function getHari(dateStr: string): string {
  if (!dateStr) return '';
  // dateStr is ISO format YYYY-MM-DD from date input
  const [year, month, day] = dateStr.split('-').map(Number);
  const d = new Date(year, month - 1, day);
  return HARI_MAP[d.getDay()] ?? '';
}

function formatTanggal(dateStr: string): string {
  if (!dateStr) return '';
  const MONTHS = [
    'Januari','Februari','Maret','April','Mei','Juni',
    'Juli','Agustus','September','Oktober','November','Desember',
  ];
  const [year, month, day] = dateStr.split('-').map(Number);
  return `${day} ${MONTHS[month - 1]} ${year}`;
}

interface Props {
  rows: Row[];
  onChange: (rows: Row[]) => void;
}

export default function AttendanceTable({ rows, onChange }: Props) {
  const addRow = () => {
    const newRow: Row = {
      id: crypto.randomUUID(),
      tanggal: '',
      hari: '',
      jam_datang: '',
      jam_pulang: '',
      keterangan: '',
    };
    onChange([...rows, newRow]);
  };

  const deleteRow = (id: string) => {
    onChange(rows.filter((r) => r.id !== id));
  };

  const updateRow = (id: string, field: keyof Row, value: string) => {
    onChange(
      rows.map((r) => {
        if (r.id !== id) return r;
        const updated = { ...r, [field]: value };
        if (field === 'tanggal') {
          updated.hari = getHari(value);
          updated.tanggal = value; // keep raw for date input
        }
        return updated;
      })
    );
  };

  const inputBase =
    'w-full bg-transparent text-sm text-gray-800 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded px-1 py-0.5';

  return (
    <div className="space-y-3">
      <div className="overflow-x-auto rounded-lg border border-gray-200 shadow-sm">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-100 text-gray-700 uppercase text-xs font-semibold tracking-wider">
            <tr>
              <th className="px-3 py-3 text-center w-36">Tanggal</th>
              <th className="px-3 py-3 text-center w-28">Hari</th>
              <th className="px-3 py-3 text-center w-28">Jam Datang</th>
              <th className="px-3 py-3 text-center w-28">Jam Pulang</th>
              <th className="px-3 py-3 text-left">Keterangan</th>
              <th className="px-3 py-3 text-center w-12"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100 bg-white">
            {rows.length === 0 && (
              <tr>
                <td colSpan={6} className="text-center py-8 text-gray-400 italic">
                  Belum ada data. Klik &ldquo;Tambah Baris&rdquo; untuk mulai.
                </td>
              </tr>
            )}
            {rows.map((row, idx) => (
              <tr
                key={row.id}
                className={idx % 2 === 0 ? 'bg-white' : 'bg-blue-50/30'}
              >
                {/* TANGGAL */}
                <td className="px-2 py-2">
                  <input
                    type="date"
                    value={row.tanggal}
                    onChange={(e) => updateRow(row.id, 'tanggal', e.target.value)}
                    className={inputBase + ' text-center'}
                  />
                  {row.tanggal && (
                    <div className="text-xs text-gray-400 text-center mt-0.5">
                      {formatTanggal(row.tanggal)}
                    </div>
                  )}
                </td>

                {/* HARI — auto-filled */}
                <td className="px-2 py-2">
                  <div className="text-center font-medium text-gray-700 text-sm">
                    {row.hari || <span className="text-gray-300">—</span>}
                  </div>
                </td>

                {/* JAM DATANG */}
                <td className="px-2 py-2">
                  <input
                    type="time"
                    value={row.jam_datang}
                    onChange={(e) => updateRow(row.id, 'jam_datang', e.target.value)}
                    className={inputBase + ' text-center'}
                  />
                </td>

                {/* JAM PULANG */}
                <td className="px-2 py-2">
                  <input
                    type="time"
                    value={row.jam_pulang}
                    onChange={(e) => updateRow(row.id, 'jam_pulang', e.target.value)}
                    className={inputBase + ' text-center'}
                  />
                </td>

                {/* KETERANGAN */}
                <td className="px-2 py-2">
                  <input
                    type="text"
                    value={row.keterangan}
                    onChange={(e) => updateRow(row.id, 'keterangan', e.target.value)}
                    placeholder="Keterangan..."
                    className={inputBase}
                  />
                </td>

                {/* DELETE */}
                <td className="px-2 py-2 text-center">
                  <button
                    onClick={() => deleteRow(row.id)}
                    title="Hapus baris"
                    className="p-1.5 rounded-md text-red-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <Trash2 size={15} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <button
        type="button"
        onClick={addRow}
        className="flex items-center gap-2 px-4 py-2 rounded-lg border-2 border-dashed border-blue-300 text-blue-600 hover:bg-blue-50 hover:border-blue-400 transition-all text-sm font-medium"
      >
        <span className="text-lg leading-none">+</span> Tambah Baris
      </button>
    </div>
  );
}
