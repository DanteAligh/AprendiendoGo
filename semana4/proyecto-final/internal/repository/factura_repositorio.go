package repository

import (
	"errors"

	"aprendiendo-go/semana4/proyecto-final/internal/models"
)

// FacturaRepository define las operaciones de persistencia para Factura.
// A propósito NO incluye "ValidarStock" ni "CalcularTotal": esa es lógica
// de negocio y vive en internal/service, que además necesita coordinar
// FacturaRepository CON ProductoRepository (para chequear/descontar
// stock). Por eso el service de facturas va a recibir AMBOS repositorios.
type FacturaRepository interface {
	Crear(factura models.Factura) (models.Factura, error)
	ObtenerPorID(id int) (models.Factura, error)
	Listar() ([]models.Factura, error)
}

// ErrFacturaNoEncontrada: mismo patrón que ErrProductoNoEncontrado.
var ErrFacturaNoEncontrada = errors.New("factura no encontrada")
