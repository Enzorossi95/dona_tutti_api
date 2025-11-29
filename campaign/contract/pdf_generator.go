package contract

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// pdfGenerator implements the PDFGenerator interface
type pdfGenerator struct{}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator() PDFGenerator {
	return &pdfGenerator{}
}

// Generate creates a PDF contract and returns the PDF bytes and SHA256 hash
func (g *pdfGenerator) Generate(data ContractData) ([]byte, string, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add title
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(190, 10, "CONTRATO LEGAL SIMPLIFICADO", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add subtitle
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 8, "Plataforma de Donaciones Dona Tutti", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add date
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(190, 6, fmt.Sprintf("Fecha de generacion: %s", data.GeneratedAt.Format("02/01/2006 15:04:05")), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Add campaign information section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, "INFORMACION DE LA CAMPANA", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 6, fmt.Sprintf("Titulo de la campana: %s", data.CampaignTitle), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Objetivo de recaudacion: $%.2f", data.CampaignGoal), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("ID de la campana: %s", data.CampaignID.String()), "", "L", false)
	pdf.Ln(5)

	// Add organizer information section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, "INFORMACION DEL ORGANIZADOR", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 6, fmt.Sprintf("Nombre: %s", data.OrganizerName), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Email: %s", data.OrganizerEmail), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Telefono: %s", data.OrganizerPhone), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Direccion: %s", data.OrganizerAddress), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("ID del Organizador: %s", data.OrganizerID.String()), "", "L", false)
	pdf.Ln(10)

	// Add terms and conditions section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, "TERMINOS Y CONDICIONES", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	pdf.SetFont("Arial", "", 10)

	// Term 1
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "1. COMPROMISO DE VERACIDAD", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "El organizador declara bajo juramento que toda la informacion proporcionada en esta campana es veridica y precisa. Cualquier informacion falsa o enganosa podra resultar en la suspension inmediata de la campana y acciones legales correspondientes.", "", "J", false)
	pdf.Ln(3)

	// Term 2
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "2. USO DE FONDOS", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "El organizador se compromete a utilizar los fondos recaudados exclusivamente para el proposito descrito en la campana. Cualquier desviacion de fondos sera considerada fraude y sera reportada a las autoridades competentes.", "", "J", false)
	pdf.Ln(3)

	// Term 3
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "3. TRANSPARENCIA Y RENDICION DE CUENTAS", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "El organizador acepta proporcionar actualizaciones regulares sobre el progreso de la campana y el uso de los fondos. Al finalizar la campana, debera presentar un informe detallado de como se utilizaron los fondos recaudados.", "", "J", false)
	pdf.Ln(3)

	// Term 4
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "4. COMISIONES Y TARIFAS", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "El organizador reconoce y acepta que la plataforma Dona Tutti puede cobrar comisiones por los servicios prestados. Estas comisiones seran deducidas automaticamente de los fondos recaudados segun las politicas vigentes de la plataforma.", "", "J", false)
	pdf.Ln(3)

	// Term 5
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "5. PROCEDIMIENTO EN CASO DE DENUNCIA", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "En caso de recibir denuncias sobre la campana, el organizador acepta cooperar plenamente con la investigacion. La plataforma se reserva el derecho de suspender la campana y retener fondos hasta que se resuelva la investigacion. El organizador acepta que cualquier decision tomada por la plataforma en este contexto sera vinculante.", "", "J", false)
	pdf.Ln(3)

	// Term 6
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "6. PROPIEDAD INTELECTUAL", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "El organizador garantiza que todo el contenido publicado en la campana (imagenes, textos, videos) es de su propiedad o cuenta con los permisos necesarios para su uso. El organizador asume toda responsabilidad por cualquier violacion de derechos de autor o propiedad intelectual.", "", "J", false)
	pdf.Ln(3)

	// Term 7
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "7. PRIVACIDAD Y PROTECCION DE DATOS", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "El organizador acepta que sus datos personales seran procesados de acuerdo con la politica de privacidad de Dona Tutti y las leyes de proteccion de datos vigentes. El organizador consiente el uso de su informacion para fines relacionados con la gestion de la campana.", "", "J", false)
	pdf.Ln(3)

	// Term 8
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(190, 6, "8. RESPONSABILIDAD LEGAL", "", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "El organizador libera a Dona Tutti de cualquier responsabilidad legal derivada del contenido de la campana, el uso de los fondos, o cualquier disputa con donantes o terceros. El organizador es el unico responsable ante la ley por todas las acciones relacionadas con su campana.", "", "J", false)
	pdf.Ln(10)

	// Add acceptance section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, "DECLARACION DE ACEPTACION", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 6, "Al firmar digitalmente este documento, el organizador declara:", "", "L", false)
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(190, 5, "[ ] He leido y comprendido todos los terminos y condiciones de este contrato.", "", "L", false)
	pdf.MultiCell(190, 5, "[ ] Acepto cumplir con todas las obligaciones establecidas en este documento.", "", "L", false)
	pdf.MultiCell(190, 5, "[ ] Acepto las condiciones del sistema en caso de denuncia o investigacion.", "", "L", false)
	pdf.MultiCell(190, 5, "[ ] Comprendo que el incumplimiento de estos terminos puede resultar en acciones legales.", "", "L", false)
	pdf.Ln(10)

	// Add signature section
	pdf.SetFont("Arial", "I", 9)
	pdf.MultiCell(190, 5, "Este documento sera firmado digitalmente mediante la aceptacion en la plataforma Dona Tutti. La firma digital incluira la fecha, hora, direccion IP y metadatos del navegador del organizador como evidencia de aceptacion.", "", "J", false)
	pdf.Ln(5)

	// Add footer
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(190, 5, "Documento generado automaticamente por Dona Tutti", "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, fmt.Sprintf("Hash del documento: Se generara al momento de la aceptacion"), "", 1, "C", false, 0, "")

	// Generate PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	pdfBytes := buf.Bytes()

	// Calculate SHA256 hash
	hash := sha256.Sum256(pdfBytes)
	hashString := fmt.Sprintf("%x", hash)

	return pdfBytes, hashString, nil
}

