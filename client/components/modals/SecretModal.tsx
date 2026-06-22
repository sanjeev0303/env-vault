import React from "react";
import { Project, Environment, Secret } from "../../utils/api";
import { ChevronRight } from "lucide-react";

interface SecretModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (e: React.FormEvent) => void;
  selectedProject: Project | null;
  selectedEnvironment: Environment | null;
  editingSecret: Secret | null;
  secretKey: string;
  setSecretKey: (key: string) => void;
  secretValue: string;
  setSecretValue: (val: string) => void;
  isError: boolean;
  errorMessage?: string;
  resetMutations: () => void;
}

export default function SecretModal({
  isOpen,
  onClose,
  onSubmit,
  selectedProject,
  selectedEnvironment,
  editingSecret,
  secretKey,
  setSecretKey,
  secretValue,
  setSecretValue,
  isError,
  errorMessage,
  resetMutations,
}: SecretModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-2xl border border-vanilla-300">
        <h2 className="font-serif text-3xl font-bold mb-6">{editingSecret ? "Edit Secret" : "Add Secret"}</h2>
        <div className="flex gap-2 items-center mb-6 text-sm bg-vanilla-200 p-3 rounded border border-vanilla-300">
          <span className="font-bold text-noir-800">{selectedProject?.name}</span>
          <ChevronRight className="w-4 h-4 text-noir-600" />
          <span className="font-bold text-noir-800">{selectedEnvironment?.name}</span>
        </div>

        <form onSubmit={onSubmit}>
          {isError && (
            <div className="mb-4 p-3 bg-red-100 border border-red-300 text-red-700 text-sm rounded">
              {errorMessage || "A secret with this key already exists."}
            </div>
          )}
          <div className="mb-4">
            <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Key</label>
            <input
              type="text"
              value={secretKey}
              onChange={(e) => setSecretKey(e.target.value)}
              disabled={!!editingSecret} // Cannot change key while editing
              className={`w-full font-mono bg-white border border-vanilla-400 p-3 rounded text-noir-900 focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all ${
                editingSecret ? "opacity-60 cursor-not-allowed bg-vanilla-200" : ""
              }`}
              placeholder="e.g. DATABASE_URL"
              required
            />
          </div>
          <div className="mb-6">
            <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Value</label>
            <textarea
              value={secretValue}
              onChange={(e) => setSecretValue(e.target.value)}
              className="w-full font-mono bg-white border border-vanilla-400 p-3 rounded text-noir-900 focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all h-32 custom-scrollbar"
              placeholder="e.g. postgresql://user:pass@host:5432/db"
              required
            />
          </div>
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => {
                onClose();
                resetMutations();
              }}
              className="px-5 py-2 font-medium text-noir-600 hover:bg-vanilla-200 rounded transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-5 py-2 font-medium bg-noir-800 text-vanilla-100 rounded hover:bg-noir-900 transition-colors shadow-md border border-noir-600"
            >
              {editingSecret ? "Save Changes" : "Add Secret"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
