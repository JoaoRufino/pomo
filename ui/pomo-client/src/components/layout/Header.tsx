// components/layout/Header.tsx
import React from 'react';
import Link from 'next/link';

const Header: React.FC = () => {
  return (
    <header className="bg-blue-600 p-4 text-white">
      <nav className="container mx-auto flex justify-between">
        <div className="text-lg font-bold">
          <Link href="/">Pomo Client</Link>
        </div>
        <div>
          <Link href="/api/tasks" className="ml-4">Tasks</Link>
        </div>
      </nav>
    </header>
  );
};

export default Header;

