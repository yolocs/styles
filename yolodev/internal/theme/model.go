package theme

const (
	Format  = "yolodev-terminal-theme"
	Version = 1
)

// Theme is the complete terminal-theme v1 document.
type Theme struct {
	Format  string         `toml:"format"`
	Version int            `toml:"version"`
	Theme   Metadata       `toml:"theme"`
	Colors  SemanticColors `toml:"colors"`
	Palette Palette        `toml:"palette"`
}

type Metadata struct {
	Name        string `toml:"name"`
	Variant     string `toml:"variant"`
	Appearance  string `toml:"appearance"`
	Description string `toml:"description"`
}

type SemanticColors struct {
	Background          string `toml:"background"`
	Foreground          string `toml:"foreground"`
	Cursor              string `toml:"cursor"`
	CursorText          string `toml:"cursor_text"`
	SelectionBackground string `toml:"selection_background"`
	SelectionForeground string `toml:"selection_foreground"`
}

type Palette struct {
	Normal ANSIColors `toml:"normal"`
	Bright ANSIColors `toml:"bright"`
}

type ANSIColors struct {
	Black   string `toml:"black"`
	Red     string `toml:"red"`
	Green   string `toml:"green"`
	Yellow  string `toml:"yellow"`
	Blue    string `toml:"blue"`
	Magenta string `toml:"magenta"`
	Cyan    string `toml:"cyan"`
	White   string `toml:"white"`
}

// Values returns palette colors in ANSI index order.
func (c ANSIColors) Values() []string {
	return []string{c.Black, c.Red, c.Green, c.Yellow, c.Blue, c.Magenta, c.Cyan, c.White}
}
