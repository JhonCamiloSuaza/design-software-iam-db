import React, { useState } from 'react';
import { authApi } from '../services/api';

export function RecoveryPage({ onSuccess, onMessage }) {
  const [resetToken, setResetToken] = useState('');
  const forgot = async (event) => {
    event.preventDefault();
    try { const data = await authApi.forgotPassword(new FormData(event.currentTarget).get('email')); setResetToken(data.resetToken || ''); onMessage(data.message); } catch (error) { onMessage(error.message); }
  };
  const reset = async (event) => {
    event.preventDefault();
    try { const data = await authApi.resetPassword(resetToken, new FormData(event.currentTarget).get('password')); onMessage(data.message); setResetToken(''); onSuccess(); } catch (error) { onMessage(error.message); }
  };
  return <>
    <form onSubmit={forgot}><label>Correo de la cuenta<input type="email" name="email" required /></label><button className="primary">Generar recuperacion</button></form>
    {resetToken && <form className="reset" onSubmit={reset}><p>Token de demostracion: <code>{resetToken}</code></p><label>Nueva contrasena<input type="password" name="password" minLength="8" required /></label><button className="primary">Cambiar contrasena</button></form>}
  </>;
}
