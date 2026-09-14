'use client';

import { useState } from 'react';
import { Header, Row, Signers } from '@/types';
import { buildPayload, exportBeritaAcara } from '@/services/berita-acara-service';

export function useBeritaAcaraExport() {
  const [loadingDocx, setLoadingDocx] = useState(false);
  const [loadingPdf, setLoadingPdf] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const clearError = () => setError(null);

  const handleExport = async (
    type: 'docx' | 'pdf',
    header: Header,
    rows: Row[],
    signers: Signers
  ) => {
    setError(null);
    const payload = buildPayload(header, rows, signers);
    const setSelf = type === 'docx' ? setLoadingDocx : setLoadingPdf;

    setSelf(true);
    try {
      await exportBeritaAcara(type, payload);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Terjadi kesalahan. Pastikan backend berjalan.');
    } finally {
      setSelf(false);
    }
  };

  return {
    loadingDocx,
    loadingPdf,
    error,
    clearError,
    handleExport,
  };
}
