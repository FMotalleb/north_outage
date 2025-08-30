export interface DataItem {
  id: number;
  unique_hash: string;
  city: string;
  address: string;
  start: string;
  end: string;
  created_at: string;
}

export interface FilterState {
  city: string;
  date: string;
  address: string;
}