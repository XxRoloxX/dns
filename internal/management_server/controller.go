package managementserver

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewRecordController struct {
	service *RecordsService
}

type DeleteRecordController struct {
	service *RecordsService
}

type GetRecordsController struct {
	service *RecordsService
}

func NewNewRecordController(service *RecordsService) *NewRecordController {
	return &NewRecordController{
		service: service,
	}
}

func NewDeleteRecordController(service *RecordsService) *DeleteRecordController {
	return &DeleteRecordController{
		service: service,
	}
}

func NewGetRecordsController(service *RecordsService) *GetRecordsController {
	return &GetRecordsController{
		service: service,
	}
}

type NewRecordParams struct {
	Name  string                `json:"name" form:"name" xml:"name" binding:"required"`
	Type  ManagedDNSRecordType  `json:"type"`
	Class ManagedDNSRecordClass `json:"class"`
	Data  string                `json:"data" form:"data" xml:"data" binding:"required"`
}

func (c *NewRecordController) Handle(g *gin.Context) {

	var record NewRecordParams

	if err := g.ShouldBindJSON(&record); err != nil {
		slog.Warn("Couldn't bind")
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the Type and Class fields
	if !isValidRecordType(record.Type) {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	if !isValidRecordClass(record.Class) {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record class"})
		return
	}

	if err := c.service.CreateRecord(&ManagedDNSResourceRecord{
		Name:  record.Name,
		Type:  record.Type,
		Class: record.Class,
		Data:  record.Data,
	}); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusCreated, gin.H{"message": "Record created successfully"})
}

func (c *DeleteRecordController) Handle(g *gin.Context) {

	idStr := g.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.DeleteRecord(id); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

func (c *GetRecordsController) Handle(g *gin.Context) {

	records, err := c.service.GetRecords()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, records)
}

func isValidRecordType(recordType ManagedDNSRecordType) bool {
	validRecordTypes := map[ManagedDNSRecordType]bool{
		ManagedDNSRecordType_A:     true,
		ManagedDNSRecordType_AAAA:  true,
		ManagedDNSRecordType_MX:    true,
		ManagedDNSRecordType_TXT:   true,
		ManagedDNSRecordType_CNAME: true,
		ManagedDNSRecordType_NS:    true,
		ManagedDNSRecordType_SOA:   true,
	}

	return validRecordTypes[recordType]
}

func isValidRecordClass(recordClass ManagedDNSRecordClass) bool {
	validRecordClasses := map[ManagedDNSRecordClass]bool{
		ManagedDNSRecordClass_IN: true,
		ManagedDNSRecordClass_CH: true,
		ManagedDNSRecordClass_CS: true,
		ManagedDNSRecordClass_HS: true,
	}

	return validRecordClasses[recordClass]
}
