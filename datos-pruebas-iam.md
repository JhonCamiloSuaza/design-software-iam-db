# Datos, comandos y SQL para las pruebas

## PUNTO: DOCKER - Carpeta del proyecto

```text
C:\Users\usuario\Downloads\design-software-iam-db-main
```

```powershell
Get-Location
```

## PUNTO: DOCKER - Reinicio limpio, borra los datos de prueba

```powershell
docker compose down -v
docker compose up -d --build
docker compose ps -a
```

```text
postgres: Up o healthy
liquibase: Exited (0)
backend: Up o healthy
frontend: Up
```

## PUNTO: DOCKER - Reinicio conservando los datos

```powershell
docker compose down
docker compose up -d --build
docker compose ps -a
```

## PUNTO: BD - Conexion en pgAdmin

```text
Pestana General
Name: IAM DB Docker
Server group: Servers

Pestana Connection
Host name/address: localhost
Port: 5436
Maintenance database: iam_db
Username: iam_user
Password: iam_pass
Save password: activado
Connect now?: activado
```

## PUNTO: BD - Conteos iniciales

```sql
SELECT
  (SELECT COUNT(*)
   FROM information_schema.schemata
   WHERE schema_name IN ('identity','rbac_catalog','rbac','session','identity_audit')
  ) AS schemas,
  (SELECT COUNT(*) FROM rbac_catalog.module) AS modules,
  (SELECT COUNT(*) FROM rbac_catalog.feature) AS features,
  (SELECT COUNT(*) FROM rbac.role) AS roles;
```

```text
schemas = 5
modules = 10
features = 62
roles = 7
```

## PUNTO: BACKEND - Estado de la API

```text
http://localhost:8080/api/health
```

```json
{"service":"iam-security-go","status":"ok"}
```

## PUNTO: FRONTEND - Aplicacion

```text
http://localhost:5173
```

## PUNTO: FRONTEND - Registro del aprendiz

```text
Nombres: Jhon
Apellidos: Prueba
Correo: pepito.suaza.prueba@sena.edu.co
Contrasena: Aprendiz2026
```

## PUNTO: BD - Comprobar el registro

```sql
SELECT u.email, u.actor_type, u.is_active, r.name AS role
FROM identity."user" u
JOIN rbac.user_role ur ON ur.user_id = u.id
JOIN rbac.role r ON r.id = ur.role_id
ORDER BY u.created_at DESC
LIMIT 1;
```

```text
email = pepito.suaza.prueba@sena.edu.co
actor_type = LEARNER
role = LEARNER
is_active = true
```

## PUNTO: FRONTEND - Login del aprendiz

```text
Correo: pepito.suaza.prueba@sena.edu.co
Contrasena: Aprendiz2026
```

## PUNTO: BD - Auditoria del login

```sql
SELECT email_attempted, outcome, attempted_at
FROM identity_audit.audit_login
ORDER BY attempted_at DESC
LIMIT 3;
```

```text
email_attempted = pepito.suaza.prueba@sena.edu.co
outcome = SUCCESS
```

## PUNTO: FRONTEND - Recuperacion del aprendiz

```text
Correo: pepito.suaza.prueba@sena.edu.co
Nueva contrasena: NuevaAprendiz2026
```

## PUNTO: FRONTEND - Login despues de la recuperacion

```text
Correo: pepito.suaza.prueba@sena.edu.co
Contrasena: NuevaAprendiz2026
```

## PUNTO: BD - Comprobar la recuperacion

```sql
SELECT u.email, pr.is_used, pr.requested_at, pr.expires_at
FROM session.password_reset_request pr
JOIN identity."user" u ON u.id = pr.user_id
ORDER BY pr.requested_at DESC
LIMIT 1;
```

```text
email = pepito.suaza.prueba@sena.edu.co
is_used = true
expires_at = aproximadamente una hora despues de requested_at
```

## PUNTO: FRONTEND - Recuperar administrador

```text
Correo: admin@sena.edu.co
Nueva contrasena: AdminSegura2026
```

## PUNTO: FRONTEND - Login del administrador

```text
Correo: admin@sena.edu.co
Contrasena: AdminSegura2026
```

```text
Tipo visible junto al correo: USER
Rol visible en Permisos asignados: SYSTEM_ADMIN
Opcion visible: Administrar usuarios
```

## PUNTO: FRONTEND - Crear usuario desde el CRUD

```text
Nombres: Usuario
Apellidos: CRUD
Correo: usuario.crud.prueba@sena.edu.co
Contrasena: Usuario2026
Tipo: LEARNER
Rol: Aprendiz
```

## PUNTO: FRONTEND - Editar y desactivar

```text
Nuevo nombre: Usuario Editado
Accion final: Desactivar
```

## PUNTO: BD - Comprobar el CRUD

```sql
SELECT email, first_name, last_name, actor_type, is_active, updated_at
FROM identity."user"
ORDER BY updated_at DESC
LIMIT 5;
```

```text
email = usuario.crud.prueba@sena.edu.co
first_name = Usuario Editado
actor_type = LEARNER
is_active = false
```

## PUNTO: BACKEND - Prueba automatica de la API

```powershell
powershell -ExecutionPolicy Bypass -File scripts\demo-api.ps1
```

```text
[OK] RBAC: el aprendiz fue bloqueado con HTTP 403.
[OK] Recuperacion y login con la nueva contrasena.
[OK] CRUD Create y Read.
[OK] CRUD Update.
[OK] CRUD Delete logico.
TODAS LAS PRUEBAS TERMINARON CORRECTAMENTE.
```

## PUNTO: BACKEND - Rutas principales

```text
GET    /api/health
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/forgot-password
POST   /api/auth/reset-password
GET    /api/me
GET    /api/catalog/summary
GET    /api/users
POST   /api/users
PUT    /api/users/{id}
DELETE /api/users/{id}
GET    /api/roles
```
