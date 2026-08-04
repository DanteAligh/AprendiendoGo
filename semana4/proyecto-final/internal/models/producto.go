// Package models contiene las entidades del dominio de este mini-ERP:
// datos puros, sin lógica de negocio ni de persistencia. Esta separación es
// a propósito (repasá el día 22 y el día 26): un modelo solo describe
// "qué forma tienen los datos", nunca "qué se puede o no se puede hacer
// con ellos" (eso vive en internal/service) ni "dónde se guardan" (eso
// vive en internal/repository).
package models

// Producto es uno de los dos recursos principales de este ERP.
//
// Los tags `json:"..."` controlan cómo se serializa/deserializa este
// struct al hablar HTTP con JSON (encoding/json, semana 3, día 17). El
// nombre en JSON suele ir en minúsculas por convención de APIs REST.
type Producto struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
	Stock  int     `json:"stock"`
}

// TODO (te toca a vos, guiado por ejercicios.md días 27-28):
// Vas a necesitar decidir y escribir en la capa de SERVICE (no acá) las
// reglas de negocio para un Producto válido. Pensalo desde ya: ¿puede un
// producto tener precio negativo? ¿puede tener nombre vacío? ¿puede tener
// stock negativo? Estas preguntas se responden con código en
// internal/service, no modificando este struct.
