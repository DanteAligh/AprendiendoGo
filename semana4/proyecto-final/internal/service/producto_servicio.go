// Package service es la capa de REGLAS DE NEGOCIO (día 22). Un service
// depende de una o más interfaces de repository (nunca de una
// implementación concreta como RepositorioProductoMemoria — eso rompería
// la posibilidad de cambiar de persistencia sin tocar esta capa) y NO sabe
// nada de HTTP: no importa net/http, no arma respuestas JSON, no conoce
// códigos de estado. Eso es trabajo de internal/handler.
package service

import (
	"errors"

	"aprendiendo-go/semana4/proyecto-final/internal/models"
	"aprendiendo-go/semana4/proyecto-final/internal/repository"
)

// ErrProductoInvalido es un error genérico de validación de negocio. En tu
// implementación real probablemente quieras errores más específicos
// (ej. ErrNombreVacio, ErrPrecioInvalido) para que el handler pueda decidir
// mensajes distintos — pero eso es una decisión de diseño que te toca a
// vos tomar y justificar.
var ErrProductoInvalido = errors.New("producto inválido")

// ProductoService define las operaciones de negocio disponibles para
// Producto. Notá que la firma es parecida a ProductoRepository, pero NO es
// lo mismo: acá es donde agregás validaciones ANTES de llamar al
// repositorio.
type ProductoService interface {
	Crear(producto models.Producto) (models.Producto, error)
	ObtenerPorID(id int) (models.Producto, error)
	Listar() ([]models.Producto, error)
	Actualizar(producto models.Producto) (models.Producto, error)
	Eliminar(id int) error
}

// productoServicio es la implementación concreta. Guarda una referencia a
// la INTERFAZ ProductoRepository (no al struct concreto), así este service
// funciona igual sin importar si por debajo hay un repositorio en memoria
// o uno respaldado por una base de datos real.
type productoServicio struct {
	repositorio repository.ProductoRepository
}

// NewProductoService construye el service inyectando su dependencia
// (repository.ProductoRepository) desde afuera — este patrón se llama
// "inyección de dependencias" y es lo que hace posible, en cmd/api/main.go,
// armar el mismo service con distintos repositorios sin cambiar código acá.
func NewProductoService(repositorio repository.ProductoRepository) ProductoService {
	return &productoServicio{repositorio: repositorio}
}

// Crear debe:
//  1. Validar el producto (¿nombre no vacío?, ¿precio > 0?, ¿stock >= 0?).
//     Si no es válido, devolver un error (envolvé ErrProductoInvalido con
//     fmt.Errorf y %w, como practicaste en la semana 2, para dar contexto
//     sin perder la posibilidad de usar errors.Is más arriba).
//  2. Si es válido, delegar al repositorio: s.repositorio.Crear(producto).
func (s *productoServicio) Crear(producto models.Producto) (models.Producto, error) {
	// TODO: implementar validaciones + delegar al repositorio.
	return models.Producto{}, errors.New("TODO: productoServicio.Crear no implementado")
}

func (s *productoServicio) ObtenerPorID(id int) (models.Producto, error) {
	// TODO: implementar (probablemente sea casi un pasamano directo al
	// repositorio, sin mucha validación adicional — pensá si hay algo que
	// SÍ valga la pena validar acá, como que el id sea positivo).
	return models.Producto{}, errors.New("TODO: productoServicio.ObtenerPorID no implementado")
}

func (s *productoServicio) Listar() ([]models.Producto, error) {
	// TODO: implementar.
	return nil, errors.New("TODO: productoServicio.Listar no implementado")
}

// Actualizar: pensá qué validaciones aplican acá además de las de Crear.
// ¿Debería poder cambiarse el ID de un producto? ¿Qué pasa si el nuevo
// stock es negativo?
func (s *productoServicio) Actualizar(producto models.Producto) (models.Producto, error) {
	// TODO: implementar.
	return models.Producto{}, errors.New("TODO: productoServicio.Actualizar no implementado")
}

func (s *productoServicio) Eliminar(id int) error {
	// TODO: implementar. Reto extra: ¿debería poder eliminarse un producto
	// que ya está referenciado por alguna factura existente? No hay una
	// única respuesta correcta — es una decisión de negocio real que vas a
	// tener que tomar (y documentar por qué) en muchos sistemas ERP.
	return errors.New("TODO: productoServicio.Eliminar no implementado")
}
