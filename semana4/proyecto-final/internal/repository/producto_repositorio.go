// Package repository es la capa de PERSISTENCIA: sabe cómo guardar y leer
// datos, pero no sabe nada de HTTP ni de reglas de negocio (repasá el día
// 22). Definimos interfaces acá para que la capa de service (que SÍ
// depende de repository) dependa de un CONTRATO, no de una implementación
// concreta. Esto es lo que te permite, el día de mañana, cambiar
// "en memoria" por SQLite o Postgres sin tocar ni una línea de
// internal/service ni internal/handler: solo escribís una nueva
// implementación de esta misma interfaz.
package repository

import (
	"errors"

	"aprendiendo-go/semana4/proyecto-final/internal/models"
)

// ProductoRepository define las operaciones de persistencia disponibles
// para Producto. Fijate que ninguna de ellas sabe qué es un "stock
// insuficiente" ni valida nada de negocio: eso es responsabilidad de
// internal/service.
type ProductoRepository interface {
	Crear(producto models.Producto) (models.Producto, error)
	ObtenerPorID(id int) (models.Producto, error)
	Listar() ([]models.Producto, error)
	Actualizar(producto models.Producto) (models.Producto, error)
	Eliminar(id int) error
}

// ErrProductoNoEncontrado es el error que una implementación de
// ProductoRepository debería devolver cuando ObtenerPorID/Actualizar/
// Eliminar reciben un ID que no existe. Se declara acá (a nivel del
// paquete, junto a la interfaz) para que el resto del código pueda usar
// errors.Is(err, repository.ErrProductoNoEncontrado) sin importar cuál
// implementación concreta esté usando.
//
// TODO: este error está declarado pero todavía no se usa en ningún lado.
// Vas a tener que devolverlo vos desde producto_memoria.go, y detectarlo
// con errors.Is en la capa de service o handler para responder 404.
var ErrProductoNoEncontrado = errors.New("producto no encontrado")
