package main

import (
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/cairo"
    "fmt"
    "strings"
    "math"
    . "gedcom_browser/token"
    "gedcom_browser/widgets"
)

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
    individListWidgetSelected := 0
    individListWidget.ScrollWidget.SetSizeRequest(300, 0)
    mainCont.Add(individListWidget.ScrollWidget)
    
    drawPerson := func(cr *cairo.Context, token *Token, x float64, y float64, highlight ...bool) {
        if len(highlight) == 1 && highlight[0] {
            cr.SetSourceRGB(1, 1, 1)
            cr.Rectangle(x, y, 45, 35)
            cr.SetLineWidth(0.3)
            cr.Stroke()
        }

        // This happens for example when a person's husband/wife is not known
        // Draw a placeholder
        if token == nil {
            cr.SetSourceRGB(0.25, 0.25, 0.25)
            cr.Rectangle(x, y, 45, 35)
            cr.Fill()
            return
        }

        gender := token.GetFirstChildWithTagValueOr(TAG_SEX, "U")
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
        cr.Rectangle(x, y, 45, 35)
        cr.Fill()

        cr.SetSourceRGB(1.0, 1.0, 1.0)
        cr.MoveTo(x, y+5)
        cr.SetFontSize(4)
        name := indi_list_widget.ParseName(token.GetFirstChildWithTagValueOr(TAG_NAME, ""))
        cr.ShowText(name.FirstName+" "+name.LastName)

        cr.SetFontSize(3)
        cr.MoveTo(x, y+8)
        birthYear := "???"
        birthToken := token.GetFirstChildWithTag(TAG_BIRT)
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
        deathToken := token.GetFirstChildWithTag(TAG_DEAT)
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
    }

    treeWidget, _ := gtk.DrawingAreaNew()
    mainCont.Add(treeWidget)
    drawCb := func(widget *gtk.DrawingArea, cr *cairo.Context) bool {
        ww, wh := float64(treeWidget.GetAllocatedWidth()), float64(treeWidget.GetAllocatedHeight())
        scale := math.Min(ww/100, wh/100)
        cr.Scale(scale, scale)

        cr.SetSourceRGB(0.2, 0.2, 0.2)
        cr.Paint()

        people := tree.GetTokensWithTag(TAG_INDI)
        person := people[individListWidgetSelected]

        family := tree.LookUpPointer(person.GetFirstChildWithTagValueOr(TAG_FAMS, ""))
        var husb *Token
        var wife *Token
        if family != nil {
            husbPtr := family.GetFirstChildWithTag(TAG_HUSB)
            if husbPtr != nil {
                husb = tree.LookUpPointer(husbPtr.LineVal.GetValueOr(""))
            }
            wifePtr := family.GetFirstChildWithTag(TAG_WIFE)
            if wifePtr != nil {
                wife = tree.LookUpPointer(wifePtr.LineVal.GetValueOr(""))
            }
        }

        if husb == nil && wife == nil { // If the person wasn't married, draw only them
            drawPerson(cr, person, 30+22.5, 30, true)
        } else {
            drawPerson(cr, husb, 30, 30, person==husb)
            drawPerson(cr, wife, 30+45+5, 30, person==wife)
        }

        widget.QueueDraw()
        return false
    }
    treeWidget.Connect("draw", drawCb)

    individListWidget.ListWidget.Connect("cursor-changed", func(widg *gtk.TreeView) bool {
        sel, _ := widg.GetSelection()
        model, iter, _ := sel.GetSelected()
        selVal, _ := model.ToTreeModel().GetValue(iter, 0)
        selectedIndex, _ := selVal.GoValue()
        individListWidgetSelected = selectedIndex.(int)
        //fmt.Println("Selection changed to ", individListWidgetSelected)
        return false
    })

    win.Add(mainCont)
    win.ShowAll()
    gtk.Main()
}
