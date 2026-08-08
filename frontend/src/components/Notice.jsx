import React from 'react';

export function Notice({ message, onClose }) {
  if (!message) return null;
  return <div className="notice">{message}<button onClick={onClose}>x</button></div>;
}
