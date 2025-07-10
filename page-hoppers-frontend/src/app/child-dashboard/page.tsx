'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

interface BookResult {
  key: string;
  title: string;
  author_name?: string[];
  cover_i?: number;
}

interface ReadingLog {
  id: number;
  title: string;
  author?: string;
  status: 'started' | 'completed';
  date: string;
  open_library_key?: string;
  cover_id?: number;
  created_at: string;
}

export default function ChildDashboard() {
  const [childName, setChildName] = useState('');
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState<BookResult[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState('');
  const [showLogModal, setShowLogModal] = useState(false);
  const [selectedBook, setSelectedBook] = useState<BookResult | null>(null);
  const [logType, setLogType] = useState<'completed' | 'started'>('completed');
  const [selectedDate, setSelectedDate] = useState('');
  const [logSuccess, setLogSuccess] = useState('');
  const [logError, setLogError] = useState('');
  const [showAddBookModal, setShowAddBookModal] = useState(false);
  const [manualBookTitle, setManualBookTitle] = useState('');
  const [manualBookAuthor, setManualBookAuthor] = useState('');
  const [readingLogs, setReadingLogs] = useState<ReadingLog[]>([]);
  const [logsLoading, setLogsLoading] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const childToken = localStorage.getItem('childToken');
    const name = localStorage.getItem('childName');
    if (!childToken) {
      router.push('/');
      return;
    }
    setChildName(name || 'Child');
    setLoading(false);
    // Set default date to today
    setSelectedDate(new Date().toISOString().split('T')[0]);
    fetchReadingLogs();
  }, [router]);

  const fetchReadingLogs = async () => {
    const childToken = localStorage.getItem('childToken');
    if (!childToken) return;

    setLogsLoading(true);
    try {
      const res = await fetch('http://localhost:8080/api/reading-logs', {
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

  const handleLogout = () => {
    localStorage.removeItem('childToken');
    localStorage.removeItem('childId');
    localStorage.removeItem('childName');
    router.push('/');
  };

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    setSearchLoading(true);
    setSearchError('');
    setSearchResults([]);
    try {
      const res = await fetch(`https://openlibrary.org/search.json?q=${encodeURIComponent(searchTerm)}`);
      if (!res.ok) throw new Error('Failed to fetch books');
      const data = await res.json();
      setSearchResults(data.docs.slice(0, 10)); // Show top 10 results
    } catch (err) {
      setSearchError('Could not fetch books.');
    } finally {
      setSearchLoading(false);
    }
  };

  const handleLogBook = (book: BookResult) => {
    setSelectedBook(book);
    setLogType('completed');
    setSelectedDate(new Date().toISOString().split('T')[0]);
    setLogSuccess('');
    setLogError('');
    setShowLogModal(true);
  };

  const handleLogSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedBook) return;
    
    const childToken = localStorage.getItem('childToken');
    if (!childToken) {
      setLogError('Not authenticated');
      return;
    }

    try {
      const res = await fetch('http://localhost:8080/api/reading-logs', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${childToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: selectedBook.title,
          author: selectedBook.author_name?.join(', '),
          status: logType,
          date: selectedDate,
          open_library_key: selectedBook.key,
          cover_id: selectedBook.cover_i,
        }),
      });

      if (res.ok) {
        const action = logType === 'completed' ? 'completed reading' : 'started reading';
        setLogSuccess(`"${selectedBook.title}" logged as ${action} on ${selectedDate}!`);
        setShowLogModal(false);
        setSelectedBook(null);
        fetchReadingLogs(); // Refresh the reading logs
      } else {
        const errorData = await res.json();
        setLogError(errorData.message || 'Failed to log book');
      }
    } catch (err) {
      setLogError('Failed to log book. Please try again.');
    }
  };

  const handleAddManualBook = () => {
    setManualBookTitle('');
    setManualBookAuthor('');
    setLogType('completed');
    setSelectedDate(new Date().toISOString().split('T')[0]);
    setLogSuccess('');
    setLogError('');
    setShowAddBookModal(true);
  };

  const handleManualBookSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!manualBookTitle.trim()) return;
    
    const childToken = localStorage.getItem('childToken');
    if (!childToken) {
      setLogError('Not authenticated');
      return;
    }

    try {
      const res = await fetch('http://localhost:8080/api/reading-logs', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${childToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: manualBookTitle,
          author: manualBookAuthor,
          status: logType,
          date: selectedDate,
        }),
      });

      if (res.ok) {
        const action = logType === 'completed' ? 'completed reading' : 'started reading';
        setLogSuccess(`"${manualBookTitle}" logged as ${action} on ${selectedDate}!`);
        setShowAddBookModal(false);
        setManualBookTitle('');
        setManualBookAuthor('');
        fetchReadingLogs(); // Refresh the reading logs
      } else {
        const errorData = await res.json();
        setLogError(errorData.message || 'Failed to log book');
      }
    } catch (err) {
      setLogError('Failed to log book. Please try again.');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-powder">
        <div className="text-2xl text-bubblegum">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-powder relative">
      {/* Logout button in top-right corner */}
      <button
        onClick={handleLogout}
        className="absolute top-4 right-4 px-4 py-2 bg-coral text-white font-bold rounded-xl hover:bg-bubblegum hover:text-white transition border-2 border-coral"
      >
        Logout
      </button>
      
      <div className="w-full max-w-4xl bg-lavender rounded-2xl shadow-lg p-8 flex flex-col items-center">
        <h1 className="text-4xl font-extrabold mb-4 text-bubblegum drop-shadow">Welcome, {childName}!</h1>
        <p className="mb-6 text-lg text-charcoal">This is your reading dashboard.</p>
        
        <div className="w-full grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Book Search Section */}
          <div className="bg-white/80 p-6 rounded-xl shadow border-2 border-sky">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold text-bubblegum">Search for Books</h2>
              <button
                onClick={handleAddManualBook}
                className="px-3 py-1 bg-lemon text-charcoal font-bold rounded-xl hover:bg-bubblegum hover:text-white transition border-2 border-bubblegum text-sm"
              >
                Add Book Manually
              </button>
            </div>
            <form onSubmit={handleSearch} className="flex gap-2 mb-4">
              <input
                type="text"
                value={searchTerm}
                onChange={e => setSearchTerm(e.target.value)}
                placeholder="Title, author, or keyword"
                className="flex-1 border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
                required
              />
              <button
                type="submit"
                className="bg-bubblegum text-white font-bold rounded-xl px-4 py-2 hover:bg-lemon hover:text-charcoal transition border-2 border-sky"
                disabled={searchLoading}
              >
                {searchLoading ? 'Searching...' : 'Search'}
              </button>
            </form>
            {searchError && <div className="text-coral text-sm mb-2">{searchError}</div>}
            {searchResults.length > 0 && (
              <ul className="space-y-4 max-h-96 overflow-y-auto">
                {searchResults.map(book => (
                  <li key={book.key} className="flex items-center gap-4 bg-lemon/30 p-3 rounded-xl border-2 border-lemon">
                    {book.cover_i ? (
                      <img
                        src={`https://covers.openlibrary.org/b/id/${book.cover_i}-S.jpg`}
                        alt={book.title}
                        className="w-12 h-16 rounded shadow"
                      />
                    ) : (
                      <div className="w-12 h-16 bg-coolgray rounded flex items-center justify-center text-white font-bold text-xl">
                        ?
                      </div>
                    )}
                    <div className="flex-1">
                      <div className="font-bold text-charcoal">{book.title}</div>
                      <div className="text-sm text-coolgray">{book.author_name?.join(', ')}</div>
                    </div>
                    <button
                      onClick={() => handleLogBook(book)}
                      className="px-3 py-1 bg-sky text-charcoal font-bold rounded-xl hover:bg-bubblegum hover:text-white transition border-2 border-bubblegum"
                    >
                      Log Book
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          {/* Reading Logs Section */}
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
        </div>

        {/* Log Book Modal */}
        {showLogModal && selectedBook && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-2xl p-6 max-w-md w-full mx-4 border-2 border-bubblegum">
              <h2 className="text-2xl font-bold text-bubblegum mb-4">
                Log "{selectedBook.title}"
              </h2>
              <form onSubmit={handleLogSubmit} className="space-y-4">
                <div>
                  <label className="block text-charcoal font-medium mb-2">What would you like to log?</label>
                  <div className="flex flex-col gap-2">
                    <label className="flex items-center gap-2">
                      <input
                        type="radio"
                        name="logType"
                        value="completed"
                        checked={logType === 'completed'}
                        onChange={() => setLogType('completed')}
                        required
                      />
                      Whole book read
                    </label>
                    <label className="flex items-center gap-2">
                      <input
                        type="radio"
                        name="logType"
                        value="started"
                        checked={logType === 'started'}
                        onChange={() => setLogType('started')}
                      />
                      Book started
                    </label>
                  </div>
                </div>
                <div>
                  <label className="block text-charcoal font-medium mb-2">Date:</label>
                  <input
                    type="date"
                    value={selectedDate}
                    onChange={e => setSelectedDate(e.target.value)}
                    className="w-full border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
                    required
                  />
                </div>
                <div className="flex gap-3">
                  <button
                    type="button"
                    onClick={() => setShowLogModal(false)}
                    className="flex-1 px-4 py-2 bg-coolgray text-white font-bold rounded-xl hover:bg-charcoal transition"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="flex-1 px-4 py-2 bg-bubblegum text-white font-bold rounded-xl hover:bg-lemon hover:text-charcoal transition"
                  >
                    Log Book
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Add Manual Book Modal */}
        {showAddBookModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-2xl p-6 max-w-md w-full mx-4 border-2 border-bubblegum">
              <h2 className="text-2xl font-bold text-bubblegum mb-4">
                Add Book Manually
              </h2>
              <form onSubmit={handleManualBookSubmit} className="space-y-4">
                <div>
                  <label className="block text-charcoal font-medium mb-2">Book Title:</label>
                  <input
                    type="text"
                    value={manualBookTitle}
                    onChange={e => setManualBookTitle(e.target.value)}
                    placeholder="Enter book title"
                    className="w-full border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
                    required
                  />
                </div>
                <div>
                  <label className="block text-charcoal font-medium mb-2">Author (optional):</label>
                  <input
                    type="text"
                    value={manualBookAuthor}
                    onChange={e => setManualBookAuthor(e.target.value)}
                    placeholder="Enter author name"
                    className="w-full border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
                  />
                </div>
                <div>
                  <label className="block text-charcoal font-medium mb-2">What would you like to log?</label>
                  <div className="flex flex-col gap-2">
                    <label className="flex items-center gap-2">
                      <input
                        type="radio"
                        name="manualLogType"
                        value="completed"
                        checked={logType === 'completed'}
                        onChange={() => setLogType('completed')}
                        required
                      />
                      Whole book read
                    </label>
                    <label className="flex items-center gap-2">
                      <input
                        type="radio"
                        name="manualLogType"
                        value="started"
                        checked={logType === 'started'}
                        onChange={() => setLogType('started')}
                      />
                      Book started
                    </label>
                  </div>
                </div>
                <div>
                  <label className="block text-charcoal font-medium mb-2">Date:</label>
                  <input
                    type="date"
                    value={selectedDate}
                    onChange={e => setSelectedDate(e.target.value)}
                    className="w-full border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
                    required
                  />
                </div>
                <div className="flex gap-3">
                  <button
                    type="button"
                    onClick={() => setShowAddBookModal(false)}
                    className="flex-1 px-4 py-2 bg-coolgray text-white font-bold rounded-xl hover:bg-charcoal transition"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="flex-1 px-4 py-2 bg-bubblegum text-white font-bold rounded-xl hover:bg-lemon hover:text-charcoal transition"
                  >
                    Add Book
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {logSuccess && <div className="text-leaf text-lg mt-4">{logSuccess}</div>}
        {logError && <div className="text-coral text-lg mt-4">{logError}</div>}
        
        <button
          onClick={() => router.push('/dashboard')}
          className="mt-6 px-6 py-2 bg-sky text-charcoal font-bold rounded-xl hover:bg-lemon hover:text-bubblegum transition border-2 border-bubblegum"
        >
          Back to Parent Dashboard
        </button>
      </div>
    </div>
  );
} 