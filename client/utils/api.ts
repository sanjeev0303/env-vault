export const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

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

export interface User {
  id: string;
  email: string;
  name: string;
  is_superadmin: boolean;
}

export interface AuthResp {
  access_token: string;
  expires_in: number;
  user: User;
}

let csrfToken: string | null = null;

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  // Attach token from localStorage
  const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;
  const method = options?.method || "GET";

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...options?.headers as Record<string, string>,
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  if (["POST", "PUT", "DELETE", "PATCH"].includes(method.toUpperCase())) {
    if (!csrfToken) {
      try {
        const csrfRes = await fetch(`${API_BASE}/auth/csrf`, { credentials: "include" });
        if (csrfRes.ok) {
          const data = await csrfRes.json();
          csrfToken = data.csrf_token;
        }
      } catch (e) {
        console.error("Failed to fetch CSRF token", e);
      }
    }
    if (csrfToken) {
      headers["X-CSRF-Token"] = csrfToken;
    }
  }

  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
    credentials: "include",
  });

  if (!res.ok) {
    let errorMessage = "An error occurred";
    try {
      const errorData = await res.json();
      errorMessage = errorData.error || errorMessage;
    } catch (_) {
      errorMessage = res.statusText;
    }

    // Clear token on 401
    if (res.status === 401 && typeof window !== 'undefined') {
      localStorage.removeItem('access_token');
      // Trigger a custom event to tell the app we're logged out
      window.dispatchEvent(new Event('auth-error'));
    }

    throw new Error(errorMessage);
  }

  // Handle 204 No Content
  if (res.status === 204) {
    return {} as T;
  }

  return res.json();
}

// Auth API
export const authApi = {
  login: (credentials: Record<string, string>) => fetchAPI<AuthResp>("/auth/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  }),
  register: (data: Record<string, string>) => fetchAPI<User>("/auth/register", {
    method: "POST",
    body: JSON.stringify(data),
  }),
  me: () => fetchAPI<User>("/auth/me"),
  reauth: (password: string) => fetchAPI<{ reauth_token: string }>("/auth/reauth", {
    method: "POST",
    body: JSON.stringify({ password }),
  }),
};

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
  getSecrets: (projectId: string, envId: string, cursor?: string, limit?: number) => {
    let url = `/projects/${projectId}/environments/${envId}/secrets`;
    const params = new URLSearchParams();
    if (cursor) params.append("cursor", cursor);
    if (limit) params.append("limit", limit.toString());
    if (params.toString()) url += `?${params.toString()}`;
    return fetchAPI<{ data: Secret[]; next_cursor: string }>(url);
  },
  revealSecret: (projectId: string, envId: string, secretId: string) =>
    fetchAPI<Secret>(`/projects/${projectId}/environments/${envId}/secrets/${secretId}/reveal`),
  exportSecrets: (projectId: string, envId: string, reauthToken: string) =>
    fetchAPI<Record<string, string>>(`/projects/${projectId}/environments/${envId}/secrets/export`, {
      headers: { "X-Reauth-Token": reauthToken },
    }),
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
