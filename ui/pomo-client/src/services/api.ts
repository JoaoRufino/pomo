export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export type Task = {
  id: number;
  message: string;
  n_pomodoros: number;
  tags: string[];
};

export type TaskList = {
  count: number;
  results: Task[];
};

const headers = {
  'Content-Type': 'application/json; charset=utf-8',
  'Accept': 'application/json; charset=utf-8',
  'Authorization': `Bearer ${process.env.NEXT_PUBLIC_API_TOKEN || 'default-token'}`,
};

const handleResponse = async <T>(response: Response): Promise<T> => {
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'API Error');
  }
  return response.json();
};

export const fetchTasks = async (): Promise<TaskList> => {
  const response = await fetch(`${API_URL}/tasks`, { headers });
  return handleResponse<TaskList>(response);
};

export const createTask = async (task: Omit<Task, 'id'>): Promise<Task> => {
  const response = await fetch(`${API_URL}/tasks`, {
    method: 'POST',
    headers,
    body: JSON.stringify(task),
  });
  return handleResponse<Task>(response);
};

export const deleteTask = async (id: number): Promise<void> => {
  const response = await fetch(`${API_URL}/tasks/${id}`, {
    method: 'DELETE',
    headers,
  });
  return handleResponse<void>(response);
};

