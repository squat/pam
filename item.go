package pam

type ItemType int

const (
	ItemTypeService ItemType = iota
	ItemTypeUser
	ItemTypeUserPrompt
	ItemTypeTTY
	ItemTypeRUser
	ItemTypeRHost
	ItemTypeAuthToken
	ItemTypeConverstaion
	ItemTypeFailDelay
	ItemTypeXDisplay
	ItemTypeXAuthData
	ItemTypeAuthTokenType
)
