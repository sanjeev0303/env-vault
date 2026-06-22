import { Project, Environment } from "../utils/api";
import { FolderLock, Plus, ChevronRight, FolderOpen, Trash2 } from "lucide-react";

interface SidebarProps {
  projects: Project[];
  environments: Environment[];
  selectedProject: Project | null;
  selectedEnvironment: Environment | null;
  onSelectProject: (project: Project) => void;
  onSelectEnvironment: (env: Environment) => void;
  onDeleteProject: (projectId: string) => void;
  onDeleteEnvironment: (env: Environment) => void;
  setProjectModalOpen: (open: boolean) => void;
  setEnvModalOpen: (open: boolean) => void;
}

export default function Sidebar({
  projects,
  environments,
  selectedProject,
  selectedEnvironment,
  onSelectProject,
  onSelectEnvironment,
  onDeleteProject,
  onDeleteEnvironment,
  setProjectModalOpen,
  setEnvModalOpen,
}: SidebarProps) {
  return (
    <div
      className="w-80 border-r border-noir-600 flex flex-col shadow-2xl relative"
      style={{ backgroundImage: "url('/noir_marble.jpeg')", backgroundSize: "cover", backgroundPosition: "center" }}
    >
      <div className="absolute inset-0 bg-noir-900/70 z-0 backdrop-blur-[1px]"></div>

      <div className="relative z-10 flex flex-col h-full">
        <div className="p-6 border-b border-white/10 flex items-center justify-between">
          <div className="flex items-center gap-3 text-gold-accent">
            <img src="/logo.png" alt="Env Vault Logo" className="w-8 h-8 drop-shadow-md rounded" />
            <h1 className="font-serif text-3xl font-bold tracking-wider drop-shadow-md">Env Vault</h1>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-4 custom-scrollbar">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xs uppercase tracking-widest text-vanilla-400 font-semibold drop-shadow">Projects</h2>
            <button
              onClick={() => setProjectModalOpen(true)}
              className="p-1 hover:bg-white/10 rounded text-vanilla-300 hover:text-vanilla-100 transition-colors backdrop-blur-sm"
            >
              <Plus className="w-4 h-4" />
            </button>
          </div>

          <div className="space-y-2">
            {projects.map((proj) => (
              <div key={proj.id} className="space-y-1">
                <div
                  className={`flex items-center justify-between p-3 rounded-lg cursor-pointer border transition-all backdrop-blur-sm ${
                    selectedProject?.id === proj.id
                      ? "bg-white/10 border-gold-accent shadow-md text-gold-accent"
                      : "border-transparent hover:bg-white/5 hover:border-white/10 text-vanilla-200"
                  }`}
                  onClick={() => onSelectProject(proj)}
                >
                  <div className="flex items-center gap-3 flex-1 min-w-0">
                    {selectedProject?.id === proj.id ? (
                      <FolderOpen className="w-5 h-5 flex-shrink-0 drop-shadow" />
                    ) : (
                      <FolderLock className="w-5 h-5 flex-shrink-0 drop-shadow" />
                    )}
                    <span className="font-medium drop-shadow truncate">{proj.name}</span>
                  </div>
                  {selectedProject?.id === proj.id && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        onDeleteProject(proj.id);
                      }}
                      className="text-red-400 hover:text-red-300 p-1 flex-shrink-0"
                    >
                      <Trash2 className="w-4 h-4 drop-shadow" />
                    </button>
                  )}
                </div>

                {/* Environments List */}
                {selectedProject?.id === proj.id && (
                  <div className="pl-6 pr-2 py-2 space-y-1 relative">
                    <div className="absolute left-6 top-0 bottom-0 w-px bg-white/20"></div>
                    <div className="flex items-center justify-between pl-4 mb-2">
                      <span className="text-xs uppercase tracking-widest text-vanilla-400 drop-shadow">Environments</span>
                      <button
                        onClick={() => setEnvModalOpen(true)}
                        className="text-xs flex items-center gap-1 text-gold-accent hover:text-vanilla-100 transition-colors drop-shadow"
                      >
                        <Plus className="w-3 h-3" /> Add
                      </button>
                    </div>
                    {environments.map((env) => (
                      <div
                        key={env.id}
                        onClick={() => onSelectEnvironment(env)}
                        className={`flex items-center justify-between pl-4 p-2 rounded-md cursor-pointer border text-sm transition-all backdrop-blur-sm ${
                          selectedEnvironment?.id === env.id
                            ? "bg-white/10 border-vanilla-400 text-vanilla-100"
                            : "border-transparent hover:bg-white/5 text-vanilla-300"
                        }`}
                      >
                        <div className="flex items-center gap-2 flex-1 min-w-0">
                          <ChevronRight
                            className={`w-3 h-3 flex-shrink-0 transition-transform drop-shadow ${
                              selectedEnvironment?.id === env.id ? "rotate-90 text-vanilla-400" : ""
                            }`}
                          />
                          <span className="drop-shadow truncate">{env.name}</span>
                        </div>
                        {selectedEnvironment?.id === env.id && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              onDeleteEnvironment(env);
                            }}
                            className="text-red-400 hover:text-red-300 p-1 flex-shrink-0"
                          >
                            <Trash2 className="w-3 h-3 drop-shadow" />
                          </button>
                        )}
                      </div>
                    ))}
                    {environments.length === 0 && (
                      <div className="pl-4 py-2 text-xs text-vanilla-400 italic drop-shadow">No environments yet.</div>
                    )}
                  </div>
                )}
              </div>
            ))}
            {projects.length === 0 && (
              <div className="text-center p-6 border border-dashed border-white/20 rounded-lg text-vanilla-400 text-sm backdrop-blur-sm drop-shadow">
                No projects found. Create one to get started.
              </div>
            )}
          </div>
        </div>

        {/* User / Footer Area */}
        <div className="p-4 border-t border-white/10 text-xs text-center text-vanilla-400 drop-shadow"></div>
      </div>
    </div>
  );
}
