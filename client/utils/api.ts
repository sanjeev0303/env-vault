export const API_BASE = "http://localhost:8080/api";

export interface Project {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Environment {
  id: string;
  project_id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Secret {
  id: string;
  project_id: string;
  environment_id: string;
  key: string;
  value: string;
  created_at: string;
  updated_at: string;
}

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (!res.ok) {
    let errorMessage = "An error occurred";
    try {
      const errorData = await res.json();
      errorMessage = errorData.error || errorMessage;
    } catch (e) {
      errorMessage = res.statusText;
    }
    throw new Error(errorMessage);
  }

  // Handle 204 No Content
  if (res.status === 204) {
    return {} as T;
  }

  return res.json();
}

// Project API
export const api = {
  getProjects: () => fetchAPI<Project[]>("/projects"),
  getProject: (id: string) => fetchAPI<Project>(`/projects/${id}`),
  createProject: (name: string) => fetchAPI<Project>("/projects", {
    method: "POST",
    body: JSON.stringify({ name }),
  }),
  deleteProject: (id: string) => fetchAPI<void>(`/projects/${id}`, {
    method: "DELETE",
  }),

  // Environment API
  getEnvironments: (projectId: string) => fetchAPI<Environment[]>(`/projects/${projectId}/environments`),
  createEnvironment: (projectId: string, name: string) => fetchAPI<Environment>(`/projects/${projectId}/environments`, {
    method: "POST",
    body: JSON.stringify({ name }),
  }),
  deleteEnvironment: (projectId: string, envId: string) => fetchAPI<void>(`/projects/${projectId}/environments/${envId}`, {
    method: "DELETE",
  }),

  // Secret API
  getSecrets: (projectId: string, envId: string) => fetchAPI<Secret[]>(`/projects/${projectId}/environments/${envId}/secrets`),
  createSecret: (projectId: string, envId: string, key: string, value: string) => fetchAPI<Secret>(`/projects/${projectId}/environments/${envId}/secrets`, {
    method: "POST",
    body: JSON.stringify({ key, value }),
  }),
  updateSecret: (projectId: string, envId: string, secretId: string, value: string) => fetchAPI<Secret>(`/projects/${projectId}/environments/${envId}/secrets/${secretId}`, {
    method: "PUT",
    body: JSON.stringify({ value }),
  }),
  deleteSecret: (projectId: string, envId: string, secretId: string) => fetchAPI<void>(`/projects/${projectId}/environments/${envId}/secrets/${secretId}`, {
    method: "DELETE",
  }),
};
