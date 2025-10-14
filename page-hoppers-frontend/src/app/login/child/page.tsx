import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';

interface ChildProfile {
  id: number;
  username: string;
  avatar?: string;
}

export default function ChildLogin() {
  const [children, setChildren] = useState<ChildProfile[]>([]);
  const [selectedChild, setSelectedChild] = useState<number | null>(null);
  const [pin, setPin] = useState('');
  const [showPin, setShowPin] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  useEffect(() => {
    // Fetch children associated with the logged-in parent
    const fetchChildren = async () => {
      try {
        const parentToken = localStorage.getItem('parentToken');
        if (!parentToken) {
          router.push('/login'); // Redirect to parent login if not authenticated
          return;
        }

        const response = await fetch('/api/children', {
          headers: {
            'Authorization': `Bearer ${parentToken}`,
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch children');
        }

        const data = await response.json();
        setChildren(data);
      } catch (err) {
        setError('Could not load children profiles');
      }
    };

    fetchChildren();
  }, [router]);

  const handleChildLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedChild) return;

    try {
      const response = await fetch('/api/auth/child/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          childId: selectedChild,
          pin: pin,
        }),
      });

      if (!response.ok) {
        throw new Error('Invalid PIN');
      }

      const data = await response.json();
      localStorage.setItem('childToken', data.token);
      router.push('/child-dashboard');
    } catch (err) {
      setError('Invalid PIN');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Who's Reading Today?
          </h2>
        </div>

        {error && (
          <div className="text-red-500 text-sm text-center">{error}</div>
        )}

        <div className="grid grid-cols-2 gap-4 mt-6">
          {children.map((child) => (
            <button
              key={child.id}
              onClick={() => setSelectedChild(child.id)}
              className={`p-4 rounded-lg border-2 ${
                selectedChild === child.id
                  ? 'border-indigo-500 bg-indigo-50'
                  : 'border-gray-200'
              }`}
            >
              {child.avatar ? (
                <Image
                  src={child.avatar}
                  alt={child.username}
                  width={64}
                  height={64}
                  className="mx-auto rounded-full"
                />
              ) : (
                <div className="w-16 h-16 bg-indigo-100 rounded-full flex items-center justify-center mx-auto">
                  <span className="text-2xl">{child.username[0].toUpperCase()}</span>
                </div>
              )}
              <p className="mt-2 text-center font-medium">{child.username}</p>
            </button>
          ))}
        </div>

        {selectedChild && (
          <form onSubmit={handleChildLogin} className="mt-6">
            <div>
              <label htmlFor="pin" className="sr-only">
                Enter your PIN
              </label>
              <div style={{ position: "relative" }}>
                <input
                  id="pin"
                  name="pin"
                  type={showPin ? "text" : "password"}
                  value={pin}
                  onChange={e => setPin(e.target.value)}
                  placeholder="PIN"
                  maxLength={4}
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-center text-2xl"
                />
                <button
                  type="button"
                  onClick={() => setShowPin((prev) => !prev)}
                  style={{
                    position: "absolute",
                    right: 8,
                    top: "50%",
                    transform: "translateY(-50%)",
                    background: "none",
                    border: "none",
                    cursor: "pointer"
                  }}
                  aria-label={showPin ? "Hide PIN" : "Show PIN"}
                >
                  {showPin ? (
                    // Eye-off icon (SVG)
                    <svg width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                      <path d="M17.94 17.94A10.94 10.94 0 0 1 12 19c-5 0-9.27-3.11-11-7.5a11.05 11.05 0 0 1 5.17-5.61M1 1l22 22" />
                    </svg>
                  ) : (
                    // Eye icon (SVG)
                    <svg width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                      <circle cx="12" cy="12" r="3" />
                      <path d="M2.05 12C3.81 7.61 7.88 4.5 12 4.5s8.19 3.11 9.95 7.5c-1.76 4.39-5.83 7.5-9.95 7.5S3.81 16.39 2.05 12z" />
                    </svg>
                  )}
                </button>
              </div>
            </div>
            <button
              type="submit"
              className="mt-4 w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Let's Read!
            </button>
          </form>
        )}
      </div>
    </div>
  );
} 