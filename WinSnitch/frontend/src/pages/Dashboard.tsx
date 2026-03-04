import { ShieldAlert, Network, History } from 'lucide-react';
import { useLiveEvents } from '../hooks/useLiveEvents';
import { DecisionModal } from '../components/DecisionModal';

export function Dashboard() {
  const events = useLiveEvents();

  return (
    <main className="min-h-screen bg-zinc-950 text-zinc-100 p-8">
      <header className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">WinSnitch</h1>
          <p className="text-zinc-400">OpenSnitch for Windows with WFP + Tauri 2</p>
        </div>
      </header>

      <section className="grid grid-cols-1 gap-4 md:grid-cols-3">
        {[{label:'Pending',value:events.filter(e=>e.state==='pending').length,icon:ShieldAlert},{label:'Total Events',value:events.length,icon:Network},{label:'History',value:500,icon:History}].map((card) => (
          <article key={card.label} className="rounded-xl border border-zinc-800 bg-zinc-900 p-4">
            <card.icon className="mb-2 h-5 w-5 text-cyan-400" />
            <p className="text-sm text-zinc-400">{card.label}</p>
            <p className="text-2xl font-semibold">{card.value}</p>
          </article>
        ))}
      </section>

      <section className="mt-6 rounded-xl border border-zinc-800 bg-zinc-900">
        <div className="border-b border-zinc-800 px-4 py-3 font-medium">Live connections</div>
        <div className="max-h-[500px] overflow-auto">
          {events.map((e) => (
            <div key={e.id} className="grid grid-cols-12 gap-2 border-b border-zinc-800 px-4 py-2 text-sm">
              <span className="col-span-4 truncate">{e.process_path}</span>
              <span className="col-span-3 truncate">{e.domain || e.dst_ip}</span>
              <span className="col-span-2">{e.protocol.toUpperCase()}</span>
              <span className="col-span-1">{e.dst_port}</span>
              <span className="col-span-2 uppercase text-cyan-300">{e.state}</span>
            </div>
          ))}
        </div>
      </section>
      <DecisionModal event={events.find((e) => e.state === 'pending')} />
    </main>
  );
}
