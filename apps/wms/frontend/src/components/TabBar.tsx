import { useTabStore } from '@/stores/tabs';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';

export function TabBar() {
  const { tabs, activeTabId, setActiveTab, closeTab } = useTabStore();

  if (tabs.length === 0) return null;

  return (
    <div className="flex items-center gap-0 bg-gray-100 border-b border-gray-200 overflow-x-auto shrink-0">
      {tabs.map((tab) => {
        const isActive = tab.id === activeTabId;
        return (
          <div
            key={tab.id}
            className={cn(
              'group flex items-center gap-1.5 px-3 py-2 text-sm cursor-pointer border-r border-gray-200',
              'min-w-0 max-w-[200px] transition-colors',
              isActive
                ? 'bg-white text-blue-700 font-medium border-t-2 border-t-blue-600'
                : 'bg-gray-50 text-gray-600 hover:bg-gray-100 border-t-2 border-t-transparent',
            )}
            onClick={() => setActiveTab(tab.id)}
          >
            <span className="truncate flex-1">{tab.label}</span>
            <button
              className={cn(
                'shrink-0 rounded p-0.5 hover:bg-gray-300 transition-colors',
                isActive ? 'opacity-100' : 'opacity-0 group-hover:opacity-100',
              )}
              onClick={(e) => {
                e.stopPropagation();
                closeTab(tab.id);
              }}
              title="关闭标签页"
            >
              <X size={12} />
            </button>
          </div>
        );
      })}
    </div>
  );
}
