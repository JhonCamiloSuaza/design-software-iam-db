const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

async function request(path, options = {}) {
  const response = await fetch(`${API_URL}${path}`, {
    method: options.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {})
    },
    ...(options.body ? { body: JSON.stringify(options.body) } : {})
  });
  const data = await response.json();
  if (!response.ok) throw new Error(data.message || 'Ocurrio un problema.');
  return data;
}

export const authApi = {
  login: (credentials) => request('/auth/login', { method: 'POST', body: credentials }),
  register: (user) => request('/auth/register', { method: 'POST', body: user }),
  forgotPassword: (email) => request('/auth/forgot-password', { method: 'POST', body: { email } }),
  resetPassword: (token, password) => request('/auth/reset-password', { method: 'POST', body: { token, password } }),
  profile: (token) => request('/me', { token }),
  summary: (token) => request('/catalog/summary', { token })
};

export const usersApi = {
  list: (token) => request('/users', { token }),
  roles: (token) => request('/roles', { token }),
  create: (token, user) => request('/users', { method: 'POST', token, body: user }),
  update: (token, id, user) => request(`/users/${id}`, { method: 'PUT', token, body: user }),
  deactivate: (token, id) => request(`/users/${id}`, { method: 'DELETE', token })
};
