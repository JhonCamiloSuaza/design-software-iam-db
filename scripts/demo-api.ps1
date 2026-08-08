$ErrorActionPreference = "Stop"
$apiUrl = "http://localhost:8080/api"
$timestamp = [DateTime]::UtcNow.ToString("yyyyMMddHHmmss")
$learnerEmail = "aprendiz.$timestamp@demo.local"
$managedEmail = "gestionado.$timestamp@demo.local"

function Show-Ok([string]$message) {
    Write-Host "[OK] $message" -ForegroundColor Green
}

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Path,
        [hashtable]$Body,
        [string]$Token
    )

    $parameters = @{
        Method      = $Method
        Uri         = "$apiUrl$Path"
        ContentType = "application/json"
    }
    if ($Body) { $parameters.Body = $Body | ConvertTo-Json }
    if ($Token) { $parameters.Headers = @{ Authorization = "Bearer $Token" } }
    Invoke-RestMethod @parameters
}

try {
    Write-Host "`nPRUEBA AUTOMATICA DEL MODULO IAM`n" -ForegroundColor Cyan

    $health = Invoke-Api -Method GET -Path "/health"
    if ($health.status -ne "ok") { throw "La API no esta saludable." }
    Show-Ok "API activa y conectada con PostgreSQL."

    $registered = Invoke-Api -Method POST -Path "/auth/register" -Body @{
        email = $learnerEmail; password = "Aprendiz2026";
        firstName = "Prueba"; lastName = "Aprendiz"
    }
    Show-Ok "Registro publico: $learnerEmail"

    $profile = Invoke-Api -Method GET -Path "/me" -Token $registered.token
    if ($profile.user.actorType -ne "LEARNER") { throw "El registro no recibio actor LEARNER." }
    Show-Ok "Perfil protegido: actor_type LEARNER y token JWT valido."

    try {
        Invoke-Api -Method GET -Path "/users" -Token $registered.token | Out-Null
        throw "La ruta administrativa permitio el acceso del aprendiz."
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -ne 403) { throw }
        Show-Ok "RBAC: el aprendiz fue bloqueado con HTTP 403."
    }

    $recovery = Invoke-Api -Method POST -Path "/auth/forgot-password" -Body @{ email = $learnerEmail }
    Invoke-Api -Method POST -Path "/auth/reset-password" -Body @{
        token = $recovery.resetToken; password = "NuevaAprendiz2026"
    } | Out-Null
    Invoke-Api -Method POST -Path "/auth/login" -Body @{
        email = $learnerEmail; password = "NuevaAprendiz2026"
    } | Out-Null
    Show-Ok "Recuperacion: token utilizado y login con la nueva contrasena."

    $adminRecovery = Invoke-Api -Method POST -Path "/auth/forgot-password" -Body @{ email = "admin@sena.edu.co" }
    Invoke-Api -Method POST -Path "/auth/reset-password" -Body @{
        token = $adminRecovery.resetToken; password = "AdminSegura2026"
    } | Out-Null
    $adminSession = Invoke-Api -Method POST -Path "/auth/login" -Body @{
        email = "admin@sena.edu.co"; password = "AdminSegura2026"
    }
    Show-Ok "Administrador autenticado con permiso de gestion."

    Invoke-Api -Method POST -Path "/users" -Token $adminSession.token -Body @{
        email = $managedEmail; password = "Gestionado2026";
        firstName = "Usuario"; lastName = "Gestionado";
        actorType = "INSTRUCTOR"; roleName = "INSTRUCTOR"
    } | Out-Null
    $users = Invoke-Api -Method GET -Path "/users" -Token $adminSession.token
    $managedUser = $users | Where-Object { $_.email -eq $managedEmail }
    if (-not $managedUser) { throw "El usuario administrativo no aparece en el listado." }
    Show-Ok "CRUD Create y Read: usuario creado y listado."

    Invoke-Api -Method PUT -Path "/users/$($managedUser.id)" -Token $adminSession.token -Body @{
        firstName = "Usuario"; lastName = "Modificado";
        actorType = "INSTRUCTOR"; isActive = $true
    } | Out-Null
    Show-Ok "CRUD Update: usuario modificado."

    Invoke-Api -Method DELETE -Path "/users/$($managedUser.id)" -Token $adminSession.token | Out-Null
    $updatedUsers = Invoke-Api -Method GET -Path "/users" -Token $adminSession.token
    $inactiveUser = $updatedUsers | Where-Object { $_.id -eq $managedUser.id }
    if ($inactiveUser.isActive -ne $false) { throw "El usuario no quedo inactivo." }
    Show-Ok "CRUD Delete logico: usuario desactivado y conservado para auditoria."

    Write-Host "`nTODAS LAS PRUEBAS TERMINARON CORRECTAMENTE.`n" -ForegroundColor Cyan
    Write-Host "Aprendiz creado: $learnerEmail"
    Write-Host "Administrador: admin@sena.edu.co / AdminSegura2026"
} catch {
    Write-Host "`n[ERROR] $($_.Exception.Message)`n" -ForegroundColor Red
    exit 1
}
