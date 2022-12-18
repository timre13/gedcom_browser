package main

import (
    "github.com/gotk3/gotk3/gtk"
    "fmt"
    "strings"
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

    win.Add(indi_list_widget.NewIndiListWidget(&tree))

    win.ShowAll()
    gtk.Main()
}
