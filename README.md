# 🐹 golang_basics: Arquitectura, Punteros y Composición

Este repositorio es una guía práctica y austera para entender la esencia de Go, desde sus tipos básicos hasta patrones de arquitectura limpia y composición.

---

## 🏗️ Estructura y Arquitectura (Master Class)

### 📂 Jerarquía de Carpetas
```text
proyecto/
├── cmd/
│   └── server/
│       └── main.go       🎼 Orquestador (cablea repo➔service➔handler)
├── internal/             🔒 Firewall (nadie externo importa)
│   ├── domain/           📁 Core business: structs + errors + interfaces
│   │   ├── user.go       📋 struct User
│   │   └── errors.go     📋 ErrNotFound, ValidationError, Temporary
│   ├── service/          📁 Lógica pura (valida + envuelve)
│   │   ├── user_service.go      ⚙️ interface UserRepository + lógica pura
│   │   └── user_service_test.go 🧪 mock + tests
│   ├── handler/          📁 Adapta HTTP ↔ Service
│   │   └── user_handler.go      ⚗️ HTTP ↔ Service (desenvuelve errores)
│   └── repository/       📁 Implementación concreta
│       └── user_repository.go   🏗️ PostgresUserRepository (genera errores)
├── pkg/                  🧰 Reutilizable (logger, utils)
└── go.mod                📑 Identidad (Nombre del módulo + versión Go)
```

### 📍 FLUJO DE DATOS (Composition)
`main CABLEA` ➔ `handler USA` ➔ `service USA` ➔ `domain`
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`↓ (via interfaz)`
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`repository IMPLEMENTA`

> **Regla de Oro:**
> * **Service** define QUÉ necesita (interfaz).
> * **Repository** dice CÓMO lo hace (concreto).
> * **Main** conecta uno con otro (inyección).

---

## 🛠️ 1. Arquitectura y Manejo de Errores
Diseñado siguiendo el flujo: **Handler ➔ Service ➔ Repository**.

*   **Sentinel Errors** 🚩: `errors.Is(err, domain.ErrNotFound)` para errores específicos.
*   **Error Types** 🏷️: `errors.As(err, &valErr)` para capturar contextos adicionales.
*   **Behavior** 🔄: Interfaces como `Temporary` para decisiones basadas en comportamiento.
*   **Wrapping** 🎁: Uso de `%w` para no perder la traza del error original.

---

## 🔢 2. Tipos y Constantes
*   **Built-in** 🧱: `int`, `float64`, `string` (inmutable), `bool`.
*   **Alias** 🎭: `byte` (uint8) y `rune` (int32/Unicode).
*   **iota** 🔢: Autoincremento inteligente para Enums y Flags.
*   **Untyped Constants** 🎈: Flexibilidad extrema y alta precisión en tiempo de compilación.

---

## 📍 3. Structs y Punteros (Memoria)
*   **Padding** 📏: El orden de los campos importa. Agrupar tipos pequeños ahorra memoria.
*   **Stack vs Heap** 🏔️:
    *   **Value Semantics**: Copia en el **Stack** (rápido).
    *   **Pointer Semantics**: Dirección de memoria. Si retorna un puntero local, "escapa" al **Heap**.
*   **Escape Analysis** 🔍: El compilador decide dónde vive cada dato.
*   **Receiver Types** 📥: `(u User)` para copia, `(u *User)` para modificar el original.

---

## 🪟 4. Arrays y Slices
*   **Mechanical Sympathy** 🏎️: Los arrays son contiguos. El CPU ama las *Cache Lines*.
*   **Slice Anatomy** 🧬: Una "ventana" de 24 bytes (`ptr`, `len`, `cap`).
*   **Sharing Memory** 🤝: Los sub-slices apuntan al mismo array. ¡Cuidado con los efectos secundarios!
*   **append** 🪄: Crea un nuevo array solo si la capacidad (cap) se agota.

---

## 🔌 5. Interfaces y Composición
*   **Duck Typing** 🦆: Si camina y grazna como pato, ES un pato. Sin palabras clave `implements`.
*   **Method Sets** 📐: Los métodos con puntero *requieren* que pases un puntero a la interfaz.
*   **Embedding** 🧩: Composición sobre herencia. Los campos y métodos internos se "promocionan".
*   **Shadowing** 👤: El struct externo puede ocultar métodos del interno sin borrarlos.

---

## 🛡️ 6. Visibilidad y Encapsulamiento
*   **Mayúscula/Minúscula** 🔠: La única regla de visibilidad.
*   **Factory Functions** 🏗️: Constructores `NewUser()` para asegurar estados válidos.
*   **Tell, Don't Ask** 🗣️: No pidas datos para decidir, dile al objeto que haga la acción (`CheckPassword`).

---

### 🚀 Cómo leer este código
Cada archivo en `internal/education` y `internal/domain` contiene comentarios didácticos que referencian estas reglas. ¡Úsalos junto a tus apuntes!

---

## 🎓 El Laboratorio: `internal/education/*`

Si quieres entender la **sintaxis y el comportamiento de Go** desde cero, esta es tu zona. Estos archivos son piezas didácticas diseñadas para ser leídas de forma aislada:

*   **`types.go`** 🔢: Tipos básicos, tipado fuerte y **alias** (`byte`, `rune`).
*   **`constants.go`** 🔢: El contador mágico `iota` y las **untyped constants** de alta precisión.
*   **`structs_pointers.go`** 📍: **Padding** (memoria), Stack vs Heap y **Escape Analysis**.
*   **`arrays_slices.go`** 🪟: Los slices como "mandos a distancia" y la mecánica del `append`.
*   **`interfaces_methods.go`** 🔌: **Duck Typing**, method sets y cómo abrir cajas `any`.
*   **`embedding.go`** 🧩: Composición con piezas de **Lego** vs herencia tradicional.
