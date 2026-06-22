import { Project, Environment } from "../utils/api";
import { ChevronRight } from "lucide-react";

interface TopBarProps {
  selectedProject: Project | null;
  selectedEnvironment: Environment | null;
  onLogout: () => void;
}

export default function TopBar({ selectedProject, selectedEnvironment, onLogout }: TopBarProps) {
  return (
    <div className="h-20 border-b border-vanilla-300 px-8 flex items-center justify-between bg-vanilla-200 shadow-sm z-10">
      <div className="flex items-center gap-3 flex-1 min-w-0">
        {selectedProject ? (
          <>
            <h2 className="font-serif text-3xl font-bold truncate max-w-[200px] sm:max-w-xs md:max-w-sm lg:max-w-md">
              {selectedProject.name}
            </h2>
            {selectedEnvironment && (
              <>
                <ChevronRight className="w-6 h-6 text-vanilla-400 flex-shrink-0" />
                <span className="font-sans text-xl font-medium text-noir-600 bg-vanilla-300 px-3 py-1 rounded-full truncate max-w-[150px] sm:max-w-[200px] md:max-w-xs lg:max-w-sm">
                  {selectedEnvironment.name}
                </span>
              </>
            )}
          </>
        ) : (
          <h2 className="font-serif text-3xl font-bold text-noir-600 italic">Select a project to view secrets</h2>
        )}
      </div>

      <button
        onClick={onLogout}
        className="px-4 py-2 text-sm font-semibold bg-noir-800 hover:bg-noir-900 rounded text-vanilla-100 transition-colors shadow-md"
      >
        Logout
      </button>
    </div>
  );
}
