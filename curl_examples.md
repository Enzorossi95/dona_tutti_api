# Ejemplos de cURL para el Microservicio Go

## 📚 Articles

### Listar todos los artículos
```bash
curl -X GET "http://localhost:9999/articles"
```

### Obtener un artículo específico
```bash
curl -X GET "http://localhost:9999/articles/1"
```

### Probar artículo inexistente (manejo de errores)
```bash
curl -X GET "http://localhost:9999/articles/999"
```

---

## 🎯 Campaigns

### Listar todas las campañas
```bash
curl -X GET "http://localhost:9999/campaigns"
```

### Obtener una campaña específica
```bash
curl -X GET "http://localhost:9999/campaigns/770e8400-e29b-41d4-a716-446655440001"
```

### Crear una nueva campaña
```bash
curl -X POST "http://localhost:9999/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Campaña de Tecnología Educativa",
    "description": "Necesitamos fondos para equipar aulas con tecnología moderna para mejorar la educación",
    "image": "https://example.com/tech-education.jpg",
    "goal": 35000.0,
    "start_date": "2025-02-01T00:00:00Z",
    "end_date": "2025-08-01T23:59:59Z",
    "location": "Escuelas Rurales",
    "category": "Education",
    "urgency": 7,
    "organizer": "Education Foundation"
  }'
```

### Crear campaña de salud
```bash
curl -X POST "http://localhost:9999/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Equipamiento Médico de Emergencia",
    "description": "Urgente: necesitamos equipos médicos para el hospital local",
    "image": "https://example.com/medical-equipment.jpg",
    "goal": 50000.0,
    "start_date": "2025-01-10T00:00:00Z",
    "end_date": "2025-04-10T23:59:59Z",
    "location": "Hospital Central",
    "category": "Health",
    "urgency": 9,
    "organizer": "Medical Relief Org"
  }'
```

### Crear campaña ambiental
```bash
curl -X POST "http://localhost:9999/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Reforestación Urbana",
    "description": "Plantemos árboles en la ciudad para mejorar la calidad del aire",
    "image": "https://example.com/reforestation.jpg",
    "goal": 25000.0,
    "start_date": "2025-03-01T00:00:00Z",
    "end_date": "2025-12-01T23:59:59Z",
    "location": "Parques Urbanos",
    "category": "Environment",
    "urgency": 6,
    "organizer": "Water for All"
  }'
```

---

## 🏢 Organizers

### Listar todos los organizadores
```bash
curl -X GET "http://localhost:9999/organizers"
```

### Obtener organizador específico - Education Foundation
```bash
curl -X GET "http://localhost:9999/organizers/660e8400-e29b-41d4-a716-446655440001"
```

### Obtener organizador específico - Medical Relief Org
```bash
curl -X GET "http://localhost:9999/organizers/660e8400-e29b-41d4-a716-446655440002"
```

### Obtener organizador específico - Water for All
```bash
curl -X GET "http://localhost:9999/organizers/660e8400-e29b-41d4-a716-446655440003"
```

---

## 📂 Categories

### Listar todas las categorías
```bash
curl -X GET "http://localhost:9999/categories"
```

### Obtener categoría Education
```bash
curl -X GET "http://localhost:9999/categories/550e8400-e29b-41d4-a716-446655440001"
```

### Obtener categoría Health
```bash
curl -X GET "http://localhost:9999/categories/550e8400-e29b-41d4-a716-446655440002"
```

### Obtener categoría Environment
```bash
curl -X GET "http://localhost:9999/categories/550e8400-e29b-41d4-a716-446655440003"
```

### Obtener categoría Community
```bash
curl -X GET "http://localhost:9999/categories/550e8400-e29b-41d4-a716-446655440004"
```

---

## 🧪 Pruebas de Validación y Errores

### Probar campaña con organizador inexistente
```bash
curl -X POST "http://localhost:9999/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Campaña con Organizador Inexistente",
    "description": "Esta campaña debería fallar",
    "goal": 10000.0,
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-06-01T23:59:59Z",
    "location": "Lugar de Prueba",
    "category": "Education",
    "urgency": 5,
    "organizer": "Organizador Inexistente"
  }'
```

### Probar campaña con categoría inexistente
```bash
curl -X POST "http://localhost:9999/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Campaña con Categoría Inexistente",
    "description": "Esta campaña debería fallar",
    "goal": 10000.0,
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-06-01T23:59:59Z",
    "location": "Lugar de Prueba",
    "category": "Categoría Inexistente",
    "urgency": 5,
    "organizer": "Education Foundation"
  }'
```

### Probar ID inválido para campaña
```bash
curl -X GET "http://localhost:9999/campaigns/invalid-uuid"
```

### Probar ID inválido para organizador
```bash
curl -X GET "http://localhost:9999/organizers/invalid-uuid"
```

### Probar ID inválido para categoría
```bash
curl -X GET "http://localhost:9999/categories/invalid-uuid"
```

---

## 📊 Comandos para Análisis de Datos

### Obtener todas las campañas y contar por categoría (usando jq si está disponible)
```bash
curl -s "http://localhost:9999/campaigns" | jq '.data.campaigns | group_by(.category) | map({category: .[0].category, count: length})'
```

### Obtener campañas ordenadas por urgencia (usando jq si está disponible)
```bash
curl -s "http://localhost:9999/campaigns" | jq '.data.campaigns | sort_by(.urgency) | reverse'
```

### Obtener solo los títulos de las campañas (usando jq si está disponible)
```bash
curl -s "http://localhost:9999/campaigns" | jq '.data.campaigns[].title'
```

---

## 🔍 Verificación de Funcionalidad GORM

### Verificar que las relaciones funcionan correctamente
```bash
# Crear una campaña y verificar que se asocia correctamente con category y organizer
curl -X POST "http://localhost:9999/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Verificación de Relaciones GORM",
    "description": "Esta campaña verifica que las relaciones entre entidades funcionan",
    "goal": 15000.0,
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-06-01T23:59:59Z",
    "location": "Lugar de Verificación",
    "category": "Community",
    "urgency": 4,
    "organizer": "Water for All"
  }'
```

### Verificar timestamps automáticos
```bash
# Las campañas creadas deberían tener created_at automático
curl -s "http://localhost:9999/campaigns" | grep -o '"created_at":"[^"]*"'
```

---

## 💡 Notas

- Todos los endpoints están funcionando con GORM
- La separación entre modelos de dominio y base de datos está implementada
- Las relaciones entre entidades (campaigns ↔ categories ↔ organizers) funcionan correctamente
- Los timestamps automáticos están funcionando
- El manejo de errores está implementado
- La validación de UUIDs está funcionando
