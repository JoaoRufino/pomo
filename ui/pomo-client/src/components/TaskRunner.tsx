import React, { useState, useEffect } from 'react';
import { Task } from '../services/api';

interface TaskRunnerProps {
  task: Task;
}

enum TaskState {
  READY = 'READY',
  RUNNING = 'RUNNING',
  PAUSED = 'PAUSED',
  BREAKING = 'BREAKING',
  COMPLETE = 'COMPLETE',
}

const TaskRunner: React.FC<TaskRunnerProps> = ({ task }) => {
  const [state, setState] = useState(TaskState.READY);
  const [remainingTime, setRemainingTime] = useState(task.duration);
  const [count, setCount] = useState(0);

  useEffect(() => {
    let timer: NodeJS.Timeout;

    if (state === TaskState.RUNNING) {
      timer = setInterval(() => {
        setRemainingTime((prev) => {
          if (prev <= 1000) {
            setState(TaskState.BREAKING);
            clearInterval(timer);
            return task.duration;
          }
          return prev - 1000;
        });
      }, 1000);
    }

    return () => clearInterval(timer);
  }, [state, task.duration]);

  const startTask = () => {
    setState(TaskState.RUNNING);
  };

  const pauseTask = () => {
    setState(TaskState.PAUSED);
  };

  const resumeTask = () => {
    setState(TaskState.RUNNING);
  };

  const completeTask = () => {
    setState(TaskState.COMPLETE);
  };

  return (
    <div className="task-runner">
      <h2 className="text-xl font-semibold">{task.message}</h2>
      <p className="text-gray-500">Pomodoros: {task.n_pomodoros}</p>
      <p className="text-gray-500">
        Remaining Time: {new Date(remainingTime).toISOString().substr(11, 8)}
      </p>
      <p className="text-gray-500">State: {state}</p>
      {state === TaskState.READY && (
        <button
          className="bg-green-500 text-white p-2 rounded"
          onClick={startTask}
        >
          Start
        </button>
      )}
      {state === TaskState.RUNNING && (
        <button
          className="bg-yellow-500 text-white p-2 rounded"
          onClick={pauseTask}
        >
          Pause
        </button>
      )}
      {state === TaskState.PAUSED && (
        <button
          className="bg-blue-500 text-white p-2 rounded"
          onClick={resumeTask}
        >
          Resume
        </button>
      )}
      {state === TaskState.BREAKING && (
        <button
          className="bg-red-500 text-white p-2 rounded"
          onClick={() => setCount(count + 1)}
        >
          Next Pomodoro
        </button>
      )}
      {count >= task.n_pomodoros && (
        <button
          className="bg-purple-500 text-white p-2 rounded"
          onClick={completeTask}
        >
          Complete Task
        </button>
      )}
    </div>
  );
};

export default TaskRunner;

