import { useState, useEffect } from 'react';
import { deleteTask, Task } from '../services/api';
import { PlusIcon } from '@heroicons/react/24/solid';
import Modal from './Modal';
import TaskForm from './TaskForm';
import TaskRunner from './TaskRunner';

interface TaskListProps {
  tasks: Task[];
  onTaskCreated: (task: Task) => void;
}

const TaskList: React.FC<TaskListProps> = ({ tasks, onTaskCreated }) => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleDelete = async (id: number) => {
    try {
      await deleteTask(id);
    } catch (error) {
      console.error('Failed to delete task', error);
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
          <TaskForm onTaskCreated={onTaskCreated} />
        </Modal>
      )}
    </div>
  );
};

export default TaskList;

