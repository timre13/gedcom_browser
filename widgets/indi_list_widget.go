package indi_list_widget

import (
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
    "gedcom_browser/token"
    "strings"
)

const (
    ColumnIndex = iota
    ColumnFirstName
    ColumnLastName
)

type Name struct {
    FirstName string
    LastName  string
}

func ParseName(name string) Name {
    output := Name{}
    if strings.Count(name, "/") >= 2 {
        startsWithSep := strings.HasPrefix(name, "/")
        for i, part := range strings.Split(name, "/") {
            j := i
            if startsWithSep {
                j++
            }
            if j % 2 == 1 {
                output.LastName += strings.Trim(part, "/") + " "
            } else {
                output.FirstName += part + " "
            }
        }
    } else {
        names := strings.Split(name, " ")
        for _, name := range names[:len(names)-1] {
            output.FirstName += name + " "
        }
        output.LastName = names[len(names)-1]
    }
    output.FirstName = strings.TrimSpace(output.FirstName)
    output.LastName = strings.TrimSpace(output.LastName)
    if output.FirstName == "" {
        output.FirstName = "???"
    }
    if output.LastName == "" {
        output.LastName = "???"
    }
    return output
}

type IndiListWidget struct {
    ScrollWidget        *gtk.ScrolledWindow
    ListWidget          *gtk.TreeView
    ListStore           *gtk.ListStore
}

func IndiListWidgetNew(tree *token.Gedcom) *IndiListWidget {
    widget := IndiListWidget{}

    widget.ScrollWidget, _ = gtk.ScrolledWindowNew(nil, nil)
    widget.ListStore, _ = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
    widget.ListWidget, _ = gtk.TreeViewNewWithModel(widget.ListStore)

    addCol := func(text string, coli int) {
        rend, _ := gtk.CellRendererTextNew()
        tvCol, _ := gtk.TreeViewColumnNewWithAttribute(text, rend, "text", coli)
        tvCol.SetResizable(true)
        tvCol.SetSortColumnID(coli)
        tvCol.SetExpand(true)
        if coli == ColumnIndex {
            tvCol.SetVisible(false)
        }
        widget.ListWidget.AppendColumn(tvCol)
    }
    addCol("#", ColumnIndex)
    addCol("First Name", ColumnFirstName)
    addCol("Last Name", ColumnLastName)

    i := 0
    for _, tok := range tree.GetTokensWithTag(token.TAG_INDI) {
        nameStr := tok.GetFirstChildWithTagValueOr(token.TAG_NAME, "")
        name := ParseName(nameStr)

        iter := widget.ListStore.Append()
        widget.ListStore.SetValue(iter, ColumnIndex, i)
        widget.ListStore.SetValue(iter, ColumnFirstName, name.FirstName)
        widget.ListStore.SetValue(iter, ColumnLastName, name.LastName)
        i++
    }

    widget.ScrollWidget.Add(widget.ListWidget)
    return &widget
}
