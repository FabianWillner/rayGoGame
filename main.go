package main

import rl "github.com/gen2brain/raylib-go/raylib"

const WINDOW_HEIGHT = 450
const WINDOW_WIDTH = 800

const TILE_HEIGHT = 32
const TILE_WIDTH = 32

const OFFSET_X = TILE_WIDTH/2

var isoBlock = rl.LoadImage("./sprites/isoBlock.png")
//var isoBlock = rl.LoadTexture("./sprites/isoBlock.png")

func gridToScreenCoord(grid rl.Vector2) rl.Vector2 {
	var x = (grid.X * TILE_WIDTH / 2 + grid.Y * -1 * TILE_WIDTH / 2) - OFFSET_X  + WINDOW_WIDTH/2
	var y = (grid.X * 0.5 * TILE_HEIGHT /2 + grid.Y * 0.5 * TILE_HEIGHT/2)  + WINDOW_HEIGHT/2

	return rl.Vector2{X: x, Y: y}
}


func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "raygogame")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	var isoText = rl.LoadTextureFromImage(isoBlock)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				rl.DrawTextureV(isoText, gridToScreenCoord(rl.Vector2{X: float32(x), Y: float32(y)}), rl.RayWhite)
			} 
			
		}
		
		rl.DrawFPS(0,0)
		rl.EndDrawing()
	}
}