import { ConnectionEvent } from '../lib/types';

type Props = { event?: ConnectionEvent };

export function DecisionModal({ event }: Props) {
  if (!event || event.state !== 'pending') return null;

  return (
    <div className="fixed right-6 top-6 w-96 rounded-xl border border-zinc-700 bg-zinc-900 p-4 shadow-2xl">
      <h2 className="text-lg font-semibold">Connection request</h2>
      <p className="mt-2 text-sm text-zinc-300">{event.process_path} → {event.domain}:{event.dst_port}</p>
      <div className="mt-4 grid grid-cols-2 gap-2 text-sm">
        <button className="rounded bg-emerald-600 px-3 py-2">Allow once</button>
        <button className="rounded bg-emerald-800 px-3 py-2">Allow always</button>
        <button className="rounded bg-rose-700 px-3 py-2">Deny once</button>
        <button className="rounded bg-rose-900 px-3 py-2">Deny always</button>
      </div>
    </div>
  );
}
