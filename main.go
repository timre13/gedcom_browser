package main

import (
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/gdk"
    "github.com/gotk3/gotk3/cairo"
    "fmt"
    "strings"
    //"math"
    . "gedcom_browser/token"
    "gedcom_browser/widgets"
)

type TreeWidgetRect struct {
    XPos        float64
    YPos        float64
    Width       float64
    Height      float64
    FontSize    float64
    Highlight   bool
    Token       *Token
    isHovered   bool
}

func (this *TreeWidgetRect) Draw(cr *cairo.Context) {
    if this.Highlight {
        cr.SetSourceRGB(1, 1, 1)
        cr.Rectangle(this.XPos, this.YPos, this.Width, this.Height)
        cr.SetLineWidth(0.3)
        cr.Stroke()
    }

    // This happens for example when a person's husband/wife is not known
    // Draw a placeholder
    if this.Token == nil {
        cr.SetSourceRGB(0.25, 0.25, 0.25)
        cr.Rectangle(this.XPos, this.YPos, this.Width, this.Height)
        cr.Fill()
        return
    }

    gender := this.Token.GetFirstChildWithTagValueOr(TAG_SEX, "U")
    switch (gender) {
    case "M": // Male
        cr.SetSourceRGB(0.35, 0.4, 0.5)
    case "F": // Female
        cr.SetSourceRGB(0.45, 0.3, 0.4)
    case "X": // Other
        cr.SetSourceRGB(0.5, 0.5, 0.4)
    case "U": // Unknown
        fallthrough
    default:
        cr.SetSourceRGB(0.5, 0.5, 0.5)
    }
    cr.Rectangle(this.XPos, this.YPos, this.Width, this.Height)
    cr.Fill()

    cr.SetSourceRGB(1.0, 1.0, 1.0)
    cr.MoveTo(this.XPos, this.YPos+4)
    cr.SetFontSize(this.FontSize)
    name := indi_list_widget.ParseName(this.Token.GetFirstChildWithTagValueOr(TAG_NAME, ""))
    cr.ShowText(name.FirstName+" "+name.LastName)

    cr.SetFontSize(this.FontSize*0.6)
    cr.MoveTo(this.XPos, this.YPos+7)
    birthYear := "???"
    birthToken := this.Token.GetFirstChildWithTag(TAG_BIRT)
    if birthToken != nil {
        birthDateToken := birthToken.GetFirstChildWithTag(TAG_DATE)
        if birthDateToken != nil {
            birthDate := birthDateToken.ParseToDate()
            if birthDate != nil {
                birthYear = fmt.Sprint(birthDate.Year)
            }
        }
    }

    deathYear := ""
    deathToken := this.Token.GetFirstChildWithTag(TAG_DEAT)
    if deathToken != nil {
        deathDateToken := deathToken.GetFirstChildWithTag(TAG_DATE)
        if deathDateToken != nil {
            deathDate := deathDateToken.ParseToDate()
            if deathDate != nil {
                deathYear = fmt.Sprint(deathDate.Year)
            } else {
                deathYear = "???"
            }
        } else {
            deathYear = "???"
        }
    }
    cr.ShowText(birthYear+" - "+deathYear)

    if this.isHovered {
        cr.SetSourceRGBA(1, 1, 1, 0.15)
        cr.Rectangle(this.XPos, this.YPos, this.Width, this.Height)
        cr.Fill()
    }
}

func (this *TreeWidgetRect) IsPointInside(x float64, y float64) bool {
    return x >= this.XPos && x < this.XPos+this.Width &&
           y >= this.YPos && y < this.YPos+this.Height
}

func (this *TreeWidgetRect) OnMouseEnter() {
    this.isHovered = true
}

func (this *TreeWidgetRect) OnMouseLeave() {
    this.isHovered = false
}

func main() {
    gtk.Init(nil)

    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        panic(err)
    }
    win.SetTitle("GEDCOM Browser")
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })
    win.SetDefaultSize(1500, 1000)
    win.ShowAll()

    path := "./private/Marta20160827.ged"
    tokens := LoadTokensFromFile(path)
    //fmt.Println("-------- Tokens --------")
    //for i, token := range tokens {
    //    fmt.Printf("%d : %s\n", i+1, token.String())
    //}

    tree := BuildTreeFromTokens(tokens)
    //fmt.Println("-------- Tree --------")
    //PrintTree(tree.Tokens, 0)

    getValueOr := func(token *Token, def string) string {
        if token != nil {
            return token.LineVal.GetValueOr(def)
        }
        return def
    }

    countTokensWithtag := func(tag Tag) int {
        count := 0
        for _, token := range tokens {
            if token.Tag == tag {
                count++
            }
        }
        return count
    }

    fmt.Println("-------- Info --------")
    fmt.Printf("Format version: %s\n", getValueOr(tree.GetTokenByPath([]Tag{TAG_HEAD, TAG_GEDC, TAG_VERS}), "???"))
    fmt.Printf("Source: %s\n", getValueOr(tree.GetTokenByPath([]Tag{TAG_HEAD, TAG_SOUR, TAG_NAME}), "???"))
    fmt.Printf("Language: %s\n", getValueOr(tree.GetTokenByPath([]Tag{TAG_HEAD, TAG_LANG}), "???"))
    fmt.Printf("File date: %s\n", getValueOr(tree.GetTokenByPath([]Tag{TAG_HEAD, TAG_DATE}), "???"))
    fmt.Printf("Individual count: %d\n", countTokensWithtag(TAG_INDI));
    fmt.Printf("Family count: %d\n", countTokensWithtag(TAG_FAM));

    win.SetTitle(fmt.Sprintf("GEDCOM Browser [%s]",
        strings.TrimSuffix(strings.TrimSuffix(path[strings.LastIndex(path, "/")+1:], ".ged"), ".GED")))

    mainCont, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)

    individListWidget := indi_list_widget.IndiListWidgetNew(&tree)
    individListWidget.ScrollWidget.SetSizeRequest(300, 0)
    mainCont.Add(individListWidget.BoxWidget)

    var treeWidgetItems []*TreeWidgetRect

    type PersonLevel int
    const (
        PERSON_LEVEL_ANCESTOR PersonLevel = iota
        PERSON_LEVEL_NORMAL
        PERSON_LEVEL_CHILD
    )
    
    // Level ancestor
    const PERSON_LA_RECT_W = 20
    const PERSON_LA_RECT_H = 12
    const PERSON_LA_RECT_PAD = 1
    const PERSON_LA_RECT_FONTS = 2
    // Level normal
    const PERSON_LN_RECT_W = 30
    const PERSON_LN_RECT_H = 20
    const PERSON_LN_RECT_PAD = 2
    const PERSON_LN_RECT_FONTS = 3
    // Level child
    const PERSON_LC_RECT_W = 20
    const PERSON_LC_RECT_H = 12
    const PERSON_LC_RECT_PAD = 1
    const PERSON_LC_RECT_FONTS = 2
    addPersonToTreeWidget := func(token *Token, x float64, y float64, level PersonLevel, highlight ...bool) {
        var person TreeWidgetRect
        switch level {
        case PERSON_LEVEL_ANCESTOR:
            person.Width    = PERSON_LA_RECT_W
            person.Height   = PERSON_LA_RECT_H
            person.FontSize = PERSON_LA_RECT_FONTS
        case PERSON_LEVEL_NORMAL:
            person.Width    = PERSON_LN_RECT_W
            person.Height   = PERSON_LN_RECT_H
            person.FontSize = PERSON_LN_RECT_FONTS
        case PERSON_LEVEL_CHILD:
            person.Width    = PERSON_LC_RECT_W
            person.Height   = PERSON_LC_RECT_H
            person.FontSize = PERSON_LC_RECT_FONTS
        }
        person.XPos = x
        person.YPos = y
        if len(highlight) == 1 { person.Highlight = highlight[0] }
        // TODO: Maybe store info about the person instead of the token
        // so we don't have to look up substructures at every redraw.
        person.Token = token

        treeWidgetItems = append(treeWidgetItems, &person)
    }

    treeWidgetEventBox, _ := gtk.EventBoxNew()
    mainCont.Add(treeWidgetEventBox)
    treeWidget, _ := gtk.DrawingAreaNew()
    treeWidgetEventBox.Add(treeWidget)
    drawCb := func(widget *gtk.DrawingArea, cr *cairo.Context) bool {
        //ww, wh := float64(treeWidget.GetAllocatedWidth()), float64(treeWidget.GetAllocatedHeight())
        ww, _ := float64(treeWidget.GetAllocatedWidth()), float64(treeWidget.GetAllocatedHeight())
        //scale := math.Min(ww/100, wh/100)
        scale := ww/100
        cr.Scale(scale, scale)

        cr.SetSourceRGB(0.2, 0.2, 0.2)
        cr.Paint()

        for _, item := range treeWidgetItems {
            item.Draw(cr)
        }

        widget.QueueDraw()
        return false
    }
    treeWidget.Connect("draw", drawCb)

    setCurrentPerson := func(person *Token) {
        treeWidgetItems = []*TreeWidgetRect{}

        //fmt.Println("Looking up FAMS")
        sfamily := tree.LookUpPointer(person.GetFirstChildWithTagValueOr(TAG_FAMS, ""))
        var husband, wife *Token
        var children []*Token
        if sfamily != nil {
            husbPtr := sfamily.GetFirstChildWithTag(TAG_HUSB)
            if husbPtr != nil {
                husband = tree.LookUpPointer(husbPtr.LineVal.GetValueOr(""))
            }
            wifePtr := sfamily.GetFirstChildWithTag(TAG_WIFE)
            if wifePtr != nil {
                wife = tree.LookUpPointer(wifePtr.LineVal.GetValueOr(""))
            }
            childPtrs := sfamily.GetChildrenWithTag(TAG_CHIL)
            for _, childPtr := range childPtrs {
                child := tree.LookUpPointer(childPtr.LineVal.GetValueOr(""))
                if child != nil {
                    children = append(children, child)
                }
            }
        }

        if husband == nil && wife == nil { // If the person wasn't married, draw only them
            addPersonToTreeWidget(person, 50-PERSON_LN_RECT_W/2, 50-PERSON_LN_RECT_H/2, PERSON_LEVEL_NORMAL, true)
        } else {
            addPersonToTreeWidget(husband, 50-PERSON_LN_RECT_W-PERSON_LN_RECT_PAD/2, 50-PERSON_LN_RECT_H/2, PERSON_LEVEL_NORMAL, person==husband)
            addPersonToTreeWidget(wife, 50+PERSON_LN_RECT_PAD/2, 50-PERSON_LN_RECT_H/2, PERSON_LEVEL_NORMAL, person==wife)
        }

        startX := float64(50-float64(len(children))/2.0*(PERSON_LC_RECT_W+PERSON_LC_RECT_PAD))+PERSON_LC_RECT_PAD/2.0
        for i, child := range children {
            xPos := startX+(PERSON_LC_RECT_W+PERSON_LC_RECT_PAD)*float64(i)
            yPos := float64(50-PERSON_LN_RECT_H/2+PERSON_LN_RECT_H+PERSON_LC_RECT_PAD)
            addPersonToTreeWidget(child, xPos, yPos, PERSON_LEVEL_CHILD)
        }

        //fmt.Println("Looking up FAMC")
        cfamily := tree.LookUpPointer(person.GetFirstChildWithTagValueOr(TAG_FAMC, ""))
        var father, mother *Token
        if cfamily != nil {
            fatherPtr := cfamily.GetFirstChildWithTag(TAG_HUSB)
            if fatherPtr != nil {
                father = tree.LookUpPointer(fatherPtr.LineVal.GetValueOr(""))
            }
            motherPtr := cfamily.GetFirstChildWithTag(TAG_WIFE)
            if motherPtr != nil {
                mother = tree.LookUpPointer(motherPtr.LineVal.GetValueOr(""))
            }
        }

        parentOffsX := 0.0
        if husband == nil && wife == nil {
            parentOffsX = 50-PERSON_LN_RECT_W/2.0
        } else if person == husband {
            parentOffsX = 50-PERSON_LN_RECT_W-PERSON_LN_RECT_PAD/2.0
        } else if person == wife {
            parentOffsX = 50+PERSON_LN_RECT_PAD/2.0
        }

        fx := parentOffsX+PERSON_LN_RECT_W/2.0-PERSON_LA_RECT_W-PERSON_LA_RECT_PAD/2.0
        fy := 50-PERSON_LN_RECT_H/2.0-PERSON_LA_RECT_H-PERSON_LA_RECT_PAD
        addPersonToTreeWidget(father, fx, fy, PERSON_LEVEL_ANCESTOR)
        mx := parentOffsX+PERSON_LN_RECT_W/2.0+PERSON_LA_RECT_PAD/2.0
        my := fy
        addPersonToTreeWidget(mother, mx, my, PERSON_LEVEL_ANCESTOR)
    }

    treeWidgetEventBox.Connect("button-release-event", func(widget *gtk.EventBox, event *gdk.Event){
        eventBtn := gdk.EventButtonNewFromEvent(event)
        ww, _ := float64(treeWidget.GetAllocatedWidth()), float64(treeWidget.GetAllocatedHeight())
        scale := ww/100
        x := eventBtn.X()/scale
        y := eventBtn.Y()/scale

        if eventBtn.Button() == gdk.BUTTON_PRIMARY {
            for _, item := range treeWidgetItems {
                if item.Token != nil && item.IsPointInside(x, y) {
                    setCurrentPerson(item.Token)
                    break
                }
            }
        }
    })

    treeWidgetEventBox.AddEvents((int)(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK))
    treeWidgetEventBox.Connect("motion-notify-event", func(widget *gtk.EventBox, event *gdk.Event){
        eventMot := gdk.EventMotionNewFromEvent(event)
        eventX, eventY := eventMot.MotionVal()

        ww, _ := float64(treeWidget.GetAllocatedWidth()), float64(treeWidget.GetAllocatedHeight())
        scale := ww/100
        x := eventX/scale
        y := eventY/scale

        for _, item := range treeWidgetItems {
            if item.IsPointInside(x, y) {
                item.OnMouseEnter()
            } else {
                item.OnMouseLeave()
            }
        }
    })

    individListWidget.ListWidget.Connect("cursor-changed", func(widg *gtk.TreeView) bool {
        sel, _ := widg.GetSelection()
        // Don't trigger when the selection changes
        // because the searching rebuilds the list
        if !individListWidget.IsSearching && sel.CountSelectedRows() == 1 {
            model, iter, _ := sel.GetSelected()
            selVal, _ := model.ToTreeModel().GetValue(iter, 0)
            selectedIndex, _ := selVal.GoValue()

            people := tree.GetTokensWithTag(TAG_INDI)
            person := people[selectedIndex.(int)]
            setCurrentPerson(person)

        }
        return false
    })

    win.Add(mainCont)
    win.ShowAll()
    gtk.Main()
}
