import React from 'react';
import { useAuth } from '../hooks/useAuth';

export function RegisterPage({ onSuccess, onMessage }) {
  const { register } = useAuth();
  const submit = async (event) => {
    event.preventDefault();
    const user = Object.fromEntries(new FormData(event.currentTarget));
    try { onMessage(await register(user)); onSuccess(); } catch (error) { onMessage(error.message); }
  };
  return <form onSubmit={submit}>
    <div className="grid"><label>Nombres<input name="firstName" required /></label><label>Apellidos<input name="lastName" required /></label></div>
    <label>Correo<input type="email" name="email" required /></label>
    <label>Contrasena<input type="password" name="password" minLength="8" required /></label>
    <button className="primary">Registrarme</button>
  </form>;
}
