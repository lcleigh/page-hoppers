'use client';
import { useEffect, useState } from 'react';

interface ReadingSummary {
  currentBook?: {
    title: string;
    author?: string;
    cover_id?: number;
  };
  lastCompletedBook?: {
    title: string;
    author?: string;
    cover_id?: number;
  };
  totalBooksReadThisMonth: number;
  totalBooksReadThisYear: number;
  totalUncompletedBooks: number;
  totalComplatedBooks: number;
}

const apiUrl = process.env.NEXT_PUBLIC_API_URL;

export default function ReadingSummaryCard() {
    const [summary, setSummary] = useState<ReadingSummary | null>(null);
    const [loading, setLoading] = useState(false);
  
    const fetchSummary = async () => {
      const childToken = localStorage.getItem('childToken');
      if (!childToken) return;
      
      const childId = localStorage.getItem('childId');
      setLoading(true);
      try {
        const res = await fetch(`${apiUrl}/children/${childId}/summary`, {
          headers: {
            'Authorization': `Bearer ${childToken}`,
            'Content-Type': 'application/json',
          },
        });
        if (res.ok) {
          const data = await res.json();
          console.log("API summary response:", data); // ✅ log the full response
          console.log("Current book:", data.currentBook); // ✅ log currentBook specifically
          setSummary(data);
        }
      } catch (err) {
        console.error('Failed to fetch reading summary:', err);
      } finally {
        setLoading(false);
      }
    };

    useEffect(() => {
    fetchSummary();
  }, []);

  const renderBook = (book?: { title: string; author?: string; cover_id?: number }) => {
    if (!book) return <span className="text-sm text-coolgray">None</span>;
    return (
      <div className="flex items-center gap-2">
        {book.cover_id ? (
          <img
            src={`https://covers.openlibrary.org/b/id/${book.cover_id}-S.jpg`}
            alt={book.title}
            className="w-10 h-14 rounded shadow"
          />
        ) : (
          <div className="w-10 h-14 bg-coolgray rounded flex items-center justify-center text-white font-bold text-lg">
            ?
          </div>
        )}
        <div>
          <div className="font-bold text-charcoal">{book.title}</div>
          {book.author && <div className="text-sm text-coolgray">{book.author}</div>}
        </div>
      </div>
    );
  };

  return (
    <div className="bg-white/80 p-6 rounded-xl shadow border-2 border-sky">
      <h2 className="text-xl font-semibold mb-4 text-bubblegum">Reading Summary</h2>
      {loading ? (
        <div className="text-center text-charcoal">Loading summary...</div>
      ) : summary ? (
        <div className="space-y-4">
          <div>
            <h3 className="font-semibold text-charcoal">Currently Reading:</h3>
            {renderBook(summary.currentBook)}
          </div>
          <div>
            <h3 className="font-semibold text-charcoal">Last Book Read:</h3>
            {renderBook(summary.lastCompletedBook)}
          </div>
          <div className="grid grid-cols-2 gap-4 mt-2">
            <div className="bg-lemon/30 p-3 rounded-xl border-2 border-lemon text-center">
              <div className="font-bold text-charcoal text-lg">{summary.totalBooksReadThisMonth}</div>
              <div className="text-sm text-coolgray">Books Read This Month</div>
            </div>
            <div className="bg-lemon/30 p-3 rounded-xl border-2 border-lemon text-center">
              <div className="font-bold text-charcoal text-lg">{summary.totalBooksReadThisYear}</div>
              <div className="text-sm text-coolgray">Books Read This Year</div>
            </div>
            <div className="bg-lemon/30 p-3 rounded-xl border-2 border-lemon text-center col-span-2">
              <div className="font-bold text-charcoal text-lg">{summary.totalUncompletedBooks}</div>
              <div className="text-sm text-coolgray">Total Books Unfinished</div>
            </div>
            <div className="bg-lemon/30 p-3 rounded-xl border-2 border-lemon text-center col-span-2">
              <div className="font-bold text-charcoal text-lg">{summary.totalComplatedBooks}</div>
              <div className="text-sm text-coolgray">Total Books Read</div>
            </div>
          </div>
        </div>
      ) : (
        <div className="text-center text-charcoal">
          <p>No reading data available.</p>
        </div>
      )}
    </div>
  );
}