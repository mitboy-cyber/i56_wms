import { useEffect, useState } from 'react';
import { Palette, Sun, Moon, Waves, Building2 } from 'lucide-react';

const themes = [
  { id: 'default', label: '翡翠', icon: Waves, color: '#059669' },
  { id: 'dark', label: '暗夜', icon: Moon, color: '#374151' },
  { id: 'blue', label: '深海', icon: Waves, color: '#2563eb' },
  { id: 'enterprise', label: '企业', icon: Building2, color: '#d97706' },
];

export default function ThemeSwitcher() {
  const [theme, setTheme] = useState(() => localStorage.getItem('i56-theme') || 'default');
  const [open, setOpen] = useState(false);

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('i56-theme', theme);
  }, [theme]);

  const current = themes.find(t => t.id === theme) || themes[0];

  return (
    <div className="relative">
      <button onClick={() => setOpen(!open)} className="flex items-center gap-2 px-3 py-2 text-sm rounded-md w-full"
        style={{ color: 'var(--color-neutral)' }}
        title="切换主题"
      >
        <Palette size={16} /> 主题
      </button>
      {open && (
        <>
          <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
          <div className="absolute bottom-full left-0 mb-1 w-40 bg-white border rounded-md shadow-sm z-20"
            style={{ background: 'var(--sidebar-bg)', borderColor: 'var(--border)' }}>
            {themes.map(t => (
              <button key={t.id} onClick={() => { setTheme(t.id); setOpen(false); }}
                className={`flex items-center gap-2 w-full px-3 py-2 text-sm hover:opacity-80 ${t.id === theme ? 'font-medium' : ''}`}
                style={{ color: t.id === theme ? 'var(--color-accent)' : 'var(--color-muted)' }}
              >
                <t.icon size={14} color={t.color} />
                {t.label}
                {t.id === theme && <span className="ml-auto text-xs">✓</span>}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
