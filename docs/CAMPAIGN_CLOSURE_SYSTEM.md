# Sistema de Cierre de Campañas y Auditoría

Este documento describe el sistema de cierre de campañas con generación de reportes de auditoría y puntuación de transparencia implementado para la HU-012.

## Tabla de Contenidos

1. [Resumen](#resumen)
2. [Arquitectura](#arquitectura)
3. [Base de Datos](#base-de-datos)
4. [Endpoints](#endpoints)
5. [Puntuación de Transparencia](#puntuación-de-transparencia)
6. [Flujo de Cierre](#flujo-de-cierre)
7. [Ejemplos de Uso](#ejemplos-de-uso)

---

## Resumen

El sistema permite:
- Cerrar campañas manualmente (por admin) o automáticamente (por meta alcanzada o fecha vencida)
- Generar un reporte de auditoría con métricas financieras y de transparencia
- Calcular una puntuación de transparencia (0-100) basada en múltiples criterios
- Generar un PDF de auditoría descargable
- Bloquear nuevas donaciones en campañas cerradas

---

## Arquitectura

### Estructura de Archivos

```
campaign/closure/
├── closure.go       # Entidades de dominio y tipos
├── model.go         # Modelo GORM para PostgreSQL
├── repository.go    # Acceso a datos y queries de métricas
├── service.go       # Lógica de negocio y cálculo de transparencia
├── handlers.go      # Endpoints HTTP
└── pdf_generator.go # Generador de PDF de auditoría
```

### Dependencias

El servicio de closure depende de:
- `campaign.Service` - Para obtener info de campaña y actualizar estado
- `organizer.Service` - Para obtener nombre del organizador
- `contract.Service` - Para verificar si existe contrato firmado
- `s3client.Client` - Para subir el PDF de auditoría

---

## Base de Datos

### Tabla: `campaign_closure_reports`

Almacena el reporte de cierre de cada campaña.

```sql
CREATE TABLE campaign_closure_reports (
    id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL UNIQUE REFERENCES campaigns(id),

    -- Tipo de cierre
    closure_type VARCHAR(50) NOT NULL,  -- 'goal_reached', 'end_date', 'manual'
    closure_reason TEXT,                 -- Justificación (obligatorio si manual)
    closed_by UUID REFERENCES users(id), -- Admin que cerró

    -- Métricas financieras
    total_raised DECIMAL(12,2),
    total_donors INTEGER,
    total_donations INTEGER,
    campaign_goal DECIMAL(12,2),
    goal_percentage DECIMAL(5,2),

    -- Gastos
    total_expenses DECIMAL(12,2),
    total_receipts INTEGER,
    receipts_with_documents INTEGER,

    -- Actividades
    total_activities INTEGER,

    -- Transparencia
    transparency_score DECIMAL(5,2),      -- 0-100
    transparency_breakdown JSONB,          -- Detalle del cálculo

    -- Alertas (placeholder)
    alerts_count INTEGER,
    alerts_resolved INTEGER,

    -- PDF
    report_pdf_url TEXT,
    report_hash VARCHAR(64),

    -- Timestamps
    closed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE
);
```

### Tabla: `campaign_alerts` (Placeholder)

Preparada para futuro sistema de alertas/denuncias.

```sql
CREATE TABLE campaign_alerts (
    id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL REFERENCES campaigns(id),
    alert_type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    severity VARCHAR(20) DEFAULT 'medium',
    reported_by UUID REFERENCES users(id),
    resolved_by UUID REFERENCES users(id),
    resolution_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE
);
```

---

## Endpoints

### 1. Cerrar Campaña (Admin)

**POST** `/api/campaigns/:id/close`

Cierra una campaña y genera el reporte de auditoría.

**Autenticación:** Requiere rol `admin`

**Request Body:**
```json
{
    "closure_type": "manual",
    "reason": "Campaña completada exitosamente antes del plazo previsto"
}
```

**Tipos de cierre válidos:**
- `goal_reached` - La campaña alcanzó su meta de recaudación
- `end_date` - La campaña llegó a su fecha límite
- `manual` - Cierre manual por el administrador (requiere `reason`)

**Response (200):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "campaign_id": "123e4567-e89b-12d3-a456-426614174000",
    "closure_type": "manual",
    "closure_reason": "Campaña completada exitosamente antes del plazo previsto",
    "closed_by": "789e0123-e45b-67d8-a901-234567890abc",
    "total_raised": 15000.00,
    "total_donors": 45,
    "total_donations": 52,
    "campaign_goal": 20000.00,
    "goal_percentage": 75.00,
    "total_expenses": 12000.00,
    "total_receipts": 8,
    "receipts_with_documents": 6,
    "total_activities": 5,
    "transparency_score": 78.5,
    "transparency_breakdown": {
        "documentation_score": 22.5,
        "activity_score": 20.0,
        "goal_progress_score": 15.0,
        "timeliness_score": 12.0,
        "alerts_deduction_score": 0,
        "bonus_score": 9.0
    },
    "alerts_count": 0,
    "alerts_resolved": 0,
    "report_pdf_url": null,
    "closed_at": "2024-01-15T10:30:00Z",
    "created_at": "2024-01-15T10:30:00Z"
}
```

**Errores:**
- `400` - Tipo de cierre inválido o razón faltante para cierre manual
- `400` - Campaña no está en estado `active` o `paused`
- `400` - Campaña ya tiene un reporte de cierre
- `401` - No autenticado
- `403` - No tiene rol de admin

---

### 2. Obtener Reporte de Cierre (Admin)

**GET** `/api/campaigns/:id/closure-report`

Obtiene el reporte completo de cierre con todas las métricas.

**Autenticación:** Requiere rol `admin`

**Response (200):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "campaign_id": "123e4567-e89b-12d3-a456-426614174000",
    "closure_type": "goal_reached",
    "total_raised": 20000.00,
    "total_donors": 78,
    "total_donations": 95,
    "campaign_goal": 20000.00,
    "goal_percentage": 100.00,
    "total_expenses": 18500.00,
    "total_receipts": 15,
    "receipts_with_documents": 14,
    "total_activities": 12,
    "transparency_score": 92.3,
    "transparency_breakdown": {
        "documentation_score": 28.0,
        "activity_score": 25.0,
        "goal_progress_score": 20.0,
        "timeliness_score": 15.0,
        "alerts_deduction_score": 0,
        "bonus_score": 4.3
    },
    "report_pdf_url": "https://bucket.s3.amazonaws.com/audits/123e.../audit-report-1705312200.pdf",
    "report_hash": "a1b2c3d4e5f6...",
    "closed_at": "2024-01-15T10:30:00Z"
}
```

**Errores:**
- `404` - No existe reporte de cierre para esta campaña
- `401` - No autenticado
- `403` - No tiene rol de admin

---

### 3. Obtener Reporte Público de Auditoría (Público)

**GET** `/api/campaigns/:id/audit`

Obtiene una versión simplificada del reporte para donantes.

**Autenticación:** No requiere

**Response (200):**
```json
{
    "campaign_id": "123e4567-e89b-12d3-a456-426614174000",
    "campaign_title": "Ayuda para Hospital Infantil",
    "organizer_name": "Fundación Esperanza",
    "closed_at": "2024-01-15T10:30:00Z",
    "total_raised": 20000.00,
    "campaign_goal": 20000.00,
    "goal_percentage": 100.00,
    "total_donors": 78,
    "total_expenses": 18500.00,
    "transparency_score": 92.3,
    "report_pdf_url": "https://bucket.s3.amazonaws.com/audits/123e.../audit-report.pdf"
}
```

**Errores:**
- `404` - No existe reporte de auditoría para esta campaña

---

### 4. Descargar PDF de Auditoría (Público)

**GET** `/api/campaigns/:id/audit/download`

Redirige al PDF de auditoría para descarga.

**Autenticación:** No requiere

**Response (302):** Redirección al URL del PDF en S3

**Errores:**
- `404` - No existe reporte o el PDF aún está siendo generado

---

## Puntuación de Transparencia

La puntuación de transparencia (0-100) se calcula automáticamente basándose en múltiples criterios:

### Desglose de Puntuación

| Criterio | Puntos Máx | Descripción |
|----------|------------|-------------|
| **Documentación** | 30 | % de receipts con documento adjunto |
| **Actividades** | 25 | Ratio de actividades registradas vs meses de campaña |
| **Progreso Meta** | 20 | % de la meta alcanzada |
| **Puntualidad** | 15 | Frecuencia promedio de actualizaciones |
| **Alertas** | -10 | Deducción por alertas no resueltas |
| **Bonus** | 10 | Contrato firmado, >10 donantes, gastos documentados |

### Fórmulas Detalladas

**1. Documentación (0-30 pts)**
```
Si hay receipts: (receipts_con_documento / total_receipts) * 30
Si no hay receipts: 0 pts
```

**2. Actividades (0-25 pts)**
```
meses_campaña = ceil((end_date - start_date) / 30 días)
actividades_esperadas = meses_campaña
ratio = total_actividades / actividades_esperadas
Si ratio > 1: bonus de hasta 50% extra
Puntos = min(ratio * 25 / 1.5, 25)
```

**3. Progreso Meta (0-20 pts)**
```
ratio = total_recaudado / meta
Si ratio >= 100%: 20 pts
Si ratio >= 75%: 15 + (ratio - 0.75) * 20
Si ratio >= 50%: 10 + (ratio - 0.50) * 20
Si ratio < 50%: ratio * 20
```

**4. Puntualidad (0-15 pts)**
```
promedio_dias_entre_actividades:
  <= 7 días: 15 pts
  <= 14 días: 12 pts
  <= 30 días: 8 pts
  > 30 días: 5 pts
```

**5. Alertas (0 a -10 pts)**
```
alertas_sin_resolver = alerts_count - alerts_resolved
deducción = alertas_sin_resolver * -2
Máximo: -10 pts
```

**6. Bonus (0-10 pts)**
```
+ 3 pts: Tiene contrato firmado
+ 2 pts: Más de 10 donantes
+ 3 pts: Gastos documentados >= 80% de lo recaudado
+ 2 pts: Cerrada antes de fecha límite (cierre manual)
```

### Interpretación del Score

| Rango | Calificación |
|-------|--------------|
| 80-100 | Excelente |
| 60-79 | Bueno |
| 40-59 | Regular |
| 0-39 | Necesita mejoras |

---

## Flujo de Cierre

```
┌─────────────────┐
│  Campaña Activa │
│   o Pausada     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ POST /close     │◄── Admin ejecuta cierre
│ con closure_type│    con tipo y razón (si manual)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Validaciones:   │
│ - Estado válido │
│ - No cerrada    │
│ - Razón (manual)│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Recopilar       │
│ Métricas:       │
│ - Donaciones    │
│ - Receipts      │
│ - Actividades   │
│ - Alertas       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Calcular Score  │
│ Transparencia   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Guardar Reporte │
│ en BD           │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Actualizar      │
│ Campaña a       │
│ "completed"     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Generar PDF     │◄── Asíncrono (goroutine)
│ y subir a S3    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Donaciones      │
│ BLOQUEADAS      │
└─────────────────┘
```

---

## Ejemplos de Uso

### Cerrar campaña por meta alcanzada

```bash
curl -X POST http://localhost:9999/api/campaigns/123e4567-e89b-12d3-a456-426614174000/close \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "closure_type": "goal_reached"
  }'
```

### Cerrar campaña manualmente con justificación

```bash
curl -X POST http://localhost:9999/api/campaigns/123e4567-e89b-12d3-a456-426614174000/close \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "closure_type": "manual",
    "reason": "El beneficiario ha recibido tratamiento alternativo y ya no requiere los fondos adicionales. Los fondos recaudados serán utilizados para cubrir gastos médicos parciales."
  }'
```

### Consultar reporte público (como donante)

```bash
curl http://localhost:9999/api/campaigns/123e4567-e89b-12d3-a456-426614174000/audit
```

### Descargar PDF de auditoría

```bash
curl -L http://localhost:9999/api/campaigns/123e4567-e89b-12d3-a456-426614174000/audit/download \
  -o reporte-auditoria.pdf
```

### Intentar donar a campaña cerrada (error esperado)

```bash
curl -X POST http://localhost:9999/api/campaigns/123e4567-e89b-12d3-a456-426614174000/donations \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100,
    "is_anonymous": true,
    "payment_method_id": 1
  }'

# Response: 400 Bad Request
# {"error": "campaign is not accepting donations: campaign has been closed"}
```

---

## Notas Importantes

1. **PDF Asíncrono**: El PDF se genera de forma asíncrona después del cierre. Puede tardar unos segundos en estar disponible.

2. **Estado Terminal**: Una vez que una campaña está en estado `completed`, no puede volver a estados anteriores.

3. **Alertas Placeholder**: La tabla `campaign_alerts` está preparada para futura implementación. Por ahora, `alerts_count` siempre será 0.

4. **Bloqueo de Donaciones**: Las donaciones se bloquean automáticamente cuando la campaña pasa a estado `completed`.

5. **Un Reporte por Campaña**: Solo puede existir un reporte de cierre por campaña (constraint `UNIQUE` en `campaign_id`).
