import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { IChatMessage } from '../types';
import { Assistant as BotIcon } from '@mui/icons-material';

interface ChatMessageProps {
  message: IChatMessage;
  children?: React.ReactNode;
}

interface BotChatMessageComponent extends React.FC<ChatMessageProps> {
  create: (message: string, type?: 'bot' | 'error') => IChatMessage;
}

const BotChatMessage: BotChatMessageComponent = ({ message, children }) => {
  const timestamp = message.payload?.timestamp
    ? new Date(message.payload.timestamp).toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit',
      })
    : '';

  return (
    <div className="react-chatbot-kit-chat-bot-message-container">
      <div className="react-chatbot-kit-chat-bot-avatar-container">
        <BotIcon />
      </div>
      <div className="bot-chat-message-bubble">
        <div className="markdown-content">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{message.message}</ReactMarkdown>
        </div>
        {children}
        {timestamp && <div className="chat-message-timestamp">{timestamp}</div>}
      </div>
    </div>
  );
};

BotChatMessage.create = (message: string, type: 'bot' | 'error' = 'bot'): IChatMessage => ({
  message,
  type,
  id: Date.now(),
  payload: { timestamp: new Date().toISOString() },
});

export default BotChatMessage;
