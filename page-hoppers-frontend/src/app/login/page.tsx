// app/login/page.tsx
'use client'; // This marks the component as a client-side component

import { useState } from "react";
import { useRouter } from 'next/navigation';

const Login = () => {
    const [role, setRole] = useState("parent");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
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
            router.push('/dashboard'); // Redirect to parent dashboard
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
                        <div>
                            <label htmlFor="password" className="sr-only">
                                Password
                            </label>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                required
                                className="appearance-none rounded-b-xl relative block w-full px-3 py-3 border-2 border-sky placeholder-coolgray text-charcoal focus:outline-none focus:border-bubblegum focus:ring-2 focus:ring-bubblegum focus:z-10 text-lg"
                                placeholder="Password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
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
