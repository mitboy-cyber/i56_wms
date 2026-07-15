import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';

/**
 * Export an array of objects to CSV and trigger download.
 * Uses current column keys to match API response fields.
 */
export function exportToCSV<T extends Record<string, unknown>>(
  data: T[],
  columns: { key: string; label: string }[],
  filename: string,
): void {
  const headers = columns.map((c) => c.label);
  const keys = columns.map((c) => c.key);
  const rows = data.map((row) => keys.map((k) => formatCell(row[k])));
  const csv = [headers, ...rows].map((r) => r.map(quoteCSV).join(',')).join('\n');

  const bom = '\uFEFF';
  const blob = new Blob([bom + csv], { type: 'text/csv;charset=utf-8' });
  saveAs(blob, `${filename}.csv`);
}

/**
 * Export an array of objects to XLSX and trigger download.
 */
export function exportToXLSX<T extends Record<string, unknown>>(
  data: T[],
  columns: { key: string; label: string }[],
  filename: string,
): void {
  const headers = columns.map((c) => c.label);
  const keys = columns.map((c) => c.key);
  const rows = data.map((row) => keys.map((k) => formatCell(row[k])));

  const ws = XLSX.utils.aoa_to_sheet([headers, ...rows]);
  // Set column widths
  ws['!cols'] = columns.map(() => ({ wch: 20 }));
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
  const buf = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
  const blob = new Blob([buf], { type: 'application/octet-stream' });
  saveAs(blob, `${filename}.xlsx`);
}

function formatCell(value: unknown): string {
  if (value === null || value === undefined) return '';
  if (typeof value === 'object') return JSON.stringify(value);
  return String(value);
}

function quoteCSV(cell: string): string {
  if (cell.includes(',') || cell.includes('"') || cell.includes('\n')) {
    return `"${cell.replace(/"/g, '""')}"`;
  }
  return cell;
}
