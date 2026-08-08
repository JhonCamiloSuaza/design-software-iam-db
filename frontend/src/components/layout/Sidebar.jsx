import React from 'react';
import { useAuth } from '../../hooks/useAuth';

export function Sidebar({ view, onNavigate }) {
  const { token, canManageUsers } = useAuth();
  const items = [['login', 'Iniciar sesion'], ['register', 'Crear usuario'], ['recovery', 'Recuperar acceso']];
  if (token) items.push(['dashboard', 'Mi panel']);
  if (canManageUsers) items.push(['users', 'Administrar usuarios']);

  return <aside>
    <div className="brand"><span>IAM</span><p>Modulo de seguridad</p></div>
    <nav>{items.map(([id, label]) => <button key={id} className={view === id ? 'active' : ''} onClick={() => onNavigate(id)}>{label}</button>)}</nav>
    <div className="status"><i></i> React + Go + PostgreSQL</div>
  </aside>;
}
