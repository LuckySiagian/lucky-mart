import { api } from './axios';

const MAX_SIDE = 1600;

/**
 * Shrinks a phone photo before upload (max 1600px, JPEG) so it uploads fast
 * on mobile data and stays under the server's 5 MB limit. Falls back to the
 * original file when the browser cannot decode it.
 */
export async function compressImage(file: File): Promise<Blob> {
  if (!file.type.startsWith('image/') || file.type === 'image/gif') return file;
  try {
    const bitmap = await createImageBitmap(file);
    const scale = Math.min(1, MAX_SIDE / Math.max(bitmap.width, bitmap.height));
    if (scale === 1 && file.size < 1_000_000) return file;

    const canvas = document.createElement('canvas');
    canvas.width = Math.round(bitmap.width * scale);
    canvas.height = Math.round(bitmap.height * scale);
    const ctx = canvas.getContext('2d');
    if (!ctx) return file;
    ctx.fillStyle = '#fff'; // transparent PNG areas become white, not black
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    ctx.drawImage(bitmap, 0, 0, canvas.width, canvas.height);
    const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/jpeg', 0.85));
    return blob && blob.size < file.size ? blob : file;
  } catch {
    return file;
  }
}

/** Uploads a photo and returns its URL (e.g. /uploads/images/2026/10/x.jpg). */
export async function uploadImage(file: File): Promise<string> {
  const blob = await compressImage(file);
  const form = new FormData();
  form.append('file', blob, blob === file ? file.name : 'foto.jpg');
  const res = await api.post<{ data: { url: string } }>('/uploads/images', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 60000,
  });
  return res.data.data.url;
}
