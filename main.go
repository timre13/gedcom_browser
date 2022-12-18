package main

import (
    "fmt"
    "strings"
    . "gedcom_browser/token"
)

func main() {
    //file := loadFile("./samples/MINIMAL555.GED")
    file := loadFile("./samples/555SAMPLE.GED")
    tokens := loadTokensFromFile("./samples/555SAMPLE.GED")
    fmt.Println("-------- Tokens --------")
    for i, token := range tokens {
        fmt.Printf("%d : %s\n", i+1, token.String())
    }

    tree := BuildTreeFromTokens(tokens)
    fmt.Printf("There are %d parent tokens\n", len(tree.Tokens))
    fmt.Println("-------- Tree --------")
    PrintTree(tree.Tokens, 0)

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
}
