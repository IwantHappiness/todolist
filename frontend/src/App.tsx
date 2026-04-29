import { useCallback, useEffect, useState, type FormEventHandler } from "react";
import "./App.css";
import { taskApi } from "./services/taskApi";
// NOTE: API calls use relative paths (e.g. "/tasks" and "/metrics").
// In development Vite proxies these endpoints to the backend (see vite.config.ts).
// This avoids cross-origin preflight and removes the need for VITE_API_BASE.
import type { Task } from "./services/models";

type TaskFilter = "all" | "completed" | "not_completed";

type NewTaskState = {
  title: string;
  description: string;
};

function formatDateTime(value: string | null): string {
  if (!value) {
    return "—";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString("ru-RU");
}

function App() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [taskFilter, setTaskFilter] = useState<TaskFilter>("all");
  const [isTasksLoading, setIsTasksLoading] = useState(false);
  const [tasksError, setTasksError] = useState<string | null>(null);

  const [isCreatingTask, setIsCreatingTask] = useState(false);
  const [newTask, setNewTask] = useState<NewTaskState>({
    title: "",
    description: "",
  });

  const [taskIdToFind, setTaskIdToFind] = useState("");
  const [foundTask, setFoundTask] = useState<Task | null>(null);
  const [isTaskSearchLoading, setIsTaskSearchLoading] = useState(false);
  const [taskSearchError, setTaskSearchError] = useState<string | null>(null);

  const [taskActionId, setTaskActionId] = useState<number | null>(null);
  const [isTaskActionLoading, setIsTaskActionLoading] = useState(false);

  const loadTasks = useCallback(async () => {
    setIsTasksLoading(true);
    setTasksError(null);

    try {
      let completedFilter: boolean | undefined = undefined;
      if (taskFilter === "completed") {
        completedFilter = true;
      } else if (taskFilter === "not_completed") {
        completedFilter = false;
      }

      const loadedTasks = await taskApi.getAllTasks(completedFilter);
      setTasks(loadedTasks);
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message);
      } else {
        setTasksError("Failed to load tasks");
      }
    } finally {
      setIsTasksLoading(false);
    }
  }, [taskFilter]);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      void loadTasks();
    }, 0);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [loadTasks]);

  const handleCreateTask: FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault();

    const title = newTask.title.trim();
    const description = newTask.description.trim();

    if (title.length === 0) {
      setTasksError("Название задачи обязательно");
      return;
    }

    setIsCreatingTask(true);
    setTasksError(null);

    try {
      await taskApi.createTask({
        title,
        description,
      });

      setNewTask({
        title: "",
        description: "",
      });
      await loadTasks();
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message);
      } else {
        setTasksError("Failed to create task");
      }
    } finally {
      setIsCreatingTask(false);
    }
  };

  const handleFindTask: FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault();

    const trimmedId = taskIdToFind.trim();
    if (trimmedId.length === 0) {
      setTaskSearchError("Введи ID задачи");
      return;
    }

    const taskId = parseInt(trimmedId, 10);
    if (Number.isNaN(taskId)) {
      setTaskSearchError("ID должен быть числом");
      return;
    }

    setIsTaskSearchLoading(true);
    setTaskSearchError(null);

    try {
      const task = await taskApi.getTaskById(taskId);
      setFoundTask(task);
    } catch (error) {
      setFoundTask(null);
      if (error instanceof Error) {
        setTaskSearchError(error.message);
      } else {
        setTaskSearchError("Failed to find task");
      }
    } finally {
      setIsTaskSearchLoading(false);
    }
  };

  const handleToggleTaskCompletion = async (task: Task) => {
    setTasksError(null);
    setTaskActionId(task.id);
    setIsTaskActionLoading(true);

    try {
      const updatedTask = await taskApi.completeTask(task.id, !task.completed);

      if (foundTask?.id === task.id) {
        setFoundTask(updatedTask);
      }
      await loadTasks();
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message);
      } else {
        setTasksError("Failed to update task status");
      }
    } finally {
      setTaskActionId(null);
      setIsTaskActionLoading(false);
    }
  };

  const handleDeleteTask = async (taskId: number) => {
    setTasksError(null);
    setTaskActionId(taskId);
    setIsTaskActionLoading(true);

    try {
      await taskApi.deleteTask(taskId);
      if (foundTask?.id === taskId) {
        setFoundTask(null);
      }
      await loadTasks();
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message);
      } else {
        setTasksError("Failed to delete task");
      }
    } finally {
      setTaskActionId(null);
      setIsTaskActionLoading(false);
    }
  };

  const handleDeleteAllTasks = async () => {
    setTasksError(null);
    try {
      await taskApi.deleteAllTasks();
      setFoundTask(null);
      await loadTasks();
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message);
      } else {
        setTasksError("Failed to delete all tasks");
      }
    }
  };

  const isCurrentTaskActionLoading = (taskId: number): boolean =>
    isTaskActionLoading && taskActionId === taskId;

  return (
    <>
      <main className="app">
        <header className="panel">
          <div className="panel-head">
            <h1>TodoList</h1>
          </div>
        </header>
        <section className="panel">
          <h1>Создать задачу</h1>
          <form className="task-form" onSubmit={handleCreateTask}>
            <label>
              Название
              <input
                type="text"
                value={newTask.title}
                onChange={(event) =>
                  setNewTask((current) => ({
                    ...current,
                    title: event.target.value,
                  }))
                }
                required
              />
            </label>
            <label>
              Описание
              <textarea
                value={newTask.description}
                onChange={(event) =>
                  setNewTask((current) => ({
                    ...current,
                    description: event.target.value,
                  }))
                }
                rows={3}
              />
            </label>
            <button type="submit" disabled={isCreatingTask}>
              {isCreatingTask ? "Создание..." : "Создать задачу"}
            </button>
          </form>
        </section>

        <section className="panel">
          <h2>Найти задачу по ID</h2>
          <form className="inline-form" onSubmit={handleFindTask}>
            <input
              type="number"
              min={0}
              placeholder="Например, 1"
              value={taskIdToFind}
              onChange={(event) => setTaskIdToFind(event.target.value)}
            />
            <button type="submit" disabled={isTaskSearchLoading}>
              {isTaskSearchLoading ? "Поиск..." : "Найти"}
            </button>
          </form>
          {taskSearchError ? <p className="error">{taskSearchError}</p> : null}
          {foundTask ? (
            <article
              className={`task-card search-task-result ${foundTask.completed ? "done" : ""
                }`}
            >
              <h3>
                #{foundTask.id} {foundTask.title || "Без названия"}
              </h3>
              <p>{foundTask.description || "Без описания"}</p>
              <p>
                <strong>Создана:</strong> {formatDateTime(foundTask.created_at)}
              </p>
              <p>
                <strong>Завершена:</strong>{" "}
                {formatDateTime(foundTask.completed_at)}
              </p>
              <p>
                <strong>Статус:</strong>{" "}
                {foundTask.completed ? "Выполнена" : "Не выполнена"}
              </p>
              <div className="actions">
                <button
                  type="button"
                  className="secondary"
                  disabled={isCurrentTaskActionLoading(foundTask.id)}
                  onClick={() => void handleToggleTaskCompletion(foundTask)}
                >
                  {isCurrentTaskActionLoading(foundTask.id)
                    ? "Сохранение..."
                    : foundTask.completed
                      ? "Отметить невыполненной"
                      : "Отметить выполненной"}
                </button>
                <button
                  type="button"
                  className="danger"
                  disabled={isCurrentTaskActionLoading(foundTask.id)}
                  onClick={() => void handleDeleteTask(foundTask.id)}
                >
                  {isCurrentTaskActionLoading(foundTask.id)
                    ? "Удаление..."
                    : "Удалить задачу"}
                </button>
              </div>
            </article>
          ) : null}
        </section>

        <section className="panel">
          <div className="panel-head">
            <h2>Все задачи</h2>
            <div className="actions">
              <select
                value={taskFilter}
                onChange={(event) =>
                  setTaskFilter(event.target.value as TaskFilter)
                }
              >
                <option value="all">Все</option>
                <option value="completed">Только выполненные</option>
                <option value="not_completed">Только невыполненные</option>
              </select>
              <button
                type="button"
                className="secondary"
                onClick={() => void loadTasks()}
              >
                Обновить
              </button>
              <button
                type="button"
                className="danger"
                onClick={() => void handleDeleteAllTasks()}
              >
                Удалить все
              </button>
            </div>
          </div>
          {tasksError ? <p className="error">{tasksError}</p> : null}
          {isTasksLoading ? <p>Загрузка задач...</p> : null}
          {!isTasksLoading && tasks.length === 0 ? (
            <p>Задач пока нет.</p>
          ) : null}
          {!isTasksLoading && tasks.length > 0 ? (
            <ul className="task-list">
              {tasks.map((task) => (
                <li
                  key={task.id}
                  className={`task-card ${task.completed ? "done" : ""}`}
                >
                  <h3>
                    #{task.id} {task.title || "Без названия"}
                  </h3>
                  <p>{task.description || "Без описания"}</p>
                  <p>
                    <strong>Создана:</strong> {formatDateTime(task.created_at)}
                  </p>
                  <p>
                    <strong>Завершена:</strong>{" "}
                    {formatDateTime(task.completed_at)}
                  </p>
                  <p>
                    <strong>Статус:</strong>{" "}
                    {task.completed ? "Выполнена" : "Не выполнена"}
                  </p>
                  <div className="actions">
                    <button
                      type="button"
                      className="secondary"
                      disabled={isCurrentTaskActionLoading(task.id)}
                      onClick={() => void handleToggleTaskCompletion(task)}
                    >
                      {isCurrentTaskActionLoading(task.id)
                        ? "Сохранение..."
                        : task.completed
                          ? "Отметить невыполненной"
                          : "Отметить выполненной"}
                    </button>
                    <button
                      type="button"
                      className="danger"
                      disabled={isCurrentTaskActionLoading(task.id)}
                      onClick={() => void handleDeleteTask(task.id)}
                    >
                      {isCurrentTaskActionLoading(task.id)
                        ? "Удаление..."
                        : "Удалить"}
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          ) : null}
        </section>
      </main >
    </>
  );
}

export default App;
