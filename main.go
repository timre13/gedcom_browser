package main

import (
    . "gedcom_browser/optional"
    "fmt"
    "strconv"
    "strings"
    "unicode"
    "unicode/utf8"
    "os"
    "bufio"
)

type Gedcom struct {
    Tokens []Token
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
    case "SEX": return TAG_SEX
    case "STAE": return TAG_STAE
    case "SURN": return TAG_SURN
    case "TIME": return TAG_TIME
    case "TITL": return TAG_TITL
    case "WIFE": return TAG_WIFE
    case "WWW": return TAG_WWW 
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
    }[tag]
}

type Token struct {
    Level       int
    Xref        Optional[string]
    Tag         Tag
    LineVal     Optional[string]
}

func (this *Token) String() string {
    return fmt.Sprintf("lvl=%d, xref='%s', tag=%s, val=\"%s\"",
        this.Level, this.Xref.GetValueOr(""), tagToStr(this.Tag), this.LineVal.GetValueOr(""))
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
    return len(value) >= 3 && value[0] == '@' && isAllRunesOfString(value[1:utf8.RuneCountInString(value)-1], isTagChar) && value[utf8.RuneCountInString(value)-1] == '@'
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

func genTokenFromLine(line string) Token {
    // Remove BOM if found
    if getRuneFromStr(line, 0) == 0xfeff {
        line = line[utf8.RuneLen(getRuneFromStr(line, 0)):]
    }

    output := Token{}
    fields := strings.Split(line, " ")
    output.Level, _ = strconv.Atoi(fields[0])
    fmt.Printf("Line: \"%s\"\n", line)
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
            return Token{}
        }
        output.LineVal.SetValue(lineVal)
    }
    return output
}

func loadFile(path string) Gedcom {
    file, err := os.Open(path)
    if err != nil { panic(err) }
    defer file.Close()

    output := Gedcom{}

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        output.Tokens = append(output.Tokens, genTokenFromLine(scanner.Text()))
    }

    return output
}

func main() {
    //file := loadFile("./samples/MINIMAL555.GED")
    file := loadFile("./samples/555SAMPLE.GED")
    fmt.Println("-------- Tokens --------")
    for i, token := range file.Tokens {
        fmt.Printf("%d : %s\n", i+1, token.String())
    }
}