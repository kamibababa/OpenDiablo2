package d2player

import (
	"log"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2resource"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2ui"
)

const (
	miniPanelX = 325
	miniPanelY = 526

	panelOffsetLeft  = 130
	panelOffsetRight = 130
)

const (
	containerOffsetX = -75
	containerOffsetY = -49

	buttonOffsetX = -72
	buttonOffsetY = -52
)

type miniPanelContent struct {
	buttonType d2ui.ButtonType
	onActivate func()
	tooltip    string
}

type miniPanelActions struct {
	characterToggle func()
	inventoryToggle func()
	skilltreeToggle func()
	partyToggle     func()
	automapToggle   func()
	messageToggle   func()
	questToggle     func()
	menuToggle      func()
}

type miniPanel struct {
	ui             *d2ui.UIManager
	asset          *d2asset.AssetManager
	container      *d2ui.Sprite
	sprite         *d2ui.Sprite
	isOpen         bool
	isSinglePlayer bool
	movedLeft      bool
	movedRight     bool
	panelGroup     *d2ui.WidgetGroup
	tooltipGroup   *d2ui.WidgetGroup
}

func newMiniPanel(asset *d2asset.AssetManager, uiManager *d2ui.UIManager, isSinglePlayer bool) *miniPanel {
	return &miniPanel{
		ui:             uiManager,
		asset:          asset,
		isOpen:         false,
		isSinglePlayer: isSinglePlayer,
	}
}

func (m *miniPanel) load(actions *miniPanelActions) {
	var err error

	m.sprite, err = m.ui.NewSprite(d2resource.MinipanelButton, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
		return
	}

	m.createWidgets(actions)
}

func (m *miniPanel) createWidgets(actions *miniPanelActions) {
	var err error

	m.panelGroup = m.ui.NewWidgetGroup(d2ui.RenderPriorityMinipanel)
	m.panelGroup.SetPosition(miniPanelX, miniPanelY)

	m.tooltipGroup = m.ui.NewWidgetGroup(d2ui.RenderPriorityForeground)

	miniPanelContainerPath := d2resource.Minipanel
	if m.isSinglePlayer {
		miniPanelContainerPath = d2resource.MinipanelSmall
	}

	m.container, err = m.ui.NewSprite(miniPanelContainerPath, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
		return
	}

	if err = m.container.SetCurrentFrame(0); err != nil {
		log.Print(err)
		return
	}

	// nolint:golint,gomnd // divide by 2 does not need a magic number
	x, y := screenWidth/2+containerOffsetX, screenHeight+containerOffsetY
	m.container.SetPosition(x, y)
	m.panelGroup.AddWidget(m.container)

	buttonWidth, buttonHeight, err := m.sprite.GetFrameSize(0)
	if err != nil {
		log.Print(err)
		return
	}

	buttonWidth++

	// nolint:golint,gomnd // divide by 2 does not need a magic number
	x, y = screenWidth/2+buttonOffsetX, screenHeight+buttonOffsetY-buttonHeight
	buttonsFirst := []miniPanelContent{
		{d2ui.ButtonTypeMinipanelCharacter,
			actions.characterToggle,
			m.asset.TranslateString("minipanelchar"),
		},
		{d2ui.ButtonTypeMinipanelInventory,
			actions.inventoryToggle,
			m.asset.TranslateString("minipanelinv"),
		},
		{d2ui.ButtonTypeMinipanelSkill,
			actions.skilltreeToggle,
			m.asset.TranslateString("minipaneltree"),
		},
	}

	for i := range buttonsFirst {
		btn := m.createButton(buttonsFirst[i], x+(i*buttonWidth), y, buttonHeight)
		m.panelGroup.AddWidget(btn)
	}

	idxOffset := len(buttonsFirst)

	if !m.isSinglePlayer {
		partyContent := miniPanelContent{d2ui.ButtonTypeMinipanelParty,
			actions.partyToggle,
			m.asset.TranslateString("minipanelparty"),
		}
		btn := m.createButton(partyContent, x+(3*buttonWidth), y, buttonHeight)
		m.panelGroup.AddWidget(btn)
		idxOffset++
	}

	buttonsLast := []miniPanelContent{
		{d2ui.ButtonTypeMinipanelAutomap,
			actions.automapToggle,
			m.asset.TranslateString("minipanelautomap"),
		},
		{d2ui.ButtonTypeMinipanelMessage,
			actions.messageToggle,
			m.asset.TranslateString("minipanelmessage"),
		},
		{d2ui.ButtonTypeMinipanelQuest,
			actions.questToggle,
			m.asset.TranslateString("minipanelquest"),
		},
		{d2ui.ButtonTypeMinipanelMen,
			actions.menuToggle,
			m.asset.TranslateString("minipanelmenubtn"),
		},
	}

	for i := range buttonsLast {
		idx := i + idxOffset
		btn := m.createButton(buttonsLast[i], x+(idx*buttonWidth), y, buttonHeight)
		m.panelGroup.AddWidget(btn)
	}

	m.panelGroup.SetVisible(false)
}

func (m *miniPanel) createButton(content miniPanelContent, x, y, buttonHeight int) *d2ui.Button {
	// Tooltip
	tt := m.ui.NewTooltip(d2resource.Font16, d2resource.PaletteSky, d2ui.TooltipXCenter, d2ui.TooltipYTop)
	tt.SetPosition(x, y-buttonHeight)
	tt.SetText(content.tooltip)
	tt.SetVisible(false)
	m.tooltipGroup.AddWidget(tt)

	// Button
	btn := m.ui.NewButton(content.buttonType, "")
	btn.SetPosition(x, y)
	btn.OnActivated(content.onActivate)
	btn.SetTooltip(tt)

	return btn
}

func (m *miniPanel) IsOpen() bool {
	return m.isOpen
}

func (m *miniPanel) Toggle() {
	if m.isOpen {
		m.Close()
	} else {
		m.Open()
	}
}

func (m *miniPanel) Open() {
	m.panelGroup.SetVisible(true)
	m.isOpen = true
}

func (m *miniPanel) Close() {
	m.panelGroup.SetVisible(false)
	m.isOpen = false
}

func (m *miniPanel) IsInRect(px, py int) bool {
	return m.panelGroup.Contains(px, py)
}

func (m *miniPanel) moveRight() {
	m.panelGroup.OffsetPosition(panelOffsetRight, 0)
	m.tooltipGroup.OffsetPosition(panelOffsetRight, 0)
}

func (m *miniPanel) undoMoveRight() {
	m.panelGroup.OffsetPosition(-panelOffsetRight, 0)
	m.tooltipGroup.OffsetPosition(-panelOffsetRight, 0)
}

func (m *miniPanel) moveLeft() {
	m.panelGroup.OffsetPosition(-panelOffsetLeft, 0)
	m.tooltipGroup.OffsetPosition(-panelOffsetLeft, 0)
}

func (m *miniPanel) undoMoveLeft() {
	m.panelGroup.OffsetPosition(panelOffsetLeft, 0)
	m.tooltipGroup.OffsetPosition(panelOffsetLeft, 0)
}

func (m *miniPanel) SetMovedLeft(moveLeft bool) {
	if m.movedLeft == moveLeft {
		return
	}

	if m.movedRight {
		if moveLeft {
			m.undoMoveRight()
			m.panelGroup.SetVisible(false)
		} else {
			m.moveRight()
			m.panelGroup.SetVisible(true)
		}
	} else {
		if moveLeft {
			m.moveLeft()
		} else {
			m.undoMoveLeft()
		}
	}

	m.movedLeft = moveLeft
}

func (m *miniPanel) SetMovedRight(moveRight bool) {
	if m.movedRight == moveRight {
		return
	}

	if m.movedLeft {
		if moveRight {
			m.undoMoveLeft()
			m.panelGroup.SetVisible(false)
		} else {
			m.moveLeft()
			m.panelGroup.SetVisible(true)
		}
	} else {
		if moveRight {
			m.moveRight()
		} else {
			m.undoMoveRight()
		}
	}

	m.movedRight = moveRight
}
