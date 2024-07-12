package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_FACTOR = 100
const WINDOW_WIDTH = 16 * SCREEN_FACTOR
const WINDOW_HEIGHT = 9 * SCREEN_FACTOR

const TILE_MULTIPLICATOR = 2
const TILE_SIZE = 32 * TILE_MULTIPLICATOR
const TILE_HEIGHT = TILE_SIZE
const TILE_WIDTH = TILE_SIZE

const OFFSET_X = WINDOW_WIDTH/2 - TILE_WIDTH/2
const OFFSET_Y = WINDOW_HEIGHT/2
var OFFSET_VEC = createVector2(OFFSET_X, OFFSET_Y)
const USE_NEAREST = true

const i_x float32 = 1;
const i_y float32 = 0.5;
const j_x float32 = -1;
const j_y float32 = 0.5;


const (
	origin_default = iota
	origin_center = iota
	origin_down = iota
)


type Vector2 struct {
	rl.Vector2
}


//go:embed sprites/isoBlock.png
var isoBytes []byte

//go:embed sprites/FireSamurai.png
var samurai []byte


func (grid Vector2) gridToScreenCoord() Vector2 {
	var res = createVector2(
		grid.X * i_x * 0.5 * TILE_WIDTH + grid.Y * j_x * 0.5 * TILE_WIDTH, 
		grid.X * i_y * 0.5 * TILE_HEIGHT + grid.Y * j_y * 0.5 * TILE_HEIGHT,
	) 

	return res.addV(OFFSET_VEC)
}

func gridToScreenCoord(x int, y int) Vector2 {
	var v = createVector2(
		float32(x) * i_x * TILE_WIDTH/2  + float32(y) * j_x * TILE_WIDTH/2, 
		float32(x) * i_y * 0.5 * TILE_HEIGHT + float32(y) * j_y * 0.5 * TILE_HEIGHT,
	)

	return v.addV(OFFSET_VEC)
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
	
	return tmp.subV(offset)
}

func LoadImageFromMemory(data []byte) *rl.Image {
    // Get a pointer to the image data and its size
    dataSize := len(data)

    return rl.LoadImageFromMemory(".png", data, int32(dataSize))
}

func (v1 Vector2) addV(v2 Vector2) Vector2 {
	return Vector2{rl.Vector2Add(v1.Vector2, v2.Vector2)}
}

func (v1 Vector2) subV(v2 Vector2) Vector2 {
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

func (v Vector2) mul(amount float32) Vector2 {
	v.X = v.X * amount
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

type Sprite struct {
	Texture rl.Texture2D
	FrameRect rl.Rectangle
	length int
	isFlipped bool
	originPixel int
}


var spriteMap map[string]Sprite = make(map[string]Sprite)



func createSpriteFromMemory(data []byte, len int, originPixel int) Sprite {
	var tex rl.Texture2D = rl.LoadTextureFromImage(LoadImageFromMemory(data))
	var frameRec rl.Rectangle = rl.Rectangle{ X: 0, Y: 0, Width: float32(tex.Width/int32(len)), Height: float32(tex.Height) };
	return Sprite{Texture: tex, FrameRect: frameRec, length: len, isFlipped: false, originPixel: originPixel}
}

func (spr Sprite) origin(origin_type int) Vector2 {
	// origin of the image (anchor)
	if origin_type == origin_down {
		if spr.isFlipped {
			return createVector2((spr.FrameRect.Width - float32(spr.originPixel))*TILE_MULTIPLICATOR, spr.FrameRect.Height * TILE_MULTIPLICATOR)
		}
	
		return createVector2(float32(spr.originPixel)*TILE_MULTIPLICATOR, spr.FrameRect.Height * TILE_MULTIPLICATOR)
	}

	if origin_type == origin_center {
		return createVector2(spr.FrameRect.Width*TILE_MULTIPLICATOR/2, spr.FrameRect.Height * TILE_MULTIPLICATOR/2)
	}	

	return createVector2(0, 0)
}

func (spr *Sprite) selectSprite(num int) {
	if num < spr.length {
		if spr.isFlipped {
			spr.FrameRect.X = float32((spr.length - num - 1) * int(spr.Texture.Width) / spr.length)
		} else {
			spr.FrameRect.X = float32(num * int(spr.Texture.Width) / spr.length)
		}
	} else {
		fmt.Println("Could not select Sprite")
	}
}

func (spr *Sprite) flip() {
	var im = rl.LoadImageFromTexture(spr.Texture)
	rl.ImageFlipHorizontal(im)
	spr.Texture = rl.LoadTextureFromImage(im)
	spr.isFlipped = !spr.isFlipped
}

func (spr Sprite) draw(x int, y int, origin_type int) {
	var coords = gridToScreenCoord(x, y) 
	if origin_type == origin_down {
		coords = coords.addY(TILE_SIZE/4).addX(TILE_SIZE/2)
	}
	var dest = rl.NewRectangle(coords.X, coords.Y,  spr.FrameRect.Width * TILE_MULTIPLICATOR, spr.FrameRect.Height * TILE_MULTIPLICATOR)
	rl.DrawTexturePro(spr.Texture, spr.FrameRect, dest, spr.origin(origin_type).Vector2, 0, rl.White)
}

func loadSprites() {
	spriteMap["tile1"] = createSpriteFromMemory(isoBytes, 1, 0)
	spriteMap["samurai"] = createSpriteFromMemory(samurai, 4, 12)
}




func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "raygogame")
	defer rl.CloseWindow()

	loadSprites()

	var sam = spriteMap["samurai"]
	var isoBlock = spriteMap["tile1"]

	var framesCounter = 0
	var framesSpeed = 8
	var currentFrame = 0


	var posX = 5
	var posY = 5

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		framesCounter++;

        if framesCounter >= (60/framesSpeed) {
            framesCounter = 0;
            currentFrame++;

            if currentFrame > 3 {
				currentFrame = 0
			}
			sam.selectSprite(currentFrame)
        }

		if rl.IsKeyPressed(rl.KeyRight) { 
			framesSpeed++ 
		} else if rl.IsKeyPressed(rl.KeyLeft) {
			framesSpeed--
		}

		if rl.IsKeyPressed(rl.KeyR) {
			sam.flip()
		}

		if rl.IsKeyPressed(rl.KeyW) {
			posX -=1
		}

		if rl.IsKeyPressed(rl.KeyS) {
			posX +=1
		}

		if rl.IsKeyPressed(rl.KeyD) {
			posY -=1
		}

		if rl.IsKeyPressed(rl.KeyA) {
			posY +=1
		}

		if framesSpeed <= 0 {
			framesSpeed = 1
		}

		rl.BeginDrawing()

		

		rl.ClearBackground(color.RGBA{0x18, 0x18, 0x18, 1})


		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				//rl.DrawTextureV(isoBlock, gridToScreenCoord(x, y).Vector2, rl.RayWhite)
				isoBlock.draw(x, y, origin_default)
			} 
			
		}

		//rl.DrawTextureV(samurai_anim[currentFrame], , rl.RayWhite)

		// for i := 0; i < 4; i++ {
		// 	mySprite.selectSprite(i)
		// 	rl.DrawTextureRec(mySprite.Texture, mySprite.FrameRect, createVector2Int(WINDOW_WIDTH/2 + i * TILE_WIDTH + TILE_WIDTH, 300).Vector2, rl.White);
		// }

		// var scrn = gridToScreenCoord(5, 5).addY(TILE_SIZE/4).addX(TILE_SIZE/2)
		// var dest = rl.NewRectangle(scrn.X, scrn.Y,  mySprite.FrameRect.Width * TILE_MULTIPLICATOR, mySprite.FrameRect.Height * TILE_MULTIPLICATOR)
		// rl.DrawTexturePro(mySprite.Texture, mySprite.FrameRect, dest, mySprite.origin().Vector2, 0, rl.White)

		sam.draw(posX, posY, origin_down)
		
		//drawHelperGrid()
		//rl.DrawLine(WINDOW_WIDTH/2, 0, WINDOW_WIDTH/2, WINDOW_HEIGHT, rl.Red)
		
		rl.DrawFPS(0,0)
		var inc_offsett int32 = 20
		var inc int32 = 20
		rl.DrawText("Controls: Left: Framespeed--; Right: Framespeed++; WASD: Movement; R: Flip Sprite", 100, 0, 20, rl.Green)
		inc += inc_offsett
		var mouseVec = Vector2{ rl.GetMousePosition() }
		rl.DrawText(fmt.Sprintf("Mouse X: %.2f - Mouse Y: %.2f", mouseVec.X, mouseVec.Y), 0, inc, 20, rl.Green)
		inc += inc_offsett
		var newMouseVec = screenToGridCoord(mouseVec)
		drawHelperGridMouse(newMouseVec, 3)
		rl.DrawText(fmt.Sprintf("Mouse X: %.2f - Mouse Y: %.2f", newMouseVec.X, newMouseVec.Y), 0, inc, 20, rl.Green)
		inc += inc_offsett
		rl.DrawText(fmt.Sprintf("Mouse X: %d - Mouse Y: %d", toNearestGrid(newMouseVec.X), toNearestGrid(newMouseVec.Y)), 0, inc, 20, rl.Green)
		inc += inc_offsett

		rl.DrawText(fmt.Sprintf("Current Frame: %d", currentFrame), 0, inc, 20, rl.Green)
		inc += inc_offsett

		rl.EndDrawing()
	}
}