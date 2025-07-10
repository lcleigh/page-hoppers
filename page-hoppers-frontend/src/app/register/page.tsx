// app/register/page.tsx
'use client';
import { useState } from "react";
import { useRouter } from 'next/navigation';
import Link from 'next/link';

const Register = () => {
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");
    const router = useRouter();

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setSuccess("");

        if (password !== confirmPassword) {
            setError("Passwords do not match");
            return;
        }

        if (password.length < 6) {
            setError("Password must be at least 6 characters long");
            return;
        }

        try {
            const res = await fetch(`http://localhost:8080/api/auth/parent/register`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    name,
                    email,
                    password,
                }),
            });

            if (!res.ok) {
                const errorData = await res.json();
                throw new Error(errorData.message || "Registration failed");
            }

            setSuccess("Registration successful! You can now log in.");
            setTimeout(() => {
                router.push('/login');
            }, 2000);
        } catch (err: any) {
            setError(err.message || "Failed to register");
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-powder">
            <div className="max-w-md w-full space-y-8 p-8 bg-lavender rounded-2xl shadow-lg flex flex-col items-center">
                <div>
                    <h2 className="mt-6 text-center text-4xl font-extrabold text-bubblegum drop-shadow">
                        Create Account
                    </h2>
                    <p className="mt-2 text-center text-lg text-charcoal">
                        Join Page Hoppers to track your family's reading journey
                    </p>
                </div>
                <form className="mt-8 space-y-6 w-full" onSubmit={handleRegister}>
                    {error && (
                        <div className="text-coral text-sm text-center font-semibold">{error}</div>
                    )}
                    {success && (
                        <div className="text-leaf text-sm text-center font-semibold">{success}</div>
                    )}
                    <div className="space-y-4">
                        <div>
                            <label htmlFor="name" className="sr-only">
                                Full Name
                            </label>
                            <input
                                id="name"
                                name="name"
                                type="text"
                                required
                                className="appearance-none relative block w-full px-3 py-3 border-2 border-sky placeholder-coolgray text-charcoal focus:outline-none focus:border-bubblegum focus:ring-2 focus:ring-bubblegum focus:z-10 text-lg rounded-xl"
                                placeholder="Full Name"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                            />
                        </div>
                        <div>
                            <label htmlFor="email" className="sr-only">
                                Email address
                            </label>
                            <input
                                id="email"
                                name="email"
                                type="email"
                                required
                                className="appearance-none relative block w-full px-3 py-3 border-2 border-sky placeholder-coolgray text-charcoal focus:outline-none focus:border-bubblegum focus:ring-2 focus:ring-bubblegum focus:z-10 text-lg rounded-xl"
                                placeholder="Email address"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                        </div>
                        <div>
                            <label htmlFor="password" className="sr-only">
                                Password
                            </label>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                required
                                className="appearance-none relative block w-full px-3 py-3 border-2 border-sky placeholder-coolgray text-charcoal focus:outline-none focus:border-bubblegum focus:ring-2 focus:ring-bubblegum focus:z-10 text-lg rounded-xl"
                                placeholder="Password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                        </div>
                        <div>
                            <label htmlFor="confirmPassword" className="sr-only">
                                Confirm Password
                            </label>
                            <input
                                id="confirmPassword"
                                name="confirmPassword"
                                type="password"
                                required
                                className="appearance-none relative block w-full px-3 py-3 border-2 border-sky placeholder-coolgray text-charcoal focus:outline-none focus:border-bubblegum focus:ring-2 focus:ring-bubblegum focus:z-10 text-lg rounded-xl"
                                placeholder="Confirm Password"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                            />
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            className="group relative w-full flex justify-center py-3 px-4 border-2 border-bubblegum text-lg font-bold rounded-xl text-white bg-bubblegum hover:bg-lemon hover:text-charcoal hover:border-sky focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-bubblegum transition"
                        >
                            Create Account
                        </button>
                    </div>
                </form>
                <div className="text-center">
                    <p className="text-coolgray">
                        Already have an account?{' '}
                        <Link href="/login" className="text-bubblegum hover:text-sky font-semibold">
                            Sign in here
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default Register; 