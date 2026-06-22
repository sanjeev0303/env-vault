"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Project, Environment, Secret } from "../utils/api";
import {
  useProjects,
  useCreateProject,
  useDeleteProject,
  useEnvironments,
  useCreateEnvironment,
  useDeleteEnvironment,
  useSecrets,
  useCreateSecret,
  useUpdateSecret,
  useDeleteSecret,
  useRevealSecret,
  useExportSecrets,
  useReauth
} from "../utils/hooks";

import Sidebar from "../components/Sidebar";
import TopBar from "../components/TopBar";
import EmptyState from "../components/EmptyState";
import SecretsTable from "../components/SecretsTable";
import CreateProjectModal from "../components/modals/CreateProjectModal";
import CreateEnvironmentModal from "../components/modals/CreateEnvironmentModal";
import SecretModal from "../components/modals/SecretModal";
import DeleteEnvironmentModal from "../components/modals/DeleteEnvironmentModal";
import DeleteSecretModal from "../components/modals/DeleteSecretModal";
import ExportSecretsModal from "../components/modals/ExportSecretsModal";

export default function Home() {
  const router = useRouter();
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [selectedEnvironment, setSelectedEnvironment] = useState<Environment | null>(null);

  const { data: projects = [] } = useProjects();
  const { data: environments = [] } = useEnvironments(selectedProject?.id);

  const {
    data: secretsData,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage
  } = useSecrets(selectedProject?.id, selectedEnvironment?.id);

  const secrets = secretsData?.pages.flatMap(page => page.data) || [];

  const createProject = useCreateProject();
  const deleteProject = useDeleteProject();
  const createEnvironment = useCreateEnvironment();
  const deleteEnvironment = useDeleteEnvironment();
  const createSecret = useCreateSecret();
  const updateSecret = useUpdateSecret();
  const deleteSecret = useDeleteSecret();

  // Modals state
  const [isProjectModalOpen, setProjectModalOpen] = useState(false);
  const [isEnvModalOpen, setEnvModalOpen] = useState(false);
  const [isSecretModalOpen, setSecretModalOpen] = useState(false);
  const [isExportModalOpen, setExportModalOpen] = useState(false);
  const [editingSecret, setEditingSecret] = useState<Secret | null>(null);
  const [envToDelete, setEnvToDelete] = useState<Environment | null>(null);
  const [secretToDelete, setSecretToDelete] = useState<Secret | null>(null);

  // Forms state
  const [projectName, setProjectName] = useState("");
  const [envName, setEnvName] = useState("");
  const [secretKey, setSecretKey] = useState("");
  const [secretValue, setSecretValue] = useState("");
  const [reauthPassword, setReauthPassword] = useState("");
  const [showReauthPassword, setShowReauthPassword] = useState(false);

  const [revealedSecrets, setRevealedSecrets] = useState<Record<string, string>>({});

  const revealSecretMutation = useRevealSecret();
  const exportSecretsMutation = useExportSecrets();
  const reauthMutation = useReauth();

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    if (!token) {
      router.push("/login");
    }

    const handleAuthError = () => router.push("/login");
    window.addEventListener('auth-error', handleAuthError);
    return () => window.removeEventListener('auth-error', handleAuthError);
  }, [router]);

  const handleLogout = () => {
    localStorage.removeItem("access_token");
    router.push("/login");
  };

  const handleSelectProject = (project: Project) => {
    setSelectedProject(project);
    setSelectedEnvironment(null);
  };

  const handleSelectEnvironment = (env: Environment) => {
    setSelectedEnvironment(env);
  };

  const handleCreateProject = async (e: React.FormEvent) => {
    e.preventDefault();
    createProject.mutate(projectName, {
      onSuccess: (newProj) => {
        setProjectModalOpen(false);
        setProjectName("");
        handleSelectProject(newProj);
      }
    });
  };

  const handleCreateEnvironment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProject) return;
    createEnvironment.mutate({ projectId: selectedProject.id, name: envName }, {
      onSuccess: (newEnv) => {
        setEnvModalOpen(false);
        setEnvName("");
        handleSelectEnvironment(newEnv);
      }
    });
  };

  const handleSaveSecret = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProject || !selectedEnvironment) return;

    if (editingSecret) {
      updateSecret.mutate({
        projectId: selectedProject.id,
        envId: selectedEnvironment.id,
        secretId: editingSecret.id,
        value: secretValue
      }, {
        onSuccess: () => {
          setSecretModalOpen(false);
          setSecretKey("");
          setSecretValue("");
          setEditingSecret(null);
        }
      });
    } else {
      createSecret.mutate({
        projectId: selectedProject.id,
        envId: selectedEnvironment.id,
        key: secretKey,
        value: secretValue
      }, {
        onSuccess: () => {
          setSecretModalOpen(false);
          setSecretKey("");
          setSecretValue("");
          setEditingSecret(null);
        }
      });
    }
  };

  const handleDeleteSecret = (secret: Secret) => {
    setSecretToDelete(secret);
  };

  const confirmDeleteSecret = async () => {
    if (!selectedProject || !selectedEnvironment || !secretToDelete) return;
    deleteSecret.mutate({
      projectId: selectedProject.id,
      envId: selectedEnvironment.id,
      secretId: secretToDelete.id
    }, {
      onSuccess: () => {
        setSecretToDelete(null);
      }
    });
  };

  const handleDeleteEnvironment = (env: Environment) => {
    setEnvToDelete(env);
  };

  const confirmDeleteEnvironment = async () => {
    if (!selectedProject || !envToDelete) return;
    deleteEnvironment.mutate({
      projectId: selectedProject.id,
      envId: envToDelete.id
    }, {
      onSuccess: () => {
        if (selectedEnvironment?.id === envToDelete.id) {
          setSelectedEnvironment(null);
        }
        setEnvToDelete(null);
      }
    });
  };

  const handleDeleteProject = async (projectId: string) => {
    if (!confirm("Are you sure you want to delete this project?")) return;
    deleteProject.mutate(projectId, {
      onSuccess: () => {
        if (selectedProject?.id === projectId) {
          setSelectedProject(null);
          setSelectedEnvironment(null);
        }
      }
    });
  };

  const toggleReveal = (secretId: string) => {
    if (revealedSecrets[secretId]) {
      const newRevealed = { ...revealedSecrets };
      delete newRevealed[secretId];
      setRevealedSecrets(newRevealed);
      return;
    }

    if (!selectedProject || !selectedEnvironment) return;

    revealSecretMutation.mutate({
      projectId: selectedProject.id,
      envId: selectedEnvironment.id,
      secretId
    }, {
      onSuccess: (data) => {
        setRevealedSecrets(prev => ({ ...prev, [secretId]: data.value }));
      }
    });
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const handleExportSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProject || !selectedEnvironment) return;

    reauthMutation.mutate(reauthPassword, {
      onSuccess: (data) => {
        exportSecretsMutation.mutate({
          projectId: selectedProject.id,
          envId: selectedEnvironment.id,
          reauthToken: data.reauth_token
        }, {
          onSuccess: (secretsRecord) => {
            const exportString = Object.entries(secretsRecord)
              .map(([k, v]) => `${k.toUpperCase()}=${v}`)
              .join('\n');
            copyToClipboard(exportString);
            setExportModalOpen(false);
            setReauthPassword("");
            alert("Exported secrets copied to clipboard!");
          }
        });
      }
    });
  };

  // If not authenticated, the useEffect will redirect to /login
  // We can render nothing or a loading state while redirecting
  if (typeof window !== 'undefined' && !localStorage.getItem("access_token")) {
    return null;
  }

  return (
    <div className="flex h-screen overflow-hidden bg-noir-900 text-vanilla-100 font-sans">
      <Sidebar
        projects={projects}
        environments={environments}
        selectedProject={selectedProject}
        selectedEnvironment={selectedEnvironment}
        onSelectProject={handleSelectProject}
        onSelectEnvironment={handleSelectEnvironment}
        onDeleteProject={handleDeleteProject}
        onDeleteEnvironment={handleDeleteEnvironment}
        setProjectModalOpen={setProjectModalOpen}
        setEnvModalOpen={setEnvModalOpen}
      />

      <div className="flex-1 flex flex-col bg-vanilla-100 text-noir-900 relative">
        <TopBar
          selectedProject={selectedProject}
          selectedEnvironment={selectedEnvironment}
          onLogout={handleLogout}
        />

        <div className="flex-1 p-8 overflow-y-auto">
          {!selectedProject ? (
            <EmptyState type="project" />
          ) : !selectedEnvironment ? (
            <EmptyState type="environment" />
          ) : (
            <SecretsTable
              selectedEnvironment={selectedEnvironment}
              secrets={secrets}
              revealedSecrets={revealedSecrets}
              onToggleReveal={toggleReveal}
              onCopyToClipboard={copyToClipboard}
              hasNextPage={hasNextPage}
              onFetchNextPage={fetchNextPage}
              isFetchingNextPage={isFetchingNextPage}
              setSecretModalOpen={setSecretModalOpen}
              setExportModalOpen={setExportModalOpen}
              setEditingSecret={setEditingSecret}
              setSecretKey={setSecretKey}
              setSecretValue={setSecretValue}
              onDeleteSecret={handleDeleteSecret}
            />
          )}
        </div>
      </div>

      <CreateProjectModal
        isOpen={isProjectModalOpen}
        onClose={() => setProjectModalOpen(false)}
        onSubmit={handleCreateProject}
        projectName={projectName}
        setProjectName={setProjectName}
        isError={createProject.isError}
        errorMessage={createProject.error?.message}
        resetMutation={createProject.reset}
      />

      <CreateEnvironmentModal
        isOpen={isEnvModalOpen}
        onClose={() => setEnvModalOpen(false)}
        onSubmit={handleCreateEnvironment}
        selectedProject={selectedProject}
        envName={envName}
        setEnvName={setEnvName}
        isError={createEnvironment.isError}
        errorMessage={createEnvironment.error?.message}
        resetMutation={createEnvironment.reset}
      />

      <SecretModal
        isOpen={isSecretModalOpen}
        onClose={() => setSecretModalOpen(false)}
        onSubmit={handleSaveSecret}
        selectedProject={selectedProject}
        selectedEnvironment={selectedEnvironment}
        editingSecret={editingSecret}
        secretKey={secretKey}
        setSecretKey={setSecretKey}
        secretValue={secretValue}
        setSecretValue={setSecretValue}
        isError={createSecret.isError || updateSecret.isError}
        errorMessage={createSecret.error?.message || updateSecret.error?.message}
        resetMutations={() => {
          createSecret.reset();
          updateSecret.reset();
        }}
      />

      <DeleteEnvironmentModal
        envToDelete={envToDelete}
        onClose={() => setEnvToDelete(null)}
        onConfirm={confirmDeleteEnvironment}
      />

      <DeleteSecretModal
        secretToDelete={secretToDelete}
        onClose={() => setSecretToDelete(null)}
        onConfirm={confirmDeleteSecret}
      />

      <ExportSecretsModal
        isOpen={isExportModalOpen}
        onClose={() => setExportModalOpen(false)}
        onSubmit={handleExportSubmit}
        reauthPassword={reauthPassword}
        setReauthPassword={setReauthPassword}
        showReauthPassword={showReauthPassword}
        setShowReauthPassword={setShowReauthPassword}
        isPending={exportSecretsMutation.isPending || reauthMutation.isPending}
      />
    </div>
  );
}
