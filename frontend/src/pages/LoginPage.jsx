import React from 'react';
import { useAuth } from '../hooks/useAuth';

export function LoginPage({ onSuccess, onNavigate, onMessage }) {
  const { login } = useAuth();
  const submit = async (event) => {
    event.preventDefault();
    const credentials = Object.fromEntries(new FormData(event.currentTarget));
    try { onMessage(await login(credentials)); onSuccess(); } catch (error) { onMessage(error.message); }
  };
  return <form onSubmit={submit}>
    <label>Correo<input type="email" name="email" required /></label>
    <label>Contrasena<input type="password" name="password" required /></label>
    <button className="primary">Ingresar</button>
    <button type="button" className="link" onClick={() => onNavigate('recovery')}>Olvide mi contrasena</button>
  </form>;
}
