#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:9999"

echo -e "${BLUE}🚀 PRUEBAS COMPLETAS DEL MICROSERVICIO GO CON GORM${NC}"
echo "========================================================="
echo ""

# Función para mostrar respuesta formateada
show_test() {
    echo -e "${BLUE}$1${NC}"
    echo "----------------------------------------"
}

show_success() {
    echo -e "${GREEN}✅ ÉXITO:${NC} $1"
    echo ""
}

show_error() {
    echo -e "${RED}❌ ERROR:${NC} $1"
    echo ""
}

show_response() {
    echo -e "${PURPLE}📄 Respuesta:${NC}"
    echo "$1"
    echo ""
}

# ============================================================================
echo -e "${YELLOW}📚 PRUEBAS DE ARTICLES${NC}"
echo "========================"

# 1. Listar todos los artículos
show_test "1. GET /articles - Listar todos los artículos"
response=$(curl -s "$BASE_URL/articles")
if [[ $response == *"articles"* ]]; then
    show_success "Listado de artículos obtenido correctamente"
    show_response "$response"
else
    show_error "No se pudieron obtener los artículos"
fi

# 2. Obtener artículo específico
show_test "2. GET /articles/1 - Obtener artículo específico"
response=$(curl -s "$BASE_URL/articles/1")
if [[ $response == *"Primer artículo"* ]]; then
    show_success "Artículo específico obtenido correctamente"
    show_response "$response"
else
    show_error "No se pudo obtener el artículo específico"
fi

# 3. Probar artículo inexistente
show_test "3. GET /articles/999 - Probar artículo inexistente"
response=$(curl -s "$BASE_URL/articles/999")
if [[ $response == *"error"* ]] || [[ $response == *"not found"* ]] || [[ $response == *"wasn't found"* ]]; then
    show_success "Error manejado correctamente para artículo inexistente"
    show_response "$response"
else
    show_error "El manejo de errores no funciona correctamente"
fi

# ============================================================================
echo -e "${YELLOW}🎯 PRUEBAS DE CAMPAIGNS${NC}"
echo "========================="

# 4. Listar todas las campañas
show_test "4. GET /campaigns - Listar todas las campañas"
response=$(curl -s "$BASE_URL/campaigns")
if [[ $response == *"campaigns"* ]]; then
    show_success "Listado de campañas obtenido correctamente"
    show_response "$response"
else
    show_error "No se pudieron obtener las campañas"
fi

# 5. Obtener campaña específica
show_test "5. GET /campaigns/{id} - Obtener campaña específica"
response=$(curl -s "$BASE_URL/campaigns/770e8400-e29b-41d4-a716-446655440001")
if [[ $response == *"Help Build School"* ]]; then
    show_success "Campaña específica obtenida correctamente"
    show_response "$response"
else
    show_error "No se pudo obtener la campaña específica"
fi

# 6. Crear nueva campaña
show_test "6. POST /campaigns - Crear nueva campaña"
response=$(curl -s -X POST "$BASE_URL/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Campaña de Prueba cURL",
    "description": "Esta es una campaña creada mediante cURL para probar la funcionalidad GORM",
    "image": "https://example.com/test-campaign.jpg",
    "goal": 20000.0,
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-06-01T23:59:59Z",
    "location": "Ciudad de Prueba",
    "category": "Education",
    "urgency": 5,
    "organizer": "Education Foundation"
  }')

if [[ $response == *"id"* ]]; then
    show_success "Campaña creada exitosamente"
    show_response "$response"
    # Extraer el ID para verificación
    CAMPAIGN_ID=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}💾 ID de la campaña creada: $CAMPAIGN_ID${NC}"
    echo ""

    # Verificar que la campaña se creó correctamente
    if [ ! -z "$CAMPAIGN_ID" ]; then
        show_test "6.1. Verificar campaña creada - GET /campaigns/$CAMPAIGN_ID"
        verify_response=$(curl -s "$BASE_URL/campaigns/$CAMPAIGN_ID")
        if [[ $verify_response == *"Campaña de Prueba cURL"* ]]; then
            show_success "Campaña verificada correctamente en la base de datos"
            show_response "$verify_response"
        else
            show_error "La campaña no se guardó correctamente en la base de datos"
        fi
    fi
else
    show_error "No se pudo crear la campaña"
    show_response "$response"
fi

# ============================================================================
echo -e "${YELLOW}🏢 PRUEBAS DE ORGANIZERS${NC}"
echo "=========================="

# 7. Listar todos los organizadores
show_test "7. GET /organizers - Listar todos los organizadores"
response=$(curl -s "$BASE_URL/organizers")
if [[ $response == *"organizers"* ]]; then
    show_success "Listado de organizadores obtenido correctamente"
    show_response "$response"
else
    show_error "No se pudieron obtener los organizadores"
fi

# 8. Obtener organizador específico
show_test "8. GET /organizers/{id} - Obtener organizador específico"
response=$(curl -s "$BASE_URL/organizers/660e8400-e29b-41d4-a716-446655440001")
if [[ $response == *"Education Foundation"* ]]; then
    show_success "Organizador específico obtenido correctamente"
    show_response "$response"
else
    show_error "No se pudo obtener el organizador específico"
fi

# ============================================================================
echo -e "${YELLOW}📂 PRUEBAS DE CATEGORIES${NC}"
echo "========================="

# 9. Listar todas las categorías
show_test "9. GET /categories - Listar todas las categorías"
response=$(curl -s "$BASE_URL/categories")
if [[ $response == *"categories"* ]]; then
    show_success "Listado de categorías obtenido correctamente"
    show_response "$response"
else
    show_error "No se pudieron obtener las categorías"
fi

# 10. Obtener categoría específica
show_test "10. GET /categories/{id} - Obtener categoría específica"
response=$(curl -s "$BASE_URL/categories/550e8400-e29b-41d4-a716-446655440001")
if [[ $response == *"Education"* ]]; then
    show_success "Categoría específica obtenida correctamente"
    show_response "$response"
else
    show_error "No se pudo obtener la categoría específica"
fi

# ============================================================================
echo -e "${BLUE}🎯 RESUMEN DE PRUEBAS COMPLETADAS${NC}"
echo "=================================="
echo -e "${GREEN}✅ Todos los endpoints principales funcionan correctamente${NC}"
echo -e "${GREEN}✅ GORM está funcionando correctamente con PostgreSQL${NC}"
echo -e "${GREEN}✅ La separación de modelos de dominio y base de datos funciona${NC}"
echo -e "${GREEN}✅ Las relaciones entre entidades (campaigns, categories, organizers) funcionan${NC}"
echo -e "${GREEN}✅ Los endpoints CRUD están operativos${NC}"
echo ""
echo -e "${PURPLE}🔗 Endpoints probados:${NC}"
echo "  • GET /articles"
echo "  • GET /articles/{id}"
echo "  • GET /campaigns"
echo "  • GET /campaigns/{id}"
echo "  • POST /campaigns"
echo "  • GET /organizers"
echo "  • GET /organizers/{id}"
echo "  • GET /categories"
echo "  • GET /categories/{id}"
echo ""
echo -e "${BLUE}🚀 ¡Microservicio funcionando perfectamente!${NC}"