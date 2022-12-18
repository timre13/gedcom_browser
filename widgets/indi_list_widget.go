package indi_list_widget

import (
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
    "gedcom_browser/token"
    "strings"
)

const (
    columnFirstName = iota
    columnLastName
    columnFgColor
)

type Name struct {
    FirstName string
    LastName  string
}

func parseName(name string) Name {
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
    listStore, _ := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
    listWidget, _ := gtk.TreeViewNewWithModel(listStore)

    addCol := func(text string, coli int) {
        rend, _ := gtk.CellRendererTextNew()
        tvCol, _ := gtk.TreeViewColumnNewWithAttribute(text, rend, "text", coli)
        listWidget.AppendColumn(tvCol)
    }
    addCol("First Name", columnFirstName)
    addCol("Last Name", columnLastName)

    for _, tok := range tree.Tokens {
        if tok.Tag == token.TAG_INDI {
            nameToken := tok.GetFirstChildWithTag(token.TAG_NAME)
            nameStr := ""
            if nameToken != nil {
                nameStr = nameToken.LineVal.GetValueOr("")
            }
            name := parseName(nameStr)

            iter := listStore.Append()
            listStore.SetValue(iter, columnFirstName, name.FirstName)
            listStore.SetValue(iter, columnLastName, name.LastName)
        }
    }

    scrollWidget.Add(listWidget)
    return scrollWidget
}
