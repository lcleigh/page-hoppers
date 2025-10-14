import { useState } from "react";

interface ChildLoginModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (pin: string) => void;
  childName?: string;
  error?: string;
}

export default function ChildLoginModal({
  isOpen,
  onClose,
  onSubmit,
  childName,
  error,
}: ChildLoginModalProps) {
  const [pin, setPin] = useState("");
  const [showPin, setShowPin] = useState(false);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
      <div className="bg-white rounded-lg shadow-lg p-8 w-full max-w-sm relative">
        <button
          className="absolute top-2 right-2 text-gray-400 hover:text-gray-600"
          onClick={onClose}
          aria-label="Close"
        >
          ×
        </button>
        <h2 className="text-xl font-bold mb-4 text-center">
          {childName ? `Enter PIN for ${childName}` : "Enter PIN"}
        </h2>
        {error && <div className="text-red-500 text-sm mb-2 text-center">{error}</div>}
        <form
          onSubmit={e => {
            e.preventDefault();
            onSubmit(pin);
            setPin("");
          }}
        >
          <div style={{ position: "relative" }}>
            <input
              type={showPin ? "text" : "password"}
              value={pin}
              onChange={e => setPin(e.target.value)}
              placeholder="PIN"
              maxLength={4}
              className="appearance-none rounded-md block w-full px-3 py-2 pr-10 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-center text-2xl"
              autoFocus
            />
            <button
              type="button"
              onClick={() => setShowPin(prev => !prev)}
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
          <button
            type="submit"
            className="mt-4 w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
          >
            Submit
          </button>
        </form>
      </div>
    </div>
  );
}