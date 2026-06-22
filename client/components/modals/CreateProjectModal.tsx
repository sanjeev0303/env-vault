import React from "react";

interface CreateProjectModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (e: React.FormEvent) => void;
  projectName: string;
  setProjectName: (name: string) => void;
  isError: boolean;
  errorMessage?: string;
  resetMutation: () => void;
}

export default function CreateProjectModal({
  isOpen,
  onClose,
  onSubmit,
  projectName,
  setProjectName,
  isError,
  errorMessage,
  resetMutation,
}: CreateProjectModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-vanilla-300">
        <h2 className="font-serif text-3xl font-bold mb-6">Create Project</h2>
        <form onSubmit={onSubmit}>
          {isError && (
            <div className="mb-4 p-3 bg-red-100 border border-red-300 text-red-700 text-sm rounded">
              {errorMessage || "A project with this name already exists or another error occurred."}
            </div>
          )}
          <div className="mb-6">
            <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">
              Project Name
            </label>
            <input
              type="text"
              value={projectName}
              onChange={(e) => setProjectName(e.target.value)}
              className="w-full bg-white border border-vanilla-400 p-3 rounded text-noir-900 font-sans focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all"
              placeholder="e.g. My Awesome App"
              required
            />
          </div>
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => {
                onClose();
                resetMutation();
                setProjectName("");
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
