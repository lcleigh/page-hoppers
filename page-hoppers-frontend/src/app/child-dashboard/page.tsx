'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import BookSearch from '@/components/BookSearch';
import ReadingLog from '@/components/ReadingLog';

interface BookResult {
  key: string;
  title: string;
  author_name?: string[];
  cover_i?: number;
}


export default function ChildDashboard() {
  const [childName, setChildName] = useState('');
  const [loading, setLoading] = useState(true);
  const [showLogModal, setShowLogModal] = useState(false);
  const [selectedBook, setSelectedBook] = useState<BookResult | null>(null);
  const [logType, setLogType] = useState<'completed' | 'started'>('completed');
  const [selectedDate, setSelectedDate] = useState('');
  const [logSuccess, setLogSuccess] = useState('');
  const [logError, setLogError] = useState('');
  const [showAddBookModal, setShowAddBookModal] = useState(false);
  const [manualBookTitle, setManualBookTitle] = useState('');
  const [manualBookAuthor, setManualBookAuthor] = useState('');
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
  }, [router]);

  const handleLogout = () => {
    localStorage.removeItem('childToken');
    localStorage.removeItem('childId');
    localStorage.removeItem('childName');
    router.push('/');
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
        // Refresh the page to update reading logs
        window.location.reload();
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
        // Refresh the page to update reading logs
        window.location.reload();
      } else {
        const errorData = await res.json();
        setLogError(errorData.message || 'Failed to log book');
      }
    } catch (err) {
      setLogError('Failed to log book. Please try again.');
    }
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
          <BookSearch 
            onLogBook={handleLogBook}
            onAddManualBook={handleAddManualBook}
          />

          {/* Reading Logs Section */}
          <ReadingLog />
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