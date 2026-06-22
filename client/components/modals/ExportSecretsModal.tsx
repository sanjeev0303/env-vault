import React from "react";
import { ShieldAlert, EyeOff, Eye } from "lucide-react";

interface ExportSecretsModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (e: React.FormEvent) => void;
  reauthPassword: string;
  setReauthPassword: (pass: string) => void;
  showReauthPassword: boolean;
  setShowReauthPassword: (show: boolean) => void;
  isPending: boolean;
}

export default function ExportSecretsModal({
  isOpen,
  onClose,
  onSubmit,
  reauthPassword,
  setReauthPassword,
  showReauthPassword,
  setShowReauthPassword,
  isPending,
}: ExportSecretsModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-noir-900/80 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-gold-accent">
        <div className="flex items-center gap-4 mb-6">
          <ShieldAlert className="w-10 h-10 text-gold-accent" />
          <h2 className="font-serif text-3xl font-bold">Secure Export</h2>
        </div>
        <p className="text-noir-700 mb-6 font-sans leading-relaxed">
          Exporting secrets requires re-authentication. Please enter your password to proceed.
        </p>
        <form onSubmit={onSubmit}>
          <div className="mb-6">
            <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Password</label>
            <div className="relative">
              <input
                type={showReauthPassword ? "text" : "password"}
                value={reauthPassword}
                onChange={(e) => setReauthPassword(e.target.value)}
                className="w-full bg-white border border-vanilla-400 p-3 pr-10 rounded text-noir-900 focus:outline-none focus:border-gold-accent transition-all"
                required
                autoFocus
              />
              <button
                type="button"
                onClick={() => setShowReauthPassword(!showReauthPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-noir-400 hover:text-gold-accent transition-colors focus:outline-none"
              >
                {showReauthPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => {
                onClose();
                setReauthPassword("");
              }}
              className="px-5 py-2 font-medium text-noir-600 hover:bg-vanilla-200 rounded transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isPending}
              className="px-5 py-2 font-medium bg-gold-accent text-noir-900 rounded hover:bg-yellow-600 transition-colors shadow-md border border-gold-600 disabled:opacity-50"
            >
              Export & Copy
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
