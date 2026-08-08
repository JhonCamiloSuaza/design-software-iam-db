import React, { useCallback, useEffect, useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import { usersApi } from '../services/api';

export function UsersPage({ onMessage }) {
  const { token } = useAuth();
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [editing, setEditing] = useState(null);
  const load = useCallback(async () => { try { const [userData, roleData] = await Promise.all([usersApi.list(token), usersApi.roles(token)]); setUsers(userData); setRoles(roleData); } catch (error) { onMessage(error.message); } }, [token, onMessage]);
  useEffect(() => { load(); }, [load]);

  const create = async (event) => { event.preventDefault(); try { const data = await usersApi.create(token, Object.fromEntries(new FormData(event.currentTarget))); onMessage(data.message); event.currentTarget.reset(); await load(); } catch (error) { onMessage(error.message); } };
  const update = async (event) => { event.preventDefault(); const body = Object.fromEntries(new FormData(event.currentTarget)); body.isActive = body.isActive === 'true'; try { const data = await usersApi.update(token, editing.id, body); onMessage(data.message); setEditing(null); await load(); } catch (error) { onMessage(error.message); } };
  const deactivate = async (user) => { if (!window.confirm(`Desea desactivar a ${user.email}?`)) return; try { const data = await usersApi.deactivate(token, user.id); onMessage(data.message); await load(); } catch (error) { onMessage(error.message); } };

  return <div className="manager">
    <form onSubmit={create}><h2>Crear usuario administrativo</h2><div className="grid"><label>Nombres<input name="firstName" required /></label><label>Apellidos<input name="lastName" required /></label></div><label>Correo<input type="email" name="email" required /></label><label>Contrasena<input type="password" name="password" minLength="8" required /></label><div className="grid"><label>Tipo<select name="actorType" defaultValue="USER"><option>USER</option><option>INSTRUCTOR</option><option>LEARNER</option></select></label><label>Rol<select name="roleName" required>{roles.map((role) => <option value={role.name} key={role.name}>{role.displayName}</option>)}</select></label></div><button className="primary">Guardar usuario</button></form>
    <div><h2>Usuarios registrados</h2><div className="table users-table"><div className="thead"><span>Usuario</span><span>Rol</span><span>Estado</span><span>Acciones</span></div>{users.map((user) => <div className="trow" key={user.id}><span><b>{user.firstName} {user.lastName}</b><small>{user.email} · {user.actorType}</small></span><span>{user.roles || 'Sin rol'}</span><span>{user.isActive ? 'Activo' : 'Inactivo'}</span><span><button className="action" onClick={() => setEditing(user)}>Editar</button>{user.isActive && <button className="danger" onClick={() => deactivate(user)}>Desactivar</button>}</span></div>)}</div></div>
    {editing && <form onSubmit={update}><h2>Editar usuario</h2><div className="grid"><label>Nombres<input name="firstName" defaultValue={editing.firstName} required /></label><label>Apellidos<input name="lastName" defaultValue={editing.lastName} required /></label></div><div className="grid"><label>Tipo<select name="actorType" defaultValue={editing.actorType}><option>USER</option><option>INSTRUCTOR</option><option>LEARNER</option></select></label><label>Estado<select name="isActive" defaultValue={String(editing.isActive)}><option value="true">Activo</option><option value="false">Inactivo</option></select></label></div><button className="primary">Guardar cambios</button><button type="button" className="link" onClick={() => setEditing(null)}>Cancelar</button></form>}
  </div>;
}
