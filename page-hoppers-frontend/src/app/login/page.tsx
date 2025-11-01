// app/login/page.tsx
'use client'; // This marks the component as a client-side component

import { useState } from "react";
import { useRouter } from 'next/navigation';

const Login = () => {
    const [role, setRole] = useState("parent");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState("");
    const router = useRouter();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();

        const loginDetails = { email, password };

        let endpoint = "/api/auth/parent/login";
        if (role === "child") {
            endpoint = "/api/auth/child/login";
        }

        try {
            const res = await fetch(`http://localhost:8080${endpoint}`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(loginDetails),
            });

            if (!res.ok) {
                throw new Error("Invalid credentials");
            }

            const data = await res.json();
            // Store the token securely
            localStorage.setItem('parentToken', data.token);
            router.push('/parent-dashboard'); // Redirect to parent dashboard
        } catch (err: any) {
            setError("Failed to login: " + err.message);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-powder">
            <div className="max-w-md w-full space-y-8 p-8 bg-lavender rounded-2xl shadow-lg flex flex-col items-center">
                <div>
                    <h2 className="mt-6 text-center text-4xl font-extrabold text-bubblegum drop-shadow">
                        Parent Sign In
                    </h2>
                    <p className="mt-2 text-center text-lg text-charcoal">
                        Sign in to manage your children's reading progress
                    </p>
                </div>
                <form className="mt-8 space-y-6 w-full" onSubmit={handleLogin}>
                    {error && (
                        <div className="text-coral text-sm text-center font-semibold">{error}</div>
                    )}
                    <div className="rounded-xl shadow-sm -space-y-px">
                        <div>
                            <label htmlFor="email" className="sr-only">
                                Email address
                            </label>
                            <input
                                id="email"
                                name="email"
                                type="email"
                                required
                                className="appearance-none rounded-t-xl relative block w-full px-3 py-3 border-2 border-sky placeholder-coolgray text-charcoal focus:outline-none focus:border-bubblegum focus:ring-2 focus:ring-bubblegum focus:z-10 text-lg"
                                placeholder="Email address"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                            />
                        </div>
                        <div style={{ position: "relative" }}>
                            <input
                                id="password"
                                name="password"
                                type={showPassword ? "text" : "password"}
                                required
                                className="appearance-none rounded-b-xl relative block w-full px-3 py-3 border-2 border-sky placeholder-coolgray text-charcoal focus:outline-none focus:border-bubblegum focus:ring-2 focus:ring-bubblegum focus:z-10 text-lg"
                                placeholder="Password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword((prev) => !prev)}
                                style={{
                                    position: "absolute",
                                    right: 8,
                                    top: "50%",
                                    transform: "translateY(-50%)",
                                    background: "none",
                                    border: "none",
                                    cursor: "pointer"
                                }}
                                aria-label={showPassword ? "Hide password" : "Show password"}
                            >
                                {showPassword ? (
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

                    <div>
                        <button
                            type="submit"
                            className="group relative w-full flex justify-center py-3 px-4 border-2 border-bubblegum text-lg font-bold rounded-xl text-white bg-bubblegum hover:bg-lemon hover:text-charcoal hover:border-sky focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-bubblegum transition"
                        >
                            Sign in
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Login;
