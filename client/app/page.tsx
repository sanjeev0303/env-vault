"use client";

import { useEffect, useState } from "react";
import { api, Project, Environment, Secret } from "../utils/api";
import { 
  FolderLock, 
  Plus, 
  ChevronRight, 
  FolderOpen, 
  KeyRound,
  Eye,
  EyeOff,
  Copy,
  Trash2,
  Edit2,
  ShieldAlert
} from "lucide-react";

export default function Home() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment | null>(null);
  const [secrets, setSecrets] = useState<Secret[]>([]);
  
  // Modals state
  const [isProjectModalOpen, setProjectModalOpen] = useState(false);
  const [isEnvModalOpen, setEnvModalOpen] = useState(false);
  const [isSecretModalOpen, setSecretModalOpen] = useState(false);
  const [editingSecret, setEditingSecret] = useState<Secret | null>(null);
  const [envToDelete, setEnvToDelete] = useState<Environment | null>(null);
  const [secretToDelete, setSecretToDelete] = useState<Secret | null>(null);

  // Forms state
  const [projectName, setProjectName] = useState("");
  const [envName, setEnvName] = useState("");
  const [secretKey, setSecretKey] = useState("");
  const [secretValue, setSecretValue] = useState("");

  const [revealedSecrets, setRevealedSecrets] = useState<Record<string, boolean>>({});

  useEffect(() => {
    loadProjects();
  }, []);

  const loadProjects = async () => {
    try {
      const data = await api.getProjects();
      setProjects(data || []);
    } catch (e) {
      console.error(e);
    }
  };

  const loadEnvironments = async (projectId: string) => {
    try {
      const data = await api.getEnvironments(projectId);
      setEnvironments(data || []);
    } catch (e) {
      console.error(e);
    }
  };

  const loadSecrets = async (projectId: string, envId: string) => {
    try {
      const data = await api.getSecrets(projectId, envId);
      setSecrets(data || []);
    } catch (e) {
      console.error(e);
    }
  };

  const handleSelectProject = (project: Project) => {
    setSelectedProject(project);
    setSelectedEnvironment(null);
    setSecrets([]);
    loadEnvironments(project.id);
  };

  const handleSelectEnvironment = (env: Environment) => {
    setSelectedEnvironment(env);
    if (selectedProject) {
      loadSecrets(selectedProject.id, env.id);
    }
  };

  const handleCreateProject = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const newProj = await api.createProject(projectName);
      setProjects([...projects, newProj]);
      setProjectModalOpen(false);
      setProjectName("");
      handleSelectProject(newProj);
    } catch (e) {
      console.error(e);
    }
  };

  const handleCreateEnvironment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProject) return;
    try {
      const newEnv = await api.createEnvironment(selectedProject.id, envName);
      setEnvironments([...environments, newEnv]);
      setEnvModalOpen(false);
      setEnvName("");
      handleSelectEnvironment(newEnv);
    } catch (e) {
      console.error(e);
    }
  };

  const handleSaveSecret = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProject || !selectedEnvironment) return;
    
    try {
      if (editingSecret) {
        await api.updateSecret(selectedProject.id, selectedEnvironment.id, editingSecret.id, secretValue);
      } else {
        await api.createSecret(selectedProject.id, selectedEnvironment.id, secretKey, secretValue);
      }
      loadSecrets(selectedProject.id, selectedEnvironment.id);
      setSecretModalOpen(false);
      setSecretKey("");
      setSecretValue("");
      setEditingSecret(null);
    } catch (e) {
      console.error(e);
    }
  };

  const handleDeleteSecret = (secret: Secret) => {
    setSecretToDelete(secret);
  };

  const confirmDeleteSecret = async () => {
    if (!selectedProject || !selectedEnvironment || !secretToDelete) return;
    try {
      await api.deleteSecret(selectedProject.id, selectedEnvironment.id, secretToDelete.id);
      loadSecrets(selectedProject.id, selectedEnvironment.id);
      setSecretToDelete(null);
    } catch (e) {
      console.error(e);
    }
  };

  const handleDeleteEnvironment = (env: Environment) => {
    setEnvToDelete(env);
  };

  const confirmDeleteEnvironment = async () => {
    if (!selectedProject || !envToDelete) return;
    try {
      await api.deleteEnvironment(selectedProject.id, envToDelete.id);
      loadEnvironments(selectedProject.id);
      if (selectedEnvironment?.id === envToDelete.id) {
        setSelectedEnvironment(null);
        setSecrets([]);
      }
      setEnvToDelete(null);
    } catch (e) {
      console.error(e);
    }
  };

  const handleDeleteProject = async (projectId: string) => {
    if (!confirm("Are you sure you want to delete this project?")) return;
    try {
      await api.deleteProject(projectId);
      loadProjects();
      setSelectedProject(null);
      setSelectedEnvironment(null);
      setEnvironments([]);
      setSecrets([]);
    } catch (e) {
      console.error(e);
    }
  };

  const toggleReveal = (secretId: string) => {
    setRevealedSecrets(prev => ({ ...prev, [secretId]: !prev[secretId] }));
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const copyAllSecrets = () => {
    if (secrets.length === 0) return;
    const allEnvString = secrets.map(s => `${s.key.toUpperCase()}=${s.value}`).join('\n');
    copyToClipboard(allEnvString);
  };

  return (
    <div className="flex h-screen overflow-hidden bg-noir-900 text-vanilla-100 font-sans">
      
      {/* SIDEBAR: Noir Marble Image */}
      <div 
        className="w-80 border-r border-noir-600 flex flex-col shadow-2xl relative"
        style={{ backgroundImage: "url('/noir_marble.jpeg')", backgroundSize: "cover", backgroundPosition: "center" }}
      >
        <div className="absolute inset-0 bg-noir-900/70 z-0 backdrop-blur-[1px]"></div>
        
        <div className="relative z-10 flex flex-col h-full">
          <div className="p-6 border-b border-white/10 flex items-center justify-between">
            <div className="flex items-center gap-3 text-gold-accent">
              <ShieldAlert className="w-8 h-8 drop-shadow-md" />
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
                    onClick={() => handleSelectProject(proj)}
                  >
                    <div className="flex items-center gap-3">
                      {selectedProject?.id === proj.id ? <FolderOpen className="w-5 h-5 drop-shadow" /> : <FolderLock className="w-5 h-5 drop-shadow" />}
                      <span className="font-medium drop-shadow">{proj.name}</span>
                    </div>
                    {selectedProject?.id === proj.id && (
                      <button 
                        onClick={(e) => { e.stopPropagation(); handleDeleteProject(proj.id); }}
                        className="text-red-400 hover:text-red-300 p-1"
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
                          onClick={() => handleSelectEnvironment(env)}
                          className={`flex items-center justify-between pl-4 p-2 rounded-md cursor-pointer border text-sm transition-all backdrop-blur-sm ${
                            selectedEnvironment?.id === env.id
                              ? "bg-white/10 border-vanilla-400 text-vanilla-100"
                              : "border-transparent hover:bg-white/5 text-vanilla-300"
                          }`}
                        >
                          <div className="flex items-center gap-2">
                            <ChevronRight className={`w-3 h-3 transition-transform drop-shadow ${selectedEnvironment?.id === env.id ? "rotate-90 text-vanilla-400" : ""}`} />
                            <span className="drop-shadow">{env.name}</span>
                          </div>
                          {selectedEnvironment?.id === env.id && (
                            <button 
                              onClick={(e) => { e.stopPropagation(); handleDeleteEnvironment(env); }}
                              className="text-red-400 hover:text-red-300 p-1"
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
          <div className="p-4 border-t border-white/10 text-xs text-center text-vanilla-400 drop-shadow">
          </div>
        </div>
      </div>

      {/* MAIN CONTENT AREA: French Vanilla (inverted slightly for dark mode elegance) */}
      <div className="flex-1 flex flex-col bg-vanilla-100 text-noir-900 relative">
        
        {/* Top Bar */}
        <div className="h-20 border-b border-vanilla-300 px-8 flex items-center justify-between bg-vanilla-200 shadow-sm z-10">
          <div className="flex items-center gap-3">
            {selectedProject ? (
              <>
                <h2 className="font-serif text-3xl font-bold">{selectedProject.name}</h2>
                {selectedEnvironment && (
                  <>
                    <ChevronRight className="w-6 h-6 text-vanilla-400" />
                    <span className="font-sans text-xl font-medium text-noir-600 bg-vanilla-300 px-3 py-1 rounded-full">
                      {selectedEnvironment.name}
                    </span>
                  </>
                )}
              </>
            ) : (
              <h2 className="font-serif text-3xl font-bold text-noir-600 italic">Select a project to view secrets</h2>
            )}
          </div>
          
          {selectedEnvironment && (
            <div className="flex items-center gap-3">
              {secrets.length > 0 && (
                <button 
                  onClick={copyAllSecrets}
                  className="flex items-center gap-2 bg-vanilla-300 text-noir-800 px-5 py-2 rounded shadow hover:bg-vanilla-400 transition-colors font-medium border border-vanilla-400"
                  title="Copy All Environments"
                >
                  <Copy className="w-4 h-4" /> Copy All
                </button>
              )}
              <button 
                onClick={() => {
                  setEditingSecret(null);
                  setSecretKey("");
                  setSecretValue("");
                  setSecretModalOpen(true);
                }}
                className="flex items-center gap-2 bg-noir-800 text-vanilla-100 px-6 py-2 rounded shadow hover:bg-noir-900 transition-colors font-medium border border-noir-600"
              >
                <Plus className="w-5 h-5" /> Add Secret
              </button>
            </div>
          )}
        </div>

        {/* Dashboard Content */}
        <div className="flex-1 p-8 overflow-y-auto">
          {!selectedProject ? (
            <div className="h-full flex flex-col items-center justify-center text-center max-w-lg mx-auto opacity-50">
              <ShieldAlert className="w-24 h-24 mb-6 text-vanilla-400" />
              <h2 className="font-serif text-4xl mb-4 font-bold text-noir-800">Welcome to Env-Vault</h2>
              <p className="font-sans text-lg text-noir-600">Select a project from the sidebar or create a new one to start securely managing your environment variables.</p>
            </div>
          ) : !selectedEnvironment ? (
            <div className="h-full flex flex-col items-center justify-center text-center max-w-lg mx-auto opacity-50">
              <FolderOpen className="w-24 h-24 mb-6 text-vanilla-400" />
              <h2 className="font-serif text-4xl mb-4 font-bold text-noir-800">Project Selected</h2>
              <p className="font-sans text-lg text-noir-600">Now, select an environment (like 'development' or 'production') to view its secrets.</p>
            </div>
          ) : (
            <div className="bg-white rounded-xl shadow-lg border border-vanilla-300 overflow-hidden">
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
                      secrets.map((secret) => (
                        <tr key={secret.id} className="border-b border-vanilla-200 hover:bg-vanilla-100 transition-colors group">
                          <td className="p-4 pl-6">
                            <div className="flex items-center gap-3">
                              <KeyRound className="w-4 h-4 text-gold-accent" />
                              <span className="font-mono font-bold text-noir-800 bg-vanilla-200 px-2 py-1 rounded">
                                {secret.key}
                              </span>
                            </div>
                          </td>
                          <td className="p-4">
                            <div className="flex items-center gap-3">
                              <span className="font-mono text-noir-700 break-all">
                                {revealedSecrets[secret.id] ? secret.value : '••••••••••••••••'}
                              </span>
                            </div>
                          </td>
                          <td className="p-4 pr-6 text-right">
                            <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                              <button 
                                onClick={() => toggleReveal(secret.id)}
                                className="p-2 text-noir-600 hover:text-noir-900 hover:bg-vanilla-300 rounded"
                                title={revealedSecrets[secret.id] ? "Hide" : "Reveal"}
                              >
                                {revealedSecrets[secret.id] ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                              </button>
                              <button 
                                onClick={() => copyToClipboard(`${secret.key.toUpperCase()}=${secret.value}`)}
                                className="p-2 text-noir-600 hover:text-noir-900 hover:bg-vanilla-300 rounded"
                                title="Copy Key=Value"
                              >
                                <Copy className="w-4 h-4" />
                              </button>
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
                                onClick={() => handleDeleteSecret(secret)}
                                className="p-2 text-noir-600 hover:text-red-600 hover:bg-red-50 rounded"
                                title="Delete"
                              >
                                <Trash2 className="w-4 h-4" />
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* MODALS */}
      
      {/* Create Project Modal */}
      {isProjectModalOpen && (
        <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-vanilla-300">
            <h2 className="font-serif text-3xl font-bold mb-6">Create Project</h2>
            <form onSubmit={handleCreateProject}>
              <div className="mb-6">
                <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Project Name</label>
                <input 
                  type="text" 
                  value={projectName}
                  onChange={e => setProjectName(e.target.value)}
                  className="w-full bg-white border border-vanilla-400 p-3 rounded text-noir-900 font-sans focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all"
                  placeholder="e.g. My Awesome App"
                  required
                />
              </div>
              <div className="flex justify-end gap-3">
                <button 
                  type="button" 
                  onClick={() => setProjectModalOpen(false)}
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
      )}

      {/* Create Environment Modal */}
      {isEnvModalOpen && (
        <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-vanilla-300">
            <h2 className="font-serif text-3xl font-bold mb-6">Create Environment</h2>
            <p className="text-noir-600 mb-6 text-sm">For project: <span className="font-bold">{selectedProject?.name}</span></p>
            <form onSubmit={handleCreateEnvironment}>
              <div className="mb-6">
                <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Environment Name</label>
                <input 
                  type="text" 
                  value={envName}
                  onChange={e => setEnvName(e.target.value)}
                  className="w-full bg-white border border-vanilla-400 p-3 rounded text-noir-900 font-sans focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all"
                  placeholder="e.g. production, staging, development"
                  required
                />
              </div>
              <div className="flex justify-end gap-3">
                <button 
                  type="button" 
                  onClick={() => setEnvModalOpen(false)}
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
      )}

      {/* Secret Modal (Create / Edit) */}
      {isSecretModalOpen && (
        <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-2xl border border-vanilla-300">
            <h2 className="font-serif text-3xl font-bold mb-6">{editingSecret ? 'Edit Secret' : 'Add Secret'}</h2>
            <div className="flex gap-2 items-center mb-6 text-sm bg-vanilla-200 p-3 rounded border border-vanilla-300">
              <span className="font-bold text-noir-800">{selectedProject?.name}</span>
              <ChevronRight className="w-4 h-4 text-noir-600" />
              <span className="font-bold text-noir-800">{selectedEnvironment?.name}</span>
            </div>
            
            <form onSubmit={handleSaveSecret}>
              <div className="mb-4">
                <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Key</label>
                <input 
                  type="text" 
                  value={secretKey}
                  onChange={e => setSecretKey(e.target.value)}
                  disabled={!!editingSecret} // Cannot change key while editing
                  className={`w-full font-mono bg-white border border-vanilla-400 p-3 rounded text-noir-900 focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all ${editingSecret ? 'opacity-60 cursor-not-allowed bg-vanilla-200' : ''}`}
                  placeholder="e.g. DATABASE_URL"
                  required
                />
              </div>
              <div className="mb-6">
                <label className="block text-sm font-semibold uppercase tracking-wider text-noir-600 mb-2">Value</label>
                <textarea 
                  value={secretValue}
                  onChange={e => setSecretValue(e.target.value)}
                  className="w-full font-mono bg-white border border-vanilla-400 p-3 rounded text-noir-900 focus:outline-none focus:border-gold-accent focus:ring-1 focus:ring-gold-accent transition-all h-32 custom-scrollbar"
                  placeholder="e.g. postgresql://user:pass@host:5432/db"
                  required
                />
              </div>
              <div className="flex justify-end gap-3">
                <button 
                  type="button" 
                  onClick={() => setSecretModalOpen(false)}
                  className="px-5 py-2 font-medium text-noir-600 hover:bg-vanilla-200 rounded transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  className="px-5 py-2 font-medium bg-noir-800 text-vanilla-100 rounded hover:bg-noir-900 transition-colors shadow-md border border-noir-600"
                >
                  {editingSecret ? 'Save Changes' : 'Add Secret'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Delete Environment Confirmation Modal */}
      {envToDelete && (
        <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-red-300">
            <div className="flex items-center gap-4 mb-4 text-red-600">
              <ShieldAlert className="w-10 h-10" />
              <h2 className="font-serif text-3xl font-bold">Delete Environment</h2>
            </div>
            <p className="text-noir-700 mb-6 font-sans leading-relaxed">
              Are you sure you want to delete the <span className="font-bold text-noir-900">'{envToDelete.name}'</span> environment? This action cannot be undone and will destroy all secrets within it.
            </p>
            <div className="flex justify-end gap-3">
              <button 
                type="button" 
                onClick={() => setEnvToDelete(null)}
                className="px-5 py-2 font-medium text-noir-600 hover:bg-vanilla-200 rounded transition-colors"
              >
                Cancel
              </button>
              <button 
                type="button"
                onClick={confirmDeleteEnvironment}
                className="px-5 py-2 font-medium bg-red-600 text-white rounded hover:bg-red-700 transition-colors shadow-md border border-red-800"
              >
                Delete Forever
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Delete Secret Confirmation Modal */}
      {secretToDelete && (
        <div className="fixed inset-0 bg-noir-900/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-vanilla-100 text-noir-900 rounded-xl shadow-2xl p-8 w-full max-w-md border border-red-300">
            <div className="flex items-center gap-4 mb-4 text-red-600">
              <ShieldAlert className="w-10 h-10" />
              <h2 className="font-serif text-3xl font-bold">Delete Secret</h2>
            </div>
            <p className="text-noir-700 mb-6 font-sans leading-relaxed">
              Are you sure you want to delete the secret <span className="font-mono font-bold text-noir-900 bg-vanilla-200 px-1 rounded">{secretToDelete.key}</span>? This action cannot be undone.
            </p>
            <div className="flex justify-end gap-3">
              <button 
                type="button" 
                onClick={() => setSecretToDelete(null)}
                className="px-5 py-2 font-medium text-noir-600 hover:bg-vanilla-200 rounded transition-colors"
              >
                Cancel
              </button>
              <button 
                type="button"
                onClick={confirmDeleteSecret}
                className="px-5 py-2 font-medium bg-red-600 text-white rounded hover:bg-red-700 transition-colors shadow-md border border-red-800"
              >
                Delete Forever
              </button>
            </div>
          </div>
        </div>
      )}

    </div>
  );
}
