package main

import (
	"image/color"
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDTH   = 800
	HEIGHT  = 800
	vRATIO  = 4
	vWIDTH  = WIDTH / vRATIO
	vHEIGHT = HEIGHT / vRATIO
)

var (
	MIN_DIST = math.MaxFloat64
	MAX_DIST = -math.MaxFloat64
)

type Point = struct {
	x  float64
	y  float64
	dx float64
	dy float64
}

func main() {
	rl.InitWindow(WIDTH, HEIGHT, "")
	rl.SetWindowPosition(200, 200)
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	var grid [vHEIGHT][vWIDTH]float64
	for j := 0; j < vHEIGHT; j++ {
		for i := 0; i < vWIDTH; i++ {
			grid[j][i] = 0.0
		}
	}

	var voronoi []Point = make([]Point, 8)
	for i := 0; i < len(voronoi); i++ {
		voronoi[i] = Point{rand.Float64(), rand.Float64(), rand.Float64()*2 - 1, rand.Float64()*2 - 1}
	}

	for j := 0; j < vHEIGHT; j++ {
		for i := 0; i < vWIDTH; i++ {
			grid[j][i] = get_distance(i, j, voronoi)
		}
	}

	for j := 0; j < vHEIGHT; j++ {
		for i := 0; i < vWIDTH; i++ {
			grid[j][i] = (grid[j][i] - MIN_DIST) / (MAX_DIST - MIN_DIST)
		}
	}

	for !rl.WindowShouldClose() {
		for i := 0; i < len(voronoi); i++ {
			voronoi[i].x += voronoi[i].dx * 0.1 * float64(rl.GetFrameTime())
			voronoi[i].y += voronoi[i].dy * 0.1 * float64(rl.GetFrameTime())

			if voronoi[i].x > 1 || voronoi[i].x < 0 {
				voronoi[i].dx *= -1
			}
			if voronoi[i].y > 1 || voronoi[i].y < 0 {
				voronoi[i].dy *= -1
			}
		}

		for j := 0; j < vHEIGHT; j++ {
			for i := 0; i < vWIDTH; i++ {
				grid[j][i] = get_distance(i, j, voronoi)
			}
		}

		for j := 0; j < vHEIGHT; j++ {
			for i := 0; i < vWIDTH; i++ {
				v := (grid[j][i] - MIN_DIST) / (MAX_DIST - MIN_DIST)
				grid[j][i] = math.Pow(1-v, 4)
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(0x00_00_00_ff))
		for j := 0; j < vHEIGHT; j++ {
			for i := 0; i < vWIDTH; i++ {
				x := int32(i * vRATIO)
				y := int32(j * vRATIO)
				col := get_gray(int(grid[j][i] * 255.0))

				draw(x, y, rl.GetColor(col))
			}
		}
		rl.EndDrawing()
	}
}

func get_distance(x, y int, points []Point) float64 {
	dist := math.MaxFloat64

	for _, p := range points {
		norm_x := float64(x) / float64(vWIDTH)
		norm_y := float64(y) / float64(vHEIGHT)

		dist = min(dist, distance(norm_x, norm_y, p.x, p.y))
	}

	MIN_DIST = min(MIN_DIST, dist)
	MAX_DIST = max(MAX_DIST, dist)

	return dist
}

func distance(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2

	return math.Abs(dx) + math.Abs(dy)
	// return max(math.Abs(dx), math.Abs(dy))
	// return math.Sqrt(float64(dx*dx + dy*dy))
}

func get_gray(val int) uint {
	return uint(0x01010100*val | 0xff)
}

func draw(x, y int32, col color.RGBA) {
	for j := 0; j < vRATIO; j++ {
		for i := 0; i < vRATIO; i++ {
			rl.DrawPixel(x+int32(i), y+int32(j), col)
		}
	}
}
