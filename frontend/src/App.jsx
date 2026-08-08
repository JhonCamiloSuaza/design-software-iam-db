import React, { useEffect, useState } from 'react';
import { Notice } from './components/Notice';
import { Sidebar } from './components/layout/Sidebar';
import { useAuth } from './hooks/useAuth';
import { DashboardPage } from './pages/DashboardPage';
import { LoginPage } from './pages/LoginPage';
import { RecoveryPage } from './pages/RecoveryPage';
import { RegisterPage } from './pages/RegisterPage';
import { UsersPage } from './pages/UsersPage';

const titles = { login: 'Bienvenido', register: 'Crear una cuenta', recovery: 'Recuperar contrasena', dashboard: 'Panel de usuario', users: 'Administrar usuarios' };

export default function App() {
  const { token, canManageUsers } = useAuth();
  const [view, setView] = useState(token ? 'dashboard' : 'login');
  const [message, setMessage] = useState('');

  useEffect(() => {
    if (!token && (view === 'dashboard' || view === 'users')) setView('login');
    if (view === 'users' && !canManageUsers) setView('dashboard');
  }, [token, canManageUsers, view]);

  return <main>
    <Sidebar view={view} onNavigate={setView} />
    <section className="content">
      <header><p>Gestion de identidad y acceso</p><h1>{titles[view]}</h1></header>
      <Notice message={message} onClose={() => setMessage('')} />
      {view === 'login' && <LoginPage onSuccess={() => setView('dashboard')} onNavigate={setView} onMessage={setMessage} />}
      {view === 'register' && <RegisterPage onSuccess={() => setView('dashboard')} onMessage={setMessage} />}
      {view === 'recovery' && <RecoveryPage onSuccess={() => setView('login')} onMessage={setMessage} />}
      {view === 'dashboard' && <DashboardPage onLogout={() => setView('login')} />}
      {view === 'users' && canManageUsers && <UsersPage onMessage={setMessage} />}
    </section>
  </main>;
}
