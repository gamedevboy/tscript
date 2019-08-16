package debug

type Component struct {
    line     int
    filePath string
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

type Info interface {
    GetLine() int
    SetLine(line int)

    GetFilePath() string
    SetFilePath(filePath string)
}
