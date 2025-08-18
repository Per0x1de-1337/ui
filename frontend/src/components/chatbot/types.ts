// frontend/src/components/chatbot/types.ts

export interface IChatMessage {
  message: string;
  type: string;
  id: number;
  loading?: boolean;
  widget?: string;
  payload?: {
    timestamp: string;
    sources?: string[];
  };
}

export interface IMessageProps {
  message: IChatMessage;
  [key: string]: unknown;
}
