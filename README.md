# Go-on-Go

**Go-on-Go** es una implementación del juego de mesa Go desarrollada en el lenguaje Go, con énfasis en modelado de estado, reglas del juego y arquitectura modular.

## Arquitectura

El proyecto está estructurado en componentes claramente definidos:

- **Representación del tablero**  
  Modelado del grid, estado de intersecciones y almacenamiento eficiente de piedras.

- **Motor de reglas**  
  Implementación de:
  - Validación de jugadas
  - Detección y captura de grupos
  - Cálculo de libertades
  - Prevención de jugadas ilegales (suicidio, ko si aplica)

- **Gestión de estado de partida**  
  Control de turnos, historial de movimientos y condiciones de finalización.

- **Separación de responsabilidades**  
  División clara entre lógica del dominio (reglas del juego) y capa de interacción.

## Enfoque técnico

El proyecto prioriza:

- Modelado explícito de estructuras de datos
- Inmutabilidad o control estricto del estado cuando es posible
- Uso idiomático de Go (paquetes, tipos, métodos asociados)
- Código testeable y extensible

## Objetivo

Servir como ejercicio de:

- Diseño de motores de reglas
- Implementación de algoritmos sobre grafos implícitos
- Gestión de estado en sistemas deterministas
- Desarrollo estructurado en Go
