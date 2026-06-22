package domain

type ContextKey string

const (
	CtxUserID       ContextKey = "user_id"
	CtxUserEmail    ContextKey = "user_email"
	CtxSessionID    ContextKey = "session_id"
	CtxIsSuperAdmin ContextKey = "is_superadmin"
	CtxIPAddress    ContextKey = "ip_address"
	CtxUserAgent    ContextKey = "user_agent"
)
