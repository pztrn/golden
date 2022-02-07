package widgets

import (
	"io"
)

type MainMenuWidget struct {
	mw *MenuWidget
}

func NewMainMenuWidget() *MainMenuWidget {

	mmw := new(MainMenuWidget)

	newMenuWidget := NewMenuWidget()
	newMainGroup := new(MenuGroup)
	newSetupGroup := new(MenuGroup)

	mmw.InitMain(newMainGroup)
	mmw.InitSetup(newSetupGroup)

	newMenuWidget.Add(newMainGroup)
	newMenuWidget.Add(newSetupGroup)

	mmw.mw = newMenuWidget

	return mmw
}

func (self *MainMenuWidget) InitMain(menuGroup *MenuGroup) {

	/* Home */
	menuAction1 := NewMenuAction()
	menuAction1.ID = "mainMenuHome"
	menuAction1.Link = "/"
	menuAction1.Label = "Home"
	menuGroup.Add(menuAction1)

	/* Netmail */
	menuAction2 := NewMenuAction()
	menuAction2.ID = "mainMenuDirect"
	menuAction2.Link = "/netmail"
	menuAction2.Label = "Netmail"
	menuAction2.Metric = -1
	menuGroup.Add(menuAction2)

	/* Echomail */
	menuAction3 := NewMenuAction()
	menuAction3.ID = "mainMenuEcho"
	menuAction3.Link = "/echo"
	menuAction3.Label = "Echomail"
	menuAction3.Metric = -1
	menuGroup.Add(menuAction3)

	/* Files */
	menuAction4 := NewMenuAction()
	menuAction4.ID = "mainMenuFile"
	menuAction4.Link = "/file"
	menuAction4.Label = "Files"
	menuAction4.Metric = -1
	menuGroup.Add(menuAction4)

	/* Service */
	menuAction5 := NewMenuAction()
	menuAction5.ID = "mainMenuService"
	menuAction5.Link = "/service"
	menuAction5.Label = "Service"
	menuGroup.Add(menuAction5)

	/* Address book */
	menuAction6 := NewMenuAction()
	menuAction6.ID = "mainMenuTwit"
	menuAction6.Link = "/twit"
	menuAction6.Label = "Twit"
	menuGroup.Add(menuAction6)

	/* Draft */
	menuAction7 := NewMenuAction()
	menuAction7.ID = "mainMenuDraft"
	menuAction7.Link = "/draft"
	menuAction7.Label = "Draft"
	menuGroup.Add(menuAction7)

}

func (self *MainMenuWidget) InitSetup(menuGroup *MenuGroup) {

	menuAction1 := NewMenuAction()
	menuAction1.ID = "mainMenuSetup"
	menuAction1.Link = "/setup"
	menuAction1.Label = "Setup"
	menuGroup.Add(menuAction1)

}

func (self *MainMenuWidget) Render(w io.Writer) error {
	self.mw.Render(w)
	return nil
}

func (self *MainMenuWidget) SetParam(ID string, value int) {
	for _, g := range self.mw.groups {
		for _, a := range g.actions {
			if a.ID == ID {
				a.Metric = value
			}
		}
	}
}
