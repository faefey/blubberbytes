import React, { useEffect, useState } from "react";
import "../stylesheets/notificationbox.css"; // css

const NotificationBox = ({
  message = "This is a notification",
  duration = 6000, // 6 seconds
  position = "bottom-right",
  setMessage,
  onClose,
}) => {
  useEffect(() => {
    const timer = setTimeout(() => {
      setMessage("");
      if (onClose) onClose();
    }, duration);

    return () => clearTimeout(timer);

  }, [duration, onClose]);

  return (
    <div className={`notification-box ${position}`}>
      <p>{message}</p>
      <button className="close-button" onClick={() => setMessage("")}>
        X
      </button>
    </div>
  );
};

export default NotificationBox;
