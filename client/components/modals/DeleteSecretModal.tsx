import { Secret } from "../../utils/api";
import { ShieldAlert } from "lucide-react";

interface DeleteSecretModalProps {
  secretToDelete: Secret | null;
  onClose: () => void;
  onConfirm: () => void;
}

export default function DeleteSecretModal({
  secretToDelete,
  onClose,
  onConfirm,
}: DeleteSecretModalProps) {
  if (!secretToDelete) return null;

  return (
    <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-red-300">
        <div className="flex items-center gap-4 mb-4 text-red-600">
          <ShieldAlert className="w-10 h-10" />
          <h2 className="font-serif text-3xl font-bold">Delete Secret</h2>
        </div>
        <p className="text-noir-700 mb-6 font-sans leading-relaxed">
          Are you sure you want to delete the secret{" "}
          <span className="font-mono font-bold text-noir-900 bg-vanilla-200 px-1 rounded">{secretToDelete.key}</span>?
          This action cannot be undone.
        </p>
        <div className="flex justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="px-5 py-2 font-medium text-noir-600 hover:bg-vanilla-200 rounded transition-colors"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="px-5 py-2 font-medium bg-red-600 text-white rounded hover:bg-red-700 transition-colors shadow-md border border-red-800"
          >
            Delete Forever
          </button>
        </div>
      </div>
    </div>
  );
}
