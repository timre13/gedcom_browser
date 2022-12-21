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
    BoxWidget          *gtk.Box
    ScrollWidget        *gtk.ScrolledWindow
    ListWidget          *gtk.TreeView
    listStore           *gtk.ListStore
    tree                *token.Gedcom
    SearchEntry         *gtk.Entry
    searchTerm          string
    IsSearching         bool
}

func IndiListWidgetNew(tree *token.Gedcom) *IndiListWidget {
    widget := IndiListWidget{}

    widget.BoxWidget, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

    widget.SearchEntry, _ = gtk.EntryNew()
    widget.BoxWidget.PackStart(widget.SearchEntry, false, false, 0)
    widget.SearchEntry.Connect("changed", func(){
        text, _ := widget.SearchEntry.GetText()
        widget.Search(text)
    })

    widget.ScrollWidget, _ = gtk.ScrolledWindowNew(nil, nil)
    widget.BoxWidget.PackStart(widget.ScrollWidget, true, true, 0)

    widget.listStore, _ = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
    widget.ListWidget, _ = gtk.TreeViewNewWithModel(widget.listStore)
    widget.ScrollWidget.Add(widget.ListWidget)

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

    widget.tree = tree
    widget.fetchItems()

    return &widget
}

func (this *IndiListWidget) fetchItems() {
    this.searchTerm = strings.ToLower((strings.ToLower(this.searchTerm)))
    searchTerms := strings.Split(this.searchTerm, " ")

    isSearchMatch := func(value string) bool {
        if this.searchTerm == "" {
            return true
        }

        value = strings.ToLower(value)

        for _, term := range searchTerms {
            if !strings.Contains(value, term) {
                return false
            }
        }
        return true
    }

    i := 0
    for _, tok := range this.tree.GetTokensWithTag(token.TAG_INDI) {
        nameStr := tok.GetFirstChildWithTagValueOr(token.TAG_NAME, "")
        if isSearchMatch(nameStr) {
            name := ParseName(nameStr)

            iter := this.listStore.Append()
            this.listStore.SetValue(iter, ColumnIndex, i)
            this.listStore.SetValue(iter, ColumnFirstName, name.FirstName)
            this.listStore.SetValue(iter, ColumnLastName, name.LastName)
        }
        i++
    }
}

func (this *IndiListWidget) Search(term string) {
    this.IsSearching = true
    this.searchTerm = term
    this.listStore.Clear()
    this.fetchItems()
    this.IsSearching = false
}
