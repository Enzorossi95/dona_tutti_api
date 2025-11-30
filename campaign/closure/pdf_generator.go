package closure

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// pdfGenerator implements the PDFGenerator interface
type pdfGenerator struct{}

// NewPDFGenerator creates a new PDF generator for audit reports
func NewPDFGenerator() PDFGenerator {
	return &pdfGenerator{}
}

// Generate creates a PDF audit report and returns the PDF bytes and SHA256 hash
func (g *pdfGenerator) Generate(data AuditReportData) ([]byte, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(0, 102, 204)
	pdf.CellFormat(190, 12, "DONA TUTTI", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(190, 10, "REPORTE DE AUDITORIA DE CAMPANA", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Campaign Information Section
	g.addSectionTitle(pdf, "INFORMACION DE LA CAMPANA")

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 6, fmt.Sprintf("Titulo: %s", data.CampaignTitle), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Organizador: %s", data.OrganizerName), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Fecha de inicio: %s", data.StartDate.Format("02/01/2006")), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Fecha de fin: %s", data.EndDate.Format("02/01/2006")), "", "L", false)
	pdf.MultiCell(190, 6, fmt.Sprintf("Fecha de cierre: %s", data.ClosedAt.Format("02/01/2006 15:04")), "", "L", false)

	closureTypeText := g.getClosureTypeText(data.ClosureType)
	pdf.MultiCell(190, 6, fmt.Sprintf("Tipo de cierre: %s", closureTypeText), "", "L", false)

	if data.ClosureReason != nil && *data.ClosureReason != "" {
		pdf.SetFont("Arial", "I", 10)
		pdf.MultiCell(190, 6, fmt.Sprintf("Justificacion: %s", *data.ClosureReason), "", "L", false)
	}
	pdf.Ln(5)

	// Financial Summary Section
	g.addSectionTitle(pdf, "RESUMEN FINANCIERO")

	// Progress bar visualization
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(190, 8, fmt.Sprintf("Progreso: %.1f%% de la meta alcanzada", data.GoalPercentage), "", 1, "L", false, 0, "")

	// Draw progress bar
	g.drawProgressBar(pdf, data.GoalPercentage)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(95, 6, fmt.Sprintf("Meta de recaudacion: $%.2f", data.CampaignGoal), "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Total recaudado: $%.2f", data.TotalRaised), "", 1, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Total donantes: %d", data.TotalDonors), "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Total donaciones: %d", data.TotalDonations), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Expenses Section
	g.addSectionTitle(pdf, "DETALLE DE GASTOS")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(95, 6, fmt.Sprintf("Total gastos documentados: $%.2f", data.TotalExpenses), "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Comprobantes: %d", data.TotalReceipts), "", 1, "L", false, 0, "")

	if data.TotalReceipts > 0 {
		docPercentage := float64(data.ReceiptsWithDocuments) / float64(data.TotalReceipts) * 100
		pdf.CellFormat(190, 6, fmt.Sprintf("Comprobantes con documento adjunto: %d (%.1f%%)", data.ReceiptsWithDocuments, docPercentage), "", 1, "L", false, 0, "")
	}

	// Receipts table (if there are receipts)
	if len(data.Receipts) > 0 {
		pdf.Ln(3)
		g.addReceiptsTable(pdf, data.Receipts)
	}
	pdf.Ln(5)

	// Activities Section
	g.addSectionTitle(pdf, "ACTIVIDADES REGISTRADAS")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(190, 6, fmt.Sprintf("Total actividades: %d", data.TotalActivities), "", 1, "L", false, 0, "")

	if len(data.Activities) > 0 {
		pdf.Ln(2)
		g.addActivitiesTable(pdf, data.Activities)
	}
	pdf.Ln(5)

	// Transparency Score Section
	g.addSectionTitle(pdf, "PUNTUACION DE TRANSPARENCIA")

	// Draw transparency score badge
	g.drawTransparencyBadge(pdf, data.TransparencyScore)
	pdf.Ln(20)

	// Score breakdown
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(120, 5, "Criterio", "1", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, "Puntos", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "Maximo", "1", 1, "C", false, 0, "")

	g.addScoreRow(pdf, "Documentacion de gastos", data.TransparencyBreakdown.DocumentationScore, 30)
	g.addScoreRow(pdf, "Registro de actividades", data.TransparencyBreakdown.ActivityScore, 25)
	g.addScoreRow(pdf, "Progreso hacia la meta", data.TransparencyBreakdown.GoalProgressScore, 20)
	g.addScoreRow(pdf, "Frecuencia de actualizaciones", data.TransparencyBreakdown.TimelinessScore, 15)
	g.addScoreRow(pdf, "Deduccion por alertas", data.TransparencyBreakdown.AlertsDeductionScore, 0)
	g.addScoreRow(pdf, "Bonificaciones", data.TransparencyBreakdown.BonusScore, 10)

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(120, 5, "TOTAL", "1", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, fmt.Sprintf("%.1f", data.TransparencyScore), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, "100", "1", 1, "C", false, 0, "")

	pdf.Ln(10)

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(190, 5, "Este documento fue generado automaticamente por Dona Tutti", "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, fmt.Sprintf("Fecha de generacion: %s", data.ClosedAt.Format("02/01/2006 15:04:05")), "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, fmt.Sprintf("ID de Campana: %s", data.CampaignID.String()), "", 1, "C", false, 0, "")

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

func (g *pdfGenerator) addSectionTitle(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(190, 8, title, "", 1, "L", true, 0, "")
	pdf.Ln(2)
}

func (g *pdfGenerator) getClosureTypeText(closureType ClosureType) string {
	switch closureType {
	case ClosureTypeGoalReached:
		return "Meta alcanzada"
	case ClosureTypeEndDate:
		return "Fecha limite alcanzada"
	case ClosureTypeManual:
		return "Cierre manual"
	default:
		return string(closureType)
	}
}

func (g *pdfGenerator) drawProgressBar(pdf *gofpdf.Fpdf, percentage float64) {
	x, y := pdf.GetXY()
	barWidth := 150.0
	barHeight := 8.0

	// Background
	pdf.SetFillColor(220, 220, 220)
	pdf.Rect(x, y, barWidth, barHeight, "F")

	// Progress
	progressWidth := barWidth * percentage / 100
	if progressWidth > barWidth {
		progressWidth = barWidth
	}

	if percentage >= 100 {
		pdf.SetFillColor(76, 175, 80) // Green
	} else if percentage >= 75 {
		pdf.SetFillColor(255, 193, 7) // Yellow
	} else if percentage >= 50 {
		pdf.SetFillColor(255, 152, 0) // Orange
	} else {
		pdf.SetFillColor(244, 67, 54) // Red
	}
	pdf.Rect(x, y, progressWidth, barHeight, "F")

	// Border
	pdf.SetDrawColor(180, 180, 180)
	pdf.Rect(x, y, barWidth, barHeight, "D")
}

func (g *pdfGenerator) addReceiptsTable(pdf *gofpdf.Fpdf, receipts []ReceiptSummary) {
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(60, 5, "Proveedor", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 5, "Descripcion", "1", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, "Monto", "1", 0, "R", false, 0, "")
	pdf.CellFormat(35, 5, "Documentado", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	maxReceipts := 10
	for i, receipt := range receipts {
		if i >= maxReceipts {
			pdf.CellFormat(190, 5, fmt.Sprintf("... y %d comprobantes mas", len(receipts)-maxReceipts), "", 1, "L", false, 0, "")
			break
		}

		documented := "No"
		if receipt.HasDocument {
			documented = "Si"
		}

		provider := receipt.Provider
		if len(provider) > 25 {
			provider = provider[:22] + "..."
		}

		name := receipt.Name
		if len(name) > 25 {
			name = name[:22] + "..."
		}

		pdf.CellFormat(60, 5, provider, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 5, name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, 5, fmt.Sprintf("$%.2f", receipt.Total), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 5, documented, "1", 1, "C", false, 0, "")
	}
}

func (g *pdfGenerator) addActivitiesTable(pdf *gofpdf.Fpdf, activities []ActivitySummary) {
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(120, 5, "Titulo", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 5, "Tipo", "1", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	maxActivities := 10
	for i, activity := range activities {
		if i >= maxActivities {
			pdf.CellFormat(190, 5, fmt.Sprintf("... y %d actividades mas", len(activities)-maxActivities), "", 1, "L", false, 0, "")
			break
		}

		title := activity.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}

		pdf.CellFormat(120, 5, title, "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 5, activity.Type, "1", 1, "L", false, 0, "")
	}
}

func (g *pdfGenerator) drawTransparencyBadge(pdf *gofpdf.Fpdf, score float64) {
	x, y := pdf.GetXY()

	// Badge background based on score
	if score >= 80 {
		pdf.SetFillColor(76, 175, 80) // Green
	} else if score >= 60 {
		pdf.SetFillColor(255, 193, 7) // Yellow
	} else if score >= 40 {
		pdf.SetFillColor(255, 152, 0) // Orange
	} else {
		pdf.SetFillColor(244, 67, 54) // Red
	}

	// Draw circle badge
	pdf.Circle(x+15, y+8, 12, "F")

	// Score text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetXY(x+5, y+4)
	pdf.CellFormat(20, 8, fmt.Sprintf("%.0f", score), "", 0, "C", false, 0, "")

	// Label
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(x+35, y+4)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(60, 8, g.getScoreLabel(score), "", 0, "L", false, 0, "")
}

func (g *pdfGenerator) getScoreLabel(score float64) string {
	if score >= 80 {
		return "Excelente"
	} else if score >= 60 {
		return "Bueno"
	} else if score >= 40 {
		return "Regular"
	}
	return "Necesita mejoras"
}

func (g *pdfGenerator) addScoreRow(pdf *gofpdf.Fpdf, label string, score float64, max float64) {
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(120, 5, label, "1", 0, "L", false, 0, "")
	pdf.CellFormat(35, 5, fmt.Sprintf("%.1f", score), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 5, fmt.Sprintf("%.0f", max), "1", 1, "C", false, 0, "")
}
