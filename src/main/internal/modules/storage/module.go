package storage

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/storage/handler"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/types"
	"go.uber.org/fx"
)






// NewStorageHandler creates handler with injected service
func NewStorageHandler() *handler.StorageHandler {
	return handler.NewStorageHandler()
}



// RegisterStorageRoutes registers routes with DI
func RegisterStorageRoutes(rg *types.RouterGroup, h *handler.StorageHandler) {
	RegisterRoutesWithDI(rg.RouterGroup, h)
}


// Module provides all storage module dependencies
var Module = fx.Module("storage",
	fx.Provide(
		NewStorageHandler,
	),
	fx.Invoke(RegisterStorageRoutes),
)
