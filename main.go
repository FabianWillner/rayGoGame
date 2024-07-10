package main

import (
	_ "embed"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_FACTOR = 100
const WINDOW_WIDTH = 16 * SCREEN_FACTOR
const WINDOW_HEIGHT = 9 * SCREEN_FACTOR


const TILE_HEIGHT = 64
const TILE_WIDTH = 64

const OFFSET_X = WINDOW_WIDTH/2 - TILE_WIDTH/2
const OFFSET_Y = WINDOW_HEIGHT/2
var OFFSET_VEC = createVector2(OFFSET_X, OFFSET_Y)
const USE_NEAREST = false

const i_x float32 = 1;
const i_y float32 = 0.5;
const j_x float32 = -1;
const j_y float32 = 0.5;

type Vector2 struct {
	rl.Vector2
}

//go:embed sprites/isoBlock.png
var isoBytes []byte

func (grid Vector2) gridToScreenCoord() Vector2 {
	var res = createVector2(
		grid.X * i_x * 0.5 * TILE_WIDTH + grid.Y * j_x * 0.5 * TILE_WIDTH, 
		grid.X * i_y * 0.5 * TILE_HEIGHT + grid.Y * j_y * 0.5 * TILE_HEIGHT,
	) 

	return Vector2Add(res, OFFSET_VEC)
}

func gridToScreenCoord(x float32, y float32) Vector2 {
	var v = createVector2(
		x * i_x * TILE_WIDTH/2  + y * j_x * TILE_WIDTH/2, 
		x * i_y * 0.5 * TILE_HEIGHT + y * j_y * 0.5 * TILE_HEIGHT,
	)

	return Vector2Add(v, OFFSET_VEC)
}

func invert_matrix(a float32, b float32, c float32, d float32) (float32, float32, float32, float32) {
	// Determinant 
	var det = (1 / (a * d - b * c));
	
	return det * d,  det * -b, det * -c, det * a
}

func screenToGridCoord(screen Vector2) Vector2 {
	//screen = rl.Vector2{X: screen.X + OFFSET_X, Y: screen.Y + OFFSET_Y}
	var a = float32(i_x * 0.5 * TILE_WIDTH);
	var b = float32(j_x * 0.5 * TILE_WIDTH);
	var c = float32(i_y * 0.5 * TILE_HEIGHT);
	var d = float32(j_y * 0.5 * TILE_HEIGHT);
	
	a, b, c, d = invert_matrix(a, b, c, d);

	var tmp = createVector2(
		screen.X * a + screen.Y * b, 
		screen.X * c + screen.Y * d,
	)
		

	var offset = createVector2(
		WINDOW_WIDTH/2 * a + WINDOW_HEIGHT/2 * b, 
		WINDOW_WIDTH/2 * c + WINDOW_HEIGHT/2 * d,
	)
	
	return Vector2Subtract(tmp, offset)
  }

func loadTile(_ string) rl.Texture2D {
	// Tiles cant be loaded globally and must be loaded after rl.InitWindow
	//var im = rl.LoadImage(path)
	var im = LoadImageFromMemory(isoBytes)
	if im.Height != TILE_HEIGHT || im.Width != TILE_WIDTH {
		if USE_NEAREST {
			rl.ImageResizeNN(im, TILE_WIDTH, TILE_HEIGHT)
		} else {
			rl.ImageResize(im, TILE_WIDTH, TILE_HEIGHT)
		}
	}

	return rl.LoadTextureFromImage(im)
}

func LoadImageFromMemory(data []byte) *rl.Image {
    // Get a pointer to the image data and its size
    dataSize := len(data)

    // Call Raylib's LoadImageFromMemory function
    img := rl.LoadImageFromMemory(".png", data, int32(dataSize))
    return img
}

func Vector2Add(v1 Vector2, v2 Vector2) Vector2 {
	return Vector2{rl.Vector2Add(v1.Vector2, v2.Vector2)}
}

func Vector2Subtract(v1 Vector2, v2 Vector2) Vector2 {
	return Vector2{rl.Vector2Subtract(v1.Vector2, v2.Vector2)}
}

func (v Vector2) add(amount float32) Vector2 {
	v.X = v.X + amount
	v.Y = v.Y + amount
	return v
}

func (v Vector2) addX(amount float32) Vector2 {
	v.X = v.X + amount
	return v
}

func (v Vector2) addY(amount float32) Vector2 {
	v.Y = v.Y + amount
	return v
}

func createVector2(x float32, y float32) Vector2 {
	return Vector2{rl.Vector2{X: x, Y: y}}
}

func createVector2Int(x int, y int) Vector2 {
	return createVector2(float32(x), float32(y))
}

func drawHelperGrid() {
	var x_0 = -5
	var y_0 = -5
	for x := x_0; x < 20; x++ {
		for y := y_0; y < 20; y++ {
			var start = createVector2Int(x, y).gridToScreenCoord().addX(TILE_WIDTH/2)
			var end = createVector2Int(20, y).gridToScreenCoord().addX(TILE_WIDTH/2)
			var end2 = createVector2Int(x, 20).gridToScreenCoord().addX(TILE_WIDTH/2)
			if x != x_0 {
				rl.DrawLineV(start.Vector2, end2.Vector2, rl.Lime)
			}
			if y != y_0 {
				rl.DrawLineV(start.Vector2, end.Vector2, rl.Lime)
			}
			
			
		}
	}
}

func drawHelperGridMouse(mouse Vector2, size int) {
	var mouse_x = toNearestGrid(mouse.X)
	var mouse_y = toNearestGrid(mouse.Y)

	var x_0 = mouse_x - size
	var y_0 = mouse_y - size
	for x := x_0; x < mouse_x + size; x++ {
		for y := y_0; y < mouse_y + size; y++ {
			var start = createVector2Int(x, y).gridToScreenCoord().addX(TILE_WIDTH/2)
			var end = createVector2Int(mouse_x + size, y).gridToScreenCoord().addX(TILE_WIDTH/2)
			var end2 = createVector2Int(x, mouse_y + size).gridToScreenCoord().addX(TILE_WIDTH/2)
			if x != x_0 {
				rl.DrawLineV(start.Vector2, end2.Vector2, rl.Lime)
			}
			if y != y_0 {
				rl.DrawLineV(start.Vector2, end.Vector2, rl.Lime)
			}
			
			
		}
	}
}

func toNearestGrid(a float32) int {
	return int(math.Floor(float64(a)))
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
				rl.DrawTextureV(isoBlock, gridToScreenCoord(float32(x), float32(y)).Vector2, rl.RayWhite)
			} 
			
		}

		
		//drawHelperGrid()
		
		rl.DrawFPS(0,0)
		var mouseVec = Vector2{ rl.GetMousePosition() }
		rl.DrawText(fmt.Sprintf("Mouse X: %.2f - Mouse Y: %.2f", mouseVec.X, mouseVec.Y), 0, 50, 20, rl.Green)
		var newMouseVec = screenToGridCoord(mouseVec)
		drawHelperGridMouse(newMouseVec, 3)
		rl.DrawText(fmt.Sprintf("Mouse X: %.2f - Mouse Y: %.2f", newMouseVec.X, newMouseVec.Y), 0, 100, 20, rl.Green)
		rl.DrawText(fmt.Sprintf("Mouse X: %d - Mouse Y: %d", toNearestGrid(newMouseVec.X), toNearestGrid(newMouseVec.Y)), 0, 150, 20, rl.Green)

		rl.EndDrawing()
	}
}