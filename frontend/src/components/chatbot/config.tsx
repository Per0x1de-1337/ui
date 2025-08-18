// frontend/src/components/chatbot/config.ts
import { createChatBotMessage } from 'react-chatbot-kit';
import BotChatMessage from './components/BotChatMessage';
import UserChatMessage from './components/UserChatMessage';
import CustomChatInput from './components/CustomChatInput';
import { IMessageProps } from './types';
import { Assistant as BotIcon } from '@mui/icons-material';

const config = {
  initialMessages: [
    createChatBotMessage(`Hello! I'm your KubeStellar assistant. How can I help you today?`, {
      delay: 1000,
      widget: 'response',
    }),
  ],
  botName: 'KubeStellar Assistant',
  customComponents: {
    botAvatar: (props: any) => <BotIcon {...props} />,
    botChatMessage: (props: IMessageProps) => <BotChatMessage message={props.message} />,
    userChatMessage: (props: IMessageProps) => <UserChatMessage message={props.message} />,
    chatInput: (props: any) => <CustomChatInput {...props} />,
    botTypingIndicator: () => (
      <div className="typing-indicator">
        <div className="dot"></div>
        <div className="dot"></div>
        <div className="dot"></div>
      </div>
    ),
  },
};

export default config;