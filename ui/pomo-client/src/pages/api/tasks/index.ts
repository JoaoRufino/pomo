import type { NextApiRequest, NextApiResponse } from 'next';
import { API_URL } from '../../../services/api';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const url = `${API_URL}/api/tasks`;

  try {
    const response = await fetch(url, {
      method: req.method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: req.method === 'GET' ? null : JSON.stringify(req.body),
    });

    if (!response.ok) {
      const errorData = await response.text();
      console.error(`Error from Golang server: ${errorData}`);
      res.status(response.status).json({ error: `Golang server error: ${errorData}` });
      return;
    }

    const data = await response.json();
    res.status(response.status).json(data);
  } catch (error) {
    console.error('Error in Next.js API route:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
}

