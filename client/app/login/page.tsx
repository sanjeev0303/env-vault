"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useLogin } from "../../utils/hooks";
import { ShieldAlert, Eye, EyeOff } from "lucide-react";

export default function LoginPage() {
  const router = useRouter();
  const [authEmail, setAuthEmail] = useState("");
  const [authPassword, setAuthPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const loginMutation = useLogin();

  useEffect(() => {
    // If already authenticated, redirect to home
    const token = localStorage.getItem("access_token");
    if (token) {
      router.push("/");
    }
  }, [router]);

  const handleAuthSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    loginMutation.mutate({ email: authEmail, password: authPassword }, {
      onSuccess: () => {
        router.push("/");
      }
    });
  };

  return (
    <div
      className="flex h-screen items-center justify-center bg-noir-900 text-vanilla-100 font-sans relative"
      style={{ backgroundImage: "url('/noir_marble.jpeg')", backgroundSize: "cover", backgroundPosition: "center" }}
    >
      <div className="absolute inset-0 bg-noir-900/70 z-0 backdrop-blur-[1px]"></div>

      <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-vanilla-300 relative z-10">
        <div className="flex justify-center mb-6">
          <img src="/logo.png" alt="Env Vault Logo" className="w-16 h-16 drop-shadow-md" />
        </div>
        <h2 className="font-serif text-3xl font-bold text-center mb-6">Welcome to Env Vault</h2>
        <form onSubmit={handleAuthSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Email</label>
            <input 
              type="email" 
              value={authEmail}
              onChange={e => setAuthEmail(e.target.value)}
              className="w-full bg-white border border-vanilla-400 p-3 rounded text-noir-900 focus:outline-none focus:border-gold-accent transition-all"
              required
            />
          </div>
          <div className="mb-6">
            <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Password</label>
            <div className="relative">
              <input 
                type={showPassword ? "text" : "password"}
                value={authPassword}
                onChange={e => setAuthPassword(e.target.value)}
                className="w-full bg-white border border-vanilla-400 p-3 pr-10 rounded text-noir-900 focus:outline-none focus:border-gold-accent transition-all"
                required
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-noir-400 hover:text-gold-accent transition-colors focus:outline-none"
              >
                {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>
          <button 
            type="submit"
            disabled={loginMutation.isPending}
            className="w-full py-3 mb-4 font-bold bg-noir-800 text-vanilla-100 rounded hover:bg-noir-900 transition-colors shadow-md disabled:opacity-50"
          >
            Login
          </button>
          <div className="text-center text-sm text-noir-600">
            Don't have an account?{" "}
            <Link 
              href="/register"
              className="font-bold text-gold-accent hover:underline"
            >
              Register here
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
