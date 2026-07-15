import { useState, useCallback } from 'react';
import type { FieldDef } from '@/types';
import { cn } from '@/lib/utils';

interface DynamicFormBuilderProps {
  fields: FieldDef[];
  initialData?: Record<string, unknown>;
  onSubmit: (data: Record<string, unknown>) => void | Promise<void>;
  submitLabel?: string;
  cancelLabel?: string;
  onCancel?: () => void;
  loading?: boolean;
  /** Show field-level validation errors */
  validateOnBlur?: boolean;
}

interface FieldError {
  [key: string]: string;
}

export function DynamicFormBuilder({
  fields,
  initialData,
  onSubmit,
  submitLabel = '保存',
  cancelLabel = '取消',
  onCancel,
  loading = false,
  validateOnBlur = true,
}: DynamicFormBuilderProps) {
  const [form, setForm] = useState<Record<string, unknown>>(() => {
    const init: Record<string, unknown> = {};
    fields.forEach((f) => {
      init[f.key] = initialData?.[f.key] ?? getDefaultValue(f);
    });
    return init;
  });
  const [errors, setErrors] = useState<FieldError>({});
  const [touched, setTouched] = useState<Set<string>>(new Set());

  const validateField = useCallback(
    (field: FieldDef, value: unknown): string => {
      if (field.required && (value === '' || value === undefined || value === null)) {
        return `${field.label} 为必填项`;
      }
      if (field.type === 'number') {
        const num = Number(value);
        if (field.required && isNaN(num)) return `${field.label} 必须为数字`;
        if (!isNaN(num)) {
          if (field.min !== undefined && num < field.min) return `${field.label} 最小值为 ${field.min}`;
          if (field.max !== undefined && num > field.max) return `${field.label} 最大值为 ${field.max}`;
        }
      }
      if (field.pattern && typeof value === 'string' && value) {
        try {
          const re = new RegExp(field.pattern);
          if (!re.test(value)) return field.helperText ?? `${field.label} 格式不正确`;
        } catch {
          // invalid regex, skip
        }
      }
      return '';
    },
    [],
  );

  const validateAll = useCallback((): boolean => {
    const newErrors: FieldError = {};
    let valid = true;
    fields.forEach((f) => {
      const err = validateField(f, form[f.key]);
      if (err) {
        newErrors[f.key] = err;
        valid = false;
      }
    });
    setErrors(newErrors);
    const allKeys = new Set(fields.map((f) => f.key));
    setTouched(allKeys);
    return valid;
  }, [fields, form, validateField]);

  const handleChange = (key: string, value: unknown) => {
    setForm((prev) => ({ ...prev, [key]: value }));
    // Clear error on change
    if (errors[key]) {
      setErrors((prev) => {
        const next = { ...prev };
        delete next[key];
        return next;
      });
    }
  };

  const handleBlur = (field: FieldDef) => {
    if (!validateOnBlur) return;
    setTouched((prev) => new Set(prev).add(field.key));
    const err = validateField(field, form[field.key]);
    if (err) {
      setErrors((prev) => ({ ...prev, [field.key]: err }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateAll()) return;
    await onSubmit(form);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4" noValidate>
      {fields.map((field) => {
        const hasError = touched.has(field.key) && !!errors[field.key];
        return (
          <div key={field.key} className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">
              {field.label}
              {field.required && <span className="text-red-500 ml-0.5">*</span>}
            </label>

            {renderField(field, form[field.key], (v) => handleChange(field.key, v), () => handleBlur(field))}

            {field.helperText && !hasError && (
              <p className="text-xs text-gray-400">{field.helperText}</p>
            )}
            {hasError && (
              <p className="text-xs text-red-500">{errors[field.key]}</p>
            )}
          </div>
        );
      })}

      <div className="flex justify-end gap-3 pt-2">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
          >
            {cancelLabel}
          </button>
        )}
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 text-sm text-white bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
        >
          {loading ? '提交中...' : submitLabel}
        </button>
      </div>
    </form>
  );
}

function getDefaultValue(field: FieldDef): unknown {
  switch (field.type) {
    case 'number':
      return '';
    case 'switch':
      return false;
    case 'select':
      return field.options?.[0]?.value ?? '';
    default:
      return '';
  }
}

function renderField(
  field: FieldDef,
  value: unknown,
  onChange: (v: unknown) => void,
  onBlur: () => void,
): React.ReactNode {
  const baseClass =
    'w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-colors';

  switch (field.type) {
    case 'textarea':
      return (
        <textarea
          value={String(value ?? '')}
          onChange={(e) => onChange(e.target.value)}
          onBlur={onBlur}
          placeholder={field.placeholder}
          rows={4}
          className={baseClass + ' resize-y min-h-[80px]'}
        />
      );

    case 'number':
      return (
        <input
          type="number"
          value={String(value ?? '')}
          onChange={(e) => onChange(e.target.value === '' ? '' : Number(e.target.value))}
          onBlur={onBlur}
          placeholder={field.placeholder}
          min={field.min}
          max={field.max}
          className={baseClass}
        />
      );

    case 'select':
      return (
        <select
          value={String(value ?? '')}
          onChange={(e) => onChange(e.target.value)}
          onBlur={onBlur}
          className={baseClass}
        >
          <option value="">{field.placeholder ?? `请选择${field.label}`}</option>
          {field.options?.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
      );

    case 'date':
      return (
        <input
          type="date"
          value={String(value ?? '')}
          onChange={(e) => onChange(e.target.value)}
          onBlur={onBlur}
          className={baseClass}
        />
      );

    case 'switch':
      return (
        <label className="relative inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            checked={!!value}
            onChange={(e) => onChange(e.target.checked)}
            onBlur={onBlur}
            className="sr-only peer"
          />
          <div className="w-10 h-5 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-blue-600" />
        </label>
      );

    default: // text
      return (
        <input
          type="text"
          value={String(value ?? '')}
          onChange={(e) => onChange(e.target.value)}
          onBlur={onBlur}
          placeholder={field.placeholder}
          className={baseClass}
        />
      );
  }
}
