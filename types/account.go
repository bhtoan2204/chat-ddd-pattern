package types

type AccountType string

const (
	AccountTypeUser  AccountType = "user"
	AccountTypeAdmin AccountType = "admin"
)

type AccountRole string

const (
	AccountRoleUser  AccountRole = "user"
	AccountRoleAdmin AccountRole = "admin"
)

type AccountStatus int

const (
	AccountStatusActive AccountStatus = iota + 1
	AccountStatusInactive
)
