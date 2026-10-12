import React, { useRef, useState } from 'react';
import { toast } from 'sonner';
import { Camera, Loader2, X } from 'lucide-react';
import { uploadImage } from '../../lib/image';
import { getErrorMessage } from '../../lib/utils';

interface ImageUploadProps {
  value: string;
  onChange: (url: string) => void;
  label?: string;
}

/** Pick or take a photo; it is uploaded right away and its URL passed up. */
export const ImageUpload: React.FC<ImageUploadProps> = ({ value, onChange, label = 'Pilih / ambil foto' }) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);

  const pick = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = '';
    if (!file) return;
    setUploading(true);
    try {
      onChange(await uploadImage(file));
    } catch (err) {
      toast.error(getErrorMessage(err, 'Gagal mengunggah foto. Coba lagi.'));
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="flex items-center gap-3">
      {value ? (
        <div className="relative h-24 w-24 flex-shrink-0 overflow-hidden rounded-2xl border border-slate-200 dark:border-slate-700">
          <img src={value} alt="Foto terpilih" className="h-full w-full object-cover" />
          <button
            type="button"
            onClick={() => onChange('')}
            className="absolute right-1 top-1 rounded-full bg-slate-900/70 p-1 text-white"
            aria-label="Hapus foto"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </div>
      ) : null}
      <button
        type="button"
        onClick={() => inputRef.current?.click()}
        disabled={uploading}
        className="flex items-center gap-2 rounded-xl border-2 border-dashed border-slate-300 px-4 py-3 text-sm font-semibold text-slate-600 hover:border-primary-500 hover:text-primary-600 disabled:opacity-60 dark:border-slate-600 dark:text-slate-300"
      >
        {uploading ? <Loader2 className="h-5 w-5 animate-spin" /> : <Camera className="h-5 w-5" />}
        {uploading ? 'Mengunggah...' : value ? 'Ganti foto' : label}
      </button>
      <input ref={inputRef} type="file" accept="image/jpeg,image/png,image/webp,image/*" className="hidden" onChange={pick} />
    </div>
  );
};
