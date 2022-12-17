package common

type DocumentType string

const (
	DATE_LAYOUT = "2006-01-02"

	// Balanço Patrimonial Ativo
	BPA DocumentType = "BPA"

	// Balanço Patrimonial Passivo
	BPP DocumentType = "BPP"

	// Demonstração de Fluxo de Caixa - Método Direto
	DFCMD DocumentType = "DFC_MD"

	// Demonstração de Fluxo de Caixa - Método Indireto
	DFCMI DocumentType = "DFC_MI"

	// Demonstração das Mutações do Patrimônio Líquido
	DMPL DocumentType = "DMPL"

	// Demonstração de Resultado Abrangente
	DRA DocumentType = "DRA"

	// Demonstração de Resultado
	DRE DocumentType = "DRE"

	// Demonstração de Valor Adicionado
	DVA DocumentType = "DVA"
)

const (
	CONSOLIDATED = "_con"
	INDIVIDUAL   = "_ind"
)

type ContentProcessor interface {
	ProcessContent([][]string) error
	PerformMigration() error
}
