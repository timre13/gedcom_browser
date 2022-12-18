package token

import (
    . "gedcom_browser/optional"
    "strconv"
    "unicode"
    "unicode/utf8"
    "os"
    "bufio"
    "html"
    "fmt"
    "strings"
)

const MAX_LEVEL_COUNT = 16

type Gedcom struct {
    Tokens []Token
}

func (this *Gedcom) GetTokenByPath(tags []Tag) *Token {
    var currToken *Token
    for _, token := range this.Tokens {
        if token.Tag == tags[0] {
            currToken = &token
            break
        }
    }
    if currToken == nil {
        return nil
    }

    for _, tag := range tags[1:] {
        match := currToken.GetFirstChildWithTag(tag)
        if match == nil {
            return nil
        }
        currToken = match
    }
    return currToken
}

type Tag int
const (
    TAG_INVALID Tag = iota
    TAG_HEAD
    TAG_GEDC
    TAG_VERS
    TAG_FORM
    TAG_CHAR
    TAG_SOUR
    TAG_NAME
    TAG_TRLR
    TAG_SUBM
    TAG_ABBR
    TAG_ADDR
    TAG_ADOP
    TAG_ADR1
    TAG_AGNC
    TAG_BIRT
    TAG_BURI
    TAG_CALN
    TAG_CHIL
    TAG_CITY
    TAG_CORP
    TAG_CTRY
    TAG_DATA
    TAG_DATE
    TAG_DEAT
    TAG_EVEN
    TAG_FAMC
    TAG_FAMS
    TAG_FILE
    TAG_GIVN
    TAG_HUSB
    TAG_LANG
    TAG_MARR
    TAG_PAGE
    TAG_PEDI
    TAG_PHON
    TAG_PLAC
    TAG_POST
    TAG_REPO
    TAG_RESI
    TAG_SEX
    TAG_STAE
    TAG_SURN
    TAG_TIME
    TAG_TITL
    TAG_WIFE
    TAG_WWW
    TAG_ADR2
    TAG_CAUS
    TAG_CONC
    TAG_CONT
    TAG_DEST
    TAG_DIV
    TAG_EDUC
    TAG_EMAIL
    TAG_FAM
    TAG_INDI
    TAG_NOTE
    TAG_OCCU
    TAG_RELI
    TAG_RIN
    TAG_TYPE
)

func strToTag(str string) Tag {
    switch str {
    case "HEAD": return TAG_HEAD
    case "GEDC": return TAG_GEDC
    case "VERS": return TAG_VERS
    case "FORM": return TAG_FORM
    case "CHAR": return TAG_CHAR
    case "SOUR": return TAG_SOUR
    case "NAME": return TAG_NAME
    case "TRLR": return TAG_TRLR
    case "SUBM": return TAG_SUBM
    case "ABBR": return TAG_ABBR
    case "ADDR": return TAG_ADDR
    case "ADOP": return TAG_ADOP
    case "ADR1": return TAG_ADR1
    case "AGNC": return TAG_AGNC
    case "BIRT": return TAG_BIRT
    case "BURI": return TAG_BURI
    case "CALN": return TAG_CALN
    case "CHIL": return TAG_CHIL
    case "CITY": return TAG_CITY
    case "CORP": return TAG_CORP
    case "CTRY": return TAG_CTRY
    case "DATA": return TAG_DATA
    case "DATE": return TAG_DATE
    case "DEAT": return TAG_DEAT
    case "EVEN": return TAG_EVEN
    case "FAMC": return TAG_FAMC
    case "FAMS": return TAG_FAMS
    case "FILE": return TAG_FILE
    case "GIVN": return TAG_GIVN
    case "HUSB": return TAG_HUSB
    case "LANG": return TAG_LANG
    case "MARR": return TAG_MARR
    case "PAGE": return TAG_PAGE
    case "PEDI": return TAG_PEDI
    case "PHON": return TAG_PHON
    case "PLAC": return TAG_PLAC
    case "POST": return TAG_POST
    case "REPO": return TAG_REPO
    case "RESI": return TAG_RESI
    case "SEX":  return TAG_SEX
    case "STAE": return TAG_STAE
    case "SURN": return TAG_SURN
    case "TIME": return TAG_TIME
    case "TITL": return TAG_TITL
    case "WIFE": return TAG_WIFE
    case "WWW":  return TAG_WWW
    case "ADR2": return TAG_ADR2
    case "CAUS": return TAG_CAUS
    case "CONC": return TAG_CONC
    case "CONT": return TAG_CONT
    case "DEST": return TAG_DEST
    case "DIV":  return TAG_DIV
    case "EDUC": return TAG_EDUC
    case "EMAIL":return TAG_EMAIL
    case "FAM":  return TAG_FAM
    case "INDI": return TAG_INDI
    case "NOTE": return TAG_NOTE
    case "OCCU": return TAG_OCCU
    case "RELI": return TAG_RELI
    case "RIN":  return TAG_RIN
    case "TYPE": return TAG_TYPE
    }
    return TAG_INVALID
}

func tagToStr(tag Tag) string {
    return [...]string{
        "<INVALID>",
        "HEAD",
        "GEDC",
        "VERS",
        "FORM",
        "CHAR",
        "SOUR",
        "NAME",
        "TRLR",
        "SUBM",
        "ABBR",
        "ADDR",
        "ADOP",
        "ADR1",
        "AGNC",
        "BIRT",
        "BURI",
        "CALN",
        "CHIL",
        "CITY",
        "CORP",
        "CTRY",
        "DATA",
        "DATE",
        "DEAT",
        "EVEN",
        "FAMC",
        "FAMS",
        "FILE",
        "GIVN",
        "HUSB",
        "LANG",
        "MARR",
        "PAGE",
        "PEDI",
        "PHON",
        "PLAC",
        "POST",
        "REPO",
        "RESI",
        "SEX",
        "STAE",
        "SURN",
        "TIME",
        "TITL",
        "WIFE",
        "WWW",
        "ADR2",
        "CAUS",
        "CONC",
        "CONT",
        "DEST",
        "DIV",
        "EDUC",
        "EMAIL",
        "FAM",
        "INDI",
        "NOTE",
        "OCCU",
        "RELI",
        "RIN",
        "TYPE",
    }[tag]
}

type Token struct {
    Level       int
    Xref        Optional[string]
    Tag         Tag
    LineVal     Optional[string]
    Subitems    []Token
}

func (this *Token) String() string {
    return fmt.Sprintf("lvl=%d, xref='%s', tag=%s, val=\"%s\", subs=%d",
        this.Level, this.Xref.GetValueOr(""), tagToStr(this.Tag), this.LineVal.GetValueOr(""), len(this.Subitems))
}

func (this *Token) GetFirstChildWithTag(tag Tag) *Token {
    for _, child := range this.Subitems {
        if child.Tag== tag {
            return &child
        }
    }
    return nil
}

func isUcLetter(char rune) bool {
    return unicode.IsUpper(char) && unicode.IsLetter(char)
}

func isDigit(char rune) bool {
    return unicode.IsDigit(char)
}

func isTagChar(char rune) bool {
    return char == '_' || isDigit(char) || isUcLetter(char);
}

func isNonAt(char rune) bool {
    return char == 0x09 || (char >= 0x20 && char <= 0x3f) || (char >= 0x41 && char <= 0x10ffff)
}

func isNonEol(char rune) bool {
    return char == 0x09 || (char >= 0x20 && char <= 0x10ffff)
}

func getRuneFromStr(str string, index int) rune {
    char, _ := utf8.DecodeRuneInString(str[index:])
    return char
}

func isAllRunesOfString(str string, cond func(rune) bool) bool {
    for _, char := range str {
        if !cond(char) {
            return false
        }
    }
    return true
}

func isLineStr(str string) bool {
    startsWithNonAt := isNonAt(getRuneFromStr(str, 0))
    startsWithDoubleAt := !isNonAt(getRuneFromStr(str, 0)) && !isNonAt(getRuneFromStr(str, 1))
    if !startsWithNonAt && !startsWithDoubleAt { return false }
    
    startByteI := 0
    if startsWithDoubleAt {
        startByteI = 2
    }
    return isAllRunesOfString(str[startByteI:], func(char rune) bool { return isNonEol(char) })
}

/*
 * Check if the string is in the format `atsign 1*tagchar atsign`.
 */
func isReference(value string) bool {
    return len(value) >= 3 && value[0] == '@' &&
        isAllRunesOfString(value[1:utf8.RuneCountInString(value)-1], isTagChar) &&
        value[utf8.RuneCountInString(value)-1] == '@'
}

func isPointer(value string) bool {
    return value == "@VOID@" || isReference(value)
}

func isValidLineVal(value string) bool {
    return isPointer(value) || isLineStr(value)
}

func concatSliceWithSpaces(slice []string) string {
    if len(slice) == 0 { return "" }

    output := ""
    for _, val := range slice {
        output += " "+val
    }
    return output[1:]
}

func genTokenFromLine(line string) *Token {
    // Remove BOM if found
    if getRuneFromStr(line, 0) == 0xfeff {
        line = line[utf8.RuneLen(getRuneFromStr(line, 0)):]
    }

    output := Token{}
    fields := strings.Split(line, " ")
    output.Level, _ = strconv.Atoi(fields[0])
    if isReference(fields[1]) {
        xref := fields[1][1:utf8.RuneCountInString(fields[1])-1]
        output.Xref.SetValue(xref)
        fmt.Printf("\tReference: \"%s\"\n", output.Xref.GetValue())
        output.Tag = strToTag(fields[2])
        fmt.Printf("\tTag: \"%s\" -> %s\n", fields[2], tagToStr(output.Tag))
    } else {
        output.Tag = strToTag(fields[1])
        fmt.Printf("\tTag: \"%s\" -> %s\n", fields[1], tagToStr(output.Tag))

        lineVal := concatSliceWithSpaces(fields[2:])
        if !isValidLineVal(lineVal) {
            fmt.Printf("\tInvalid line value: \"%s\"\n", lineVal)
            return nil
        }
        // TODO: Remove HTML tags
        output.LineVal.SetValue(html.UnescapeString(lineVal))
    }
    return &output
}

func LoadTokensFromFile(path string) []Token {
    fmt.Println("Loading tokens from file...")

    file, err := os.Open(path)
    if err != nil { panic(err) }
    defer file.Close()

    tokens := []Token{}

    scanner := bufio.NewScanner(file)
    lineNum := 1
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Printf("Line #%d: \"%s\"\n", lineNum, line)
        token := genTokenFromLine(line)
        if token != nil {
            tokens = append(tokens, *token)
        }
        lineNum++
    }

    //for i:=0; i < 300; i++ {
    //    for _, token := range tokens {
    //        fmt.Println(token.Level)
    //    }
    //}

    return tokens
}

func BuildTreeFromTokens(tokens []Token) Gedcom {
    fmt.Println("Building tree...")

    output := Gedcom{}
    var path [MAX_LEVEL_COUNT]*[]Token
    path[0] = &output.Tokens
    level := 0

    addSibling := func(token *Token) {
        *path[level] = append(*path[level], *token)
        path[level+1] = &(*path[level])[len((*path[level]))-1].Subitems
    }

    addChild := func(token *Token) {
        level++
        addSibling(token)
    }

    for _, token := range tokens {
        fmt.Println(token.String())

        if token.Level == level+1 {
            fmt.Printf("Child branch with level %d\n", token.Level)
            addChild(&token)
        } else if token.Level == level {
            fmt.Printf("Sibling branch with level %d\n", token.Level)
            addSibling(&token)
        } else if token.Level < level {
            fmt.Printf("Upper branch with level %d\n", token.Level)
            level = token.Level
            addSibling(&token)
        } else {
            panic("Token skipped a level")
        }
    }

    return output
}

func PrintTree(tree []Token, level int) {
    for _, token := range tree {
        fmt.Print(strings.Repeat("  ", level))
        fmt.Printf("(%s)\n", token.String())
        PrintTree(token.Subitems, level+1)
    }
}
