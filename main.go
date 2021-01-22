package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	// Settings for the game screen
	screenWidth  = 960
	screenHeight = 540
)

var (
	leftSprite      *ebiten.Image
	rightSprite     *ebiten.Image
	idleSprite      *ebiten.Image
	backgroundImage *ebiten.Image
)

func Init() {
	// Preload main character images for each direction
	// Right Image
	var err error
	rightSprite, _, err = ebitenutil.NewImageFromFile("Assets/right.png")
	if err != nil {
		log.Fatal(err)
	}

	// Left Image
	leftSprite, _, err = ebitenutil.NewImageFromFile("Assets/left.png")
	if err != nil {
		log.Fatal(err)
	}

	// Middle Image
	idleSprite, _, err = ebitenutil.NewImageFromFile("Assets/middle.png")
	if err != nil {
		log.Fatal(err)
	}

	// Background Image
	backgroundImage, _, err = ebitenutil.NewImageFromFile("Assets/background.png")
	if err != nil {
		log.Fatal(err)
	}

}

const (
	unit    = 160
	groundY = 400
)

type char struct {
	x  int
	y  int
	vx int // Horizontal velocity
	vy int // Vertical velocity
}

// Velocity is equivalent to a specification of an object's speed and direction of motion (e.g. 60 km/h to the North)
// vx represents the velocity on the horizontal axis and is how quickly the object’s x value is changing in value moving left to right.
// vy represents the velocity in the vertical axis, and is how quickly the object’s y changes in value moving up and down.

func (c *char) tryJump() {
	// To make the character jump

	c.vy = -5 * unit // Jump level	[- Increase, + Decrease] 			// Measurement in pixel per second
}

func (c *char) update() {
	c.x += c.vx // "+=" Add and assign
	c.y += c.vy
	if c.y > groundY*unit {
		c.y = groundY * unit
	}
	if c.vx > 0 {
		c.vx -= 16 // Stop speed level       // "-=" Subtract and assign
	} else if c.vx < 0 {
		c.vx += 16
	}
	if c.vy < 20*unit {
		c.vy += 16 // Jump Height [+ Decrease, - Increase]
	}
}

// Draw() draws the screen based on the current game state
func (c *char) draw(screen *ebiten.Image) {
	s := idleSprite
	switch {
	case c.vx > 0:
		s = rightSprite
	case c.vx < 0:
		s = leftSprite
	}

	// GeoM is used to rotate, scale and move an image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)                                 // Scale matrix by
	op.GeoM.Translate(float64(c.x)/unit, float64(c.y)/unit) // (0,0) is the top-left corner
	screen.DrawImage(s, op)
}

// Game implements the ebiten.Game interface
type Game struct {
	mario *char
}

// Update() updates the game logic by 1 tick (60 ticks per second)
func (g *Game) Update() error {
	if g.mario == nil {
		g.mario = &char{x: 50 * unit, y: groundY * unit}
	}

	// Controls
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) { // "||" OR Operator
		g.mario.vx = -2 * unit
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.mario.vx = 2 * unit
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.mario.tryJump()
	}
	g.mario.update()
	return nil
}

// Draw() draws the screen based on the current game state
func (g *Game) Draw(screen *ebiten.Image) {

	// Draws Background Image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(backgroundImage, op)

	// Draws the Mario
	g.mario.draw(screen)

	// Show the message
	msg := fmt.Sprintf("TPS: %0.2f\nPress the space key to jump.", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

// Layout() gets the window size and return the game logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Interface to run the game
func main() {
	Init()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mario Game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

// Game Loop : It’s an infinite loop that updates and redraws to animate gameplay
