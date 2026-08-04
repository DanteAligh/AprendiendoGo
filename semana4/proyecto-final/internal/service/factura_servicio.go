package service

import (
	"errors"

	"aprendiendo-go/semana4/proyecto-final/internal/models"
	"aprendiendo-go/semana4/proyecto-final/internal/repository"
)

// ErrStockInsuficiente y ErrProductoNoExiste: mismo concepto que en el
// ejemplo del día 26 (ejemplos/dia26_modelado_erp.go). Volvé a mirar ese
// archivo si no te acordás del patrón de ValidarStock — ahí está resuelto
// en un contexto sin capas ni HTTP; acá el reto es aplicar la MISMA idea
// dentro de la arquitectura de capas.
var (
	ErrStockInsuficiente = errors.New("stock insuficiente")
	ErrProductoNoExiste  = errors.New("producto no existe")
)

// FacturaService es intencionalmente el service MÁS interesante de este
// proyecto: para crear una factura válida necesita LEER productos (para
// conocer precio y stock actuales) Y escribir tanto en FacturaRepository
// (guardar la factura) como en ProductoRepository (descontar stock). Por
// eso recibe ambos repositorios.
type FacturaService interface {
	Crear(factura models.Factura) (models.Factura, error)
	ObtenerPorID(id int) (models.Factura, error)
	Listar() ([]models.Factura, error)
}

type facturaServicio struct {
	facturaRepo  repository.FacturaRepository
	productoRepo repository.ProductoRepository
}

// NewFacturaService recibe DOS repositorios distintos. Este es el patrón
// que vas a repetir constantemente en un ERP real: casi ninguna operación
// de negocio interesante toca un solo recurso; casi todas cruzan varios
// (una factura toca productos, un pedido de compra toca proveedores y
// productos, etc.).
func NewFacturaService(facturaRepo repository.FacturaRepository, productoRepo repository.ProductoRepository) FacturaService {
	return &facturaServicio{facturaRepo: facturaRepo, productoRepo: productoRepo}
}

// Crear es el corazón de este proyecto integrador. Tu implementación
// tiene que, en orden:
//
//  1. Para cada ItemFactura en factura.Items:
//     a. Buscar el Producto correspondiente con s.productoRepo.ObtenerPorID.
//     Si no existe, devolver un error que envuelva ErrProductoNoExiste
//     (fmt.Errorf("...: %w", ErrProductoNoExiste)).
//     b. Verificar que producto.Stock >= item.Cantidad. Si no alcanza,
//     devolver un error que envuelva ErrStockInsuficiente.
//     c. Fijar item.PrecioUnitario = producto.Precio (el precio se congela
//     al momento de facturar, no se recalcula después — repasá el
//     comentario de PrecioUnitario en internal/models/factura.go).
//
//  2. Recién si TODOS los items pasaron la validación (nunca a medias):
//     para cada item, descontar producto.Stock -= item.Cantidad y
//     guardar el producto actualizado con s.productoRepo.Actualizar.
//
//  3. Guardar la factura con s.facturaRepo.Crear y marcarla Confirmada = true.
//
// Pensalo bien: ¿qué pasa si falla el paso 2 para el item 3 de 5, después
// de haber descontado stock de los items 1 y 2? En un sistema con base de
// datos real, esto se resolvería con una transacción (todo o nada). Con
// almacenamiento en memoria no tenés transacciones reales — es un buen
// tema para investigar como reto extra en ejercicios.md.
func (s *facturaServicio) Crear(factura models.Factura) (models.Factura, error) {
	// TODO: implementar los 3 pasos descritos arriba.
	return models.Factura{}, errors.New("TODO: facturaServicio.Crear no implementado")
}

func (s *facturaServicio) ObtenerPorID(id int) (models.Factura, error) {
	// TODO: implementar.
	return models.Factura{}, errors.New("TODO: facturaServicio.ObtenerPorID no implementado")
}

func (s *facturaServicio) Listar() ([]models.Factura, error) {
	// TODO: implementar.
	return nil, errors.New("TODO: facturaServicio.Listar no implementado")
}

// TODO opcional (reto extra): agregá un método Total(factura models.Factura)
// float64 a este service (sumar item.Cantidad * item.PrecioUnitario de
// cada item, tal como en el ejemplo del día 26) y usalo en el handler para
// devolver el total junto con la factura en la respuesta JSON.
