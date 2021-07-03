package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"github.com/tdewolff/canvas/text"
)

var fontFamily *canvas.FontFamily

func main() {
	fontFamily = canvas.NewFontFamily("times")
	if err := fontFamily.LoadLocalFont("DejaVuSerif", canvas.FontRegular); err != nil {
		panic(err)
	}

	c := canvas.New(200, 230)
	ctx := canvas.NewContext(c)
	ctx.SetFillColor(canvas.White)
	ctx.DrawPath(0, 0, canvas.Rectangle(c.W, c.H))
	// draw(ctx)
	draw2(ctx)
	renderers.Write("out.png", c, canvas.DPMM(5.0))
}

var lorem = []string{
	`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla malesuada fringilla libero vel ultricies. Phasellus eu lobortis lorem. Phasellus eu cursus mi. Sed enim ex, ornare et velit vitae, sollicitudin volutpat dolor. Sed aliquam sit amet nisi id sodales. Aliquam erat volutpat. In hac habitasse platea dictumst. Pellentesque luctus varius nibh sit amet porta. Vivamus tempus, enim ut sodales aliquet, magna massa viverra eros, nec gravida risus ipsum a erat. Etiam dapibus sem augue, at porta nisi dictum non. Vestibulum quis urna ut ligula dapibus mollis eu vel nisl. Vestibulum lorem dolor, eleifend lacinia fringilla eu, pulvinar vitae metus.`,
	`Morbi dapibus purus vel erat auctor, vehicula tempus leo maximus. Aenean feugiat vel quam sit amet iaculis. Fusce et justo nec arcu maximus porttitor. Cras sed aliquam ipsum. Sed molestie mauris nec dui interdum sollicitudin. Nulla id egestas massa. Fusce congue ante. Interdum et malesuada fames ac ante ipsum primis in faucibus. Praesent faucibus tellus eu viverra blandit. Vivamus mi massa, hendrerit in commodo et, luctus vitae felis.`,
	`Quisque semper aliquet augue, in dignissim eros cursus eu. Pellentesque suscipit consequat nibh, sit amet ultricies risus. Suspendisse blandit interdum tortor, consectetur tristique magna aliquet eu. Aliquam sollicitudin eleifend sapien, in pretium nisi. Sed tempor eleifend velit quis vulputate. Donec condimentum, lectus vel viverra pharetra, ex enim cursus metus, quis luctus est urna ut purus. Donec tempus gravida pharetra. Sed leo nibh, cursus at hendrerit at, ultricies a dui. Maecenas eget elit magna. Quisque sollicitudin odio erat, sed consequat libero tincidunt in. Nullam imperdiet, neque quis consequat pellentesque, metus nisl consectetur eros, ut vehicula dui augue sed tellus.`,
	` Vivamus varius ex sed nisi vestibulum, sit amet tincidunt ante vestibulum. Nullam et augue blandit dolor accumsan tempus. Quisque at dictum elit, id ullamcorper dolor. Nullam feugiat mauris eu aliquam accumsan.`,
}

var y = 205.0

func drawText(c *canvas.Context, x float64, text *canvas.Text) {
	h := text.Bounds().H
	c.DrawText(x, y, text)
	y -= h + 10.0
}

func draw(c *canvas.Context) {
	c.SetFillColor(canvas.Black)

	headerFace := fontFamily.Face(28.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	textFace := fontFamily.Face(12.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	drawText(c, 30.0, canvas.NewTextBox(headerFace, "Document Example", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[0], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))

	lenna, err := os.Open("lenna.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(lenna)
	if err != nil {
		panic(err)
	}
	imgDPMM := 15.0
	imgWidth := float64(img.Bounds().Max.X) / imgDPMM
	imgHeight := float64(img.Bounds().Max.Y) / imgDPMM
	c.DrawImage(170.0-imgWidth, y-imgHeight, img, canvas.DPMM(imgDPMM))

	drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[1], 140.0-imgWidth-10.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))
	// drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[2], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))
	// drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[3], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))
}

func drawText2(c *canvas.Context, x, y float64, face *canvas.FontFace, rich *canvas.RichText) {
	metrics := face.Metrics()
	width, height := 90.0, 32.0

	text := rich.ToText(width, height, canvas.Justify, canvas.Top, 0.0, 0.0)

	c.SetFillColor(color.RGBA{192, 0, 64, 255})
	c.DrawPath(x, y, text.Bounds().ToPath())
	c.SetFillColor(color.RGBA{51, 51, 51, 51})
	c.DrawPath(x, y, canvas.Rectangle(width, -metrics.LineHeight))
	c.SetFillColor(color.RGBA{0, 0, 0, 51})
	c.DrawPath(x, y+metrics.CapHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.CapHeight-metrics.Descent))
	c.DrawPath(x, y+metrics.XHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.XHeight))

	c.SetFillColor(canvas.Black)
	c.DrawPath(x, y, canvas.Rectangle(width, -height).Stroke(0.2, canvas.RoundCap, canvas.RoundJoin))
	c.DrawText(x, y, text)
}

func draw2(c *canvas.Context) {
	// Draw a comprehensive text box
	pt := 14.0
	face := fontFamily.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	rt := canvas.NewRichText(face)
	rt.Add(face, "Lorem dolor ipsum ")
	rt.Add(fontFamily.Face(pt, canvas.White, canvas.FontBold, canvas.FontNormal), "confiscator")
	rt.Add(face, " cur\u200babitur ")
	rt.Add(fontFamily.Face(pt, canvas.Black, canvas.FontItalic, canvas.FontNormal), "mattis")
	rt.Add(face, " dui ")
	rt.Add(fontFamily.Face(pt, canvas.Black, canvas.FontBold|canvas.FontItalic, canvas.FontNormal), "tellus")
	rt.Add(face, " vel. Proin ")
	rt.Add(fontFamily.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontUnderline), "sodales")
	rt.Add(face, " eros vel ")
	rt.Add(fontFamily.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontSineUnderline), "nibh")
	rt.Add(face, " fringilla pellen\u200btesque eu cillum. ")

	face = fontFamily.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	face.Language = "ru"
	face.Script = text.Cyrillic
	rt.Add(face, "дёжжэнтиюнт холст ")

	face = fontFamily.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	face.Language = "ar"
	face.Script = text.Arabic
	face.Direction = text.RightToLeft
	rt.Add(face, "تسجّل يتكلّم ")

	face = fontFamily.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	face.Language = "hi"
	face.Script = text.Devanagari
	rt.Add(face, "हालाँकि प्र ")

	drawText2(c, 5, 95, face, rt)

	// Draw the word Stroke being stroked
	face = fontFamily.Face(80.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	p, _, _ := face.ToPath("Stroke")
	c.DrawPath(100, 5, p.Stroke(0.75, canvas.RoundCap, canvas.RoundJoin))

	// Draw an elliptic arc being dashed
	ellipse, err := canvas.ParseSVG(fmt.Sprintf("A10 30 30 1 0 30 0z"))
	if err != nil {
		panic(err)
	}
	c.SetFillColor(canvas.Whitesmoke)
	c.DrawPath(110, 40, ellipse)

	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(0.75)
	c.SetStrokeCapper(canvas.RoundCap)
	c.SetStrokeJoiner(canvas.RoundJoin)
	c.SetDashes(0.0, 2.0, 4.0, 2.0, 2.0, 4.0, 2.0)
	//ellipse = ellipse.Dash(0.0, 2.0, 4.0, 2.0).Stroke(0.5, canvas.RoundCap, canvas.RoundJoin)
	c.DrawPath(110, 40, ellipse)
	c.SetStrokeColor(canvas.Transparent)
	c.SetDashes(0.0)

	// Draw a raster image
	lenna, err := os.Open("lenna.png")
	if err != nil {
		panic(err)
	}
	img, err := canvas.NewPNGImage(lenna)
	if err != nil {
		panic(err)
	}
	c.Rotate(5)
	c.DrawImage(50.0, 0.0, img, 15)
	c.SetView(canvas.Identity.Translate(0, 0).Scale(0.5, 0.5))

	// Draw an closed set of points being smoothed
	polyline := &canvas.Polyline{}
	polyline.Add(0.0, 0.0)
	polyline.Add(30.0, 0.0)
	polyline.Add(30.0, 15.0)
	polyline.Add(0.0, 30.0)
	polyline.Add(0.0, 0.0)
	c.SetFillColor(canvas.Seagreen)
	c.FillColor.R = byte(float64(c.FillColor.R) * 0.25)
	c.FillColor.G = byte(float64(c.FillColor.G) * 0.25)
	c.FillColor.B = byte(float64(c.FillColor.B) * 0.25)
	c.FillColor.A = byte(float64(c.FillColor.A) * 0.25)
	c.SetStrokeColor(canvas.Seagreen)
	c.DrawPath(155, 35, polyline.Smoothen())

	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(0.5)
	c.DrawPath(155, 35, polyline.ToPath())
	c.SetStrokeWidth(0.75)
	for _, coord := range polyline.Coords() {
		c.DrawPath(155, 35, canvas.Circle(2.0).Translate(coord.X, coord.Y))
	}

	// Draw a open set of points being smoothed
	polyline = &canvas.Polyline{}
	polyline.Add(0.0, 0.0)
	polyline.Add(20.0, 10.0)
	polyline.Add(40.0, 30.0)
	polyline.Add(60.0, 40.0)
	polyline.Add(80.0, 20.0)
	c.SetStrokeColor(canvas.Dodgerblue)
	c.DrawPath(10, 15, polyline.Smoothen())
	c.SetStrokeColor(canvas.Black)
	for _, coord := range polyline.Coords() {
		c.DrawPath(10, 15, canvas.Circle(2.0).Translate(coord.X, coord.Y))
	}
}
