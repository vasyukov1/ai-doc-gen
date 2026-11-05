package documents

import "strings"

// FKNMetadataDefaults supplies institution labels, never personal approvals or dates.
func FKNMetadataDefaults(language string) map[string]string {
	if language == "ru" {
		return map[string]string{
			"organization": "ПРАВИТЕЛЬСТВО РОССИЙСКОЙ ФЕДЕРАЦИИ\nФЕДЕРАЛЬНОЕ ГОСУДАРСТВЕННОЕ АВТОНОМНОЕ\nОБРАЗОВАТЕЛЬНОЕ УЧРЕЖДЕНИЕ ВЫСШЕГО ОБРАЗОВАНИЯ\nНАЦИОНАЛЬНЫЙ ИССЛЕДОВАТЕЛЬСКИЙ УНИВЕРСИТЕТ\n«ВЫСШАЯ ШКОЛА ЭКОНОМИКИ»",
			"faculty":      "Факультет компьютерных наук", "program": "Образовательная программа «Программная инженерия»"}
	}
	return map[string]string{"organization": "NATIONAL RESEARCH UNIVERSITY\nHIGHER SCHOOL OF ECONOMICS", "faculty": "Faculty of Computer Science", "program": "Software Engineering degree programme"}
}

func manualTODO(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToUpper(value), "TODO") {
		return value
	}
	return "TODO: " + value
}

// Normalize legacy placeholders only in the rendered copy, not saved prose.
func renderTODOs(value string) string {
	return strings.NewReplacer("[Не предоставлено:", "[TODO:", "[Not provided:", "[TODO:").Replace(value)
}

func fknLabels(language string) map[string]string {
	if language == "ru" {
		return map[string]string{
			"student_group": "Студент группы", "page_numbers": "Номера листов (страниц)", "total_pages": "Всего листов (страниц в докум.)", "document_number": "№ документа", "incoming_document": "Входящий № сопроводительного докум. и дата", "change": "Изм.", "changed": "Измененных", "replaced": "Замененных", "new_pages": "Новых", "cancelled_pages": "Аннулированных", "signature_short": "Подп.", "date": "Дата", "sheet": "Лист", "inventory_original": "Инв. № подл.", "inventory_replacement": "Взам. Инв. №", "inventory_duplicate": "Инв. № дубл.", "signature_todo": "Подпись ____________", "date_todo": "Дата ____________"}
	}
	return map[string]string{"student_group": "Student group", "page_numbers": "Sheet (page) numbers", "total_pages": "Total sheets (pages)", "document_number": "Document no.", "incoming_document": "Incoming document no. and date", "change": "Change", "changed": "Changed", "replaced": "Replaced", "new_pages": "New", "cancelled_pages": "Cancelled", "signature_short": "Signature", "date": "Date", "sheet": "Sheet", "inventory_original": "Original inventory no.", "inventory_replacement": "Replacement inventory no.", "inventory_duplicate": "Duplicate inventory no.", "signature_todo": "Signature ____________", "date_todo": "Date ____________"}
}
