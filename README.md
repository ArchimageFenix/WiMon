# 🍋 WiMon – Windows Lemon Monitor

WiMon (**Wi**ndows **Mon**itor, 🍋) es una herramienta escrita en **Go** que muestra las **conexiones TCP salientes en estado ESTABLISHED** de tu equipo Windows, todo a través de un **dashboard web moderno** servido desde un único ejecutable.

> Backend en Go + HTML embebido con `embed` + API JSON = un monitor ligero, local y sin dependencias externas pesadas.

---

## ✨ Características

- 🟢 **Go nativo**: todo el backend en un solo binario.
- 🌐 **Servidor web embebido** (`net/http`) escuchando en `http://localhost:8080`.
- 📊 **Dashboard HTML** (en `templates/index.html`) con estilo moderno tipo “lemon dark theme”.
- 🔍 **Conexiones reales TCP ESTABLISHED** usando `github.com/shirou/gopsutil/v3/net`.
- ⏱️ Cálculo de **tiempo conectado** por conexión:
  - Se guarda el momento en que WiMon ve la conexión por primera vez.
  - Se muestra una duración legible (`10s`, `2m 15s`, `1h 3m`, etc.).
- 🔁 **Actualizable fácilmente** desde el frontend mediante llamadas a `/api/connections`.

> En esta versión, los campos `Country`, `Range` y `ASN` están definidos en el modelo pero aún se devuelven vacíos (pensados para futuras mejoras).

---

## 🧩 Arquitectura básica

### 🖥 Backend (Go)

Archivo principal: `main.go`

- Embebe plantillas HTML:

  ```go
  //go:embed templates/*
  var templatesFS embed.FS
