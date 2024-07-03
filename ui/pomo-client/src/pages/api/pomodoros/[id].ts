import type { NextApiRequest, NextApiResponse } from 'next';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const { method, query, body } = req;
  const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const url = `${baseUrl}/api/pomodoros/${query.id || ''}`;

  try {
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${process.env.NEXT_PUBLIC_API_TOKEN || 'default-token'}`,
      },
      body: method === 'POST' || method === 'PUT' ? JSON.stringify(body) : null,
    });

    const data = await response.json();
    res.status(response.status).json(data);
  } catch (error) {
    res.status(500).json({ error: 'Internal Server Error' });
  }
}

