package repository

import (
	"errors"
	"sync"

	"aprendiendo-go/semana4/proyecto-final/internal/models"
)

// RepositorioProductoMemoria implementa ProductoRepository guardando todo
// en un map en memoria (se pierde al reiniciar el proceso — para
// persistencia real necesitarías escribir otra implementación, por
// ejemplo con SQLite, que satisfaga la misma interfaz ProductoRepository).
//
// El mutex (mu) es OBLIGATORIO acá: net/http atiende cada petición en su
// propia goroutine (semana 3, día 15), así que dos requests concurrentes
// (ej. dos POST /productos al mismo tiempo) podrían leer/escribir el mapa
// al mismo tiempo sin el mutex, causando una data race o directamente un
// panic de Go ("concurrent map writes").
type RepositorioProductoMemoria struct {
	mu          sync.Mutex
	productos   map[int]models.Producto
	siguienteID int
}

// NewRepositorioProductoMemoria construye el repositorio ya listo para
// usarse (con el mapa inicializado — un mapa nil no se puede escribir).
// Este constructor SÍ está implementado por vos como andamiaje: la parte
// que te toca escribir son los métodos de abajo.
func NewRepositorioProductoMemoria() *RepositorioProductoMemoria {
	return &RepositorioProductoMemoria{
		productos:   make(map[int]models.Producto),
		siguienteID: 1,
	}
}

// Crear debe:
//  1. Tomar el lock (r.mu.Lock() y defer r.mu.Unlock()).
//  2. Asignarle un ID nuevo al producto usando r.siguienteID, y luego
//     incrementar r.siguienteID.
//  3. Guardarlo en r.productos.
//  4. Devolver el producto ya con su ID asignado.
func (r *RepositorioProductoMemoria) Crear(producto models.Producto) (models.Producto, error) {
	// TODO: implementar (ver los pasos arriba).
	return models.Producto{}, errors.New("TODO: RepositorioProductoMemoria.Crear no implementado")
}

// ObtenerPorID debe devolver el producto si existe, o
// ErrProductoNoEncontrado si no.
func (r *RepositorioProductoMemoria) ObtenerPorID(id int) (models.Producto, error) {
	// TODO: implementar.
	return models.Producto{}, errors.New("TODO: RepositorioProductoMemoria.ObtenerPorID no implementado")
}

// Listar debe devolver TODOS los productos guardados, como un slice.
// Pista: el orden de un map en Go no está garantizado; si te importa un
// orden estable (ej. por ID), vas a tener que ordenarlo vos (paquete
// "sort" o "slices").
func (r *RepositorioProductoMemoria) Listar() ([]models.Producto, error) {
	// TODO: implementar.
	return nil, errors.New("TODO: RepositorioProductoMemoria.Listar no implementado")
}

// Actualizar debe reemplazar el producto existente con el mismo ID. Si el
// ID no existe, debe devolver ErrProductoNoEncontrado.
func (r *RepositorioProductoMemoria) Actualizar(producto models.Producto) (models.Producto, error) {
	// TODO: implementar.
	return models.Producto{}, errors.New("TODO: RepositorioProductoMemoria.Actualizar no implementado")
}

// Eliminar debe borrar el producto con ese ID. Si no existe, debe devolver
// ErrProductoNoEncontrado.
func (r *RepositorioProductoMemoria) Eliminar(id int) error {
	// TODO: implementar.
	return errors.New("TODO: RepositorioProductoMemoria.Eliminar no implementado")
}

// Verificación en tiempo de compilación de que RepositorioProductoMemoria
// satisface la interfaz ProductoRepository. Si te olvidás de implementar
// algún método (o le cambiás la firma sin querer), esta línea hace que el
// error aparezca ACÁ, con un mensaje claro, en vez de en otro archivo
// lejano que intenta usar la interfaz.
var _ ProductoRepository = (*RepositorioProductoMemoria)(nil)
