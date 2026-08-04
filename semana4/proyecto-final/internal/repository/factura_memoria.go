package repository

import (
	"errors"
	"sync"

	"aprendiendo-go/semana4/proyecto-final/internal/models"
)

// RepositorioFacturaMemoria: mismo patrón que RepositorioProductoMemoria,
// aplicado a Factura. Repetir el patrón (mutex + map + contador de ID) es
// intencional: una vez que lo entendés para un recurso, ya sabés
// replicarlo para cualquier otro recurso de tu futuro ERP.
type RepositorioFacturaMemoria struct {
	mu          sync.Mutex
	facturas    map[int]models.Factura
	siguienteID int
}

func NewRepositorioFacturaMemoria() *RepositorioFacturaMemoria {
	return &RepositorioFacturaMemoria{
		facturas:    make(map[int]models.Factura),
		siguienteID: 1,
	}
}

// Crear: mismo contrato que ProductoRepository.Crear (asignar ID, guardar,
// devolver con ID asignado).
func (r *RepositorioFacturaMemoria) Crear(factura models.Factura) (models.Factura, error) {
	// TODO: implementar.
	return models.Factura{}, errors.New("TODO: RepositorioFacturaMemoria.Crear no implementado")
}

func (r *RepositorioFacturaMemoria) ObtenerPorID(id int) (models.Factura, error) {
	// TODO: implementar. Devolver ErrFacturaNoEncontrada si no existe.
	return models.Factura{}, errors.New("TODO: RepositorioFacturaMemoria.ObtenerPorID no implementado")
}

func (r *RepositorioFacturaMemoria) Listar() ([]models.Factura, error) {
	// TODO: implementar.
	return nil, errors.New("TODO: RepositorioFacturaMemoria.Listar no implementado")
}

var _ FacturaRepository = (*RepositorioFacturaMemoria)(nil)
