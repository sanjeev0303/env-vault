import React from "react";
import { Project } from "../../utils/api";

interface CreateEnvironmentModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (e: React.FormEvent) => void;
  selectedProject: Project | null;
  envName: string;
  setEnvName: (name: string) => void;
  isError: boolean;
  errorMessage?: string;
  resetMutation: () => void;
}

export default function CreateEnvironmentModal({
  isOpen,
  onClose,
  onSubmit,
  selectedProject,
  envName,
  setEnvName,
  isError,
  errorMessage,
  resetMutation,
}: CreateEnvironmentModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-vanilla-300">
        <h2 className="font-serif text-3xl font-bold mb-6">Create Environment</h2>
        <p className="text-noir-600 mb-6 text-sm">
          For project: <span className="font-bold">{selectedProject?.name}</span>
        </p>
        <form onSubmit={onSubmit}>
          {isError && (
            <div className="mb-4 p-3 bg-red-100 border border-red-300 text-red-700 text-sm rounded">
              {errorMessage || "An environment with this name already exists."}
            </div>
          )}
          <div className="mb-6">
            <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">
              Environment Name
            </label>
            <input
              type="text"
              value={envName}
              onChange={(e) => setEnvName(e.target.value)}
              className="w-full bg-white border border-vanilla-400 p-3 rounded text-noir-900 font-sans focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all"
              placeholder="e.g. production, staging, development"
              required
            />
          </div>
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => {
                onClose();
                resetMutation();
                setEnvName("");
              }}
              className="px-5 py-2 font-medium text-noir-600 hover:bg-vanilla-200 rounded transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-5 py-2 font-medium bg-noir-800 text-vanilla-100 rounded hover:bg-noir-900 transition-colors shadow-md border border-noir-600"
            >
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
