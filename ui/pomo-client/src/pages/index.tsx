import type { NextPage } from 'next';
import { GetServerSideProps } from 'next';
import TaskList from '../components/TaskList';
import { Task, fetchTasks } from '../services/api';
import { useState } from 'react';

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
      <TaskList tasks={tasks} />
    </div>
  );
};

export const getServerSideProps: GetServerSideProps = async () => {
  try {
    const data = await fetchTasks();
    return { props: { initialTasks: data.results } };
  } catch (error) {
    return { props: { initialTasks: [] } };
  }
};

export default Home;

