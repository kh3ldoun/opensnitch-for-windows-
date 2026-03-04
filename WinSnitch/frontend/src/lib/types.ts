export interface ConnectionEvent {
  id: string;
  node_id: string;
  process_path: string;
  domain: string;
  dst_ip: string;
  dst_port: number;
  protocol: string;
  ipv6: boolean;
  timestamp: string;
  state: string;
}
