import { Header, Row, Signers, BeritaAcaraPayload } from '@/types';
import { formatTanggal } from '@/utils/date-utils';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

export function buildPayload(header: Header, rows: Row[], signers: Signers): BeritaAcaraPayload {
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

export async function downloadFile(url: string, payload: BeritaAcaraPayload, filename: string): Promise<void> {
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

export async function exportBeritaAcara(
  type: 'docx' | 'pdf',
  payload: BeritaAcaraPayload,
  customFilename?: string
): Promise<void> {
  const ext = type === 'docx' ? 'docx' : 'pdf';
  const safeName = (payload.header.nama || 'BeritaAcara').replace(/\s+/g, '_');
  const filename = customFilename ?? `BeritaAcara_${safeName}.${ext}`;
  const url = `${API_BASE_URL}/api/generate/${type}`;

  await downloadFile(url, payload, filename);
}
