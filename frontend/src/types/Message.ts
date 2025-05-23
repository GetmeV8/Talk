export interface Message {
  type: string;
  content: string;
  sender: string;
  timestamp?: number;
} 