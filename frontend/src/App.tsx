import { useCallback, useEffect, useState, type FormEventHandler } from 'react'
import './App.css'

const API_BASE = '/api'

type Task = {
  id: number
  title: string
  created: string
  completed_at: string | null
  completed: boolean
  description: string
}

type TasksResponse = Record<string, Task> | Task[]

type MetricsResponse = {
  memory_usage_mb: number
  uptime_seconds: number
  goroutines: number
  cpu_cores: number
}

type TaskFilter = 'all' | 'completed' | 'not_completed'

type NewTaskState = {
  title: string
  description: string
}

function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

function getErrorMessage(payload: unknown, status: number): string {
  if (isObject(payload) && typeof payload.message === 'string') {
    return payload.message
  }
  return `Request failed with status ${status}`
}

function formatDateTime(value: string | null): string {
  if (!value) {
    return '—'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toLocaleString('ru-RU')
}

async function apiRequest<T>(path: string, init?: RequestInit): Promise<T | null> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
  })

  const text = await response.text()
  const payload = text.length > 0 ? (JSON.parse(text) as unknown) : null

  if (!response.ok) {
    throw new Error(getErrorMessage(payload, response.status))
  }

  return payload as T | null
}

function normalizeTasks(payload: TasksResponse | null): Task[] {
  if (!payload) {
    return []
  }

  if (Array.isArray(payload)) {
    return payload
  }

  return Object.values(payload)
}

function App() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [taskFilter, setTaskFilter] = useState<TaskFilter>('all')
  const [isTasksLoading, setIsTasksLoading] = useState(false)
  const [tasksError, setTasksError] = useState<string | null>(null)

  const [isCreatingTask, setIsCreatingTask] = useState(false)
  const [newTask, setNewTask] = useState<NewTaskState>({
    title: '',
    description: '',
  })

  const [taskIdToFind, setTaskIdToFind] = useState('')
  const [foundTask, setFoundTask] = useState<Task | null>(null)
  const [isTaskSearchLoading, setIsTaskSearchLoading] = useState(false)
  const [taskSearchError, setTaskSearchError] = useState<string | null>(null)

  const [taskActionId, setTaskActionId] = useState<number | null>(null)
  const [isTaskActionLoading, setIsTaskActionLoading] = useState(false)

  const [metrics, setMetrics] = useState<MetricsResponse | null>(null)
  const [isMetricsLoading, setIsMetricsLoading] = useState(false)
  const [metricsError, setMetricsError] = useState<string | null>(null)

  const loadTasks = useCallback(async () => {
    setIsTasksLoading(true)
    setTasksError(null)

    try {
      const query =
        taskFilter === 'all'
          ? ''
          : `?completed=${String(taskFilter === 'completed')}`
      const payload = await apiRequest<TasksResponse>(`/task${query}`)
      const normalizedTasks = normalizeTasks(payload).sort((a, b) => a.id - b.id)
      setTasks(normalizedTasks)
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message)
      } else {
        setTasksError('Failed to load tasks')
      }
    } finally {
      setIsTasksLoading(false)
    }
  }, [taskFilter])

  const loadMetrics = useCallback(async () => {
    setIsMetricsLoading(true)
    setMetricsError(null)

    try {
      const payload = await apiRequest<MetricsResponse>('/metrics')
      setMetrics(payload)
    } catch (error) {
      if (error instanceof Error) {
        setMetricsError(error.message)
      } else {
        setMetricsError('Failed to load metrics')
      }
    } finally {
      setIsMetricsLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadTasks()
  }, [loadTasks])

  useEffect(() => {
    void loadMetrics()
  }, [loadMetrics])

  const handleCreateTask: FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault()

    const title = newTask.title.trim()
    const description = newTask.description.trim()

    if (title.length === 0) {
      setTasksError('Название задачи обязательно')
      return
    }

    setIsCreatingTask(true)
    setTasksError(null)

    try {
      await apiRequest<Task>('/task', {
        method: 'POST',
        body: JSON.stringify({
          title,
          description,
        }),
      })

      setNewTask({
        title: '',
        description: '',
      })
      await loadTasks()
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message)
      } else {
        setTasksError('Failed to create task')
      }
    } finally {
      setIsCreatingTask(false)
    }
  }

  const handleFindTask: FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault()

    const trimmedId = taskIdToFind.trim()
    if (trimmedId.length === 0) {
      setTaskSearchError('Введи ID задачи')
      return
    }

    setIsTaskSearchLoading(true)
    setTaskSearchError(null)

    try {
      const payload = await apiRequest<Task>(`/task/${trimmedId}`)
      setFoundTask(payload)
    } catch (error) {
      setFoundTask(null)
      if (error instanceof Error) {
        setTaskSearchError(error.message)
      } else {
        setTaskSearchError('Failed to find task')
      }
    } finally {
      setIsTaskSearchLoading(false)
    }
  }

  const handleToggleTaskCompletion = async (task: Task) => {
    setTasksError(null)
    setTaskActionId(task.id)
    setIsTaskActionLoading(true)

    try {
      const updatedTask = await apiRequest<Task>(`/task/${task.id}`, {
        method: 'PATCH',
        body: JSON.stringify({ completed: !task.completed }),
      })

      if (foundTask?.id === task.id) {
        setFoundTask(updatedTask)
      }
      await loadTasks()
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message)
      } else {
        setTasksError('Failed to update task status')
      }
    } finally {
      setTaskActionId(null)
      setIsTaskActionLoading(false)
    }
  }

  const handleDeleteTask = async (taskId: number) => {
    setTasksError(null)
    setTaskActionId(taskId)
    setIsTaskActionLoading(true)

    try {
      await apiRequest<null>(`/task/${taskId}`, { method: 'DELETE' })
      if (foundTask?.id === taskId) {
        setFoundTask(null)
      }
      await loadTasks()
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message)
      } else {
        setTasksError('Failed to delete task')
      }
    } finally {
      setTaskActionId(null)
      setIsTaskActionLoading(false)
    }
  }

  const handleDeleteAllTasks = async () => {
    setTasksError(null)
    try {
      await apiRequest<null>('/task', { method: 'DELETE' })
      setFoundTask(null)
      await loadTasks()
    } catch (error) {
      if (error instanceof Error) {
        setTasksError(error.message)
      } else {
        setTasksError('Failed to delete all tasks')
      }
    }
  }

  const isCurrentTaskActionLoading = (taskId: number): boolean =>
    isTaskActionLoading && taskActionId === taskId

  return (
    <main className="app">
      <header className="panel">
        <div className="panel-head">
          <h1>TodoList</h1>
          <button
            type="button"
            className="secondary"
            onClick={() => void loadMetrics()}
            disabled={isMetricsLoading}
          >
            {isMetricsLoading ? 'Обновление метрик...' : 'Обновить метрики'}
          </button>
        </div>
        {metricsError ? <p className="error">{metricsError}</p> : null}
        <div className="metrics-grid">
          <article className="metric">
            <span>Память</span>
            <strong>{metrics?.memory_usage_mb ?? '—'} MB</strong>
          </article>
          <article className="metric">
            <span>Uptime</span>
            <strong>{metrics?.uptime_seconds ?? '—'} s</strong>
          </article>
          <article className="metric">
            <span>Goroutines</span>
            <strong>{metrics?.goroutines ?? '—'}</strong>
          </article>
          <article className="metric">
            <span>CPU Cores</span>
            <strong>{metrics?.cpu_cores ?? '—'}</strong>
          </article>
        </div>
      </header>

      <section className="panel">
        <h2>Создать задачу</h2>
        <form className="task-form" onSubmit={handleCreateTask}>
          <label>
            Название
            <input
              type="text"
              value={newTask.title}
              onChange={(event) =>
                setNewTask((current) => ({ ...current, title: event.target.value }))
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
            {isCreatingTask ? 'Создание...' : 'Создать задачу'}
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
            {isTaskSearchLoading ? 'Поиск...' : 'Найти'}
          </button>
        </form>
        {taskSearchError ? <p className="error">{taskSearchError}</p> : null}
        {foundTask ? (
          <article className={`task-card search-task-result ${foundTask.completed ? 'done' : ''}`}>
            <h3>
              #{foundTask.id} {foundTask.title || 'Без названия'}
            </h3>
            <p>{foundTask.description || 'Без описания'}</p>
            <p>
              <strong>Создана:</strong> {formatDateTime(foundTask.created)}
            </p>
            <p>
              <strong>Завершена:</strong> {formatDateTime(foundTask.completed_at)}
            </p>
            <p>
              <strong>Статус:</strong>{' '}
              {foundTask.completed ? 'Выполнена' : 'Не выполнена'}
            </p>
            <div className="actions">
              <button
                type="button"
                className="secondary"
                disabled={isCurrentTaskActionLoading(foundTask.id)}
                onClick={() => void handleToggleTaskCompletion(foundTask)}
              >
                {isCurrentTaskActionLoading(foundTask.id)
                  ? 'Сохранение...'
                  : foundTask.completed
                    ? 'Отметить невыполненной'
                    : 'Отметить выполненной'}
              </button>
              <button
                type="button"
                className="danger"
                disabled={isCurrentTaskActionLoading(foundTask.id)}
                onClick={() => void handleDeleteTask(foundTask.id)}
              >
                {isCurrentTaskActionLoading(foundTask.id)
                  ? 'Удаление...'
                  : 'Удалить задачу'}
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
              onChange={(event) => setTaskFilter(event.target.value as TaskFilter)}
            >
              <option value="all">Все</option>
              <option value="completed">Только выполненные</option>
              <option value="not_completed">Только невыполненные</option>
            </select>
            <button type="button" className="secondary" onClick={() => void loadTasks()}>
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
        {!isTasksLoading && tasks.length === 0 ? <p>Задач пока нет.</p> : null}
        {!isTasksLoading && tasks.length > 0 ? (
          <ul className="task-list">
            {tasks.map((task) => (
              <li key={task.id} className={`task-card ${task.completed ? 'done' : ''}`}>
                <h3>
                  #{task.id} {task.title || 'Без названия'}
                </h3>
                <p>{task.description || 'Без описания'}</p>
                <p>
                  <strong>Создана:</strong> {formatDateTime(task.created)}
                </p>
                <p>
                  <strong>Завершена:</strong> {formatDateTime(task.completed_at)}
                </p>
                <p>
                  <strong>Статус:</strong> {task.completed ? 'Выполнена' : 'Не выполнена'}
                </p>
                <div className="actions">
                  <button
                    type="button"
                    className="secondary"
                    disabled={isCurrentTaskActionLoading(task.id)}
                    onClick={() => void handleToggleTaskCompletion(task)}
                  >
                    {isCurrentTaskActionLoading(task.id)
                      ? 'Сохранение...'
                      : task.completed
                        ? 'Отметить невыполненной'
                        : 'Отметить выполненной'}
                  </button>
                  <button
                    type="button"
                    className="danger"
                    disabled={isCurrentTaskActionLoading(task.id)}
                    onClick={() => void handleDeleteTask(task.id)}
                  >
                    {isCurrentTaskActionLoading(task.id) ? 'Удаление...' : 'Удалить'}
                  </button>
                </div>
              </li>
            ))}
          </ul>
        ) : null}
      </section>
    </main>
  )
}

export default App
