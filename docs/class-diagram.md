# Diagrama de Clases - Dona Tutti API

## Descripción General
Este diagrama representa la estructura de clases del sistema de donaciones "Dona Tutti", mostrando las entidades principales y sus relaciones.

## Diagrama de Clases

```mermaid
classDiagram
    class User {
        +UUID id
        +String email
        +String passwordHash
        +UUID roleId
        +String firstName
        +String lastName
        +Boolean isActive
        +Boolean isVerified
        +String resetToken
        +DateTime resetTokenExpires
        +DateTime lastLogin
        +DateTime createdAt
        +DateTime updatedAt
        +login()
        +logout()
        +resetPassword()
        +updateProfile()
    }

    class Role {
        +UUID id
        +String name
        +String description
        +Boolean isActive
        +DateTime createdAt
        +DateTime updatedAt
    }

    class Permission {
        +UUID id
        +String name
        +String resource
        +String action
        +String description
        +DateTime createdAt
    }

    class Organizer {
        +UUID id
        +UUID userId
        +String name
        +String avatar
        +Boolean verified
        +String email
        +String phone
        +String website
        +String address
        +DateTime createdAt
        +createCampaign()
        +updateProfile()
        +getStatistics()
    }

    class Campaign {
        +UUID id
        +String title
        +String description
        +String image
        +Float64 goal
        +DateTime startDate
        +DateTime endDate
        +String location
        +Integer urgency
        +String status
        +UUID categoryId
        +UUID organizerId
        +String beneficiaryName
        +Integer beneficiaryAge
        +String currentSituation
        +String urgencyReason
        +DateTime createdAt
        +DateTime updatedAt
        +calculateProgress()
        +updateStatus()
        +getTotalDonations()
    }

    class CampaignCategory {
        +UUID id
        +String name
        +String description
        +DateTime createdAt
        +getCampaigns()
    }

    class Activity {
        +UUID id
        +UUID campaignId
        +String title
        +String description
        +DateTime date
        +String type
        +String author
        +DateTime createdAt
        +DateTime updatedAt
        +create()
        +update()
        +getActivitiesByCampaign()
    }

    class Receipt {
        +UUID id
        +UUID campaignId
        +String provider
        +String name
        +String description
        +Float64 total
        +Integer quantity
        +DateTime date
        +String documentURL
        +String note
        +DateTime createdAt
        +DateTime updatedAt
        +upload()
        +validateReceipt()
        +getTotalExpenses()
    }

    class Donation {
        +UUID id
        +UUID campaignId
        +UUID donorId
        +Float64 amount
        +DateTime date
        +String message
        +Boolean isAnonymous
        +Integer paymentMethodId
        +String status
        +DateTime createdAt
        +DateTime updatedAt
        +processDonation()
        +refund()
        +updateStatus()
    }

    class Donor {
        +UUID id
        +String firstName
        +String lastName
        +String email
        +String phone
        +Boolean isVerified
        +DateTime createdAt
        +DateTime updatedAt
        +makeDonation()
        +getDonationHistory()
        +updateProfile()
    }

    class PaymentMethod {
        +Integer id
        +String code
        +String name
        +Boolean isActive
        +DateTime createdAt
        +activate()
        +deactivate()
    }

    class CampaignPaymentMethod {
        +Integer id
        +UUID campaignId
        +Integer paymentMethodId
        +String instructions
        +Boolean isActive
        +DateTime createdAt
        +DateTime updatedAt
        +addTransferDetails()
        +addCashLocation()
    }

    class TransferDetail {
        +Integer id
        +Integer campaignPaymentMethodId
        +String bankName
        +String accountHolder
        +String cbu
        +String alias
        +String swiftCode
        +String additionalNotes
        +DateTime createdAt
        +DateTime updatedAt
        +validateBankDetails()
    }

    class CashLocation {
        +Integer id
        +Integer campaignPaymentMethodId
        +String locationName
        +String address
        +String contactInfo
        +String availableHours
        +String additionalNotes
        +DateTime createdAt
        +DateTime updatedAt
        +validateLocation()
    }

    class DonationStatus {
        <<enumeration>>
        COMPLETED
        PENDING
        FAILED
        REFUNDED
    }

    class CampaignStatus {
        <<enumeration>>
        ACTIVE
        INACTIVE
        COMPLETED
        CANCELLED
    }

    %% Relationships

    User "1" --> "1" Role : has
    Role "1" --> "*" Permission : has
    User "1" --> "*" Organizer : creates

    Organizer "1" --> "*" Campaign : manages
    Campaign "*" --> "1" CampaignCategory : belongs to
    Campaign "1" --> "*" Activity : has
    Campaign "1" --> "*" Receipt : has
    Campaign "1" --> "*" Donation : receives
    Campaign "1" --> "*" CampaignPaymentMethod : accepts

    Donation "*" --> "1" Donor : made by
    Donation "*" --> "1" PaymentMethod : uses
    Donation --> DonationStatus : has status

    CampaignPaymentMethod "*" --> "1" PaymentMethod : uses
    CampaignPaymentMethod "1" --> "*" TransferDetail : has
    CampaignPaymentMethod "1" --> "*" CashLocation : has

    Campaign --> CampaignStatus : has status
```

## Descripción de las Relaciones

### Relaciones Principales

1. **User - Role - Permission**: Sistema RBAC
   - Un usuario tiene un rol
   - Un rol tiene múltiples permisos

2. **User - Organizer**:
   - Un usuario puede crear múltiples organizadores
   - Un organizador pertenece a un usuario

3. **Organizer - Campaign**:
   - Un organizador puede gestionar múltiples campañas
   - Una campaña pertenece a un organizador

4. **Campaign - CampaignCategory**:
   - Una campaña pertenece a una categoría
   - Una categoría puede tener múltiples campañas

5. **Campaign - Activity**:
   - Una campaña puede tener múltiples actividades
   - Una actividad pertenece a una campaña

6. **Campaign - Receipt**:
   - Una campaña puede tener múltiples recibos/comprobantes
   - Un recibo pertenece a una campaña

7. **Campaign - Donation**:
   - Una campaña puede recibir múltiples donaciones
   - Una donación se hace a una campaña específica

8. **Donation - Donor**:
   - Un donador puede hacer múltiples donaciones
   - Una donación es hecha por un donador

9. **Donation - PaymentMethod**:
   - Una donación usa un método de pago
   - Un método de pago puede usarse en múltiples donaciones

10. **Campaign - CampaignPaymentMethod - PaymentMethod**:
    - Una campaña puede aceptar múltiples métodos de pago
    - La relación se gestiona a través de CampaignPaymentMethod

11. **CampaignPaymentMethod - TransferDetail/CashLocation**:
    - Un método de pago de campaña puede tener detalles de transferencia
    - Un método de pago de campaña puede tener ubicaciones para pago en efectivo

## Notas de Implementación

### Tipos de Datos
- **UUID**: Usado para IDs principales (User, Campaign, Donation, etc.)
- **Integer**: Usado para PaymentMethod y entidades relacionadas
- **Float64**: Usado para montos monetarios
- **DateTime**: Timestamps con createdAt/updatedAt automáticos
- **Boolean**: Flags de estado (isActive, isVerified, etc.)

### Patrones de Diseño
- **Repository Pattern**: Capa de acceso a datos
- **Service Layer**: Lógica de negocio
- **Clean Architecture**: Separación entre modelos de DB y entidades de dominio
- **DTO Pattern**: Objetos de transferencia para request/response

### Características Especiales
- **Soft Delete**: Manejo de eliminación lógica con flags de estado
- **Auditoría**: Campos createdAt/updatedAt automáticos
- **Validación**: Validaciones a nivel de modelo y servicio
- **Autorización**: Sistema RBAC integrado con JWT

## Flujos Principales

1. **Creación de Campaña**:
   User → Organizer → Campaign → CampaignPaymentMethod

2. **Proceso de Donación**:
   Donor → Donation → Campaign (con PaymentMethod)

3. **Gestión de Actividades**:
   Campaign → Activity (registro de eventos/actualizaciones)

4. **Control de Gastos**:
   Campaign → Receipt (comprobantes de gastos)

## Tecnologías Utilizadas
- **Framework**: Echo v4 (Go)
- **ORM**: GORM
- **Base de Datos**: PostgreSQL 15
- **Autenticación**: JWT
- **Documentación**: Swagger/OpenAPI