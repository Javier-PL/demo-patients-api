package services

import (
	"ccl/ccl-patients-api/models"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"net/http"
	"os"

	"github.com/jung-kurt/gofpdf"
)

//CreatePDF retrieves the PDF data and serves the pdf
func CreatePDF(w http.ResponseWriter, r *http.Request) {

	var invoices []models.Invoice
	_ = json.NewDecoder(r.Body).Decode(&invoices)

	thePDF := BuildPDF(*sortInvoicesByInvoicenumber(invoices))

	f, err := os.Open("historial.pdf")
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	thePDF.Output(w)
	//http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), f)  //this is for overwriting the local pdf file

}

//BuildPDF creates the PDF page itself
func BuildPDF(invoices []models.Invoice) *gofpdf.Fpdf {

	//breaks = calculateBreaks(tsdata)

	//var hasbreak = false

	/*
		var colCount = 6

		var rowsTotal = len(invoices)
		var colWd = 20.0*/

	var rowsTotal = len(invoices)

	const colsTotal = 8

	const (
		colWd_main     = 20.0
		coldWd_remarks = 45.0
		coldWd_total   = colWd_main * 5
		marginH        = 18.0
		marginTop      = 10.0
		marginLeft     = 10.0
		marginRight    = 10.0
		lineHt         = 2.5
		cellGap        = 2.0
		colsTotal_cus  = 2
		colCount_cus   = 2
		colWd_cus      = 20.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}

	var (
		cellList [colsTotal]cellType
		cell     cellType
	)

	pdf := gofpdf.New("Landscape", "mm", "A4", "") // 210 x 297

	columnsHeaders := [colsTotal]string{"FACTURA N°", "FACTURAR A:", "DESCRIPCIÓN", "CANTIDAD", "PRECIO UNITARIO", "TOTAL", "PAGADO", "FECHA"}
	columnsWidths := [colsTotal]float64{24, 80, 50, 20, 40, 20, 20, 20}
	//alignList := [colsTotal]string{"L", "C", "C", "R"}
	//strList := detailsData(tsdata, 0)

	pdf.SetMargins(marginLeft, marginTop, marginRight)

	pdf.AddPage()

	//Company Data
	CompanyData := []string{"Name Surname", "Col. 7242", "Avd. Alameda Sundheim, 28", "21003 Huelva", "620211388"}

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	//Address data
	var address_lineH = 5.0
	pdf.SetFont("Helvetica", "B", 10)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[0]), "L")
	pdf.SetMargins(marginLeft, marginTop, 30)
	pdf.WriteAligned(0, address_lineH, tr("FACTURA"), "R")
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 10)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[1]), "L")
	pdf.Ln(4)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[2]), "L")
	pdf.Ln(4)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[3]), "L")
	pdf.Ln(4)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[4]), "L")
	pdf.Ln(14)

	//Customer data
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(marginLeft, 30)
	pdf.Ln(-1)

	// Rows

	/*
		colWd = 20 //reinit va
			y := pdf.GetY()

			count := -1
			datacount := 0*/

	pdf.SetFont("Arial", "B", 8)

	//Headers
	for colJ := 0; colJ < colsTotal; colJ++ {
		pdf.SetFillColor(172, 181, 174)
		pdf.CellFormat(columnsWidths[colJ], 8, tr(columnsHeaders[colJ]), "1", 0, "CM", true, 0, "")

	}
	pdf.SetFont("Arial", "", 10)
	pdf.Ln(-1)

	// Rows
	y := pdf.GetY()
	count := -1

	for rowJ := 0; rowJ < rowsTotal; rowJ++ {
		//datacount = 0 //which string is gonna be put in every cell
		maxHt := lineHt

		// Cell height calculation loop
		for colJ := 0; colJ < colsTotal; colJ++ {

			colWd := columnsWidths[colJ]

			count++

			/*if count > len(detailsData(tsdata, 0)) {
				count = 1

			}*/
			if colJ == 0 { //fill breaks column

				cell.str = tr(strconv.FormatInt(invoices[rowJ].InvoiceNumber, 10))
			} else if colJ == 1 { //fill breaks column

				text := invoices[rowJ].Patient + " " + invoices[rowJ].PatientDNI
				if invoices[rowJ].IsOrg {
					text = text + " " + invoices[rowJ].PatientAddress
				}
				cell.str = tr(text)
			} else if colJ == 2 {

				description := ""
				retention := ""
				if invoices[rowJ].Retention > 0 {
					retention = "Retención " + strconv.Itoa(invoices[rowJ].Retention) + "%"
				}

				if invoices[rowJ].Description != "" {
					description = invoices[rowJ].Description + "\n" + retention
				} else {
					description = retention
				}
				cell.str = tr(description)

			} else if colJ == 3 {
				cell.str = tr(strconv.Itoa(invoices[rowJ].Units))
			} else if colJ == 4 { //Precio Unitario
				//In order to make line break \n to work, in the cell rendering SplitLines must be applied
				retention := ""
				retention_calc := invoices[rowJ].Price * float64(invoices[rowJ].Retention) / 100
				if invoices[rowJ].Retention > 0 {
					retention = "\n" + fmt.Sprintf("%.2f", retention_calc) + " €"
				}
				cell.str = tr(fmt.Sprintf("%.2f", invoices[rowJ].Price) + " €" + retention)

			} else if colJ == 5 { //Total
				retention := ""
				retention_calc := float64(invoices[rowJ].Units) * invoices[rowJ].Price * float64(invoices[rowJ].Retention) / 100
				if invoices[rowJ].Retention > 0 {
					retention = "\n-" + fmt.Sprintf("%.2f", retention_calc) + " €"
				}
				total_payed := float64(invoices[rowJ].Units) * invoices[rowJ].Price
				cell.str = tr(fmt.Sprintf("%.2f", total_payed) + " €" + retention)
			} else if colJ == 6 { //Pagado
				cell.str = tr(fmt.Sprintf("%.2f", invoices[rowJ].Payed) + " €")
			} else if colJ == 7 {
				cell.str = tr(invoices[rowJ].Date.Format("02-01-2006"))
			} else { //fallback
				cell.str = ""
			}

			/*
				if colJ == 4 {
					cell.list = pdf.SplitLines([]byte(cell.str), 5)
				} else {
					cell.list = pdf.SplitLines([]byte(cell.str), colWd) //cell.list = pdf.SplitLines([]byte(cell.str), colWd-cellGap-cellGap)

				}*/

			cell.list = pdf.SplitLines([]byte(cell.str), colWd) //cell.list = pdf.SplitLines([]byte(cell.str), colWd-cellGap-cellGap)

			//fmt.Println(float64(len(cell.list)))

			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}

			cellList[colJ] = cell

			//datacount++
		}
		// Cell render loop
		x := marginLeft
		for colJ := 0; colJ < colsTotal; colJ++ {

			//handle styles for every cell
			pdf.SetFont("Arial", "", 8)

			pdf.Rect(x, y, columnsWidths[colJ], maxHt+cellGap+cellGap, "D")

			cell = cellList[colJ]
			cellY := y + cellGap + (maxHt-cell.ht)/2

			//pdf.SetXY(x+cellGap, cellY)

			if colJ == 1 || colJ == 2 { //Special case. Columns name and description must be printed in 2 lines

				for splitJ := 0; splitJ < len(cell.list); splitJ++ {
					pdf.SetXY(x+cellGap, cellY)
					pdf.CellFormat(columnsWidths[colJ]-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, "L", false, 0, "")
					cellY += lineHt + lineHt/4
				}

			} else if colJ == 4 || colJ == 5 { //Special case. Columns price and total must be printed in 2 lines

				for splitJ := 0; splitJ < len(cell.list); splitJ++ {
					pdf.SetXY(x+cellGap, cellY)
					pdf.CellFormat(columnsWidths[colJ]-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, "CM", false, 0, "")
					cellY += lineHt + lineHt/4 //added lineHt/4 for more space between lines
				}

			} else {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(columnsWidths[colJ]-cellGap-cellGap, lineHt, cell.str, "", 0, "CM", false, 0, "")
			}

			//cellY += lineHt

			x += columnsWidths[colJ]

		}
		y += maxHt + cellGap + cellGap

		//Log.Info("Y:", y)

		if y > 180 {
			pdf.AddPage()

			y = pdf.GetY()
		}

	}

	return pdf
}

//CreatePDF retrieves the PDF data and serves the pdf
func CreatePatientInvoicePDF(w http.ResponseWriter, r *http.Request) {

	var invoice models.Invoice
	_ = json.NewDecoder(r.Body).Decode(&invoice)

	thePDF := buildPatientInvoicePDF(invoice)

	f, err := os.Open("historial.pdf")
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	thePDF.Output(w)
	//http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), f)  //this is for overwriting the local pdf file

}

func buildPatientInvoicePDF(invoice models.Invoice) *gofpdf.Fpdf {

	const colsTotal = 4

	const (
		marginTop   = 10.0
		marginLeft  = 20.0
		marginRight = 10.0
		lineHt      = 2.5
		cellGap     = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}

	var (
		cellList [colsTotal]cellType
		cell     cellType
	)

	pdf := gofpdf.New("Portrait", "mm", "A4", "") // 210 x 297

	columnsHeaders := [colsTotal]string{"DESCRIPCIÓN", "CANTIDAD", "PRECIO UNITARIO", "TOTAL"}
	columnsWidths := [colsTotal]float64{80, 30, 30, 30}

	pdf.SetMargins(marginLeft, marginTop, marginRight)

	pdf.AddPage()

	//Company Data
	CompanyData := []string{"Name Surname", "49062004Z", "Col. 7242", "Avd. Alameda Sundheim, 28", "21003 Huelva", "620211388"}

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.Image("assets/ccl_logo.png", 160, marginTop, 26, 26, false, "", 0, "")

	//Address data
	var address_lineH = 5.0
	var lines_separation = 5.0
	pdf.SetFont("Helvetica", "B", 10)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[0]), "L")
	pdf.SetMargins(marginLeft, marginTop, 60)
	pdf.WriteAligned(0, address_lineH, tr("FACTURA"), "R")
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.Ln(lines_separation)
	pdf.SetFont("Helvetica", "", 10)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[1]), "L")
	pdf.SetMargins(marginLeft, marginTop, 60)
	pdf.WriteAligned(0, address_lineH, tr("Fecha: "+invoice.Date.Format("01-02-2006")), "R")
	pdf.Ln(lines_separation)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[2]), "L")
	pdf.SetMargins(marginLeft, marginTop, 60)
	pdf.WriteAligned(0, address_lineH, tr("Factura N°: "+strconv.FormatInt(invoice.InvoiceNumber, 10)), "R")
	pdf.Ln(lines_separation)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[3]), "L")
	pdf.Ln(lines_separation)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[4]), "L")
	pdf.Ln(lines_separation)
	pdf.WriteAligned(0, address_lineH, tr(CompanyData[5]), "L")
	pdf.Ln(14)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.WriteAligned(0, address_lineH, tr("Facturar a:"), "L")
	pdf.Ln(lines_separation)
	pdf.SetFont("Helvetica", "", 10)
	pdf.WriteAligned(0, address_lineH, tr(invoice.Patient), "L")
	pdf.Ln(lines_separation)
	pdf.WriteAligned(0, address_lineH, tr(invoice.PatientDNI), "L")
	pdf.Ln(lines_separation)
	if invoice.IsOrg == true {
		pdf.WriteAligned(0, address_lineH, tr(invoice.PatientAddress), "L")
		pdf.Ln(lines_separation)
	}

	//Customer data
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(marginLeft, 80)
	pdf.Ln(-1)

	pdf.SetFont("Arial", "B", 8)

	//Headers
	for colJ := 0; colJ < colsTotal; colJ++ {
		pdf.SetFillColor(172, 181, 174)
		pdf.CellFormat(columnsWidths[colJ], 8, tr(columnsHeaders[colJ]), "1", 0, "CM", true, 0, "")

	}
	pdf.SetFont("Arial", "", 7)
	pdf.Ln(-1)

	// Rows
	y := pdf.GetY()

	maxHt := lineHt

	// Cell height calculation loop
	for colJ := 0; colJ < colsTotal; colJ++ {

		colWd := columnsWidths[colJ]

		if colJ == 0 {

			description := ""
			retention := ""
			if invoice.Retention > 0 {
				retention = "Retención " + strconv.Itoa(invoice.Retention) + "%"
			}

			if invoice.Description != "" {
				description = invoice.Description + "\n" + retention
			} else {
				description = retention
			}
			cell.str = tr(description)

		} else if colJ == 1 {
			cell.str = tr(strconv.Itoa(invoice.Units))
		} else if colJ == 2 {
			//In order to make line break \n to work, in the cell rendering SplitLines must be applied
			retention := ""
			retention_calc := invoice.Price * float64(invoice.Retention) / 100
			if invoice.Retention > 0 {
				retention = "\n-" + fmt.Sprintf("%.2f", retention_calc) + " €"
			}
			cell.str = tr(fmt.Sprintf("%.2f", invoice.Price) + " €" + retention)

		} else if colJ == 3 {
			retention := ""
			retention_calc := float64(invoice.Units) * invoice.Price * float64(invoice.Retention) / 100
			if invoice.Retention > 0 {
				retention = "\n-" + fmt.Sprintf("%.2f", retention_calc) + " €"
			}
			cell.str = tr(fmt.Sprintf("%.2f", float64(invoice.Units)*invoice.Price) + " €" + retention)

		} else { //fallback
			cell.str = ""
		}

		cell.list = pdf.SplitLines([]byte(cell.str), colWd)

		cell.ht = float64(len(cell.list))*lineHt + 50
		if cell.ht > maxHt {
			maxHt = cell.ht
		}

		cellList[colJ] = cell

	}
	// Cell render loop
	x := marginLeft
	for colJ := 0; colJ < colsTotal; colJ++ {

		//handle styles for every cell
		pdf.SetFont("Helvetica", "", 9)

		pdf.Rect(x, y, columnsWidths[colJ], maxHt+cellGap+cellGap, "D")

		cell = cellList[colJ]
		cellY := y + cellGap + (maxHt-cell.ht)/2

		if colJ == 0 { //Special case. Columns name and description must be printed in 2 lines

			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(columnsWidths[colJ]-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, "L", false, 0, "")
				cellY += lineHt + lineHt/2
			}

		} else if colJ == 2 || colJ == 3 { //Special case. Columns price and total must be printed in 2 lines

			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(columnsWidths[colJ]-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, "CM", false, 0, "")
				cellY += lineHt + lineHt/2 //added lineHt/4 for more space between lines
			}

		} else {
			pdf.SetXY(x+cellGap, cellY)
			pdf.CellFormat(columnsWidths[colJ]-cellGap-cellGap, lineHt, cell.str, "", 0, "CM", false, 0, "")
		}

		x += columnsWidths[colJ]

	}
	y += maxHt + cellGap + cellGap

	pdf.SetXY(marginLeft, 160)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetMargins(marginLeft, marginTop, 30)

	retention_calc := float64(invoice.Units) * invoice.Price * float64(invoice.Retention) / 100
	total_payed := (float64(invoice.Units) * invoice.Price) - retention_calc
	pdf.WriteAligned(0, address_lineH, tr("Total "+fmt.Sprintf("%.2f", total_payed)+" €"), "R")

	if invoice.IsOrg == false {
		pdf.WriteAligned(0, address_lineH, tr("Pagado "+fmt.Sprintf("%.2f", invoice.Payed)+" €"), "R")
	}

	return pdf
}

func sortInvoicesByInvoicenumber(invoices []models.Invoice) *[]models.Invoice {

	if len(invoices) > 0 {
		sort.Slice(invoices, func(i, j int) bool {
			return invoices[i].InvoiceNumber < invoices[j].InvoiceNumber
		})
	}

	return &invoices
}
