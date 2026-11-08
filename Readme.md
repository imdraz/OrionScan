# OrionScan

Escáner de puertos TCP rápido y liviano, escrito en Go. Forma parte del proyecto [ORION.LAB](https://imdraz.github.io/orion-lab/).

Por defecto revisa la IP de tu red local, así podés ver qué puertos de tu equipo están accesibles desde la red. También podés indicar otro equipo con `-host`.

## Características

- Un único ejecutable (`.exe` en Windows), sin instalar Python, Node ni otras dependencias.
- Escaneo en paralelo mediante goroutines; la cantidad de conexiones simultáneas es configurable.
- Acepta listas (`80,443`) y rangos (`1-1024`) de puertos.
- Detecta automáticamente la IP de la red local si no se indica un objetivo.
- Usa TCP connect, así que no necesita permisos de administrador ni drivers adicionales.

## Uso

```text
OrionScan [opciones]

  -host string   equipo a revisar (vacío = IP de tu red local)
  -p string      puertos: lista (80,443) o rango (1-1024) (default "1-1024")
  -t duration    tiempo máximo de espera por puerto (default 500ms)
  -w int         conexiones simultáneas (default 500)
```

Ejemplos:

```powershell
.\orion-scan.exe
.\orion-scan.exe -p 22,80,443,3389,8080
.\orion-scan.exe -host 192.168.1.1 -p 1-65535 -w 1000
```

## Cómo interpretar el resultado

Un puerto aparece como **abierto** cuando acepta una conexión TCP. Un puerto que no aparece puede estar cerrado o filtrado por un firewall, esta herramienta no distingue entre ambos casos.

## Compilar desde el código

Requisito: [Go](https://go.dev/dl/) instalado.

```powershell
go build -ldflags="-s -w" -o orion-scan.exe
```

Para compilar para Linux:

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o orion-scan
Remove-Item Env:GOOS, Env:GOARCH
```

## Descarga

Podés descargar el ejecutable precompilado desde la sección [Releases](../../releases).

Windows SmartScreen puede advertir que el archivo no está firmado digitalmente. Es normal en herramientas nuevas.

## Limitaciones

- Solo detecta puertos TCP. No hace escaneo UDP ni SYN.
- No identifica el servicio que corre en cada puerto.
- Un timeout muy corto (`-t`) puede omitir servicios con respuesta lenta.

## Uso responsable

La idea es que sea ejecutado desde el equipo personal o laboral que utilices.
