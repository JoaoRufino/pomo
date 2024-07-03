import React, { useState } from 'react';
import { createTask, Task } from '../services/api';

const TaskForm: React.FC<{ onTaskCreated: (task: Task) => void }> = ({ onTaskCreated }) => {
  const [message, setMessage] = useState('');
  const [nPomodoros, setNPomodoros] = useState(0);
  const [tags, setTags] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const newTask = await createTask({ message, n_pomodoros: nPomodoros, tags });
      onTaskCreated(newTask);
      setMessage('');
      setNPomodoros(0);
      setTags([]);
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mb-8">
      <h2 className="text-2xl font-bold mb-4">Create Task</h2>
      {error && <div className="text-red-500 mb-4">{error}</div>}
      <div className="mb-4">
        <label className="block mb-2">Message</label>
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          className="border p-2 w-full"
          required
        />
      </div>
      <div className="mb-4">
        <label className="block mb-2">Number of Pomodoros</label>
        <input
          type="number"
          value={nPomodoros}
          onChange={(e) => setNPomodoros(Number(e.target.value))}
          className="border p-2 w-full"
          required
        />
      </div>
      <div className="mb-4">
        <label className="block mb-2">Tags (comma separated)</label>
        <input
          type="text"
          value={tags.join(', ')}
          onChange={(e) => setTags(e.target.value.split(',').map(tag => tag.trim()))}
          className="border p-2 w-full"
        />
      </div>
      <button type="submit" className="bg-blue-500 text-white p-2 rounded">Create Task</button>
    </form>
  );
};

export default TaskForm;

