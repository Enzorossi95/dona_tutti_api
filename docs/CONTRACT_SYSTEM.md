# Campaign Legal Contract System

## Overview

Sistema completo de contratos legales para campañas con generación de PDF, firma digital simplificada y flujo de estados validado.

## Arquitectura

### Base de Datos

#### Nueva Tabla: `campaign_contracts`

```sql
CREATE TABLE campaign_contracts (
    id UUID PRIMARY KEY,
    campaign_id UUID NOT NULL UNIQUE REFERENCES campaigns(id),
    organizer_id UUID NOT NULL REFERENCES organizers(id),
    contract_pdf_url TEXT NOT NULL,
    contract_hash VARCHAR(64) NOT NULL,
    accepted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    acceptance_ip VARCHAR(45) NOT NULL,
    acceptance_user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### Actualización: Enum de Estados de Campaña

```sql
CREATE TYPE campaign_status AS ENUM (
    'draft',
    'pending_approval',
    'active',
    'paused',
    'completed',
    'rejected'
);
```

### Dominio: `campaign/contract/`

Estructura siguiendo Domain-Driven Design:

```
campaign/contract/
├── contract.go          # Entidades del dominio
├── model.go            # Modelos GORM (base de datos)
├── repository.go       # Capa de acceso a datos
├── service.go          # Lógica de negocio
├── handlers.go         # Handlers HTTP (Echo)
└── pdf_generator.go    # Generación de PDF
```

## Flujo de Trabajo

### 1. Crear Campaña (Estado: draft)

```bash
POST /api/campaigns
{
  "title": "Mi Campaña",
  "description": "Descripción...",
  "goal": 10000.00,
  ...
}
# Status inicial: "draft"
```

### 2. Generar Contrato PDF

```bash
POST /api/campaigns/{id}/contract/generate
# No requiere body - obtiene toda la información de la BD
# Genera PDF, sube a S3, retorna URL
```

### 3. Visualizar Contrato

```bash
GET /api/campaigns/{id}/contract
# Retorna contrato con URL del PDF y metadata
```

### 4. Aceptar Contrato (Firma Digital)

```bash
POST /api/campaigns/{id}/contract/accept
{
  "organizer_id": "uuid..."
}
# Registra: IP, timestamp, user agent, hash del PDF
# Cambia status: draft → pending_approval
```

### 5. Admin Revisa Comprobante

```bash
GET /api/campaigns/{id}/contract/proof
# Vista para admin con todos los detalles del contrato
```

### 6. Admin Aprueba Campaña

```bash
PUT /api/campaigns/{id}
{
  "status": "active"
}
# Transición: pending_approval → active
```

## Endpoints API

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/campaigns/:id/contract/generate` | Generar contrato PDF |
| GET | `/api/campaigns/:id/contract` | Ver contrato generado |
| POST | `/api/campaigns/:id/contract/accept` | Firmar/aceptar contrato |
| GET | `/api/campaigns/:id/contract/proof` | Ver comprobante (admin) |

## Firma Digital Simplificada

No requiere PKI ni certificados digitales complejos. La "firma digital" consiste en:

1. **Timestamp**: Fecha y hora exacta de aceptación
2. **IP Address**: Dirección IP del organizador
3. **User Agent**: Navegador/dispositivo utilizado
4. **Hash SHA256**: Hash del documento PDF aceptado

Este método es **legalmente válido** para términos de servicio y proporciona evidencia suficiente de aceptación.

## Estados de Campaña

### Transiciones Válidas

```
draft → pending_approval → active → paused
                         ↘         ↗
                          completed

draft → rejected (terminal)
pending_approval → rejected (terminal)
```

### Validaciones

- ✅ No se puede crear campaña directamente como `active`
- ✅ Para pasar a `pending_approval` se requiere contrato firmado
- ✅ Solo admin puede aprobar: `pending_approval` → `active`
- ✅ No se puede cambiar estado de campañas en `completed` o `rejected`

## Contenido del Contrato

El PDF generado incluye:

### 1. Información de la Campaña
- Título
- Objetivo de recaudación
- ID único

### 2. Información del Organizador
- Nombre completo
- Email
- Teléfono
- Dirección
- ID único

### 3. Términos y Condiciones

1. **Compromiso de Veracidad**
2. **Uso de Fondos**
3. **Transparencia y Rendición de Cuentas**
4. **Comisiones y Tarifas**
5. **Procedimiento en Caso de Denuncia**
6. **Propiedad Intelectual**
7. **Privacidad y Protección de Datos**
8. **Responsabilidad Legal**

### 4. Declaración de Aceptación

Checkboxes virtuales que el organizador acepta al firmar:
- He leído y comprendido todos los términos
- Acepto cumplir con todas las obligaciones
- Acepto las condiciones del sistema en caso de denuncia
- Comprendo las consecuencias legales del incumplimiento

## Integración S3

Los contratos PDF se almacenan en S3 con la siguiente estructura:

```
s3://bucket-name/
  contracts/
    {campaign-id}/
      contract-{timestamp}.pdf
```

**URL Pública**: Se genera automáticamente según el entorno:
- **AWS**: `https://bucket.s3.amazonaws.com/contracts/{campaign-id}/contract-{timestamp}.pdf`
- **LocalStack**: `http://localhost:4566/bucket/contracts/{campaign-id}/contract-{timestamp}.pdf`

## Middleware de Validación

### `ValidateStatusValue()`

Valida que el valor de status sea uno de los permitidos.

### `RequireContractForApproval()`

Verifica que exista un contrato firmado antes de aprobar una campaña.

### `ValidateStatusTransition()`

Placeholder para validaciones adicionales de transiciones de estado.

## Testing

### Script de Pruebas

```bash
chmod +x tests/test_contracts.sh
./tests/test_contracts.sh
```

El script prueba:

1. ✅ Autenticación de usuario
2. ✅ Creación de campaña en estado draft
3. ✅ Generación de contrato PDF (obtiene datos de BD)
4. ✅ Visualización del contrato
5. ✅ Aceptación del contrato (firma digital)
6. ✅ Verificación de estado pending_approval
7. ✅ Visualización del comprobante legal (admin)
8. ✅ Validaciones negativas

**Nota**: El endpoint de generación fue simplificado y ya no requiere request body.

## Dependencias

### Nueva Dependencia: gofpdf

```bash
go get github.com/jung-kurt/gofpdf
go mod tidy
```

Esta librería se usa para generar PDFs en Go de manera sencilla.

## Configuración

No requiere configuración adicional. El sistema utiliza:

- **S3 Client**: Ya configurado en el sistema
- **Database**: Usa la conexión GORM existente
- **Auth**: Integrado con el sistema JWT actual

## Consideraciones de Seguridad

### ✅ Implementado

- Hash SHA256 del documento para verificar integridad
- Registro de IP para auditoría
- Timestamps inmutables
- Validación de transiciones de estado
- Relación única campaign_id → contract (un contrato por campaña)

### 🔄 Para Producción

- Rate limiting en endpoints de generación de contratos
- Validación adicional de IP (geolocalización, VPN detection)
- Backup automático de contratos en múltiples regiones
- Firma con timestamp authority (TSA) para mayor validez legal
- Encriptación adicional de datos sensibles en reposo

## Mejoras Futuras

1. **Versionamiento de Contratos**: Mantener histórico si cambian términos legales
2. **Notificaciones**: Email automático al generar/firmar contrato
3. **Recordatorios**: Sistema de recordatorios si no se firma en X días
4. **Analytics**: Dashboard de contratos pendientes/firmados
5. **Firma Manuscrita**: Canvas para dibujar firma (opcional)
6. **Multi-idioma**: Generar contratos en diferentes idiomas
7. **Templates Personalizados**: Diferentes tipos de contratos según categoría

## Troubleshooting

### Error: "S3 client not initialized"

**Solución**: Configurar variables de entorno AWS:
```bash
export AWS_REGION=us-east-1
export AWS_S3_BUCKET=dona-tutti-files
export AWS_ACCESS_KEY_ID=your-key
export AWS_SECRET_ACCESS_KEY=your-secret
```

### Error: "Contract not found"

**Causa**: Intentando aceptar un contrato que no ha sido generado.

**Solución**: Seguir el flujo correcto:
1. Generar contrato primero (`/contract/generate`)
2. Luego aceptarlo (`/contract/accept`)

### Error: "Invalid status transition"

**Causa**: Intentando cambiar a un estado no permitido.

**Solución**: Verificar las transiciones válidas en la sección "Estados de Campaña".

## Autor y Mantenimiento

- **Versión**: 1.0.0
- **Fecha**: Enero 2025
- **Migración**: `20250129000000_campaign_contracts_system.sql`

---

Para más información, consultar el código fuente en `campaign/contract/` o los tests en `tests/test_contracts.sh`.

