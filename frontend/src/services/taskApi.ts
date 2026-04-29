import type {
  Task,
  TasksResponse,
  CreateTaskRequest,
  CompleteTaskRequest,
} from "./models";

// Use relative endpoints (dev proxy or same-origin). API_BASE removed so calls use the provided path directly.

async function apiRequest<T>(
  path: string,
  init?: RequestInit,
): Promise<T | null> {
  const opts: RequestInit = {
    ...init,
    headers: {
      ...(init?.headers ?? {}),
    },
  };

  // only set content-type when there is a body
  if (
    opts.body &&
    !(opts.headers && (opts.headers as Record<string, string>)["Content-Type"])
  ) {
    opts.headers = {
      ...(opts.headers as Record<string, string>),
      "Content-Type": "application/json",
    };
  }

  const res = await fetch(path, opts);

  // No Content
  if (res.status === 204) {
    return null;
  }

  const text = await res.text();
  const payload = text.length > 0 ? (JSON.parse(text) as unknown) : null;

  if (!res.ok) {
    throw new Error(getErrorMessage(payload, res.status));
  }

  return payload as T | null;
}

function getErrorMessage(payload: unknown, status: number): string {
  if (isObject(payload)) {
    const message = payload.message;
    if (typeof message === "string") {
      return message;
    }

    const error = payload.error;
    if (typeof error === "string") {
      return error;
    }
  }
  return `Request failed with status ${status}`;
}

function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function normalizeTasks(payload: TasksResponse | null): Task[] {
  if (!payload) return [];
  if (Array.isArray(payload)) return payload;
  return Object.values(payload);
}

export const taskApi = {
  // GET /tasks or GET /tasks?completed=true|false
  async getAllTasks(completed?: boolean): Promise<Task[]> {
    const query =
      completed === undefined ? "" : `?completed=${String(completed)}`;
    const resp = await apiRequest<TasksResponse>(`/tasks${query}`);
    return normalizeTasks(resp).sort((a, b) => a.id - b.id);
  },

  // GET /tasks/:id
  async getTaskById(id: number): Promise<Task> {
    const resp = await apiRequest<Task>(`/tasks/${id}`);
    if (!resp) throw new Error("Task not found");
    return resp;
  },

  // POST /tasks
  async createTask(payload: CreateTaskRequest): Promise<Task> {
    const resp = await apiRequest<Task>("/tasks", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    if (!resp) throw new Error("Failed to create task");
    return resp;
  },

  // PUT /tasks/:id  (replace/update full fields)
  async updateTask(id: number, payload: CreateTaskRequest): Promise<Task> {
    const resp = await apiRequest<Task>(`/tasks/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    });
    if (!resp) throw new Error("Failed to update task");
    return resp;
  },

  // PATCH /tasks/:id  (update completion status)
  async completeTask(id: number, completed: boolean): Promise<Task> {
    const body: CompleteTaskRequest = { completed };
    const resp = await apiRequest<Task>(`/tasks/${id}`, {
      method: "PATCH",
      body: JSON.stringify(body),
    });
    if (!resp) throw new Error("Failed to update task status");
    return resp;
  },

  // DELETE /tasks/:id
  async deleteTask(id: number): Promise<void> {
    await apiRequest<null>(`/tasks/${id}`, {
      method: "DELETE",
    });
  },

  // DELETE /tasks
  async deleteAllTasks(): Promise<void> {
    await apiRequest<null>("/tasks", {
      method: "DELETE",
    });
  },
};
