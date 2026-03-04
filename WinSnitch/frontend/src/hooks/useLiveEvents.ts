import { useEffect, useState } from 'react';
import { ConnectionEvent } from '../lib/types';

export function useLiveEvents() {
  const [events, setEvents] = useState<ConnectionEvent[]>([]);

  useEffect(() => {
    const ws = new WebSocket('ws://127.0.0.1:47777/ws');
    ws.onmessage = (msg) => {
      const event = JSON.parse(msg.data) as ConnectionEvent;
      setEvents((prev) => [event, ...prev].slice(0, 500));
    };
    return () => ws.close();
  }, []);

  return events;
}
