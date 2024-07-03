import { useState, useEffect } from 'react';
import { fetchTasks, deleteTask, createTask, Task, TaskList } from '../services/api';
import { PlusIcon } from '@heroicons/react/24/solid';
import Modal from './Modal';
import TaskForm from './TaskForm';
import TaskRunner from './TaskRunner';

const TaskList = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    try {
      const data: TaskList = await fetchTasks();
      setTasks(data.results);
    } catch (error) {
      console.error('Failed to load tasks', error);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteTask(id);
      loadTasks();
    } catch (error) {
      console.error('Failed to delete task', error);
    }
  };

  const handleCreate = async (task: Omit<Task, 'id'>) => {
    try {
      await createTask(task);
      loadTasks();
      setIsModalOpen(false);
    } catch (error) {
      console.error('Failed to create task', error);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Tasks</h1>
        <button
          className="bg-blue-500 text-white p-2 rounded-full shadow-md"
          onClick={() => setIsModalOpen(true)}
        >
          <PlusIcon className="h-6 w-6" />
        </button>
      </div>
      <ul className="space-y-4">
        {tasks.map(task => (
          <li key={task.id} className="bg-white shadow-md rounded p-4">
            <TaskRunner task={task} />
            <button
              className="bg-red-500 text-white p-2 rounded shadow-md mt-2"
              onClick={() => handleDelete(task.id)}
            >
              Delete
            </button>
          </li>
        ))}
      </ul>
      {isModalOpen && (
        <Modal onClose={() => setIsModalOpen(false)}>
          <TaskForm onSubmit={handleCreate} />
        </Modal>
      )}
    </div>
  );
};

export default TaskList;

