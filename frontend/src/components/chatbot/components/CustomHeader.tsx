import React from 'react';

interface CustomHeaderProps {
  toggleFullScreen: () => void;
  isFullScreen: boolean;
  onClose: () => void;
}

const CustomHeader: React.FC<CustomHeaderProps> = ({
  toggleFullScreen,
  isFullScreen,
  onClose,
}) => {
  return (
    <div className="custom-chatbot-header">
      <span className="custom-chatbot-header-title">KubeStellar Assistant</span>
      <div>
        <button
          onClick={toggleFullScreen}
          className="custom-chatbot-fullscreen-btn"
          aria-label={isFullScreen ? 'Exit full screen' : 'Enter full screen'}
        >
          {isFullScreen ? (
            // Collapse Icon
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path d="M8 3v3a2 2 0 0 1-2 2H3m18 0h-3a2 2 0 0 1-2-2V3m0 18v-3a2 2 0 0 1 2-2h3M3 16h3a2 2 0 0 1 2 2v3" />
            </svg>
          ) : (
            // Expand Icon
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3" />
            </svg>
          )}
        </button>
        <button
          onClick={onClose}
          className="custom-chatbot-close-btn"
          aria-label="Close chat"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
          >
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>
  );
};

export default CustomHeader;
