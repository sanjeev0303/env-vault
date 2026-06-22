import { Environment, Secret } from "../utils/api";
import { TerminalSquare, Plus, Copy, KeyRound, EyeOff, Eye, Edit2, Trash2 } from "lucide-react";

interface SecretsTableProps {
  selectedEnvironment: Environment;
  secrets: Secret[];
  revealedSecrets: Record<string, string>;
  onToggleReveal: (secretId: string) => void;
  onCopyToClipboard: (text: string) => void;
  hasNextPage: boolean;
  onFetchNextPage: () => void;
  isFetchingNextPage: boolean;
  setSecretModalOpen: (open: boolean) => void;
  setExportModalOpen: (open: boolean) => void;
  setEditingSecret: (secret: Secret | null) => void;
  setSecretKey: (key: string) => void;
  setSecretValue: (value: string) => void;
  onDeleteSecret: (secret: Secret) => void;
}

export default function SecretsTable({
  selectedEnvironment,
  secrets,
  revealedSecrets,
  onToggleReveal,
  onCopyToClipboard,
  hasNextPage,
  onFetchNextPage,
  isFetchingNextPage,
  setSecretModalOpen,
  setExportModalOpen,
  setEditingSecret,
  setSecretKey,
  setSecretValue,
  onDeleteSecret,
}: SecretsTableProps) {
  return (
    <div className="flex-1 overflow-hidden flex flex-col bg-vanilla-100 rounded-xl shadow-2xl border border-vanilla-300">
      <div className="p-6 border-b border-vanilla-300 flex justify-between items-center bg-vanilla-200">
        <h3 className="text-xl font-bold font-serif text-noir-900 flex items-center gap-2 flex-1 min-w-0">
          <TerminalSquare className="w-6 h-6 text-gold-accent flex-shrink-0" />
          <span className="truncate">{selectedEnvironment.name} Secrets</span>
        </h3>
        <div className="flex gap-2 flex-shrink-0 pl-4">
          <button
            onClick={() => {
              setEditingSecret(null);
              setSecretKey("");
              setSecretValue("");
              setSecretModalOpen(true);
            }}
            className="flex items-center gap-2 px-4 py-2 bg-noir-800 text-vanilla-100 rounded hover:bg-noir-900 transition-colors text-sm font-semibold shadow-md"
          >
            <Plus className="w-4 h-4" /> Add Secret
          </button>
          <button
            onClick={() => setExportModalOpen(true)}
            className="flex items-center gap-2 px-4 py-2 bg-gold-accent text-noir-900 rounded hover:bg-yellow-600 transition-colors text-sm font-semibold shadow-md"
          >
            <Copy className="w-4 h-4" /> Export All
          </button>
        </div>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse table-fixed">
          <thead>
            <tr className="bg-vanilla-200 text-noir-600 text-sm uppercase tracking-wider font-semibold border-b border-vanilla-300">
              <th className="p-4 pl-6 w-1/3">Key</th>
              <th className="p-4 w-1/2">Value</th>
              <th className="p-4 text-right pr-6 w-1/6">Actions</th>
            </tr>
          </thead>
          <tbody>
            {secrets.length === 0 ? (
              <tr>
                <td colSpan={3} className="p-12 text-center text-noir-600 font-medium">
                  No secrets found for this environment. Click "Add Secret" to create one.
                </td>
              </tr>
            ) : (
              <>
                {secrets.map((secret) => (
                  <tr key={secret.id} className="border-b border-vanilla-200 hover:bg-vanilla-100 transition-colors group">
                    <td className="p-4 pl-6">
                      <div className="flex items-center gap-3">
                        <KeyRound className="w-4 h-4 text-gold-accent" />
                        <span className="font-mono font-bold text-noir-800 bg-vanilla-200 px-2 py-1 rounded">
                          {secret.key}
                        </span>
                      </div>
                    </td>
                    <td className="p-4 text-noir-600">
                      <div className="flex items-center gap-2">
                        <span className={`font-mono ${!revealedSecrets[secret.id] ? "blur-sm select-none" : ""}`}>
                          {revealedSecrets[secret.id] !== undefined ? revealedSecrets[secret.id] : "••••••••"}
                        </span>
                        <button
                          onClick={() => onToggleReveal(secret.id)}
                          className="text-noir-400 hover:text-gold-accent transition-colors focus:outline-none ml-2"
                          title={revealedSecrets[secret.id] ? "Hide value" : "Reveal value"}
                        >
                          {revealedSecrets[secret.id] ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                        </button>
                        {revealedSecrets[secret.id] && (
                          <button
                            onClick={() => onCopyToClipboard(revealedSecrets[secret.id])}
                            className="text-noir-400 hover:text-gold-accent transition-colors focus:outline-none"
                            title="Copy value"
                          >
                            <Copy className="w-4 h-4" />
                          </button>
                        )}
                      </div>
                    </td>
                    <td className="p-4 pr-6 text-right">
                      <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                          onClick={() => {
                            setEditingSecret(secret);
                            setSecretKey(secret.key);
                            setSecretValue(secret.value);
                            setSecretModalOpen(true);
                          }}
                          className="p-2 text-noir-600 hover:text-blue-600 hover:bg-blue-50 rounded"
                          title="Edit"
                        >
                          <Edit2 className="w-4 h-4" />
                        </button>
                        <button
                          onClick={() => onDeleteSecret(secret)}
                          className="p-2 text-noir-600 hover:text-red-600 hover:bg-red-50 rounded"
                          title="Delete"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
                {hasNextPage && (
                  <tr>
                    <td colSpan={3} className="p-4 text-center">
                      <button
                        onClick={onFetchNextPage}
                        disabled={isFetchingNextPage}
                        className="text-noir-600 hover:text-noir-900 font-medium text-sm disabled:opacity-50"
                      >
                        {isFetchingNextPage ? "Loading more..." : "Load More Secrets"}
                      </button>
                    </td>
                  </tr>
                )}
              </>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
