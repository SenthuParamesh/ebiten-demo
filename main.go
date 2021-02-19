package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	// Settings for the ugame screen
	screenWidth  = 1500
	screenHeight = 400
)

const (
	//	unit    = 160
	wallLeftX  = 0
	wallRightX = 1470
	groundY    = 320
)

var (
	backgroundImage *ebiten.Image
	marioImage      *ebiten.Image
)

func loadAssets() {

	var err error

	// Background Image
	backgroundImage, _, err = ebitenutil.NewImageFromFile("imgs/world.png")
	if err != nil {
		log.Fatalln(err)
	}

	marioImage, _, err = ebitenutil.NewImageFromFile("imgs/mario.png")
	if err != nil {
		log.Fatalln(err)
	}

}

type Mario struct {
	Position struct {
		X float64
		Y float64
	}
	Sprite     *ebiten.Image
	Animations struct {
		Idle  *Animation
		Small struct {
			WalkRight *Animation
			WalkLeft  *Animation
			JumpRight *Animation
			JumpLeft  *Animation
		}
		Large struct {
			WalkRight *Animation
			WalkLeft  *Animation
			JumpRight *Animation
			JumpLeft  *Animation
		}
	}
	CurrentAnimation *Animation
}

func createMario(mi *ebiten.Image) *Mario {
	mario := &Mario{}
	mario.Position.X = 100
	mario.Position.Y = 285
	mario.Sprite = mi
	mario.Animations.Idle = newIdleAnimation(mario.Sprite)
	mario.Animations.Large.WalkRight = newLargeWalkRightAnimation(mario.Sprite)
	mario.CurrentAnimation = mario.Animations.Idle
	return mario
}

func (m *Mario) GetAnimationFrameImage() *ebiten.Image {
	return m.CurrentAnimation.GetFrameImage()
}

func (m *Mario) ChangePositionX(x float64) {
	m.Position.X += x
	if m.Position.X < wallLeftX {
		m.Position.X = wallLeftX
	}

	if m.Position.X > wallRightX {
		m.Position.X = wallRightX
	}
}

type Animation struct {
	CurrentFrame int
	Period       int
	Frames       []*Frame
	n            int
}

func newIdleAnimation(mi *ebiten.Image) *Animation {
	anim := &Animation{}

	idleFrame := &Frame{}
	idleFrame.Image = mi.SubImage(image.Rect(208, 50, 228, 85)).(*ebiten.Image)

	anim.Frames = append(anim.Frames, idleFrame)
	return anim
}

func newLargeWalkRightAnimation(mi *ebiten.Image) *Animation {
	anim := &Animation{}
	anim.Period = 10

	frame1 := &Frame{}
	frame1.Image = mi.SubImage(image.Rect(238, 50, 258, 85)).(*ebiten.Image)

	frame2 := &Frame{}
	frame2.Image = mi.SubImage(image.Rect(268, 50, 288, 85)).(*ebiten.Image)

	frame3 := &Frame{}
	frame3.Image = mi.SubImage(image.Rect(298, 50, 318, 85)).(*ebiten.Image)

	anim.Frames = []*Frame{frame3, frame2, frame1}
	return anim
}

func (a *Animation) Advance() {

	if a.Period == 0 {
		return
	}

	a.n++
	if a.n < a.Period {
		return
	}

	a.n = 0
	a.CurrentFrame = a.CurrentFrame + 1

	if a.CurrentFrame >= len(a.Frames) {
		a.CurrentFrame = 0
	}
}

func (a *Animation) GetFrameImage() *ebiten.Image {
	return a.Frames[a.CurrentFrame].Image
}

type Frame struct {
	Image *ebiten.Image
}

// Interface to run the game
func main() {

	loadAssets()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mario Game")

	game := &Game{}
	game.Mario = createMario(marioImage)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

	fmt.Println()
}

type Game struct {
	Mario *Mario
}

// Update() updates the game logic by 1 tick (60 ticks per second)
func (g *Game) Update() error {

	// Controls
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) { // "||" OR Operator
		g.Mario.ChangePositionX(-4)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Mario.ChangePositionX(4)
		g.Mario.CurrentAnimation = g.Mario.Animations.Large.WalkRight
	} else {
		g.Mario.CurrentAnimation = g.Mario.Animations.Idle
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	}

	// change animation

	g.Mario.CurrentAnimation.Advance()
	return nil
}

// Draw() draws the screen based on the current game state
func (g *Game) Draw(screen *ebiten.Image) {

	// Draws Background Image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	screen.DrawImage(backgroundImage, op)

	g.Mario.Draw(screen)

}

func (m *Mario) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)                           // Scale matrix by
	op.GeoM.Translate(m.Position.X, m.Position.Y) // (0,0) is the top-left corner
	screen.DrawImage(m.GetAnimationFrameImage(), op)
}

// Layout() gets the window size and return the game logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

/*
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
*/
// Game Loop : It’s an infinite loop that updates and redraws to animate gameplay

/*
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

*/
