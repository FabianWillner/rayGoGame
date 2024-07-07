package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_FACTOR = 100
const WINDOW_WIDTH = 16 * SCREEN_FACTOR
const WINDOW_HEIGHT = 9 * SCREEN_FACTOR


const TILE_HEIGHT = 64
const TILE_WIDTH = 64

const OFFSET_X = WINDOW_WIDTH/2 - TILE_WIDTH/2
const OFFSET_Y = WINDOW_HEIGHT/2
var OFFSET_VEC = rl.Vector2{X: OFFSET_X, Y: OFFSET_Y}
const USE_NEAREST = false

const i_x = 1;
const i_y = 0.5;
const j_x = -1;
const j_y = 0.5;


func gridToScreenCoord(grid rl.Vector2) rl.Vector2 {
	var res = rl.Vector2{
		X: (grid.X * i_x * 0.5 * TILE_WIDTH + grid.Y * j_x * 0.5 * TILE_WIDTH), // + OFFSET_X,
		Y: (grid.X * i_y * 0.5 * TILE_HEIGHT + grid.Y * j_y * 0.5 * TILE_HEIGHT), //+ OFFSET_Y,
	}

	return rl.Vector2Add(res, OFFSET_VEC)
}

func invert_matrix(a float32, b float32, c float32, d float32) (float32, float32, float32, float32) {
	// Determinant 
	var det = (1 / (a * d - b * c));
	
	return det * d,  det * -b, det * -c, det * a
}

func screenToGridCoord(screen rl.Vector2) rl.Vector2 {
	//screen = rl.Vector2{X: screen.X + OFFSET_X, Y: screen.Y + OFFSET_Y}
	var a = float32(i_x * 0.5 * TILE_WIDTH);
	var b = float32(j_x * 0.5 * TILE_WIDTH);
	var c = float32(i_y * 0.5 * TILE_HEIGHT);
	var d = float32(j_y * 0.5 * TILE_HEIGHT);
	
	a, b, c, d = invert_matrix(a, b, c, d);

	var tmp = rl.Vector2{
		X: (screen.X * a + screen.Y * b) ,
		Y: (screen.X * c + screen.Y * d) ,
	}

	var offset = rl.Vector2{
		X: WINDOW_WIDTH/2 * a + WINDOW_HEIGHT/2 * b,
		Y: WINDOW_WIDTH/2 * c + WINDOW_HEIGHT/2 * d,
	}
	
	return rl.Vector2Subtract(tmp, offset)
  }

func loadTile(path string) rl.Texture2D {
	// Tiles cant be loaded globally and must be loaded after rl.InitWindow
	var im = rl.LoadImage(path)
	if im.Height != TILE_HEIGHT || im.Width != TILE_WIDTH {
		if USE_NEAREST {
			rl.ImageResizeNN(im, TILE_WIDTH, TILE_HEIGHT)
		} else {
			rl.ImageResize(im, TILE_WIDTH, TILE_HEIGHT)
		}
	}

	return rl.LoadTextureFromImage(im)
}


func main() {
	fmt.Println("Test")
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "raygogame")
	defer rl.CloseWindow()

	var isoBlock = loadTile("./sprites/isoBlock.png")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		

		rl.ClearBackground(rl.Black)
		
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				rl.DrawTextureV(isoBlock, gridToScreenCoord(rl.Vector2{X: float32(x), Y: float32(y)}), rl.RayWhite)
			} 
			
		}
		
		rl.DrawFPS(0,0)
		var mouseVec = rl.GetMousePosition()
		rl.DrawText(fmt.Sprintf("Mouse X: %.2f - Mouse Y: %.2f", mouseVec.X, mouseVec.Y), 0, 50, 20, rl.Green)
		var newMouseVec = screenToGridCoord(mouseVec)
		rl.DrawText(fmt.Sprintf("Mouse X: %.2f - Mouse Y: %.2f", newMouseVec.X, newMouseVec.Y), 0, 100, 20, rl.Green)
		rl.EndDrawing()
	}
}