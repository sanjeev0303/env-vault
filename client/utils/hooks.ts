import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from "@tanstack/react-query";
import { api } from "./api";

// Projects
export const useProjects = () => {
  return useQuery({
    queryKey: ["projects"],
    queryFn: api.getProjects,
  });
};

export const useCreateProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (name: string) => api.createProject(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
    },
  });
};

export const useDeleteProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
    },
  });
};

// Environments
export const useEnvironments = (projectId: string | undefined) => {
  return useQuery({
    queryKey: ["environments", projectId],
    queryFn: () => api.getEnvironments(projectId!),
    enabled: !!projectId,
  });
};

export const useCreateEnvironment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, name }: { projectId: string; name: string }) => 
      api.createEnvironment(projectId, name),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["environments", variables.projectId] });
    },
  });
};

export const useDeleteEnvironment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, envId }: { projectId: string; envId: string }) => 
      api.deleteEnvironment(projectId, envId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["environments", variables.projectId] });
    },
  });
};

// Secrets
export const useSecrets = (projectId: string | undefined, envId: string | undefined) => {
  return useInfiniteQuery({
    queryKey: ["secrets", projectId, envId],
    queryFn: ({ pageParam = "" }) => api.getSecrets(projectId!, envId!, pageParam),
    getNextPageParam: (lastPage) => lastPage.next_cursor || undefined,
    initialPageParam: "",
    enabled: !!projectId && !!envId,
  });
};

export const useCreateSecret = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, envId, key, value }: { projectId: string; envId: string; key: string; value: string }) =>
      api.createSecret(projectId, envId, key, value),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["secrets", variables.projectId, variables.envId] });
    },
  });
};

export const useUpdateSecret = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, envId, secretId, value }: { projectId: string; envId: string; secretId: string; value: string }) =>
      api.updateSecret(projectId, envId, secretId, value),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["secrets", variables.projectId, variables.envId] });
    },
  });
};

export const useDeleteSecret = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, envId, secretId }: { projectId: string; envId: string; secretId: string }) =>
      api.deleteSecret(projectId, envId, secretId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["secrets", variables.projectId, variables.envId] });
    },
  });
};

export const useRevealSecret = () => {
  return useMutation({
    mutationFn: ({ projectId, envId, secretId }: { projectId: string; envId: string; secretId: string }) =>
      api.revealSecret(projectId, envId, secretId)
  });
};

export const useExportSecrets = () => {
  return useMutation({
    mutationFn: ({ projectId, envId, reauthToken }: { projectId: string; envId: string; reauthToken: string }) =>
      api.exportSecrets(projectId, envId, reauthToken)
  });
};

// Auth Hooks
import { authApi } from "./api";

export const useLogin = () => {
  return useMutation({
    mutationFn: (credentials: Record<string, string>) => authApi.login(credentials),
    onSuccess: (data) => {
      localStorage.setItem("access_token", data.access_token);
    }
  });
};

export const useRegister = () => {
  return useMutation({
    mutationFn: (data: Record<string, string>) => authApi.register(data)
  });
};

export const useReauth = () => {
  return useMutation({
    mutationFn: (password: string) => authApi.reauth(password)
  });
};
