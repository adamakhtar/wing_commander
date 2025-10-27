package filepicker

import (
	"fmt"
	"os"
	"sort"
	"strings"

	filewalker "github.com/adamakhtar/wing_commander/internal/file_walker"
	"github.com/adamakhtar/wing_commander/internal/ui/context"
	"github.com/adamakhtar/wing_commander/internal/ui/keys"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

//
// TYPES
//================================================

type Dimension struct {
	width int
	height int
}

type ComponentDimensions struct {
	margin int
	searchInput Dimension
	resultsPanel Dimension
	selectedPathsPanel Dimension
}

type Focused string

const (
	SearchInput Focused = "searchInput"
	SearchResults Focused = "searchResults"
	SelectedPaths Focused = "selectedPaths"
)

type Model struct {
	ctx context.Context
	allPaths []string
	resultsTable table.Model
	selectedPaths UniqueFilesSet
	searchInput textinput.Model
	currentPanel Focused
	testFilePathsLoaded bool
}

//
// BUILDERS
//================================================

func NewModel(ctx context.Context) Model {
	si := textinput.New()
	si.Placeholder = "Pikachu"
	si.CharLimit = 156
	si.Width = 0
	si.Focus();

	return Model{
		ctx: ctx,
		allPaths: []string{},
		resultsTable: table.New(table.WithColumns([]table.Column{{Title: "File", Width: 0}})),
		selectedPaths: make(UniqueFilesSet),
		searchInput: si,
		currentPanel: SearchResults,
		testFilePathsLoaded: true,
	}
}

//
// BUBBLETEA
//================================================

func (m Model) Init() tea.Cmd {
	fmt.Println("initaializing")
	return getTestFilePathsCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
		case filePathsMsg:
			m.onFilePathsLoaded(msg)
			return m, nil
    case tea.KeyMsg:
			switch {
				case key.Matches(msg, keys.FilepickerKeys.Cancel):
					return m, cancelCmd
				case key.Matches(msg, keys.FilepickerKeys.LineUp):
					m.resultsTable, cmd = m.resultsTable.Update(msg)
					return m, cmd
				case key.Matches(msg, keys.FilepickerKeys.LineDown):
					m.resultsTable, cmd = m.resultsTable.Update(msg)
					return m, cmd
				case key.Matches(msg, keys.FilepickerKeys.Select):
					m.addSelectedPath()
					return m, nil
				case key.Matches(msg, keys.FilepickerKeys.RunTests):
					return m, nil
				default:
					m.searchInput, cmd = m.searchInput.Update(msg)
					m.onSearchInputChanged(m.searchInput.Value())
					return m, cmd
			}
		case errMsg:
			fmt.Printf(msg.Error())
			return m, nil
	}

	return m, cmd
}

func (m Model) View() string {
	// return "File Picker View"
	s := strings.Builder{}

	dimensions := m.getComponentDimensions()
	_, searchInput, resultsPanel, selectedPathsPanel := dimensions.margin, dimensions.searchInput, dimensions.resultsPanel, dimensions.selectedPathsPanel

	search := NewPanel(searchInput.width, searchInput.height, m.searchInputFocussed(), m.searchInput.View())
	searchResults := NewPanel(resultsPanel.width, resultsPanel.height, m.searchResultsFocussed(), m.resultsTable.View())
	selectedPaths := NewPanel(selectedPathsPanel.width, selectedPathsPanel.height, m.selectedPathsFocussed(), m.selectedPaths.String())

	s.WriteString(search.render())
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, searchResults.render(), selectedPaths.render()))

	return s.String()
}

//
// MESSAGES & HANDLERS
//================================================

type filePathsMsg struct {
	filePaths []string
}

func (msg filePathsMsg) FilePaths() []string {
	return msg.filePaths
}

type errMsg struct {
	err error
}

type CancelMsg struct {}

func (m *Model) onFilePathsLoaded(msg filePathsMsg) {
	m.allPaths = msg.FilePaths()
	m.setResultTableRows(m.allPaths)
	m.resultsTable.Focus()
	m.testFilePathsLoaded = true
}

func (m *Model) onSearchInputChanged(value string) {
	if value == "" {
		m.setResultTableRows(m.allPaths)
		return
	}

	matchingPaths := fuzzy.FindFold(value, m.allPaths)

	if len(matchingPaths) == 0 {
		m.setResultTableRows(m.allPaths)
	} else {
		m.setResultTableRows(matchingPaths)
	}
}

//
// Commands
//===============================================================

func getTestFilePathsCmd() tea.Msg {
	filePaths, err := getTestFilePaths()

	if err != nil {
		return errMsg{err: err}
	}

	return filePathsMsg{filePaths: filePaths}
}

func (e errMsg) Error() string {
	return e.err.Error()
}

func cancelCmd() tea.Msg {
	return CancelMsg{}
}

//
// EXTERNAL FUNCTIONS
//================================================

func (m *Model) Prepare() tea.Cmd {
	m.resetPicker()
	return getTestFilePathsCmd
}

func (m *Model) UpdateContext(ctx context.Context) {
	m.ctx = ctx

	dimensions := m.getComponentDimensions()
	m.resultsTable.SetHeight(dimensions.resultsPanel.height)
	m.resultsTable.SetWidth(dimensions.resultsPanel.width)
	m.searchInput.Width = dimensions.searchInput.width
	log.Debug("filepicker.UpdateContext: dimensions", dimensions)
}


type UniqueFilesSet map[string]bool

func (s *UniqueFilesSet) Add(key string) {
		(*s)[key] = true
}

func (s *UniqueFilesSet) Keys() []string {
		keys := make([]string, 0, len(*s))
		for key := range *s {
				keys = append(keys, key)
		}
		return keys
}

func (s *UniqueFilesSet) String() string {
		keys := s.Keys()
		sort.Strings(keys)
		return strings.Join(keys, "\n")
}


//
//  METHODS
//===============================================================

func getTestFilePaths() ([]string, error) {
	fmt.Println("getTestFilePaths Start")
	cwd, err := os.Getwd()

	if err != nil {
		return nil, errMsg{err: err}
	}

	excludePatterns := []string{".git/**"}

	filePaths := filewalker.FileEntriesRecursive(cwd, []string{}, excludePatterns)

	fmt.Println("getTestFilePaths End")
	return filePaths, nil
}



func (m *Model) setResultTableRows(filePaths []string) {
	tableRows := []table.Row{}
	for _, path := range filePaths {
		tableRow := table.Row{path}
		tableRows = append(tableRows, tableRow)
	}

	dimensions := m.getComponentDimensions()

	m.resultsTable.SetColumns([]table.Column{{Title: "File", Width: dimensions.resultsPanel.width}})
	m.resultsTable.SetRows(tableRows)
	m.resultsTable.SetCursor(0)
}


func (m *Model) addSelectedPath() {
	selectedRow := m.resultsTable.SelectedRow()

	if selectedRow == nil {
		return
	}

	file := selectedRow[0]

	m.selectedPaths.Add(file)
}

func (m *Model) togglePanelFocus() {
	switch m.currentPanel {
	case SearchResults:
		m.focusSelectedPaths()
	case SelectedPaths:
		m.focusSearchResults()
	default:
		m.focusSearchResults()
	}
}

func (m *Model) focusSearchResults() {
	m.currentPanel = SearchResults
	m.resultsTable.Focus()
}

func (m *Model) focusSelectedPaths() {
	m.currentPanel = SelectedPaths
	m.resultsTable.Blur()
}

func (m Model) searchResultsFocussed() bool {
	return m.currentPanel == SearchResults
}

func (m Model) selectedPathsFocussed() bool {
	return m.currentPanel == SelectedPaths
}

func (m Model) searchInputFocussed() bool {
	return m.currentPanel == SearchInput
}

func (m Model) getComponentDimensions() ComponentDimensions {
	margin := 1
	searchInputHeight := 1

	return ComponentDimensions{
		margin: margin,
		searchInput: Dimension{m.ctx.ScreenWidth - (margin * 2), searchInputHeight},
		resultsPanel: Dimension{(m.ctx.ScreenWidth - margin * 3) / 2, m.ctx.ScreenHeight - searchInputHeight - (margin * 2)},
		selectedPathsPanel: Dimension{(m.ctx.ScreenWidth - margin * 3) / 2, m.ctx.ScreenHeight - searchInputHeight - (margin * 2)},
	}
}

func (m *Model) resetPicker() {
	m.selectedPaths = make(UniqueFilesSet)
	m.searchInput.SetValue("")
	m.setResultTableRows(m.allPaths)
	m.resultsTable.SetCursor(0)
}


// Panel Componeents

type PanelStyles struct {
	regular lipgloss.Style
	focused lipgloss.Style
}

type Panel struct {
	width int
	height int
	focused bool
	styles PanelStyles
	pContent string
}


func NewPanel(width int, height int, focused bool, pContent string) Panel {
	return Panel{
		width: width,
		height: height,
		focused: focused,
		styles: DefaultPanelStyles(),
		pContent: pContent,
	}
}

func (p Panel) render() string {
	panelStyle := p.styles.regular
	if p.focused {
		panelStyle = p.styles.focused
	}

	return panelStyle.Width(p.width).Height(p.height).Render(p.pContent)
}

func DefaultPanelStyles() PanelStyles {
	return PanelStyles{
		regular: lipgloss.NewStyle().
    	BorderStyle(lipgloss.NormalBorder()).
    	BorderForeground(lipgloss.Color("63")).
			Padding(0, 0),
		focused: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("5")).
			Padding(0, 0),
	}
}
