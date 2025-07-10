// app/page.tsx
'use client';
import Link from 'next/link';

const Home = () => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-powder">
      <div className="w-full max-w-xl bg-lavender rounded-2xl shadow-lg p-10 flex flex-col items-center">
        <h1 className="text-5xl font-extrabold mb-4 text-bubblegum drop-shadow text-center">Welcome to Page Hoppers!</h1>
        <p className="mb-8 text-xl text-charcoal text-center max-w-lg">
          Track your reading adventures in and out of school.<br />
          Parents, log in to manage your family's reading journey!
        </p>
        <div className="flex flex-col sm:flex-row gap-4 w-full max-w-sm">
          <Link
            href="/login"
            className="px-8 py-3 bg-bubblegum text-white text-xl font-bold rounded-xl shadow hover:bg-lemon hover:text-charcoal hover:border-sky border-2 border-bubblegum transition text-center"
          >
            Parent Login
          </Link>
          <Link
            href="/register"
            className="px-8 py-3 bg-sky text-charcoal text-xl font-bold rounded-xl shadow hover:bg-lemon hover:text-bubblegum border-2 border-sky transition text-center"
          >
            Register
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Home;
