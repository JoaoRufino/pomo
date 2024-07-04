import type { NextPage, GetStaticProps } from 'next';
import { useState } from 'react';
import TaskList from '../components/TaskList';
import { Task, fetchTasks } from '../services/api';

type HomeProps = {
  initialTasks: Task[];
};

const Home: NextPage<HomeProps> = ({ initialTasks }) => {
  const [tasks, setTasks] = useState<Task[]>(initialTasks);

  const handleTaskCreated = (newTask: Task) => {
    setTasks([newTask, ...tasks]);
  };

  return (
    <div className="p-8">
      <TaskList tasks={tasks} onTaskCreated={handleTaskCreated} />
    </div>
  );
};

export const getStaticProps: GetStaticProps = async () => {
  try {
    const data = await fetchTasks();
    return { props: { initialTasks: data.results } };
  } catch (error) {
    return { props: { initialTasks: [] } };
  }
};

export default Home;

