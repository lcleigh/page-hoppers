'use client';
import { useEffect, useState } from 'react';

interface ReadingLogItem {
  id: number;
  title: string;
  author?: string;
  status: 'started' | 'completed';
  date: string;
  open_library_key?: string;
  cover_id?: number;
  created_at: string;
}

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

export default function ReadingLog() {
  const [readingLogs, setReadingLogs] = useState<ReadingLogItem[]>([]);
  const [logsLoading, setLogsLoading] = useState(false);

  const fetchReadingLogs = async () => {
    const childToken = localStorage.getItem('childToken');
    if (!childToken) return;

    setLogsLoading(true);
    try {
      const res = await fetch(`${apiUrl}/reading-logs`, {
        headers: {
          'Authorization': `Bearer ${childToken}`,
          'Content-Type': 'application/json',
        },
      });
      if (res.ok) {
        const logs = await res.json();
        setReadingLogs(logs || []);
      }
    } catch (err) {
      console.error('Failed to fetch reading logs:', err);
    } finally {
      setLogsLoading(false);
    }
  };

  useEffect(() => {
    fetchReadingLogs();
  }, []);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <div className="bg-white/80 p-6 rounded-xl shadow border-2 border-sky">
      <h2 className="text-xl font-semibold mb-4 text-bubblegum">Your Reading Log</h2>
      {logsLoading ? (
        <div className="text-center text-charcoal">Loading reading logs...</div>
      ) : (readingLogs && readingLogs.length > 0) ? (
        <div className="space-y-3 max-h-96 overflow-y-auto">
          {readingLogs.map(log => (
            <div key={log.id} className="flex items-center gap-3 bg-lemon/30 p-3 rounded-xl border-2 border-lemon">
              {log.cover_id ? (
                <img
                  src={`https://covers.openlibrary.org/b/id/${log.cover_id}-S.jpg`}
                  alt={log.title}
                  className="w-12 h-16 rounded shadow"
                />
              ) : (
                <div className="w-12 h-16 bg-coolgray rounded flex items-center justify-center text-white font-bold text-xl">
                  ?
                </div>
              )}
              <div className="flex-1">
                <div className="font-bold text-charcoal">{log.title}</div>
                {log.author && <div className="text-sm text-coolgray">{log.author}</div>}
                <div className="text-xs text-coolgray">{formatDate(log.date)}</div>
              </div>
              <span className={`px-2 py-1 rounded text-xs font-bold ${
                log.status === 'completed' 
                  ? 'bg-leaf text-white' 
                  : 'bg-sky text-charcoal'
              }`}>
                {log.status === 'completed' ? 'Completed' : 'Started'}
              </span>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center text-charcoal">
          <p>No reading logs yet.</p>
          <p className="text-sm">Search for books or add them manually to get started!</p>
        </div>
      )}
    </div>
  );
}
