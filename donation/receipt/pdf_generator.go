package receipt

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

// ReceiptData contains all the information needed to generate a donation receipt
type ReceiptData struct {
	DonationID    uuid.UUID
	CampaignTitle string
	DonorName     string
	Amount        float64
	Date          time.Time
	PaymentMethod string
	IsAnonymous   bool
}

// PDFGenerator defines the interface for generating PDF receipts
type PDFGenerator interface {
	Generate(data ReceiptData) ([]byte, error)
}

// pdfGenerator implements the PDFGenerator interface
type pdfGenerator struct{}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator() PDFGenerator {
	return &pdfGenerator{}
}

// Generate creates a PDF receipt and returns the PDF bytes
func (g *pdfGenerator) Generate(data ReceiptData) ([]byte, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add header with logo/branding
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(41, 128, 185) // Blue color
	pdf.CellFormat(190, 15, "DONA TUTTI", "", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(190, 8, "Plataforma de Donaciones", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add receipt title
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(190, 10, "COMPROBANTE DE DONACION", "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// Add horizontal line
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(8)

	// Receipt details section
	pdf.SetFont("Arial", "", 11)
	
	// Receipt number
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(70, 7, "Numero de Recibo:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(120, 7, data.DonationID.String())
	pdf.Ln(7)

	// Date
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(70, 7, "Fecha y Hora:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(120, 7, data.Date.Format("02/01/2006 15:04:05"))
	pdf.Ln(10)

	// Add section separator
	pdf.SetLineWidth(0.2)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(8)

	// Campaign information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "INFORMACION DE LA CAMPANA")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(70, 7, "Campana:")
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(120, 7, data.CampaignTitle, "", "L", false)
	pdf.Ln(3)

	// Add section separator
	pdf.SetLineWidth(0.2)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(8)

	// Donor information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "INFORMACION DEL DONANTE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(70, 7, "Donante:")
	pdf.SetFont("Arial", "", 11)
	donorName := data.DonorName
	if data.IsAnonymous {
		donorName = "Donacion Anonima"
	}
	pdf.Cell(120, 7, donorName)
	pdf.Ln(10)

	// Add section separator
	pdf.SetLineWidth(0.2)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(8)

	// Payment information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "DETALLES DEL PAGO")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(70, 7, "Metodo de Pago:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(120, 7, data.PaymentMethod)
	pdf.Ln(7)

	// Amount (highlighted)
	pdf.Ln(5)
	pdf.SetFillColor(240, 240, 240) // Light grey background
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(70, 12, "Monto Total:", "", 0, "L", true, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(46, 204, 113) // Green color for amount
	pdf.CellFormat(120, 12, fmt.Sprintf("$%.2f", data.Amount), "", 1, "L", true, 0, "")
	pdf.SetTextColor(0, 0, 0) // Reset to black
	pdf.Ln(10)

	// Add section separator
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(10)

	// Thank you message
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(41, 128, 185) // Blue color
	pdf.MultiCell(190, 7, "¡Gracias por tu generosa contribucion!", "", "C", false)
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 6, "Tu apoyo hace la diferencia y ayuda a que esta campana alcance su objetivo. Este comprobante es valido como constancia de tu donacion.", "", "C", false)
	pdf.Ln(10)

	// Footer section
	pdf.Ln(10)
	pdf.SetY(-40) // Position 40mm from bottom
	pdf.SetLineWidth(0.2)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(3)

	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128) // Grey color
	pdf.CellFormat(190, 5, "Documento generado automaticamente por Dona Tutti", "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, fmt.Sprintf("Fecha de generacion: %s", time.Now().Format("02/01/2006 15:04:05")), "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, "Este comprobante es valido como constancia de donacion", "", 1, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	// Generate PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

