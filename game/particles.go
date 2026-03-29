package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// particle is a single pixel with velocity and remaining lifetime.
type particle struct {
	x, y   float32
	vx, vy float32
	life    int // frames remaining
	maxLife int
	r, g, b uint8
}

// ParticleSystem holds all active particles.
type ParticleSystem struct {
	particles []particle
}

// SpawnParticles emits count pixels of the given color bursting outward from (worldX, worldY).
func (ps *ParticleSystem) SpawnParticles(worldX, worldY float32, count int, r, g, b uint8, rng *rand.Rand) {
	for range count {
		angle := rng.Float64() * 2 * math.Pi
		speed := float32(rng.Float64()*2.0 + 0.5)
		life := 18 + rng.Intn(14)
		ps.particles = append(ps.particles, particle{
			x: worldX, y: worldY,
			vx: float32(math.Cos(angle)) * speed,
			vy: float32(math.Sin(angle)) * speed,
			life: life, maxLife: life,
			r: r, g: g, b: b,
		})
	}
}

// SpawnBlood emits count red pixels bursting outward from (worldX, worldY).
func (ps *ParticleSystem) SpawnBlood(worldX, worldY float32, count int, rng *rand.Rand) {
	ps.SpawnParticles(worldX, worldY, count, 200, 30, 30, rng)
}

// Update advances all particles and removes dead ones.
func (ps *ParticleSystem) Update() {
	alive := ps.particles[:0]
	for i := range ps.particles {
		p := &ps.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vx *= 0.92 // drag
		p.vy *= 0.92
		p.life--
		if p.life > 0 {
			alive = append(alive, *p)
		}
	}
	ps.particles = alive
}

// Draw renders all particles to the screen.
func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for _, p := range ps.particles {
		alpha := uint8(255 * p.life / p.maxLife)
		col := color.RGBA{p.r, p.g, p.b, alpha}
		vector.DrawFilledRect(screen, p.x, p.y, 2, 2, col, false)
	}
}
