package indi_list_widget

import (
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
    "gedcom_browser/token"
    "strings"
)

const (
    columnIndex = iota
    columnFirstName
    columnLastName
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

func NewIndiListWidget(tree *token.Gedcom) *gtk.ScrolledWindow {
    scrollWidget, _ := gtk.ScrolledWindowNew(nil, nil)
    listStore, _ := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
    listWidget, _ := gtk.TreeViewNewWithModel(listStore)

    addCol := func(text string, coli int) {
        rend, _ := gtk.CellRendererTextNew()
        tvCol, _ := gtk.TreeViewColumnNewWithAttribute(text, rend, "text", coli)
        tvCol.SetResizable(true)
        tvCol.SetSortColumnID(coli)
        tvCol.SetExpand(true)
        if coli == columnIndex {
            // Column 0 is only used to distinguish between items, so hide it
            tvCol.SetVisible(false)
        }
        listWidget.AppendColumn(tvCol)
    }
    addCol("#", columnIndex)
    addCol("First Name", columnFirstName)
    addCol("Last Name", columnLastName)

    i := 0
    for _, tok := range tree.GetTokensWithTag(token.TAG_INDI) {
        nameStr := tok.GetFirstChildWithTagValueOr(token.TAG_NAME, "")
        name := ParseName(nameStr)

        iter := listStore.Append()
        listStore.SetValue(iter, columnIndex, i)
        listStore.SetValue(iter, columnFirstName, name.FirstName)
        listStore.SetValue(iter, columnLastName, name.LastName)
        i++
    }

    scrollWidget.Add(listWidget)
    return scrollWidget
}
