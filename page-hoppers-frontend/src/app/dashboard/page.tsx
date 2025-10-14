// app/dashboard/page.tsx
'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { createChild } from '@/api/children';
import ChildLoginModal from '@/components/modals/ChildLoginModal';
import handleLogin  from '@/app/login/page';



interface ChildProfile {
  id: number;
  username: string;
  avatar?: string;
  name?: string;
  age?: number;
}

export default function ParentDashboard() {
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const [children, setChildren] = useState<ChildProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [newChildFirstName, setNewChildFirstName] = useState("");
  const [newChildLastName, setNewChildLastName] = useState("");
  const [newChildAge, setNewChildAge] = useState("");
  const [newChildPin, setNewChildPin] = useState("");
  const [addError, setAddError] = useState("");
  const [addSuccess, setAddSuccess] = useState("");
  
  // Child login state
  const [showChildLoginModal, setShowChildLoginModal] = useState(false);
  const [selectedChildName, setSelectedChildName] = useState<string | undefined>();
  const [childPinError, setChildPinError] = useState<string | undefined>();
  const [childPin, setChildPin] = useState("");
  const [childLoginError, setChildLoginError] = useState("");
  const [childLoginLoading, setChildLoginLoading] = useState(false);
  
  const router = useRouter();

  useEffect(() => {
    const fetchChildren = async () => {
      const parentToken = localStorage.getItem('parentToken');
      if (!parentToken) {
        router.push('/login');
        return;
      }
      try {
        const res = await fetch(`${API_URL}/api/children`, {
          headers: {
            'Authorization': `Bearer ${parentToken}`,
          },
        });
        if (!res.ok) {
          throw new Error('Failed to fetch children');
        }
        const data = await res.json();
        // Map ID to id for compatibility
        const mappedChildren = data.map((child: any) => ({
          id: child.ID ?? child.id,
          username: child.username,
          name: child.name,
          age: child.age,
          avatar: child.avatar,
        }));
        setChildren(mappedChildren);
      } catch (err) {
        setError('Could not load children');
      } finally {
        setLoading(false);
      }
    };
    fetchChildren();
  }, [router]);

  const handleAddChild = async (e: React.FormEvent) => {
    e.preventDefault();
    setAddError("");
    setAddSuccess("");
    try {
      const parentToken = localStorage.getItem('parentToken');
      if (!parentToken) {
        router.push('/login');
        return;
      }
      const data = await createChild({
        name: `${newChildFirstName} ${newChildLastName}`,
        age: Number(newChildAge),
        pin: newChildPin,
        token: parentToken,
      });
      setAddSuccess('Child added!');
      setNewChildFirstName("");
      setNewChildLastName("");
      setNewChildAge("");
      setNewChildPin("");
      setChildren(prev => [...prev, { id: data.id, username: data.username, name: `${newChildFirstName} ${newChildLastName}`, age: Number(newChildAge) }]);
    } catch (err) {
      setAddError('Could not add child');
    }
  };

  const handleChildSelect = (childName: string) => {
    setSelectedChildName(childName);
    setShowChildLoginModal(true);
    setChildPinError(undefined); // reset error
  };

  const handleChildPinSubmit = async (pin: string) => {

    // Your PIN validation logic here (API call, etc.)
    // On success:
    //   setShowChildLoginModal(false);
    //   ...navigate or update state...
    // On error:
    //   setChildPinError("Invalid PIN");
  };

  const handleChildLogin = (child: ChildProfile) => {
    setSelectedChildName(child.name || child.username);
    setShowChildLoginModal(true);
    setChildPinError(undefined); // reset error
  };

  const handleChildLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedChildName) return;

    setChildLoginLoading(true);
    setChildLoginError("");

    try {
      const res = await fetch('http://localhost:8080/api/auth/child/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          childId: selectedChildName, // Assuming selectedChildName is the ID
          pin: childPin
        }),
      });

      if (!res.ok) {
        throw new Error('Invalid PIN');
      }

      const data = await res.json();
      
      // Store child token and redirect to child dashboard
      localStorage.setItem('childToken', data.token);
      localStorage.setItem('childId', selectedChildName);
      localStorage.setItem('childName', selectedChildName);
      
      // Close modal and redirect
      setShowChildLoginModal(false);
      setSelectedChildName(undefined);
      setChildPin("");
      
      // Redirect to child dashboard (you'll need to create this)
      router.push('/child-dashboard');
      
    } catch (err) {
      setChildLoginError('Invalid PIN. Please try again.');
    } finally {
      setChildLoginLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('parentToken');
    router.push('/');
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-powder relative">
      {/* Logout button in top-right corner */}
      <button
        onClick={handleLogout}
        className="absolute top-4 right-4 px-4 py-2 bg-coral text-white font-bold rounded-xl hover:bg-bubblegum hover:text-white transition border-2 border-coral"
      >
        Logout
      </button>
      
      <div className="w-full max-w-2xl bg-lavender rounded-2xl shadow-lg p-8 flex flex-col items-center">
        <h1 className="text-4xl font-extrabold mb-4 text-bubblegum drop-shadow">Parent Dashboard</h1>
        <p className="mb-6 text-lg text-charcoal">Welcome! Here are your children:</p>
        
        <form onSubmit={handleAddChild} className="mb-8 w-full max-w-md bg-white/80 p-6 rounded-xl shadow flex flex-col gap-3 border-2 border-lavender">
          <h2 className="text-xl font-semibold mb-2 text-bubblegum">Add a Child</h2>
          <input
            type="text"
            placeholder="First Name"
            className="border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
            value={newChildFirstName}
            onChange={e => setNewChildFirstName(e.target.value)}
            required
          />
          <input
            type="text"
            placeholder="Last Name"
            className="border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
            value={newChildLastName}
            onChange={e => setNewChildLastName(e.target.value)}
            required
          />
          <input
            type="number"
            placeholder="Age"
            className="border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
            value={newChildAge}
            onChange={e => setNewChildAge(e.target.value)}
            min={1}
            required
          />
          <input
            type="password"
            placeholder="4-digit PIN"
            className="border-2 border-sky rounded px-3 py-2 focus:outline-none focus:border-bubblegum"
            value={newChildPin}
            onChange={e => setNewChildPin(e.target.value)}
            pattern="[0-9]{4}"
            maxLength={4}
            required
          />
          <button
            type="submit"
            className="bg-bubblegum text-white font-bold rounded-xl px-4 py-2 hover:bg-lemon hover:text-charcoal transition border-2 border-sky"
          >
            Add Child
          </button>
          {addError && <div className="text-coral text-sm">{addError}</div>}
          {addSuccess && <div className="text-leaf text-sm">{addSuccess}</div>}
        </form>
        
        {loading ? (
          <div>Loading...</div>
        ) : error ? (
          <div className="text-coral">{error}</div>
        ) : children.length === 0 ? (
          <div className="text-coolgray">No children found. Add a child above!</div>
        ) : (
          <ul className="w-full max-w-md space-y-4">
            {children.map(child => (
              <li key={child.id} className="p-4 bg-white rounded-xl shadow flex flex-col sm:flex-row sm:items-center sm:justify-between border-2 border-sky">
                <div className="flex items-center">
                  {child.avatar ? (
                    <img src={child.avatar} alt={child.username} className="w-10 h-10 rounded-full mr-4" />
                  ) : (
                    <div className="w-10 h-10 bg-lemon rounded-full flex items-center justify-center mr-4 border-2 border-bubblegum">
                      <span className="text-xl font-bold text-bubblegum">{child.name ? child.name[0].toUpperCase() : child.username[0].toUpperCase()}</span>
                    </div>
                  )}
                  <div className="flex flex-col">
                    <span className="font-medium text-charcoal text-lg">{child.name || child.username}</span>
                    <span className="text-coolgray text-sm">Age: {child.age || "-"}</span>
                  </div>
                </div>
                <button
                  onClick={() => handleChildLogin(child)}
                  className="mt-2 sm:mt-0 px-4 py-2 bg-sky text-charcoal font-bold rounded-xl hover:bg-lemon hover:text-bubblegum transition border-2 border-bubblegum"
                >
                  Login
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Child Login Modal */}
      <ChildLoginModal
        isOpen={showChildLoginModal}
        onClose={() => setShowChildLoginModal(false)}
        onSubmit={handleLogin}
        childName={selectedChildName}
        error={childPinError}
      />
    </div>
  );
} 