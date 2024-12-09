import React, { useEffect, useState } from "react";
import "../stylesheets/notificationbox.css"; // css

const NotificationBox = ({
  message = "This is a notification",
  duration = 5000, // 5 seconds
  position = "bottom-right",
  onClose,
}) => {
  const [visible, setVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false);
      if (onClose) onClose();
    }, duration);

    return () => clearTimeout(timer);

  }, [duration, onClose]);

  if (!visible) return null;

  return (
    <div className={`notification-box ${position}`}>
      <p>{message}</p>
      <button className="close-button" onClick={() => setVisible(false)}>
        X
      </button>
    </div>
  );
};

export default NotificationBox;
