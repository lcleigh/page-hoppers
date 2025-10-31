'use client';
import { useState } from 'react';

interface BookResult {
  key: string;
  title: string;
  author_name?: string[];
  cover_i?: number;
}

interface BookSearchProps {
  onLogBook: (book: BookResult) => void;
  onAddManualBook: () => void;
}

export default function BookSearch({ onLogBook, onAddManualBook }: BookSearchProps) {
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState<BookResult[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState('');

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

  return (
    <div className="bg-white/80 p-6 rounded-xl shadow border-2 border-sky">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-bubblegum">Search for Books</h2>
        <button
          onClick={onAddManualBook}
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
                onClick={() => onLogBook(book)}
                className="px-3 py-1 bg-sky text-charcoal font-bold rounded-xl hover:bg-bubblegum hover:text-white transition border-2 border-bubblegum"
              >
                Log Book
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
