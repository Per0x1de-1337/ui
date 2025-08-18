import React, { useState } from 'react';
import BaseChatbot from 'react-chatbot-kit';
// import 'react-chatbot-kit/build/main.css';
import './Chatbot.css';

import config from './config.tsx';
import ActionProvider from './ActionProvider';
import MessageParser from './MessageParser';
import CustomHeader from './components/CustomHeader';
import { IChatMessage } from './types';

interface ChatbotProps {
  onClose: () => void;
}

const Chatbot: React.FC<ChatbotProps> = ({ onClose }) => {
  const [isFullScreen, setIsFullScreen] = useState(false);

  const toggleFullScreen = () => setIsFullScreen(prev => !prev);

  const chatbotConfig = {
    ...config,
    customComponents: {
      ...config.customComponents,
      header: () => (
        <CustomHeader
          toggleFullScreen={toggleFullScreen}
          isFullScreen={isFullScreen}
          onClose={onClose}
        />
      ),
    },
  };

  const saveMessages = (messages: IChatMessage[]): void => {
    localStorage.setItem('chatbot_messages', JSON.stringify(messages));
  };

  const loadMessages = (): IChatMessage[] => {
    const messagesJSON = localStorage.getItem('chatbot_messages');
    if (!messagesJSON) return [];

    try {
      const parsedMessages = JSON.parse(messagesJSON);
      return Array.isArray(parsedMessages) ? parsedMessages : [];
    } catch (e) {
      console.error('Failed to parse messages from localStorage', e);
      return [];
    }
  };

  const containerClasses = isFullScreen
    ? 'chatbot-main-container fullscreen'
    : 'chatbot-main-container';

  return (
    <div className={containerClasses}>
      <div className="chatbot-container">
        <BaseChatbot
          config={chatbotConfig}
          messageHistory={loadMessages()}
          saveMessages={saveMessages}
          messageParser={MessageParser}
          actionProvider={ActionProvider}
        />
      </div>
    </div>
  );
};

export default Chatbot;
