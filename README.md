# 🐹 golang_basics: La Masterclass Definitiva de Go

> **Del Zero al Hero**: Una guía práctica y sin rodeos para dominar Go desde tipos básicos hasta Clean Architecture y concurrencia de producción.

Este repositorio no es un tutorial más. Es un **sistema de aprendizaje integrado** donde cada concepto tiene:
- 📖 **Explicación teórica** (aquí en el README)
- 💻 **Código ejecutable** (en `internal/education` y sistema completo)
- 🧪 **Tests reales** (con mocks, race detection, benchmarks)
- 🏗️ **Arquitectura de producción** (Clean Architecture + Concurrency Patterns)

---

## 🧩 ¿Qué hace esta aplicación?

Este repositorio incluye **dos aplicaciones reales y ejecutables**, no solo código de ejemplo:

### 🌐 Servidor HTTP de Usuarios (`cmd/server`)

Una API REST que gestiona usuarios con arquitectura de producción:

| Operación | Qué hace el usuario |
|-----------|---------------------|
| `POST /users` | Registra un nuevo usuario (nombre, email, contraseña) |
| `GET /users/{id}` | Consulta los datos de un usuario por su ID |

```bash
# Arrancar el servidor
go run ./cmd/server
# Servidor escuchando en http://localhost:8080

# Crear un usuario
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret"}'

# Consultar un usuario
curl http://localhost:8080/users/1
```

El servidor valida entradas, devuelve errores descriptivos (400, 404, 500) y almacena los datos en memoria.

---

### ⚙️ Procesador Batch de Usuarios (`cmd/workers`)

Un sistema de procesamiento masivo en paralelo que:

1. **Genera** N usuarios ficticios (configurable, por defecto 2000)
2. **Procesa** cada usuario en paralelo: calcula un hash SHA-256 (carga CPU) y simula un guardado en base de datos (carga I/O)
3. **Registra** toda la actividad con un logger que nunca bloquea el procesamiento
4. **Reporta** cuántos usuarios se procesaron y en cuánto tiempo

```bash
# Ejecutar con valores por defecto (8 workers, 2000 usuarios, 10s timeout)
go run ./cmd/workers

# Personalizar
go run ./cmd/workers -workers=16 -jobs=10000 -timeout=30s

# Salida de ejemplo:
# [worker-01] cpu: 1.2ms, io: 3.1ms → user-0042
# [worker-03] cpu: 1.1ms, io: 3.0ms → user-0091
# ...
# Processed: 2000 | Dropped logs: 0 | Duration: 4.3s
```

Es el mismo patrón que usan sistemas reales para:
- Procesar pedidos de e-commerce masivamente
- Enviar notificaciones push a millones de usuarios
- Indexar documentos en buscadores
- Generar reportes nocturnos

---

## 🗺️ Mapa de Aprendizaje

```
NIVEL 1: FUNDAMENTOS (Sintaxis & Memoria)
├─ 1. Types & Constants             → internal/education/types.go, constants.go
├─ 2. Structs & Pointers            → internal/education/structs_pointers.go
├─ 3. Arrays & Slices               → internal/education/arrays_slices.go
└─ 4. Interfaces & Methods          → internal/education/interfaces_methods.go, embedding.go

NIVEL 2: CLEAN ARCHITECTURE (Clean Code & Design Patterns)
├─ 5. Visibility & Encapsulation    → internal/domain/user.go
├─ 6. Error Handling                → internal/domain/errors.go
├─ 7. Dependency Injection          → internal/service/user_service.go
├─ 8. Layered Architecture          → cmd/server/main.go + handler/service/repository
└─ 9. Testing & Mocking             → internal/service/user_service_test.go

NIVEL 3: CONCURRENCY (Multithreading & Patterns)
├─ 10. Goroutines & Channels        → internal/worker/pool.go
├─ 11. Atomic & Mutex               → internal/worker/stats.go
├─ 12. Context & Cancellation       → internal/worker/pool.go (línea 48-68)
├─ 13. Drop Pattern                 → internal/platform/logger/logger.go
└─ 14. Complete System              → cmd/workers/main.go (batch processor)

NIVEL 4: ADVANCED TOOLING (Testing & Performance)
├─ 15. Advanced Testing             → internal/service/user_service_test.go
├─ 16. Benchmarking                 → internal/worker/stats_bench_test.go
└─ 17. Profiling & Tracing          → cmd/workers/main.go (pprof)
```

---

## 📚 NIVEL 1: FUNDAMENTOS

### 🔢 1. Types & Constants — Type System

> **Archivo**: `internal/education/types.go` + `constants.go`

#### 🧱 Tipos Básicos (Built-in Types)

Go es de **tipado fuerte y estático**. No hay conversiones implícitas ni coerción.

```go
// Booleanos
var b bool = true

// Numéricos (Arquitectura-Dependientes)
var i int      // 32 o 64 bits según tu CPU
var u uint     // Sin signo (solo positivos)

// Numéricos (Tamaño Fijo)
var i8 int8    // -128 a 127
var u8 uint8   // 0 a 255
var i32 int32  // -2^31 a 2^31-1
var f64 float64 // IEEE-754 de 64 bits

// Strings (Inmutables & UTF-8)
var s string = "Hola 🚀"  // No puedes hacer s[0] = 'X'

// Complex (Para matemáticas avanzadas)
var c complex128 = complex(1, 2)  // 1+2i
```

**💡 Reglas de Oro:**
- `int` vs `int32`: `int` es más rápido (se adapta a tu CPU). Úsalo por defecto.
- **Strings son inmutables**: Para modificarlos, conviértelos a `[]byte` o `[]rune`.
- **No hay conversión implícita**: `var x int32 = 10; var y int = x` → ❌ ERROR

#### 🎭 Alias (byte & rune)

```go
var b byte = 255  // Alias de uint8. Para datos binarios.
var r rune = 'A'  // Alias de int32. Para caracteres Unicode.
```

**Cuándo usar cada uno:**
- `byte`: Archivos, sockets, protocolos binarios.
- `rune`: Procesar texto carácter por carácter (especialmente no-ASCII).

#### 🔢 Constants: iota & Untyped Constants

```go
// Enums con iota (Contador Automático)
const (
    _        = iota  // 0: Lo saltamos
    Lunes             // 1
    Martes            // 2
    Miercoles         // 3
)

// Flags con Bitwise (Permisos, Estados)
const (
    Read  = 1 << iota  // 001 (1)
    Write = 1 << iota  // 010 (2)
    Exec  = 1 << iota  // 100 (4)
)

// Combinar flags: permisos := Read | Write
// Verificar: if permisos & Read != 0 { ... }
```

**🎈 Untyped Constants: La Magia de Go**

```go
const Pi = 3.14159265358979323846264338327950288419716939937510

var f32 float32 = Pi  // ✅ OK: Pi se adapta
var f64 float64 = Pi  // ✅ OK: Precisión máxima

// Las constantes tienen +256 bits de precisión.
// Solo pierden decimales al asignarse a una variable.
```

**📍 Ver en código**: `internal/education/types.go:7-32`, `constants.go:8-39`

---

### 📍 2. Structs & Pointers — Memory Management

> **Archivo**: `internal/education/structs_pointers.go`

#### 📏 Memory Padding: El Orden Importa

```go
// ❌ MAL: Ocupa 24 bytes
type BadOrder struct {
    A bool   // 1 byte + 7 padding
    B int64  // 8 bytes
    C bool   // 1 byte + 7 padding
}

// ✅ BIEN: Ocupa 16 bytes
type GoodOrder struct {
    B int64  // 8 bytes
    A bool   // 1 byte
    C bool   // 1 byte + 6 padding
}
```

**💡 Regla**: Ordena campos de **mayor a menor tamaño**. El compilador alinea datos a múltiplos de su tamaño natural (int64 → múltiplo de 8).

**🔍 Herramientas**:
```bash
# Ver el layout de memoria
go tool compile -S structs_pointers.go
```

#### 🏔️ Stack vs Heap: Escape Analysis

```go
// STACK: Rápido, memoria limitada, se limpia automáticamente
func NewPersonValue(name string) Person {
    return Person{Name: name}  // Copia → Stack
}

// HEAP: Lento, memoria grande, Garbage Collector
func NewPersonPointer(name string) *Person {
    p := Person{Name: name}
    return &p  // ⚠️ "Escapa" al Heap (sobrevive a la función)
}
```

**📊 Cuándo usar cada uno:**

| Semántica | Cuándo | Ejemplo |
|-----------|--------|---------|
| **Value** | Datos pequeños (<64 bytes), sin mutación | `time.Time`, `uuid.UUID` |
| **Pointer** | Datos grandes (>64 bytes), necesitas mutar | `User`, `Config`, `Database` |

**🔍 Ver escape analysis:**
```bash
go build -gcflags="-m" structs_pointers.go
```

#### ⬇️⬆️ Sharing Down vs Sharing Up

```go
// ✅ Sharing DOWN: Seguro (main → función)
func main() {
    p := &Person{Name: "Alice"}
    updateName(p)  // main sigue vivo → Stack OK
}

// ⚠️ Sharing UP: Escapa al Heap (función → main)
func createPerson() *Person {
    p := Person{Name: "Bob"}
    return &p  // Sobrevive a su creador → HEAP
}
```

#### 💉 Value vs Pointer Receivers

```go
type Counter struct { Count int }

// Value receiver: NO modifica el original
func (c Counter) GetValue() int {
    return c.Count  // Trabaja sobre una COPIA
}

// Pointer receiver: SÍ modifica el original
func (c *Counter) Increment() {
    c.Count++  // Modifica la memoria compartida
}
```

**🎯 Regla de oro**:
- Usa **pointer receiver** si:
  - El struct es >64 bytes
  - Necesitas mutar el estado
  - Ya tienes otros métodos con pointer receiver (consistencia)

**📍 Ver en código**: `internal/education/structs_pointers.go:7-103`

---

### 🪟 3. Arrays & Slices — Mechanical Sympathy

> **Archivo**: `internal/education/arrays_slices.go`

#### 🏎️ Arrays: Contiguos y Cache-Friendly

```go
arr := [3]int{10, 20, 30}
// Memoria: [10][20][30] → Todo junto, el CPU lo ama

// ⚠️ CUIDADO: Arrays se copian ENTEROS al pasarlos
func print(a [1000000]int) {  // Copia 8MB de memoria
    fmt.Println(a[0])
}
```

**💡 Mechanical Sympathy**: Los arrays aprovechan las **cache lines** del CPU (64 bytes). Datos contiguos = menos cache misses = más velocidad.

#### 🧬 Slices: Anatomy (24 bytes)

```go
type slice struct {
    ptr unsafe.Pointer  // 8 bytes: Puntero al array subyacente
    len int             // 8 bytes: Longitud (elementos usados)
    cap int             // 8 bytes: Capacidad (elementos disponibles)
}

s := []int{10, 20, 30}
// ptr → [10, 20, 30]
// len = 3, cap = 3
```

**🔪 Slicing comparte memoria:**

```go
s1 := []int{10, 20, 30}
s2 := s1[1:3]  // [20, 30] → Mismo array que s1

s2[0] = 999
fmt.Println(s1)  // [10, 999, 30] ← s1 también cambió!
```

#### 🪄 append: Capacity Rule

```go
s := make([]int, 0, 2)  // len=0, cap=2
s = append(s, 1)        // len=1, cap=2 (escribe en [0])
s = append(s, 2)        // len=2, cap=2 (escribe en [1])
s = append(s, 3)        // len=3, cap=4 ← CREA NUEVO ARRAY
```

**📐 Estrategia de crecimiento**:
- cap < 1024 → duplica capacidad
- cap ≥ 1024 → crece 25%

**⚠️ PELIGRO: Append rompe el vínculo**

```go
s1 := make([]int, 2, 2)  // [0, 0], cap=2
s2 := s1                 // Comparten array
s1 = append(s1, 99)      // s1 se muda a array nuevo
s1[0] = 888
fmt.Println(s2)          // [0, 0] ← s2 sigue en el viejo
```

#### 🎭 Strings: Immutable Slices

```go
s := "Hello"
// Internamente: struct { ptr *byte, len int }

s2 := s  // Copia solo 16 bytes (ptr + len), no el texto
```

#### 🔄 Range: Value vs Pointer Semantics

```go
nums := []int{1, 2, 3}

// Value: 'v' es una COPIA
for i, v := range nums {
    v = 999  // No modifica nums
}

// Pointer: Usa el índice para acceder directamente
for i := range nums {
    nums[i] = 999  // SÍ modifica nums
}
```

**📍 Ver en código**: `internal/education/arrays_slices.go:8-62`

---

### 🔌 4. Interfaces & Methods — Duck Typing

> **Archivos**: `internal/education/interfaces_methods.go` + `embedding.go`

#### 🦆 Interfaces: Implementación Implícita

```go
type Speaker interface {
    Speak() string
}

type Dog struct { Name string }

// Dog "implementa" Speaker automáticamente
func (d Dog) Speak() string {
    return "Guau!"
}

var s Speaker = Dog{Name: "Rex"}  // ✅ Funciona sin 'implements'
```

**💡 Filosofía Go**: "Si camina como pato y grazna como pato, ES un pato". No necesitas declarar `implements`.

#### 📐 Method Sets: Critical Rule

```go
type Item struct { Val int }

// Método con POINTER receiver
func (i *Item) Update() {
    i.Val++
}

// ✅ CORRECTO
var updater interface{ Update() } = &Item{}

// ❌ ERROR: 'Item' no tiene Update(), solo '*Item'
var fail interface{ Update() } = Item{}
```

**🎯 Regla**:
- **Pointer receiver** (`*T`) → Solo `*T` implementa la interfaz
- **Value receiver** (`T`) → Tanto `T` como `*T` implementan la interfaz

#### 🧩 Embedding: Composition > Inheritance

```go
type Engine struct { HP int }
func (e Engine) Start() { fmt.Println("Vroom!") }

type Car struct {
    Engine  // Sin nombre → "promoción"
    Brand string
}

c := Car{Brand: "Tesla", Engine: Engine{HP: 300}}
c.HP        // ✅ Acceso directo (promoción)
c.Start()   // ✅ Método promocionado
```

#### 👤 Shadowing: Ocultar sin Destruir

```go
func (c Car) Start() {
    fmt.Println("Silencioso...")  // Sobrescribe Engine.Start()
}

c.Start()        // → "Silencioso..."
c.Engine.Start() // → "Vroom!" (sigue disponible)
```

#### ⚠️ Composition ≠ Type Polymorphism

```go
func Repair(e Engine) { ... }

c := Car{}
Repair(c)         // ❌ ERROR: Car no es Engine
Repair(c.Engine)  // ✅ CORRECTO: Pasamos la pieza
```

#### 📦 Type Assertions: Abrir la Caja 'any'

```go
func InspectAny(v any) {
    // Type Switch (elegante)
    switch val := v.(type) {
    case int:
        fmt.Printf("Int: %d\n", val)
    case string:
        fmt.Printf("String: %s\n", val)
    default:
        fmt.Println("Unknown")
    }

    // Comma-OK Idiom (seguro)
    s, ok := v.(string)
    if !ok {
        panic("not a string!")  // Manejo explícito
    }
    fmt.Println(s)
}
```

**📍 Ver en código**: `internal/education/interfaces_methods.go:7-85`, `embedding.go:6-51`

---

## 🏗️ NIVEL 2: CLEAN ARCHITECTURE

### 🔐 5. Visibility & Encapsulation

> **Archivo**: `internal/domain/user.go`

#### 🔠 Capitalization Rule

```go
package domain

type User struct {
    ID    int     // ✅ Exportado (público)
    Name  string  // ✅ Exportado

    passwordHash string  // ❌ Unexported (privado)
}

// Solo el código en 'domain' puede tocar passwordHash
```

**💡 Consecuencia**: No hay getters/setters automáticos. Si quieres acceso, lo diseñas explícitamente.

#### 🏗️ Factory Functions: Safe Constructors

```go
// NewUser: El "constructor" oficial
func NewUser(name, email, password string) *User {
    return &User{
        Name:  name,
        Email: email,
        passwordHash: hash(password),  // Lógica de negocio
    }
}

// ❌ Sin factory, alguien podría hacer:
// u := User{passwordHash: ""}  // Estado inválido
```

#### 🗣️ Tell, Don't Ask

```go
// ❌ MAL: "Ask" (pedir el hash y validar fuera)
if user.GetPasswordHash() == hash(input) { ... }

// ✅ BIEN: "Tell" (decirle al objeto que se valide)
if user.CheckPassword(input) { ... }
```

**💡 Principio**: Los objetos deben **hacer** cosas, no ser **sacos de datos** que otros manipulan.

**📍 Ver en código**: `internal/domain/user.go:5-51`

---

### 🚨 6. Error Handling — The 3 Types

> **Archivo**: `internal/domain/errors.go`

Go no tiene excepciones. Los errores son **valores** que retornas y manejas explícitamente.

#### 🚩 1. Sentinel Errors (Valor Específico)

```go
var ErrNotFound = errors.New("user not found")

// Uso: errors.Is()
if errors.Is(err, domain.ErrNotFound) {
    return http.StatusNotFound
}
```

**Cuándo usar**: Estados conocidos y sin contexto adicional (`EOF`, `ErrNotFound`, `ErrClosed`).

#### 🏷️ 2. Type Errors (Struct con Contexto)

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}

// Uso: errors.As()
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println("Field:", valErr.Field)
}
```

**Cuándo usar**: Necesitas contexto (qué campo falló, código de error, detalles).

#### 🔄 3. Behavior Errors (Interfaz)

```go
type Temporary interface {
    IsTemporary() bool
}

// Uso: Preguntamos por capacidad
var tempErr Temporary
if errors.As(err, &tempErr) && tempErr.IsTemporary() {
    retry()  // Error temporal, reintentar
}
```

**Cuándo usar**: Decisiones basadas en **comportamiento** (¿es temporal? ¿es timeout? ¿es retryable?).

#### 🎁 Error Wrapping: %w (Go 1.13+)

```go
// Repository genera el error
return domain.ErrNotFound

// Service lo envuelve con contexto
return fmt.Errorf("service failed: %w", err)

// Handler puede desenvolverlo
if errors.Is(err, domain.ErrNotFound) {
    // ✅ Encuentra ErrNotFound aunque esté envuelto
}
```

**💡 Regla**: Usa `%w` para mantener la cadena de errores. Usa `%v` si quieres "esconder" el error original.

**📍 Ver en código**: `internal/domain/errors.go:5-30`

---

### 💉 7. Dependency Injection — Interfaces

> **Archivo**: `internal/service/user_service.go`

#### 📜 The Contract (Interface)

```go
type UserRepository interface {
    Save(user domain.User) error
    GetByID(id int) (domain.User, error)
}

type UserService struct {
    repo UserRepository  // ← Depende de la ABSTRACCIÓN
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

**💡 Ventaja**: `UserService` no sabe si `repo` es Postgres, MongoDB o un Mock. Solo sabe que cumple el contrato.

#### 🏗️ Concrete First (Filosofía Go)

```go
// 1. Escribes el tipo concreto
type PostgresUserRepository struct { ... }

// 2. Solo extraes interface SI LA NECESITAS (tests, multiple impls)
type UserRepository interface {
    Save(user domain.User) error
}
```

**⚠️ Anti-patrón**: Crear interfaces "por si acaso". **YAGNI** (You Ain't Gonna Need It).

#### 🤡 Consumer-Defined Interfaces

```go
// ❌ MAL: El proveedor define una interfaz gigante
type IUserRepository interface {
    Save(...) error
    GetByID(...) error
    Update(...) error
    Delete(...) error
    List(...) error
    // 20 métodos más...
}

// ✅ BIEN: El consumidor define lo que necesita
package service

type UserSaver interface {  // Solo Save()
    Save(user domain.User) error
}

func (s *UserService) RegisterUser(repo UserSaver) { ... }
```

**💡 Regla**: Las interfaces deben ser **pequeñas** y estar donde se **usan**, no donde se implementan.

**📍 Ver en código**: `internal/service/user_service.go:8-43`, `user_service_test.go:8-27`

---

### 🎼 8. Layered Architecture — Clean Architecture

> **Archivos**: `cmd/server/main.go` + `handler/` + `service/` + `repository/`

#### 📂 La Jerarquía de Carpetas

```
cmd/
├── server/main.go     → Orquestador (cablea todo)
└── workers/main.go    → Batch processor (concurrencia)

internal/              → Firewall (solo este proyecto puede importarlo)
├── domain/            → Core Business (sin dependencias externas)
│   ├── user.go        → Entidades de negocio
│   └── errors.go      → Errores de dominio
├── service/           → Lógica de Negocio (orquestación)
│   └── user_service.go → Valida, coordina, envuelve
├── handler/           → Adaptadores (HTTP → Service)
│   └── user_handler.go → Traduce JSON, mapea códigos HTTP
├── repository/        → Implementaciones de Persistencia
│   └── user_repository.go → Postgres, MySQL, Mock
├── platform/          → Infraestructura compartida
│   └── logger/        → Logger con Drop Pattern
└── worker/            → Sistema de concurrencia
    ├── pool.go        → Worker pool
    ├── stats.go       → Estadísticas atómicas
    └── job.go         → Unidad de trabajo

education/             → Material didáctico (no ejecutable)
```

#### 🔄 Flujo de Datos (Dependency Rule)

```
HTTP Request
    ↓
Handler (adapter)
    ↓
Service (business logic)
    ↓
Repository Interface ← Service depende de abstracción
    ↓
Repository Concrete ← Implementación depende del dominio
    ↓
Database
```

**💡 Regla de Dependencia**: Las capas internas NO conocen las externas.
- ✅ Service usa Repository (interface)
- ❌ Repository NO puede importar Service

#### 🎼 Main: The Orchestrator

```go
func main() {
    // 1. Crear capas de dentro hacia fuera
    repo := repository.NewPostgresUserRepository()

    // 2. Inyectar dependencias
    svc := service.NewUserService(repo)

    // 3. Crear adaptadores
    hdl := handler.NewUserHandler(svc)

    // 4. Configurar servidor
    http.HandleFunc("/users", hdl.CreateUser)
    http.ListenAndServe(":8080", nil)
}
```

**💡 Responsabilidad de main**: **CABLEAR** todo. Es el único lugar que conoce las implementaciones concretas.

#### 🗺️ Error Mapping by Layer

| Capa | Responsabilidad | Ejemplo |
|------|----------------|---------|
| **Repository** | GENERA errores de dominio | `return domain.ErrNotFound` |
| **Service** | ENVUELVE con contexto | `fmt.Errorf("service failed: %w", err)` |
| **Handler** | DESENVUELVE y mapea HTTP | `if errors.Is(err, ErrNotFound) → 404` |

**📍 Ver en código**: `cmd/server/main.go:11-29`, `handler/user_handler.go:22-73`, `service/user_service.go:15-43`, `repository/user_repository.go:10-34`

---

### 🧪 9. Testing & Mocking

> **Archivo**: `internal/service/user_service_test.go`

#### 🤡 Simple Mock

```go
type MockRepo struct {
    saveCalled bool
    saveError  error
}

func (m *MockRepo) Save(user domain.User) error {
    m.saveCalled = true
    return m.saveError  // Control total sobre el comportamiento
}

func TestRegisterUser(t *testing.T) {
    mock := &MockRepo{}
    svc := service.NewUserService(mock)

    err := svc.RegisterUser(domain.User{Name: "Alice"})

    if !mock.saveCalled {
        t.Error("Save() should have been called")
    }
}
```

#### 🏃 Running Tests

```bash
# Tests normales
go test ./...

# Con coverage
go test -cover ./...

# Race detector (CRÍTICO para concurrencia)
go test -race ./...

# Benchmarks
go test -bench=. -benchmem ./...
```

**📍 Ver en código**: `internal/service/user_service_test.go:8-31`, `internal/platform/logger/logger_test.go`, `internal/worker/pool_test.go`

---

## ⚡ NIVEL 3: PRODUCTION-GRADE CONCURRENCY

### 🏭 Complete System: Batch User Processor

Un **sistema real e integrado** que demuestra TODOS los conceptos de concurrencia de Go (Módulos 8-11 de Ultimate Go) trabajando juntos en un único pipeline funcional.

#### 🎯 ¿Qué hace?

Genera **N usuarios ficticios**, los **procesa en paralelo** con un pool de workers, y lo **registra todo** con un logger que nunca se bloquea (Drop Pattern).

Es lo mismo que hace cualquier sistema batch de producción:
- Procesar pedidos de e-commerce
- Indexar documentos en un buscador
- Enviar notificaciones masivas
- Generar reportes

```bash
# Ejecutar (defaults: 8 workers, 2000 jobs, 10s timeout)
go run ./cmd/workers

# Custom
go run ./cmd/workers -workers=16 -jobs=10000 -timeout=30s -mode=semaphore
```

#### 📐 System Architecture

```
                 cmd/workers/main.go
                 ══════════════════
                 Orquestador principal
                        │
          ┌─────────────┼─────────────┐
          │             │             │
          ▼             ▼             ▼
    logger.New()   GenerateJobs()  Pool.Run()
    (platform)     (2000 users)   (concurrency)
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
                    ▼                  ▼                  ▼
              Goroutines           Channels          Context
              WaitGroup            Buffered        Cancellation
              (Mod 8)             (Mod 10)         (Mod 10)
                    │                  │                  │
                    └──────────────────┼──────────────────┘
                                       │
                         ┌─────────────┴─────────────┐
                         │                           │
                         ▼                           ▼
                    Stats (Mod 9)            Logger (Mod 11)
                    atomic.Int64             Drop Pattern
                    RWMutex                  Never Blocks
```

---

### 🧵 10. Goroutines & Channels

> **Archivo**: `internal/worker/pool.go`

#### 🏊‍♂️ Worker Pool Pattern

```go
func (p *Pool) Run(ctx context.Context, jobs []Job) *Stats {
    ch := make(chan Job, p.numWorkers)  // Canal buffered

    // PRODUCER: Envía jobs al canal
    go func() {
        defer close(ch)  // Señal de "terminé"
        for _, job := range jobs {
            select {
            case ch <- job:  // Encolar job
            case <-ctx.Done():  // Cancelación
                return
            }
        }
    }()

    // WORKERS: Procesan en paralelo
    var wg sync.WaitGroup
    for i := 0; i < p.numWorkers; i++ {
        wg.Add(1)
        go func(workerID string) {
            defer wg.Done()

            // POOLING: Iterar hasta que el canal se cierre
            for job := range ch {
                p.process(ctx, workerID, job)
            }
        }(fmt.Sprintf("worker-%02d", i))
    }

    wg.Wait()  // Esperar a todos los workers
    return p.stats
}
```

**💡 Patrones demostrados:**
- **Pooling**: N goroutines fijas consumiendo de un canal
- **Fan Out Bounded**: Limitar concurrencia (vs lanzar 10,000 goroutines)
- **Signaling**: `close(ch)` para terminar limpiamente
- **Cancellation**: `ctx.Done()` para parar inmediatamente

#### 🌋 CPU-bound vs 🌊 IO-bound

```go
func (p *Pool) process(ctx context.Context, workerID string, job Job) {
    // FASE CPU-BOUND: Cálculo intensivo
    startCPU := time.Now()
    for i := 0; i < 5000; i++ {
        hash := sha256.Sum256([]byte(job.Name))
        _ = hash
    }
    cpuDur := time.Since(startCPU)

    // FASE IO-BOUND: Espera a disco/red
    startIO := time.Now()
    select {
    case <-time.After(3 * time.Millisecond):  // Simular "guardado"
    case <-ctx.Done():  // Respetar cancelación
        return
    }
    ioDur := time.Since(startIO)

    p.logger.Log(fmt.Sprintf("[%s] cpu: %v, io: %v", workerID, cpuDur, ioDur))
}
```

**🔬 Experimento:**
```bash
# 1 core → CPU-bound se serializa
GOMAXPROCS=1 go run ./cmd/workers

# 8 cores → CPU-bound se paraleliza
GOMAXPROCS=8 go run ./cmd/workers

# IO-bound NO mejora con más cores (espera externa)
```

**📍 Ver en código**: `internal/worker/pool.go:31-79` (Pool.Run), `pool.go:82-109` (process)

---

### 🔒 11. Atomic & Mutex

> **Archivo**: `internal/worker/stats.go`

#### ⚡ Atomic: Lock-Free Performance

```go
type Stats struct {
    processed atomic.Int64  // ← Operación atómica (sin locks)

    mu        sync.RWMutex
    perWorker map[string]int
}

func (s *Stats) RecordProcessed(workerID string) {
    // 1. Contador simple → Atomic (10-100x más rápido que Mutex)
    s.processed.Add(1)

    // 2. Estructura compleja → Mutex (mapas no son thread-safe)
    s.mu.Lock()
    s.perWorker[workerID]++
    s.mu.Unlock()
}
```

**💡 Cuándo usar cada uno:**

| Operación | Usar | Razón |
|-----------|------|-------|
| Contador simple | `atomic.Int64` | Lock-free, extremadamente rápido |
| Mapa | `sync.Mutex` | Mapas no son thread-safe en Go |
| Struct complejo | `sync.Mutex` | Múltiples campos relacionados |
| Lectura >> Escritura | `sync.RWMutex` | Permite múltiples lectores |

#### 🔐 RWMutex: Multiple Readers

```go
func (s *Stats) PerWorker() map[string]int {
    s.mu.RLock()  // ← Múltiples goroutines pueden leer a la vez
    defer s.mu.RUnlock()

    // IMPORTANTE: Retornar una COPIA
    copyMap := make(map[string]int, len(s.perWorker))
    for k, v := range s.perWorker {
        copyMap[k] = v
    }
    return copyMap  // El llamador recibe su propia copia
}
```

**⚠️ PELIGRO: Data Races en Mapas**

```go
// ❌ ERROR: Sin protección → panic: concurrent map writes
m := make(map[string]int)
go func() { m["a"]++ }()
go func() { m["b"]++ }()

// ✅ CORRECTO: Con Mutex
var mu sync.Mutex
go func() { mu.Lock(); m["a"]++; mu.Unlock() }()
```

#### 🏁 Race Detector

```bash
# Detectar data races
go test -race ./...

# Ver data race intencional (educativo)
go test -tags=race_example -run=TestIntentionalRace ./internal/worker/
```

**📍 Ver en código**: `internal/worker/stats.go:12-58`, `race_example_test.go`

---

### ⏱️ 12. Context & Cancellation

> **Archivo**: `internal/worker/pool.go` (líneas 39-50, 66-69, 97-102)

#### 🛑 Context: Cooperative Cancellation

```go
// main.go: Crear context con timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// pool.go: Producer respeta cancelación
for _, job := range jobs {
    select {
    case ch <- job:  // Enviar job
    case <-ctx.Done():  // ¡Tiempo agotado! Parar.
        return
    }
}

// pool.go: Worker respeta cancelación
for job := range ch {
    if ctx.Err() != nil {  // Verificar antes de trabajo pesado
        return
    }
    process(ctx, job)
}

// pool.go: Process respeta cancelación
select {
case <-time.After(3 * time.Millisecond):  // Simular IO
case <-ctx.Done():  // Cancelar durante IO
    return
}
```

**💡 Patrones:**
- **Timeout**: `context.WithTimeout(parent, 10*time.Second)`
- **Cancelación manual**: `cancel()` (limpia recursos)
- **Propagación**: Un context cancelado propaga a todos sus hijos

**📍 Ver en código**: `cmd/workers/main.go:39`, `internal/worker/pool.go:48-50,66-69,97-102`

---

### 🛡️ 13. Drop Pattern — Never Block

> **Archivo**: `internal/platform/logger/logger.go`

#### 📦 Never-Blocking Logger

```go
type Logger struct {
    ch      chan string   // Buffer de mensajes
    writer  io.Writer
    dropped atomic.Int64  // Contador de drops
    wg      sync.WaitGroup
}

func New(w io.Writer, bufferSize int) *Logger {
    l := &Logger{
        ch:     make(chan string, bufferSize),
        writer: w,
    }

    // Goroutine de escritura en background
    l.wg.Add(1)
    go func() {
        defer l.wg.Done()
        for msg := range l.ch {  // Pooling
            io.WriteString(l.writer, msg+"\n")
        }
    }()

    return l
}

func (l *Logger) Log(msg string) {
    select {
    case l.ch <- msg:  // ✅ Buffer tiene espacio
    default:           // ❌ Buffer lleno → DROP
        l.dropped.Add(1)
    }
}
```

**💡 Filosofía del Drop Pattern:**
- La **aplicación principal nunca se bloquea** por logs
- Si el disco está lento, descartamos logs (se registra el count)
- Prioridad: **Throughput > Completitud de logs**

#### 🚪 Clean Shutdown

```go
func (l *Logger) Shutdown() {
    close(l.ch)  // Señal: No más mensajes
    l.wg.Wait()  // Esperar drenaje completo del buffer
}
```

**⚠️ Prevención de Goroutine Leak:**

```go
// Si usáramos un channel de resultado:
// ❌ MAL: Si main se va antes, goroutine se queda bloqueada
done := make(chan struct{})
go func() {
    // trabajo...
    done <- struct{}{}  // ← Bloqueado si nadie escucha
}()

// ✅ BIEN: Buffer de 1 permite enviar aunque nadie escuche
done := make(chan struct{}, 1)
go func() {
    // trabajo...
    done <- struct{}{}  // ← No bloquea (buffer de 1)
}()
```

**📍 Ver en código**: `internal/platform/logger/logger.go:9-79`, `logger_test.go:10-68`

---

### 🏗️ 14. Sistema Completo en Acción

> **Archivo**: `cmd/workers/main.go`

#### 🎼 Complete Orchestration

```go
func main() {
    // 1. Configuración
    numWorkers := flag.Int("workers", 8, "Número de workers")
    numJobs := flag.Int("jobs", 2000, "Jobs a procesar")
    timeout := flag.Duration("timeout", 10*time.Second, "Timeout")
    flag.Parse()

    // 2. Componentes
    l := logger.New(os.Stdout, 256)       // Logger (Mod 11)
    jobs := worker.GenerateJobs(*numJobs) // Data source
    ctx, cancel := context.WithTimeout(context.Background(), *timeout)
    defer cancel()

    // 3. Procesamiento
    p := worker.NewPool(*numWorkers, l)
    stats := p.Run(ctx, jobs)  // 🚀 Aquí ocurre la magia

    // 4. Shutdown limpio
    l.Shutdown()  // Drenar buffer

    // 5. Reportes
    fmt.Printf("Processed: %d\n", stats.TotalProcessed())
    fmt.Printf("Dropped logs: %d\n", l.DroppedCount())
}
```

#### 🚦 Alternative Mode: Semaphore

```go
func runWithSemaphore(ctx context.Context, maxWorkers int, jobs []Job) *Stats {
    sem := make(chan struct{}, maxWorkers)  // Semáforo
    var wg sync.WaitGroup

    for _, job := range jobs {
        wg.Add(1)
        go func(j Job) {
            defer wg.Done()

            // ACQUIRE: Esperar hueco
            select {
            case sem <- struct{}{}:
            case <-ctx.Done():
                return
            }
            defer func() { <-sem }()  // RELEASE

            // Procesar...
        }(job)
    }

    wg.Wait()
    return stats
}
```

**💡 Pool vs Semaphore:**

| Pattern | Uso | Pros | Contras |
|---------|-----|------|---------|
| **Pool** | Jobs predecibles | Memoria constante | Setup inicial |
| **Semaphore** | Jobs esporádicos | Flexible | Más goroutines |

**📍 Ver en código**: `cmd/workers/main.go:16-133`

---

## 📊 Tabla de Conceptos vs Ubicación

| Concepto | Archivo | Líneas |
|----------|---------|--------|
| **Tipos básicos** | `education/types.go` | 7-32 |
| **iota & constantes** | `education/constants.go` | 8-39 |
| **Memory padding** | `education/structs_pointers.go` | 10-24 |
| **Stack vs Heap** | `education/structs_pointers.go` | 44-57 |
| **Value/Pointer receivers** | `education/structs_pointers.go` | 33-42 |
| **Slices anatomía** | `education/arrays_slices.go` | 9-34 |
| **Duck typing** | `education/interfaces_methods.go` | 26-39 |
| **Embedding** | `education/embedding.go` | 6-37 |
| **Sentinel errors** | `domain/errors.go` | 8-10 |
| **Type errors** | `domain/errors.go` | 12-22 |
| **Behavior errors** | `domain/errors.go` | 24-29 |
| **Factory functions** | `domain/user.go` | 23-29 |
| **Tell Don't Ask** | `domain/user.go` | 32-35 |
| **DI con interfaces** | `service/user_service.go` | 8-23 |
| **Error wrapping** | `service/user_service.go` | 39-42 |
| **Layered architecture** | `cmd/server/main.go` | 11-29 |
| **Mocking** | `service/user_service_test.go` | 8-27 |
| **Goroutines + WaitGroup** | `worker/pool.go` | 56-73 |
| **Channels + Pooling** | `worker/pool.go` | 36-52 |
| **CPU vs IO-bound** | `worker/pool.go` | 82-104 |
| **Context cancellation** | `worker/pool.go` | 48-50,66-69,97-102 |
| **Atomic** | `worker/stats.go` | 25-28 |
| **RWMutex** | `worker/stats.go` | 33-36,47-57 |
| **Drop Pattern** | `platform/logger/logger.go` | 44-53 |
| **Shutdown limpio** | `platform/logger/logger.go` | 58-74 |
| **Semaphore** | `cmd/workers/main.go` | 66-111 |
| **Table-Driven Tests** | `service/user_service_test.go` | 44-90 |
| **Mocking (Service)** | `service/user_service_test.go` | 14-40 |
| **Benchmarking** | `worker/stats_bench_test.go` | 15-67 |
| **Profiling (pprof)** | `cmd/workers/main.go` | 7-33 |

---

## 🚀 Cómo Estudiar Este Repositorio

### 📚 Ruta de Aprendizaje Recomendada

#### Fundamentos
1. Lee `internal/education/types.go` + ejecuta en tu mente
2. Lee `internal/education/constants.go` + crea tus propios enums
3. Lee `internal/education/structs_pointers.go` + dibuja memoria
4. **Ejercicio**: Crea un `type Product struct` con padding optimizado

#### Estructuras de Datos
1. Lee `internal/education/arrays_slices.go` + experimenta con append
2. Lee `internal/education/interfaces_methods.go` + crea tus interfaces
3. Lee `internal/education/embedding.go` + compáralo con herencia
4. **Ejercicio**: Implementa un `Stack` con slices y método `Push/Pop`

#### Clean Architecture
1. Lee `internal/domain/user.go` + `errors.go`
2. Lee `internal/service/user_service.go` + `user_service_test.go`
3. Lee `internal/handler/user_handler.go` + `repository/user_repository.go`
4. Estudia `cmd/server/main.go` (cómo se cablea)
5. **Ejercicio**: Añade `UpdateUser()` en todas las capas

#### Concurrency
1. Lee `internal/platform/logger/logger.go` + sus tests
2. Lee `internal/worker/stats.go` (atomic + mutex)
3. Lee `internal/worker/pool.go` (el corazón)
4. Lee `cmd/workers/main.go` (orquestación)
5. **Ejercicio**: Ejecuta con `GOMAXPROCS=1` vs `GOMAXPROCS=8`

#### Advanced Tooling
1. Lee `internal/service/user_service_test.go` (Advanced Testing & Mocks)
2. Lee `internal/worker/stats_bench_test.go` (Performance measurement)
3. Ejecuta `go test -bench=. ./internal/worker/` y analiza los resultados
4. Abre `http://localhost:6060/debug/pprof` mientras corres el batch processor
5. **Ejercicio**: Encuentra el número óptimo de workers para tu CPU usando el benchmark

### 🧪 Comandos Esenciales

```bash
# Tests completos
go test -v -race -cover ./...

# Ver data race educativo
go test -tags=race_example -run=TestIntentionalRace ./internal/worker/

# Batch processor
go run ./cmd/workers -workers=4 -jobs=1000

# Server HTTP
go run ./cmd/server

# Escape analysis
go build -gcflags="-m" ./internal/education/structs_pointers.go

# Benchmarks
go test -bench=. -benchmem ./internal/worker/
```

---

## 🎯 Recursos Adicionales

- **Ultimate Go**: https://www.ardanlabs.com/ultimate-go/
- **Go by Example**: https://gobyexample.com/
- **Effective Go**: https://go.dev/doc/effective_go
- **Go Memory Model**: https://go.dev/ref/mem
- **Dave Cheney (Performance)**: https://dave.cheney.net/

---

## 📜 Licencia

MIT License — Úsalo, modifícalo, enséñalo, róbalo, mejóralo.

---

**🔥 Nota Final**: Este repositorio es tu gimnasio de Go. No basta con leerlo. **Ejecuta el código**, **rompe cosas**, **añade features**, **mide con benchmarks**. Solo así pasarás de junior a senior.

¿Listo para ser un experto en Go? **Let's Go! 🚀**
