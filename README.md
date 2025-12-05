<img width="1117" height="785" alt="{A9802DD2-3342-41CD-8B40-2528B56D2B98}" src="https://github.com/user-attachments/assets/a22856e3-3426-4b96-843c-d1ee62f5bd09" /># 🍋 WiMon – Windows Lemon Monitor

WiMon (**Wi**ndows **Mon**itor) es una herramienta escrita en **Go** que muestra las **conexiones TCP salientes en estado ESTABLISHED** de un equipo Windows, a través de un **dashboard web moderno** servido desde un único ejecutable.

> Todo corre localmente.  
> Sin dependencias externas.  
> Sin instalar servidores o frameworks extra.

---


<img width="1117" height="785" alt="{A9802DD2-3342-41CD-8B40-2528B56D2B98}" src="https://github.com/user-attachments/assets/7850bec0-c25a-4144-a00f-148a448faba5" />



## ✨ Características principales

- 🟢 Backend escrito completamente en **Go**
- 🌐 Servidor web embebido (`net/http`)
- 🍋 Dashboard HTML moderno
- 📡 Obtención en tiempo real de conexiones TCP
- 🔎 Filtra solo conexiones en estado **ESTABLISHED**
- ⏱ Mide duración de cada conexión
- 🔁 API interna `/api/connections`
- ☑️ Cero dependencias de sistema externas

⚠️ *Campos como País, Rango y ASN están definidos en el modelo, pero aún no se rellenan (reservados para futura integración).*

---

# 🛠 Pasos de creación del proyecto

Estos son los pasos exactos usados para crear WiMon desde cero:

### 1. Crear carpeta del proyecto

```bash
mkdir WiMon
cd WiMon
go mod init wimon


3. Crear estructura inicial

WiMon/
 ├─ main.go
 └─ templates/
     └─ index.html


4. Instalar dependencia

go get github.com/shirou/gopsutil/v3/net


5. go run .
WiMon escuchando en http://localhost:8080



6. [
  {
    "remote_ip": "142.250.190.78",
    "country": "",
    "range": "",
    "asn": "",
    "protocol": "TCP",
    "since": "2025-12-05T16:30:20Z",
    "duration_secs": 75,
    "display_since": "16:30:20",
    "display_duration": "1m 15s"
  }
]


7. go build .

8. WiMon.exe

9. WiMon escuchando en http://localhost:8080







