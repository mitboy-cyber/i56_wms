import { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import client from '@/api/client';
import Modal from '@/components/Modal';
import { DynamicFormBuilder } from '@/components/DynamicFormBuilder';
import { exportToCSV, exportToXLSX } from '@/lib/export';
import type { ColumnDef, FieldDef, FilterDef } from '@/types';
import {
  Search,
  Filter,
  Columns,
  Download,
  Trash2,
  Pencil,
  ChevronLeft,
  ChevronRight,
  CheckSquare,
  Square,
  X,
  RotateCcw,
} from 'lucide-react';

// ── Types ──

interface GenericListPageProps {
  title: string;
  queryKey: string[];
  queryFn: () => Promise<any>;
  /** API base path for mutations, e.g. '/admin/api/clients' */
  apiBase?: string;
  /** Deprecated: old name for apiBase */
  mutateUrl?: string;
  columns: ColumnDef[];
  fields?: FieldDef[];
  getRowId: (row: Record<string, unknown>, i: number) => string;

  // ── Optional enhancements ──
  /** Enable server-side search (debounced) */
  searchable?: boolean;
  /** Search placeholder text */
  searchPlaceholder?: string;
  /** Advanced filter definitions */
  filterDefs?: FilterDef[];
  /** Enable bulk actions (select + batch delete) */
  enableBulkActions?: boolean;
  /** Enable pagination (reads `total` from API response) */
  pagination?: boolean;
  /** Default page size */
  pageSize?: number;
  /** Enable CSV/XLSX export */
  enableExport?: boolean;
  /** Custom delete mutation: (id) => Promise */
  onDelete?: (row: Record<string, unknown>) => Promise<void>;
  /** Custom batch delete: (ids: string[]) => Promise */
  onBatchDelete?: (ids: string[]) => Promise<void>;
  /** Custom create/update mutation URL override */
  createUrl?: string;
  updateUrl?: string;
  /** Enable drag-to-reorder columns (simple state-based) */
  enableColumnReorder?: boolean;
  /** Persist column prefs to localStorage under this key */
  columnPrefsKey?: string;
}

// ── Helpers ──

function useDebounce<T>(value: T, delay: number): T {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(t);
  }, [value, delay]);
  return debounced;
}

// ── Component ──

export default function GenericListPage(props: GenericListPageProps) {
  const {
    title,
    queryKey,
    queryFn,
    apiBase,
    mutateUrl,
    columns: rawColumns,
    fields,
    getRowId,
    searchable = true,
    searchPlaceholder = '搜索...',
    filterDefs,
    enableBulkActions = true,
    pagination: enablePagination = true,
    pageSize: defaultPageSize = 20,
    enableExport = true,
    onDelete: customDelete,
    onBatchDelete: customBatchDelete,
    createUrl,
    updateUrl,
    enableColumnReorder = false,
    columnPrefsKey,
  } = props;

  const qc = useQueryClient();
  // ── Auto-derive API base from query function ──
  const autoApiBase = useMemo(() => {
    if (apiBase || mutateUrl) return '';
    const fnStr = queryFn.toString();
    const m = fnStr.match(/client\.(?:get|post)\('([^']+)'(?:\)|'\))/) || fnStr.match(/client\.(?:get|post)\("([^"]+)"(?:\)|"\))/);
    return m ? m[1] : '';
  }, [queryFn, apiBase, mutateUrl]);

  const baseUrl = apiBase || mutateUrl || autoApiBase || '';

  // ── Column management (show/hide, reorder, persist) ──
  const [columnOrder, setColumnOrder] = useState<string[]>(() => {
    if (columnPrefsKey) {
      try {
        const saved = localStorage.getItem(`cols:${columnPrefsKey}`);
        if (saved) return JSON.parse(saved) as string[];
      } catch { /* ignore */ }
    }
    return rawColumns.map((c) => c.key);
  });
  const [hiddenColumns, setHiddenColumns] = useState<Set<string>>(new Set(rawColumns.filter((c) => c.hidden).map((c) => c.key)));

  // Sync columnOrder with rawColumns (new columns appear at end)
  useEffect(() => {
    setColumnOrder((prev) => {
      const allKeys = rawColumns.map((c) => c.key);
      const existing = new Set(prev);
      const added = allKeys.filter((k) => !existing.has(k));
      const valid = prev.filter((k) => allKeys.includes(k));
      return [...valid, ...added];
    });
  }, [rawColumns]);

  const saveColumnPrefs = useCallback(
    (order: string[], hidden: Set<string>) => {
      if (columnPrefsKey) {
        localStorage.setItem(`cols:${columnPrefsKey}`, JSON.stringify(order));
        localStorage.setItem(`cols:hidden:${columnPrefsKey}`, JSON.stringify([...hidden]));
      }
    },
    [columnPrefsKey],
  );

  const visibleColumns: ColumnDef[] = useMemo(() => {
    const map = new Map(rawColumns.map((c) => [c.key, c]));
    return columnOrder
      .filter((k) => !hiddenColumns.has(k) && map.has(k))
      .map((k) => map.get(k)!);
  }, [rawColumns, columnOrder, hiddenColumns]);

  const toggleColumn = (key: string) => {
    setHiddenColumns((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      saveColumnPrefs(columnOrder, next);
      return next;
    });
  };

  // ── Filters & Search ──
  const [searchInput, setSearchInput] = useState('');
  const debouncedSearch = useDebounce(searchInput, 350);
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [showFilters, setShowFilters] = useState(false);
  const [showColumnsMenu, setShowColumnsMenu] = useState(false);
  const [showExportMenu, setShowExportMenu] = useState(false);

  // ── Bulk selection ──
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const lastClickedRef = useRef<number>(-1);

  // ── Pagination ──
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(defaultPageSize);

  // ── Query ──
  const activeSearch = debouncedSearch.trim() || undefined;
  const activeFilters = Object.keys(filters).length > 0 ? filters : undefined;

  const searchQueryKey = [queryKey, { search: activeSearch, filters: activeFilters, page, pageSize }];

  const result = useQuery({
    queryKey: searchQueryKey,
    queryFn: () => queryFn(),
  });

  // API response: AxiosResponse wrapping the JSON body.
  // Body can be either a flat array or { data: [...], total?: N }
  const respBody: unknown = (result.data as any)?.data;
  const rawData: Record<string, unknown>[] = Array.isArray(respBody)
    ? respBody
    : (respBody as any)?.data ?? [];
  const total: number | undefined = Array.isArray(respBody)
    ? undefined
    : (respBody as any)?.total;
  const isLoading = (result as any).isLoading;

  const rows = rawData;

  // ── Modal state ──
  const [modalOpen, setModalOpen] = useState(false);
  const [editItem, setEditItem] = useState<Record<string, unknown> | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<{ ids: string[]; label: string } | null>(null);

  const fieldList: FieldDef[] =
    fields ??
    rawColumns.map((c) => ({
      key: c.key,
      label: c.label,
      type: 'text' as const,
    }));

  // ── Mutations ──
  const addMut = useMutation({
    mutationFn: (body: Record<string, unknown>) =>
      createUrl
        ? client.post(createUrl, body)
        : baseUrl
          ? client.post(baseUrl, body)
          : Promise.reject(new Error('No API URL configured')),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey });
      setModalOpen(false);
    },
  });

  const updateMut = useMutation({
    mutationFn: ({ id, body }: { id: string; body: Record<string, unknown> }) =>
      updateUrl
        ? client.put(updateUrl.replace(':id', id), body)
        : baseUrl
          ? client.put(`${baseUrl}/${id}`, body)
          : Promise.reject(new Error('No API URL configured')),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey });
      setModalOpen(false);
    },
  });

  const deleteMut = useMutation({
    mutationFn: async (ids: string[]) => {
      if (customBatchDelete) {
        await customBatchDelete(ids);
      } else if (customDelete && ids.length === 1) {
        // We need the row object — fall back
        const row = rows.find((r) => ids.includes(String(getRowId(r, 0))));
        if (row) await customDelete(row);
      } else if (baseUrl) {
        // Batch delete: delete each individually
        await Promise.all(ids.map((id) => client.delete(`${baseUrl}/${id}`)));
      } else {
        throw new Error('No API URL configured');
      }
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey });
      setSelectedIds(new Set());
      setDeleteConfirm(null);
    },
  });

  // ── Handlers ──
  const openAdd = () => {
    setEditItem(null);
    setModalOpen(true);
  };

  const openEdit = (row: Record<string, unknown>) => {
    setEditItem(row);
    setModalOpen(true);
  };

  const handleSave = (formData: Record<string, unknown>) => {
    if (editItem) {
      const id = String(editItem['id'] ?? getRowId(editItem, 0));
      updateMut.mutate({ id, body: formData });
    } else {
      addMut.mutate(formData);
    }
  };

  const confirmDelete = (ids: string[], label: string) => {
    setDeleteConfirm({ ids, label });
  };

  const executeDelete = () => {
    if (deleteConfirm) {
      deleteMut.mutate(deleteConfirm.ids);
    }
  };

  // ── Bulk selection helpers ──
  const toggleSelect = (id: string, index: number) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      lastClickedRef.current = index;
      return next;
    });
  };

  const toggleSelectAll = () => {
    if (selectedIds.size === rows.length) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(rows.map((r, i) => getRowId(r, i) as string)));
    }
  };

  const isAllSelected = rows.length > 0 && selectedIds.size === rows.length;
  const isPartialSelected = selectedIds.size > 0 && selectedIds.size < rows.length;

  // ── Export ──
  const handleExport = (format: 'csv' | 'xlsx') => {
    const data = selectedIds.size > 0
      ? rows.filter((r, i) => selectedIds.has(getRowId(r, i) as string))
      : rows;
    const cols = visibleColumns.map((c) => ({ key: c.key, label: c.label }));
    const filename = `${title}_${new Date().toISOString().slice(0, 10)}`;
    if (format === 'csv') exportToCSV(data, cols, filename);
    else exportToXLSX(data, cols, filename);
    setShowExportMenu(false);
  };

  // ── Filter clear ──
  const clearFilters = () => {
    setFilters({});
    setSearchInput('');
  };

  const hasActiveFilters = Object.keys(filters).length > 0 || searchInput.trim().length > 0;

  // ── Pagination helpers ──
  const totalPages = total ? Math.ceil(total / pageSize) : 0;

  // ── Reset page when search/filters change ──
  useEffect(() => {
    setPage(1);
  }, [debouncedSearch, filters]);

  // ── Render ──
  return (
    <div>
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-bold text-gray-800">{title}</h2>
        <div className="flex items-center gap-2">
          {/* Column manager button */}
          <div className="relative">
            <button
              onClick={() => setShowColumnsMenu(!showColumnsMenu)}
              className="px-3 py-2 text-sm text-gray-600 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 flex items-center gap-1.5 transition-colors"
            >
              <Columns size={15} />
              <span className="hidden sm:inline">列</span>
            </button>
            {showColumnsMenu && (
              <div className="absolute right-0 top-full mt-1 w-56 bg-white rounded-lg shadow-xl border border-gray-200 z-30 py-1 max-h-64 overflow-y-auto">
                <div className="px-3 py-1.5 text-xs font-medium text-gray-400 uppercase">显示/隐藏列</div>
                {rawColumns.map((c) => (
                  <label
                    key={c.key}
                    className="flex items-center gap-2 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={!hiddenColumns.has(c.key)}
                      onChange={() => toggleColumn(c.key)}
                      className="rounded border-gray-300 text-blue-600"
                    />
                    {c.label}
                  </label>
                ))}
              </div>
            )}
          </div>

          {/* Export button */}
          {enableExport && (
            <div className="relative">
              <button
                onClick={() => setShowExportMenu(!showExportMenu)}
                className="px-3 py-2 text-sm text-gray-600 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 flex items-center gap-1.5 transition-colors"
              >
                <Download size={15} />
                <span className="hidden sm:inline">导出</span>
              </button>
              {showExportMenu && (
                <div className="absolute right-0 top-full mt-1 w-40 bg-white rounded-lg shadow-xl border border-gray-200 z-30 py-1">
                  <button
                    onClick={() => handleExport('csv')}
                    className="w-full text-left px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
                  >
                    CSV (.csv)
                  </button>
                  <button
                    onClick={() => handleExport('xlsx')}
                    className="w-full text-left px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
                  >
                    Excel (.xlsx)
                  </button>
                </div>
              )}
            </div>
          )}

          {/* Add button */}
          <button
            onClick={openAdd}
            className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
          >
            + 添加
          </button>
        </div>
      </div>

      {/* Search & Filters bar */}
      <div className="mb-4 space-y-3">
        {/* Search */}
        {searchable && (
          <div className="flex items-center gap-2">
            <div className="relative flex-1 max-w-md">
              <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder={searchPlaceholder}
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                className="w-full pl-9 pr-8 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
              {searchInput && (
                <button
                  onClick={() => setSearchInput('')}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <X size={14} />
                </button>
              )}
            </div>

            {filterDefs && filterDefs.length > 0 && (
              <button
                onClick={() => setShowFilters(!showFilters)}
                className={`px-3 py-2 text-sm rounded-lg border flex items-center gap-1.5 transition-colors ${
                  showFilters ? 'bg-blue-50 border-blue-300 text-blue-700' : 'border-gray-300 text-gray-600 hover:bg-gray-50'
                }`}
              >
                <Filter size={15} />
                <span className="hidden sm:inline">筛选</span>
                {Object.keys(filters).length > 0 && (
                  <span className="bg-blue-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                    {Object.keys(filters).length}
                  </span>
                )}
              </button>
            )}

            {hasActiveFilters && (
              <button
                onClick={clearFilters}
                className="px-3 py-2 text-sm text-gray-500 hover:text-gray-700 flex items-center gap-1"
              >
                <RotateCcw size={14} />
                <span className="hidden sm:inline">重置</span>
              </button>
            )}
          </div>
        )}

        {/* Advanced filters */}
        {showFilters && filterDefs && filterDefs.length > 0 && (
          <div className="bg-gray-50 border border-gray-200 rounded-lg p-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
              {filterDefs.map((f) => (
                <div key={f.key}>
                  <label className="block text-xs font-medium text-gray-600 mb-1">{f.label}</label>
                  {f.type === 'date-range' ? (
                    <div className="flex gap-2">
                      <input
                        type="date"
                        value={filters[`${f.key}_from`] ?? ''}
                        onChange={(e) =>
                          setFilters((prev) => ({ ...prev, [`${f.key}_from`]: e.target.value }))
                        }
                        className="w-full px-2 py-1.5 border border-gray-300 rounded text-sm outline-none focus:ring-1 focus:ring-blue-500"
                      />
                      <input
                        type="date"
                        value={filters[`${f.key}_to`] ?? ''}
                        onChange={(e) =>
                          setFilters((prev) => ({ ...prev, [`${f.key}_to`]: e.target.value }))
                        }
                        className="w-full px-2 py-1.5 border border-gray-300 rounded text-sm outline-none focus:ring-1 focus:ring-blue-500"
                      />
                    </div>
                  ) : f.type === 'select' ? (
                    <select
                      value={filters[f.key] ?? ''}
                      onChange={(e) => {
                        const v = e.target.value;
                        setFilters((prev) => {
                          const next = { ...prev };
                          if (v) next[f.key] = v;
                          else delete next[f.key];
                          return next;
                        });
                      }}
                      className="w-full px-2 py-1.5 border border-gray-300 rounded text-sm outline-none focus:ring-1 focus:ring-blue-500"
                    >
                      <option value="">全部</option>
                      {f.options?.map((o) => (
                        <option key={o.value} value={o.value}>{o.label}</option>
                      ))}
                    </select>
                  ) : (
                    <input
                      type={f.type === 'date' ? 'date' : 'text'}
                      value={filters[f.key] ?? ''}
                      onChange={(e) => {
                        const v = e.target.value;
                        setFilters((prev) => {
                          const next = { ...prev };
                          if (v) next[f.key] = v;
                          else delete next[f.key];
                          return next;
                        });
                      }}
                      placeholder={f.placeholder}
                      className="w-full px-2 py-1.5 border border-gray-300 rounded text-sm outline-none focus:ring-1 focus:ring-blue-500"
                    />
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Bulk actions bar */}
        {enableBulkActions && selectedIds.size > 0 && (
          <div className="flex items-center gap-2 bg-blue-50 border border-blue-200 rounded-lg px-4 py-2">
            <span className="text-sm text-blue-700 font-medium">
              已选择 {selectedIds.size} 项
            </span>
            <button
              onClick={() => confirmDelete([...selectedIds], `${selectedIds.size} 条`)}
              className="px-3 py-1 text-sm text-red-600 bg-red-50 border border-red-200 rounded hover:bg-red-100 transition-colors flex items-center gap-1"
            >
              <Trash2 size={13} />
              批量删除
            </button>
            {enableExport && (
              <button
                onClick={() => handleExport('csv')}
                className="px-3 py-1 text-sm text-blue-600 bg-blue-50 border border-blue-200 rounded hover:bg-blue-100 transition-colors flex items-center gap-1"
              >
                <Download size={13} />
                导出选中
              </button>
            )}
            <button
              onClick={() => setSelectedIds(new Set())}
              className="ml-auto text-sm text-gray-500 hover:text-gray-700"
            >
              <X size={16} />
            </button>
          </div>
        )}
      </div>

      {/* Table */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-400">加载中...</div>
      ) : rows.length === 0 ? (
        <div className="text-center py-12 text-gray-400">
          {hasActiveFilters ? '未找到匹配的记录' : '暂无数据'}
        </div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-gray-50 border-b border-gray-200">
                  {/* Bulk select checkbox */}
                  {enableBulkActions && (
                    <th className="w-10 px-3 py-3">
                      <button onClick={toggleSelectAll} className="text-gray-400 hover:text-blue-600">
                        {isAllSelected ? (
                          <CheckSquare size={16} className="text-blue-600" />
                        ) : isPartialSelected ? (
                          <CheckSquare size={16} className="text-blue-400" />
                        ) : (
                          <Square size={16} />
                        )}
                      </button>
                    </th>
                  )}
                  {visibleColumns.map((c) => (
                    <th
                      key={c.key}
                      className="text-left px-4 py-3 font-medium text-gray-600 whitespace-nowrap"
                    >
                      {c.label}
                    </th>
                  ))}
                  <th className="text-right px-4 py-3 font-medium text-gray-600 w-24">操作</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row, i) => {
                  const rowId = getRowId(row, i) as string;
                  const isSelected = selectedIds.has(rowId);
                  return (
                    <tr
                      key={rowId}
                      className={`border-b border-gray-100 transition-colors ${
                        isSelected ? 'bg-blue-50' : 'hover:bg-gray-50'
                      }`}
                    >
                      {enableBulkActions && (
                        <td className="px-3 py-2.5">
                          <button
                            onClick={() => toggleSelect(rowId, i)}
                            className="text-gray-400 hover:text-blue-600"
                          >
                            {isSelected ? (
                              <CheckSquare size={16} className="text-blue-600" />
                            ) : (
                              <Square size={16} />
                            )}
                          </button>
                        </td>
                      )}
                      {visibleColumns.map((c) => (
                        <td key={c.key} className="px-4 py-2.5 text-gray-700 whitespace-nowrap">
                          {c.render
                            ? c.render(row[c.key], row)
                            : (c.key === 'status' || c.key.endsWith('_status'))
                              ? renderStatusBadge(String(row[c.key] ?? ''))
                              : formatCellValue(row[c.key])}
                        </td>
                      ))}
                      <td className="px-4 py-2.5 text-right whitespace-nowrap">
                        <button
                          onClick={() => openEdit(row)}
                          className="text-blue-600 hover:text-blue-800 mr-3 text-xs font-medium"
                        >
                          <Pencil size={14} className="inline mr-1" />
                          编辑
                        </button>
                        <button
                          onClick={() => confirmDelete([rowId], String(rowId))}
                          className="text-red-500 hover:text-red-700 text-xs font-medium"
                        >
                          <Trash2 size={14} className="inline mr-1" />
                          删除
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {enablePagination && total && totalPages > 0 && (
            <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 bg-gray-50">
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <span>
                  共 {total} 条，第 {page}/{totalPages} 页
                </span>
                <select
                  value={pageSize}
                  onChange={(e) => {
                    setPageSize(Number(e.target.value));
                    setPage(1);
                  }}
                  className="ml-2 px-2 py-1 border border-gray-300 rounded text-sm outline-none focus:ring-1 focus:ring-blue-500"
                >
                  {[10, 20, 50, 100].map((n) => (
                    <option key={n} value={n}>
                      {n} 条/页
                    </option>
                  ))}
                </select>
              </div>
              <div className="flex items-center gap-1">
                <button
                  disabled={page <= 1}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  className="px-2 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed"
                >
                  <ChevronLeft size={16} />
                </button>
                {renderPageNumbers(page, totalPages, setPage)}
                <button
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  className="px-2 py-1 text-sm border border-gray-300 rounded hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed"
                >
                  <ChevronRight size={16} />
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Add/Edit Modal */}
      <Modal
        open={modalOpen}
        title={editItem ? '编辑' : '添加'}
        onClose={() => setModalOpen(false)}
      >
        <DynamicFormBuilder
          fields={fieldList}
          initialData={editItem ?? undefined}
          onSubmit={handleSave}
          loading={addMut.isPending || updateMut.isPending}
          onCancel={() => setModalOpen(false)}
        />
      </Modal>

      {/* Delete confirmation modal */}
      <Modal
        open={!!deleteConfirm}
        title="确认删除"
        onClose={() => setDeleteConfirm(null)}
        footer={
          <div className="flex justify-end gap-3">
            <button
              onClick={() => setDeleteConfirm(null)}
              className="px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200"
            >
              取消
            </button>
            <button
              onClick={executeDelete}
              className="px-4 py-2 text-sm text-white bg-red-600 rounded-lg hover:bg-red-700"
              disabled={deleteMut.isPending}
            >
              {deleteMut.isPending ? '删除中...' : '确认删除'}
            </button>
          </div>
        }
      >
        <p className="text-gray-700">
          确定要删除 {deleteConfirm?.label ?? ''} 吗？此操作不可撤销。
        </p>
      </Modal>

      {/* Click-away handlers for dropdowns */}
      {showColumnsMenu && <ClickAwayLayer onClick={() => setShowColumnsMenu(false)} />}
      {showExportMenu && <ClickAwayLayer onClick={() => setShowExportMenu(false)} />}
    </div>
  );
}

// ── Helpers (outside component) ──

// ── Status badge color mapping (BFT56-aligned) ──
const STATUS_COLORS: Record<string, string> = {
  // Order statuses
  '待拣货': 'bg-yellow-100 text-yellow-800',
  '待装柜': 'bg-orange-100 text-orange-800',
  '已取消': 'bg-red-100 text-red-800',
  '已完成': 'bg-green-100 text-green-800',
  'pending_picking': 'bg-yellow-100 text-yellow-800',
  'pending_packing': 'bg-yellow-100 text-yellow-800',
  'pending_loading': 'bg-orange-100 text-orange-800',
  'in_transit': 'bg-blue-100 text-blue-800',
  'completed': 'bg-green-100 text-green-800',
  'shipped': 'bg-indigo-100 text-indigo-800',
  'loaded': 'bg-purple-100 text-purple-800',
  'customs_clearance': 'bg-teal-100 text-teal-800',
  // Parcel statuses
  '待打包': 'bg-yellow-100 text-yellow-800',
  '已上架': 'bg-green-100 text-green-800',
  'stored': 'bg-green-100 text-green-800',
  'packed': 'bg-blue-100 text-blue-800',
  'weighed': 'bg-teal-100 text-teal-800',
  // Work order statuses
  '待处理': 'bg-yellow-100 text-yellow-800',
  '处理中': 'bg-blue-100 text-blue-800',
  'pending': 'bg-yellow-100 text-yellow-800',
  'in_progress': 'bg-blue-100 text-blue-800',
  // General
  'active': 'bg-green-100 text-green-800',
  '启用': 'bg-green-100 text-green-800',
  '装货中': 'bg-yellow-100 text-yellow-800',
  '已发运': 'bg-green-100 text-green-800',
  '已结算': 'bg-green-100 text-green-800',
  '待结算': 'bg-yellow-100 text-yellow-800',
  '认证中': 'bg-yellow-100 text-yellow-800',
  '认证成功': 'bg-green-100 text-green-800',
  '认证失败': 'bg-red-100 text-red-800',
};

function renderStatusBadge(status: string) {
  const color = STATUS_COLORS[status] || 'bg-gray-100 text-gray-700';
  return <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${color}`}>{status.replace(/_/g, ' ')}</span>;
}

function formatCellValue(value: unknown): string {
  if (value === null || value === undefined) return '';
  if (typeof value === 'boolean') return value ? '是' : '否';
  if (typeof value === 'object') return JSON.stringify(value);
  return String(value);
}

function ClickAwayLayer({ onClick }: { onClick: () => void }) {
  // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-static-element-interactions
  return <div className="fixed inset-0 z-20" onClick={onClick} />;
}

function renderPageNumbers(
  current: number,
  total: number,
  setPage: (p: number) => void,
): React.ReactNode {
  const pages: (number | '...')[] = [];
  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i);
  } else {
    pages.push(1);
    if (current > 3) pages.push('...');
    for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
      pages.push(i);
    }
    if (current < total - 2) pages.push('...');
    pages.push(total);
  }

  return pages.map((p, i) =>
    p === '...' ? (
      <span key={`ellipsis-${i}`} className="px-2 py-1 text-sm text-gray-400">
        ...
      </span>
    ) : (
      <button
        key={p}
        onClick={() => setPage(p)}
        className={`px-2 py-1 text-sm rounded transition-colors ${
          p === current
            ? 'bg-blue-600 text-white'
            : 'border border-gray-300 hover:bg-gray-100 text-gray-700'
        }`}
      >
        {p}
      </button>
    ),
  );
}
