import { FolderOpen } from "lucide-react";

interface EmptyStateProps {
  type: "project" | "environment";
}

export default function EmptyState({ type }: EmptyStateProps) {
  if (type === "project") {
    return (
      <div className="h-full flex flex-col items-center justify-center text-center max-w-lg mx-auto opacity-50">
        <img src="/logo.png" alt="Env Vault Logo" className="w-24 h-24 mb-6 opacity-80 rounded-xl drop-shadow-lg" />
        <h2 className="font-serif text-4xl mb-4 font-bold text-noir-800">Welcome to Env-Vault</h2>
        <p className="font-sans text-lg text-noir-600">
          Select a project from the sidebar or create a new one to start securely managing your environment variables.
        </p>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col items-center justify-center text-center max-w-lg mx-auto opacity-50">
      <FolderOpen className="w-24 h-24 mb-6 text-vanilla-400" />
      <h2 className="font-serif text-4xl mb-4 font-bold text-noir-800">Project Selected</h2>
      <p className="font-sans text-lg text-noir-600">
        Now, select an environment (like 'development' or 'production') to view its secrets.
      </p>
    </div>
  );
}
