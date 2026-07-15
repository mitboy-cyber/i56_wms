// Shared types for the I56 WMS frontend
import type { ReactNode } from 'react';

/** A single column definition for GenericListPage */
export interface ColumnDef {
  key: string;
  label: string;
  render?: (value: unknown, row: Record<string, unknown>) => ReactNode;
  sortable?: boolean;
  /** Hide by default (user can reveal via column manager) */
  hidden?: boolean;
}

/** A field definition for forms */
export interface FieldDef {
  key: string;
  label: string;
  type?: 'text' | 'number' | 'select' | 'date' | 'switch' | 'textarea';
  options?: { value: string; label: string }[];
  required?: boolean;
  placeholder?: string;
  /** Validation rule: regex or custom error message */
  pattern?: string;
  min?: number;
  max?: number;
  helperText?: string;
}

/** Advanced filter definition */
export interface FilterDef {
  key: string;
  label: string;
  type: 'text' | 'select' | 'date' | 'date-range';
  options?: { value: string; label: string }[];
  placeholder?: string;
}

/** Pagination state */
export interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

/** API response wrapper */
export interface ApiListResponse<T> {
  data: T[];
  total?: number;
  page?: number;
  page_size?: number;
}

/** Tab state for multi-tab navigation */
export interface TabItem {
  id: string;
  label: string;
  href: string;
  /** Timestamp when tab was opened, for ordering */
  openedAt: number;
}

/** KPI card data */
export interface KpiCardData {
  label: string;
  value: number | string;
  change?: number; // percentage change
  changeLabel?: string;
  icon?: ReactNode;
  color?: 'blue' | 'green' | 'amber' | 'red' | 'purple';
  href?: string;
}
