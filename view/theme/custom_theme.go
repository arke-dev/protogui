package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct{}

func (c CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameDisabled:
		if variant == theme.VariantLight {
			return color.Black
		}
		return color.White
	// case theme.ColorNameForeground:
	// 	return color.Black
	// case theme.ColorNameInputBackground:
	// 	return color.White
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (c CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (c CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	}

	return theme.DefaultTheme().Size(name)
}
