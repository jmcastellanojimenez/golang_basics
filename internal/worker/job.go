package worker

import "fmt"

// Job representa una unidad de trabajo en nuestro sistema. 📝
// En este caso, simula el procesamiento de un usuario.
type Job struct {
	ID    int
	Name  string
	Email string
}

// GenerateJobs crea una lista de N trabajos ficticios para procesar. 🏭
// Es nuestra fuente de datos para el pipeline concurrente.
func GenerateJobs(n int) []Job {
	jobs := make([]Job, n)
	for i := 0; i < n; i++ {
		jobs[i] = Job{
			ID:    i,
			Name:  fmt.Sprintf("user_%04d", i),
			Email: fmt.Sprintf("user_%04d@example.com", i),
		}
	}
	return jobs
}
