import BeritaAcaraForm from '@/components/berita-acara-form';
import { FileSignature } from 'lucide-react';

export default function HomePage() {
  return (
    <main className="min-h-screen bg-linear-to-br from-slate-100 via-blue-50 to-slate-100">

      {/* Top nav bar */}
      <header className="sticky top-0 z-40 border-b border-white/60 bg-white/80 backdrop-blur-md shadow-sm">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 h-14 flex items-center gap-3">
          <div className="p-1.5 rounded-lg bg-blue-600 text-white">
            <FileSignature size={18} />
          </div>
          <div>
            <span className="font-bold text-gray-800 text-sm">Berita Acara Generator</span>
            <span className="ml-2 text-xs text-gray-400 font-normal hidden sm:inline">
              PT. Bank Negara Indonesia • Divisi Retail Digital Delivery
            </span>
          </div>
          <div className="ml-auto">
            <span className="px-2.5 py-1 rounded-full bg-green-100 text-green-700 text-xs font-semibold">
              ● BNI Internal
            </span>
          </div>
        </div>
      </header>

      {/* Main content */}
      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-2">

        {/* Page title */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Berita Acara</h1>
          <p className="text-sm text-gray-500 mt-1">
            Isi data di bawah ini, lalu export ke Word atau PDF.
          </p>
        </div>

        <BeritaAcaraForm />
      </div>

      {/* Footer */}
      <footer className="mt-12 pb-8 text-center text-xs text-gray-400">
        © {new Date().getFullYear()} PT. Bank Negara Indonesia (Persero) Tbk — Internal Tool
      </footer>
    </main>
  );
}
