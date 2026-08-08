import React from 'react';
import { useAuth } from '../hooks/useAuth';

export function DashboardPage({ onLogout }) {
  const { profile, summary, logout } = useAuth();
  if (!profile) return <p>Cargando informacion del usuario...</p>;
  const closeSession = () => { logout(); onLogout(); };
  return <>
    <div className="profile"><div>{profile.user.firstName?.[0]}{profile.user.lastName?.[0]}</div><section><h2>{profile.user.firstName} {profile.user.lastName}</h2><p>{profile.user.email} · {profile.user.actorType}</p></section><button onClick={closeSession}>Cerrar sesion</button></div>
    <div className="metrics">{summary && Object.entries(summary).map(([key, value]) => <article key={key}><strong>{value}</strong><span>{{ users: 'Usuarios', roles: 'Roles', features: 'Permisos', modules: 'Modulos' }[key]}</span></article>)}</div>
    <h2>Permisos asignados</h2>
    <div className="table"><div className="thead"><span>Rol</span><span>Funcion</span><span>Alcance</span></div>{profile.permissions.map((item, index) => <div className="trow" key={index}><span>{item.role}</span><span>{item.feature}</span><span>{item.scope}</span></div>)}</div>
  </>;
}
