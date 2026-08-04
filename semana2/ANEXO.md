# Anexo — Cómo investigar por tu cuenta (semana 2: HTTP y JSON)

Esta semana entras a backend real: un servidor HTTP que responde JSON. Son
paquetes más grandes que `fmt` o `strconv`, así que además de `go doc`
necesitas saber probar cosas en vivo, sin escribir código.

## `net/http`: busca la firma del handler primero

Todo handler HTTP en Go tiene la misma forma:

```powershell
go doc http.HandlerFunc
```

Vas a ver que es literalmente `func(ResponseWriter, *Request)`. Cualquier
función con esa firma sirve como handler. Cuando te pase algo raro con un
handler (por ejemplo un error de tipos al pasarlo a `HandleFunc`), lo
primero es comparar la firma de tu función contra esta.

Para ver qué puedes hacer con el request:

```powershell
go doc http.Request
```

Te muestra todos los campos exportados (`Method`, `URL`, `Header`, `Body`,
etc.) — ahí está la respuesta a "¿cómo leo el query param X?" o "¿cómo
leo un header?".

## `encoding/json`: los struct tags no son magia, están documentados

```powershell
go doc json.Marshal
go doc json.Unmarshal
```

Cuando algo no se serializa como esperas (un campo no aparece, o aparece
con otro nombre), el 90% de las veces es un tag `json:"..."` mal escrito
en el struct. `go doc` no te va a mostrar tus propios structs, pero sí te
recuerda la sintaxis exacta del tag (`json:"nombre,omitempty"`).

## Probar una API sin escribir código: `curl`

Antes de asumir que tu servidor está mal, pruébalo directo:

```powershell
# GET
curl http://localhost:8080/productos

# POST con JSON
curl -X POST http://localhost:8080/productos \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Teclado","precio":50}'

# ver también el código de status HTTP
curl -i http://localhost:8080/productos
```

Esto te aísla el problema: si `curl` ya te devuelve algo raro, el bug está
en el servidor, no en algo que intentaste hacer con otra herramienta.

## Regla de esta semana

Cuando algo en HTTP o JSON no funcione como esperas: primero aísla el problema
con `curl` (o con `Invoke-WebRequest`, su equivalente nativo de PowerShell),
después revisa la firma exacta con `go doc`, y solo entonces toca el código.
