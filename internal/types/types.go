package types

// StackFrame represents a single frame in a backtrace
type StackFrame struct {
	FilePath        AbsPath
	Line            int
	Function        string
	ChangeIntensity int
	ChangeReason    string
}

// NewStackFrame creates a new StackFrame
func NewStackFrame(file AbsPath, line int, function string) StackFrame {
	return StackFrame{
		FilePath:        file,
		Line:            line,
		Function:        function,
		ChangeIntensity: 0,
		ChangeReason:    "",
	}
}
