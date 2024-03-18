// internal/serverinterface/server.go

package serverinterface

import "github.com/supabase-community/supabase-go"

type Server interface {
	Sb() *supabase.Client
}
