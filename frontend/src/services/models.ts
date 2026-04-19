export type Task = {
  id: number;
  title: string;
  description: string;
  // ISO 8601 string from backend: created_at
  created_at: string;
  // ISO 8601 string or null from backend: completed_at
  completed_at: string | null;
  completed: boolean;
};

export type TasksResponse = Record<number, Task> | Task[];

export type ErrorDTO = {
  message: string;
  time: string;
};

export type CreateTaskRequest = {
  title: string;
  description: string;
};

export type CompleteTaskRequest = {
  completed: boolean;
};

export type MetricsResponse = {
  memory_usage_mb: number;
  uptime_seconds: number;
  goroutines: number;
  cpu_cores: number;
};
