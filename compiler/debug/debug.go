package debug

type Component struct {
	line, skipLine int
	filePath       string
}

func (c *Component) GetFilePath() string {
	return c.filePath
}

func (c *Component) SetFilePath(filePath string) {
	c.filePath = filePath
}

func (c *Component) GetLine() int {
	return c.line
}

func (c *Component) SetLine(line int) {
	c.line = line
}

func (c *Component) GetSkipLine() int {
	return c.skipLine
}

func (c *Component) SetSkipLine(line int) {
	c.skipLine = line
}

type Info interface {
	GetLine() int
	SetLine(line int)

	SetSkipLine(line int)
	GetSkipLine() int

	GetFilePath() string
	SetFilePath(filePath string)
}
