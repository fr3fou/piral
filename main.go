package main

import (
	"fmt"
	"math"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(300, 125, "Piral Launcher")

	fullscreen := false

	comboActive := 0
	res := [][]int32{
		{1920, 1080},
		{1366, 768},
		{800, 600},
	}
	comboText := []string{"1920x1080", "1366x768", "800x600"}

	n := float32(10_000)

	raygui.LoadGuiStyle("zahnrad.style")
	rl.SetTargetFPS(60)

	start := false
	primes := []int{}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(raygui.BackgroundColor())

		if start {

			rl.DrawText("Loading...", (300-rl.MeasureText("Loading...", 20))/2, (125-20)/2, 20, raygui.TextColor())
			rl.EndDrawing()
			primes = findPrimes(int(n))
			rl.CloseWindow()
			break
		} else {
			raygui.Label(rl.NewRectangle(40, 25, 0, 20), "Fullscreen")
			fullscreen = raygui.CheckBox(rl.NewRectangle(20, 25, 20, 20), fullscreen)

			raygui.Label(rl.NewRectangle(125, 125-20-20, 0, 20), fmt.Sprintf("%d primes", int(n)))
			n = raygui.SliderBar(rl.NewRectangle(20, 125-20-20, 100, 20), n, 10_000, 1_000_000)

			comboActive = raygui.ComboBox(rl.NewRectangle(300-100-20-30, 25, 100, 20), comboText, comboActive)

			start = raygui.Button(rl.NewRectangle(300-40-20+1, 125-20-20, 40, 20), "Start!")
		}

		rl.EndDrawing()
	}

	if start {
		piral(res[comboActive][0], res[comboActive][1], fullscreen, primes)
	}
}

func piral(screenWidth, screenHeight int32, fullscreen bool, primes []int) {

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	rl.InitWindow(screenWidth, screenHeight, "Piral")
	if fullscreen {
		rl.MaximizeWindow()
		rl.ToggleFullscreen()
	}

	i := 0
	scale := 0.1
	theta := 0.0
	delta := 4
	generating := false
	auto := true
	for !rl.WindowShouldClose() {
		scale += float64(rl.GetMouseWheelMove()) * 0.001

		if rl.IsKeyReleased(rl.KeySpace) {
			generating = !generating
		}

		if rl.IsKeyReleased(rl.KeyR) {
			i = 0
			scale = 0.1
		}

		if rl.IsKeyReleased(rl.KeyZ) {
			auto = !auto
		}

		if generating {
			if auto {
				scale = constrain(scale/1.0005, 1, 0.002)
			}

			l := len(primes)
			if i+delta <= l {
				i += delta
			} else {
				i += len(primes) - i
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		for _, prime := range primes[:i] {
			p := float64(prime)
			sin, cos := math.Sincos(p)

			vec := rl.Vector2Scale(
				rl.Vector2{
					X: float32(p * cos),
					Y: float32(p * sin),
				},
				float32(scale),
			)

			origin := rl.Vector2{
				X: float32(screenWidth / 2),
				Y: float32(screenHeight / 2),
			}

			x, y := vec.X, vec.Y
			sin, cos = math.Sincos(theta)

			vec.X = float32(float64(x)*cos - float64(y)*sin)
			vec.Y = float32(float64(x)*sin + float64(y)*cos)

			rl.DrawCircleV(rl.Vector2Add(vec, origin), 2, rl.SkyBlue)
		}

		theta += 0.001

		rl.EndDrawing()
	}
}

func findPrimes(n int) []int {
	if n < 2 {
		return []int{}
	}

	if n < 5 {
		return []int{2}
	}

	// Create a list of consecutive integers from 2 through n
	// (2, 3, 4, ..., n)
	nums := make([]int, n, n)
	primes := make([]int, 0)

	// Fill up the array with all the numbers
	for i := 2; i < n; i++ {
		nums[i] = i
	}

	// Initially, let p equal 2, the smallest prime number
	p := 2
	primes = append(primes, p)

	for {
		np := primeSieve(p, p, n, nums)

		// If np hasn't changed, our work is done
		if np == p {
			break
		}

		p = np
		primes = append(primes, p)
	}

	// Remove the last 2 (to compensate for starting from 2)
	primes = primes[:len(primes)-2]

	return primes
}

// primeSieve takes in p (current num)
// i (incrementer / original p)
// n (maximum number)
// nums (list of the numbers)
func primeSieve(p, i, n int, nums []int) int {
	// Enumerate the multiples of p by counting in increments of p from 2p to n,
	// and mark them in the list (these will be 2p, 3p, 4p, ...;
	// the p itself should not be marked).

	// If we have reached the maximum,
	// return the first number greater than p that is not marked
	if p+i >= n {
		for _, v := range nums {
			if v > i && v != -1 {
				return v
			}
		}

		return p
	}

	// Mark in the list
	nums[p] = -1

	// Continue recursively
	return primeSieve(p+i, i, n, nums)
}

func scale(value float64, minFrom float64, maxFrom float64, minTo float64, maxTo float64) float64 {
	return ((maxTo-minTo)*(value-minFrom))/(maxFrom-minFrom) + minTo
}

func constrain(n, high, low float64) float64 {
	return math.Max(math.Min(n, high), low)
}

func findNums(n int) []int {
	v := make([]int, n)
	for i := 0; i < n; i++ {
		v[i] = i
	}
	return v
}
