package ui

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"strconv"
	"strings"
	"time"
	"wechat_ui/ui/assets"
	"wechat_ui/ui/page/chat/appwidget/apptheme"
	"wechat_ui/ui/page/chat/gen"
	"wechat_ui/ui/page/chat/model"
	"wechat_ui/ui/pkg/async"
	"wechat_ui/ui/pkg/list"
	"wechat_ui/ui/pkg/ninepatch"
	"wechat_ui/ui/v"
	"wechat_ui/ui/values"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	lorem "github.com/drhodes/golorem"

	chatlayout "wechat_ui/ui/pkg/layout"

	chatwidget "wechat_ui/ui/pkg/widget"
	matchat "wechat_ui/ui/pkg/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Config struct {
	// theme to use {light,dark}.
	Theme string
	// latency specifies maximum latency (in millis) to simulate
	Latency int
	// loadSize specifies maximum number of items to load at a time.
	LoadSize int
	// bufferSize specifies how many elements to hold in memory before
	// compacting the list.
	BufferSize int
}

// th is the active theme object.
var (
	th = apptheme.NewTheme()
)

var (
	// SidebarMaxWidth specifies how large the side bar should be on
	// desktop layouts.
	SidebarMaxWidth = unit.Dp(250)
	// Breakpoint at which to switch from desktop to mobile layout.
	Breakpoint = unit.Dp(600)
)

// UI manages the state for the entire application's UI.
type UI struct {
	// Loader loads resources asynchronously.
	// Deallocates stale resources.
	// Stale is defined as "not being scheduled frequently".
	async.Loader
	// Rooms is the root of the data, containing messages chunked by
	// room.
	// It also contains interact state, rather than maintaining two
	// separate lists for the model and state.
	Rooms Rooms
	// Local user for this client.
	Local *model.User
	// Users contains user data.
	Users *model.Users
	// RoomList for the sidebar.
	RoomList widget.List
	// Modal can show widgets atop the rest of the ui.
	Modal component.ModalState
	// Bg is the background color of the content area.
	Bg color.NRGBA
	// InsideRoom if we are currently in the room view.
	// Used to decide when to render the sidebar on small viewports.
	InsideRoom bool
	// AddBtn holds click state for a button that adds a new message to
	// the current room.
	AddBtn widget.Clickable
	// DeleteBtn holds click state for a button that removes a message
	// from the current room.
	DeleteBtn widget.Clickable
	// MessageMenu is the context menu available on messages.
	MessageMenu component.MenuState
	// ContextMenuTarget tracks the message state on which the context
	// menu is currently acting.
	ContextMenuTarget *model.Message

	SearchEditor  *widget.Editor
	AddContactBtn v.IconButton
}

// loadNinePatch from the embedded resources package.
func loadNinePatch(path string) ninepatch.NinePatch {
	imgf, err := assets.Resources.Open(path)
	if err != nil {
		panic(fmt.Errorf("opening image: %w", err))
	}
	defer imgf.Close()
	img, err := png.Decode(imgf)
	if err != nil {
		panic(fmt.Errorf("decoding png: %w", err))
	}
	return ninepatch.DecodeNinePatch(img)
}

var (
	cookie = loadNinePatch("9-Patch/iap_platocookie_asset_2.png")
	hotdog = loadNinePatch("9-Patch/iap_hotdog_asset.png")
)

// NewUI constructs a UI and populates it with dummy data.
func NewUI(invalidator func(), conf Config) *UI {
	var ui UI

	switch conf.Theme {
	case "light":
		th.UsePalette(apptheme.Light)
	case "dark":
		th.UsePalette(apptheme.Dark)
	}

	ui.SearchEditor = &widget.Editor{}

	ui.Modal.VisibilityAnimation.Duration = time.Millisecond * 250

	ui.MessageMenu = component.MenuState{
		Options: []func(gtx C) D{
			component.MenuItem(th.Theme, &ui.DeleteBtn, "Delete").Layout,
		},
	}

	g := &gen.Generator{
		FetchImage: func(sz image.Point) image.Image {
			img, _ := randomImage(sz)
			return img
		},
	}

	ui.AddContactBtn = v.NewIconButton(ContentAdd, values.Gray1, th.Bg)
	ui.AddContactBtn.Size = unit.Dp(30)

	// Generate most of the model data.
	var (
		rooms = g.GenRooms(3, 10)
		users = g.GenUsers(10, 30)
		local = users.Random()
	)

	ui.Users = users
	ui.Local = local

	for _, r := range rooms.List() {
		rt := NewExampleData(users, local, g, 100)
		rt.SimulateLatency = conf.Latency
		rt.MaxLoads = conf.LoadSize
		lm := list.NewManager(conf.BufferSize,
			list.Hooks{
				// Define an allocator function that can instaniate the appropriate
				// state type for each kind of row data in our list.
				Allocator: func(data list.Element) interface{} {
					switch data.(type) {
					case model.Message:
						return &chatwidget.Row{}
					default:
						return nil
					}
				},
				// Define a presenter that can transform each kind of row data
				// and state into a widget.
				Presenter: ui.presentChatRow,
				// NOTE(jfm): awkard coupling between message data and `list.Manager`.
				Loader:      rt.Load,
				Synthesizer: synth,
				Comparator:  rowLessThan,
				Invalidator: invalidator,
			},
		)
		lm.Stickiness = list.After
		ui.Rooms.List = append(ui.Rooms.List, Room{
			Room:      r,
			Messages:  rt,
			ListState: lm,
		})
	}

	// spin up a bunch of async actors to send messages to rooms.
	for _, u := range users.List() {
		u := u
		if u.Name == local.Name {
			continue
		}
		go func() {
			for {
				var (
					respond = time.Second * time.Duration(1)
					compose = time.Second * time.Duration(1)
					room    = ui.Rooms.Random()
				)
				func() {
					time.Sleep(respond)
					room.SetComposing(u.Name, true)
					time.Sleep(compose)
					room.SetComposing(u.Name, false)
					room.Send(u.Name, lorem.Paragraph(1, 4))
				}()
			}
		}()
	}

	ui.Rooms.Select(0)
	for ii := range ui.Rooms.List {
		ui.Rooms.List[ii].List.ScrollToEnd = true
		ui.Rooms.List[ii].List.Axis = layout.Vertical
	}

	ui.Bg = th.Palette.Bg

	return &ui
}

// Layout the application UI.
func (ui *UI) Layout(gtx C) D {

	return ui.Loader.Frame(gtx, ui.layout)
}

func (ui *UI) layout(gtx C) D {
	for ii := range ui.Rooms.List {
		r := &ui.Rooms.List[ii]
		if r.Interact.Clicked() {
			ui.Rooms.Select(ii)
			ui.InsideRoom = true
			break
		}
	}

	paint.FillShape(gtx.Ops, ui.Bg, clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X = gtx.Dp(SidebarMaxWidth)
			gtx.Constraints.Min = gtx.Constraints.Constrain(gtx.Constraints.Min)
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return ui.layoutSearch(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					// 搜索结果
					if searchR := ui.layoutSearchResult(); searchR != nil {
						return searchR(gtx)
					}
					// 会话列表
					return ui.layoutRoomList(gtx)
				}),
			)
		}),
		layout.Rigid(v.SeparatorVertical(gtx.Constraints.Max.Y, 1, component.WithAlpha(th.Fg, 50)).Layout),
		layout.Flexed(1, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					gtx.Constraints.Min = gtx.Constraints.Max
					return ui.layoutChat(gtx)
				}),
				layout.Expanded(func(gtx C) D {
					return ui.layoutModal(gtx)
				}),
			)
		}),
	)
}

// layoutChat lays out the chat interface with associated controls.
func (ui *UI) layoutChat(gtx C) D {
	room := ui.Rooms.Active()
	var (
		//scrollWidth unit.Dp
		list  = &room.List
		state = room.ListState
	)
	listStyle := material.List(th.Theme, list)
	//scrollWidth = listStyle.ScrollbarStyle.Width()
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return listStyle.Layout(gtx,
				state.UpdatedLen(&list.List),
				state.Layout,
			)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.layoutEditor2(gtx)
			//return chatlayout.Background(th.Palette.BgSecondary).Layout(gtx, func(gtx C) D {
			//	if ui.AddBtn.Clicked() {
			//		active := ui.Rooms.Active()
			//		active.SendLocal(active.Editor.Text())
			//		active.Editor.SetText("")
			//	}
			//	if ui.DeleteBtn.Clicked() {
			//		serial := ui.ContextMenuTarget.Serial()
			//		ui.Rooms.Active().DeleteRow(serial)
			//	}
			//	return layout.Inset{
			//		Bottom: unit.Dp(8),
			//		Top:    unit.Dp(8),
			//	}.Layout(gtx, func(gtx C) D {
			//		gutter := chatlayout.Gutter()
			//		gutter.RightWidth = gutter.RightWidth + scrollWidth
			//		return gutter.Layout(gtx,
			//			nil,
			//			func(gtx C) D {
			//				return ui.layoutEditor(gtx)
			//			},
			//			material.IconButton(th.Theme, &ui.AddBtn, Send, "Send").Layout,
			//		)
			//	})
			//})
		}),
	)
}

// layoutRoomList lays out a list of rooms that can be clicked to view
// the messages in that room.
func (ui *UI) layoutRoomList(gtx C) D {
	return layout.Stack{}.Layout(
		gtx,
		layout.Expanded(func(gtx C) D {
			return component.Rect{
				Size: image.Point{
					X: gtx.Constraints.Min.X,
					Y: gtx.Constraints.Max.Y,
				},
				Color: th.Palette.Surface,
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			ui.RoomList.Axis = layout.Vertical
			gtx.Constraints.Min = gtx.Constraints.Max
			listL := material.List(th.Theme, &ui.RoomList)
			listL.AnchorStrategy = material.Overlay
			return listL.Layout(gtx, len(ui.Rooms.List), func(gtx C, ii int) D {
				r := ui.Rooms.Index(ii)
				latest := r.Latest()
				return apptheme.Room(th.Theme, &r.Interact, &apptheme.RoomConfig{
					Name:    r.Room.Name,
					Image:   r.Room.Image,
					Content: latest.Content,
					SentAt:  latest.SentAt,
				}).Layout(gtx)
			})
		}),
	)
}
func (ui *UI) layoutSearchResult() layout.Widget {
	searchContent := strings.TrimSpace(ui.SearchEditor.Text())
	for !ui.SearchEditor.Focused() || len(searchContent) == 0 {
		return nil
	}
	searchContent = strings.ToLower(searchContent)
	var r Room
	has := false
	for _, roomInfo := range ui.Rooms.List {
		if strings.ToLower(roomInfo.Name) == searchContent {
			has = true
			r = roomInfo
			break
		}
	}
	if !has {
		return layout.Spacer{}.Layout
	}

	latest := r.Latest()
	return apptheme.Room(th.Theme, &r.Interact, &apptheme.RoomConfig{
		Name:    r.Room.Name,
		Image:   r.Room.Image,
		Content: latest.Content,
		SentAt:  latest.SentAt,
	}).Layout
}

// layoutSearch lays out the search editor.
func (ui *UI) layoutSearch(gtx C) D {
	inset := layout.Inset{
		Top:    values.MarginPadding20,
		Bottom: values.MarginPadding10,
	}
	return inset.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly,
		}.Layout(gtx, layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X / 4 * 3
			return chatlayout.Rounded(unit.Dp(2)).Layout(gtx, func(gtx C) D {
				return chatlayout.Background(values.Gray1).Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
						for _, e := range ui.SearchEditor.Events() {
							switch e.(type) {
							case widget.SubmitEvent:
								ui.SearchEditor.SetText("")
							}
						}
						ui.SearchEditor.Submit = true
						ui.SearchEditor.SingleLine = true
						ui.SearchEditor.MaxLen = 10
						ed := material.Editor(th.Theme, ui.SearchEditor, "Search")
						return ed.Layout(gtx)
					})
				})
			})
		}),
			layout.Rigid(func(gtx C) D {
				for ui.AddContactBtn.Button.Clicked() {
					fmt.Println("点击添加好友")
				}
				return ui.AddContactBtn.Layout(gtx)
			}),
		)
	})
}

// layoutEditor lays out the message editor.
func (ui *UI) layoutEditor(gtx C) D {
	return chatlayout.Rounded(unit.Dp(8)).Layout(gtx, func(gtx C) D {
		return chatlayout.Background(th.Palette.Surface).Layout(gtx, func(gtx C) D {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
				active := ui.Rooms.Active()
				editor := &active.Editor
				for _, e := range editor.Events() {
					switch e.(type) {
					case widget.SubmitEvent:
						active.SendLocal(editor.Text())
						editor.SetText("")
					}
				}
				editor.Submit = true
				editor.SingleLine = true
				return material.Editor(th.Theme, editor, "Send a message").Layout(gtx)
			})
		})
	})
}

func (ui *UI) layoutModal(gtx C) D {
	if ui.Modal.Clicked() {
		ui.Modal.ToggleVisibility(gtx.Now)
	}
	// NOTE(jfm): scrim should be dark regardless of theme.
	// Perhaps "scrim color" could be specified on the theme.
	t := *th.Theme
	t.Fg = apptheme.Dark.Surface
	return component.Modal(&t, &ui.Modal).Layout(gtx)
}

// synth inserts date separators and unread separators
// between chat rows as a list.Synthesizer.
func synth(previous, row, next list.Element) []list.Element {
	var out []list.Element
	asMessage, ok := row.(model.Message)
	if !ok {
		out = append(out, row)
		return out
	}
	if previous == nil {
		if !asMessage.Read {
			out = append(out, model.UnreadBoundary{})
		}
		out = append(out, row)
		return out
	}
	lastMessage, ok := previous.(model.Message)
	if !ok {
		out = append(out, row)
		return out
	}
	if !asMessage.Read && lastMessage.Read {
		out = append(out, model.UnreadBoundary{})
	}
	y, m, d := asMessage.SentAt.Local().Date()
	yy, mm, dd := lastMessage.SentAt.Local().Date()
	if y == yy && m == mm && d == dd {
		out = append(out, row)
		return out
	}
	out = append(out, model.DateBoundary{Date: asMessage.SentAt}, row)
	return out
}

// rowLessThan acts as a list.Comparator, returning whether a sorts before b.
func rowLessThan(a, b list.Element) bool {
	aID := string(a.Serial())
	bID := string(b.Serial())
	aAsInt, _ := strconv.Atoi(aID)
	bAsInt, _ := strconv.Atoi(bID)
	return aAsInt < bAsInt
}

// presentChatRow returns a widget closure that can layout the given chat item.
// `data` contains managed data for this chat item, `state` contains UI defined
// interactive state.
func (ui *UI) presentChatRow(data list.Element, state interface{}) layout.Widget {
	switch data := data.(type) {
	case model.Message:
		state, ok := state.(*chatwidget.Row)
		if !ok {
			return func(C) D { return D{} }
		}
		return func(gtx C) D {
			if state.Clicked() {
				ui.Modal.Show(gtx.Now, func(gtx C) D {
					return layout.UniformInset(unit.Dp(25)).Layout(gtx, func(gtx C) D {
						return widget.Image{
							Src:      state.Image.Op(),
							Fit:      widget.ScaleDown,
							Position: layout.Center,
						}.Layout(gtx)
					})
				})
			}
			if state.ContextArea.Active() {
				// If the right-click context area for this message is activated,
				// inform the UI that this message is the target of any action
				// taken within that menu.
				ui.ContextMenuTarget = &data
			}
			return ui.row(data, state)(gtx)
		}
	case model.DateBoundary:
		return matchat.DateSeparator(th.Theme, data.Date).Layout
	case model.UnreadBoundary:
		return matchat.UnreadSeparator(th.Theme).Layout
	default:
		return func(gtx C) D { return D{} }
	}
}

// row returns either a plato.RowStyle or a chatmaterial.RowStyle based on the
// provided boolean.
func (ui *UI) row(data model.Message, state *chatwidget.Row) layout.Widget {
	user, ok := ui.Users.Lookup(data.Sender)
	if !ok {
		return func(C) D { return D{} }
	}
	np := func() *ninepatch.NinePatch {
		switch user.Theme {
		case model.ThemeHotdog:
			return &hotdog
		case model.ThemePlatoCookie:
			return &cookie
		}
		return nil
	}()
	var (
		avatar image.Image
		body   image.Image
	)
	if data.Avatar != "" {
		avatar = avatarPlaceholder
		if img := loadImage(string(data.Serial())+"-avatar", data.Avatar, &ui.Loader); img != nil {
			state.Avatar.Reload()
			avatar = img
		}
	}
	if data.Image != "" {
		body = imageMessagePlaceholder
		if img := loadImage(string(data.Serial())+"-body", data.Image, &ui.Loader); img != nil {
			state.Image.Reload()
			body = img
		}
	}
	msg := matchat.NewRow(th.Theme, state, &ui.MessageMenu, matchat.RowConfig{
		Sender:  data.Sender,
		Content: data.Content,
		SentAt:  data.SentAt,
		Avatar:  avatar,
		Image:   body,
		Local:   user.Name == ui.Local.Name,
	})
	if np != nil {
		msg.MessageStyle = msg.WithNinePatch(th.Theme, *np)
	}
	msg.MessageStyle.BubbleStyle.Color = user.Color
	for i := range msg.Content.Styles {
		msg.Content.Styles[i].Color = th.Contrast(matchat.Luminance(user.Color))
	}
	return msg.Layout
}

var (
	// placeholderColor to use for placeholder images.
	placeholderColor = color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	// avatarPlaceholder used when avatar image has not been loaded yet.
	avatarPlaceholder *image.NRGBA = placeholder(image.Pt(64, 64), placeholderColor)
	// imageMessagePlaceholder used when message image has not been loaded yet.
	imageMessagePlaceholder *image.NRGBA = placeholder(image.Pt(320, 320), placeholderColor)
)

// placeholder helper generates a rectangle image of the given size for the
// given color.
func placeholder(sz image.Point, c color.NRGBA) (ph *image.NRGBA) {
	ph = image.NewNRGBA(image.Rectangle{Max: sz})
	for xx := ph.Bounds().Min.X; xx < ph.Bounds().Max.X; xx++ {
		for yy := ph.Bounds().Min.Y; yy < ph.Bounds().Max.Y; yy++ {
			ph.SetNRGBA(xx, yy, c)
		}
	}
	return ph
}

// loadImage helper schedules an image to be downloaded and returns it if ready.
func loadImage(id, u string, l *async.Loader) image.Image {
	r := l.Schedule(id, func(_ context.Context) interface{} {
		img, err := fetch(id, u)
		if err != nil {
			log.Printf("loading image: %v", err)
		}
		return img
	})
	switch r.State {
	case async.Queued, async.Loading:
	case async.Loaded:
		if img, ok := r.Value.(image.Image); ok {
			return img
		}
	}
	return nil
}
