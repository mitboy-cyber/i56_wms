import { useState, useRef, useEffect } from 'react';
import { Send, Bot, User, Loader2, Settings } from 'lucide-react';
import client from '@/api/client';

interface Message { id: number; role: 'user' | 'assistant'; content: string; time: string; }

export default function AIChatPage() {
  const [messages, setMessages] = useState<Message[]>(() => {
    try { const s = localStorage.getItem('ai-chat-msgs'); return s ? JSON.parse(s) : []; }
    catch { return []; }
  });
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const chatEnd = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => { chatEnd.current?.scrollIntoView({ behavior: 'smooth' }); }, [messages]);
  useEffect(() => { localStorage.setItem('ai-chat-msgs', JSON.stringify(messages)); }, [messages]);

  const handleSend = async () => {
    const text = input.trim();
    if (!text || loading) return;
    setInput('');

    const userMsg: Message = { id: Date.now(), role: 'user', content: text, time: new Date().toLocaleTimeString() };
    setMessages(prev => [...prev, userMsg]);
    setLoading(true);

    try {
      const res = await client.post('/admin/api/system/ai-chat', { role: 'user', content: text, time: new Date().toISOString() });
      const aiMsg: Message = {
        id: Date.now() + 1, role: 'assistant',
        content: `收到你的消息：「${text}」\n\n作为 I56 智能助手，我可以帮你：\n• 查询订单/包裹状态\n• 分析仓库运营数据\n• 解答系统使用问题\n• 生成报表\n\n请告诉我你需要什么帮助？`,
        time: new Date().toLocaleTimeString()
      };
      setMessages(prev => [...prev, aiMsg]);
    } catch {
      const errMsg: Message = { id: Date.now() + 1, role: 'assistant', content: '抱歉，AI 服务暂时不可用，请稍后重试。', time: new Date().toLocaleTimeString() };
      setMessages(prev => [...prev, errMsg]);
    } finally {
      setLoading(false);
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSend(); }
  };

  return (
    <div className="flex flex-col h-[calc(100vh-160px)] max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h2 className="text-xl font-bold text-gray-800">AI 智能助手</h2>
          <p className="text-sm text-gray-500">基于 I56 大模型的仓库运营助手</p>
        </div>
        <a href="/admin/system/ai-settings" className="flex items-center gap-1 text-sm text-blue-600 hover:text-blue-800">
          <Settings size={14} /> AI 设置
        </a>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto bg-white rounded-xl border border-gray-200 p-4 mb-4 space-y-4">
        {messages.length === 0 && (
          <div className="flex flex-col items-center justify-center h-full text-gray-400">
            <Bot size={48} className="mb-3 opacity-50" />
            <p className="text-lg font-medium">I56 AI 助手</p>
            <p className="text-sm mt-1">输入你的问题，获取实时帮助</p>
          </div>
        )}
        {messages.map(msg => (
          <div key={msg.id} className={`flex gap-3 ${msg.role === 'user' ? 'justify-end' : ''}`}>
            {msg.role === 'assistant' && (
              <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center shrink-0 mt-0.5">
                <Bot size={16} className="text-blue-600" />
              </div>
            )}
            <div className={`max-w-[75%] rounded-2xl px-4 py-2.5 ${
              msg.role === 'user'
                ? 'bg-blue-600 text-white rounded-br-md'
                : 'bg-gray-100 text-gray-800 rounded-bl-md'
            }`}>
              <p className="text-sm whitespace-pre-wrap">{msg.content}</p>
              <p className={`text-xs mt-1 ${msg.role === 'user' ? 'text-blue-200' : 'text-gray-400'}`}>{msg.time}</p>
            </div>
            {msg.role === 'user' && (
              <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center shrink-0 mt-0.5">
                <User size={16} className="text-gray-600" />
              </div>
            )}
          </div>
        ))}
        {loading && (
          <div className="flex gap-3">
            <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center shrink-0">
              <Bot size={16} className="text-blue-600" />
            </div>
            <div className="bg-gray-100 rounded-2xl rounded-bl-md px-4 py-2.5">
              <Loader2 size={16} className="animate-spin text-gray-400" />
            </div>
          </div>
        )}
        <div ref={chatEnd} />
      </div>

      {/* Input */}
      <div className="flex gap-2">
        <input
          ref={inputRef}
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="输入你的问题，按 Enter 发送..."
          className="flex-1 px-4 py-3 rounded-xl border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
          disabled={loading}
          autoFocus
        />
        <button
          onClick={handleSend}
          disabled={loading || !input.trim()}
          className="px-5 py-3 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >
          <Send size={18} />
        </button>
      </div>
    </div>
  );
}
