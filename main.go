package main

import (
    "fmt"
    "log"
    "strconv"
)

type state int

const (
    STAT_0 state = iota
    STAT_NUM
)

type unit struct {
    eletype string
    lit string
}


func generate_units(str string) []unit {
    units := []unit{}
    word := ""
    current_state := STAT_0

    if str == "" {
        return units
    }

    i := 0
    for {
        if i >= len(str) {
            if word != "" {
                units = append(units, unit{eletype:"ET_NUM", lit:word})
                word = ""
            }
            units = append(units, unit{eletype: "ET_EOF", lit:""})
            return units
        }
        current_char := str[i]
        fmt.Printf("current_char: %c\n", current_char)
        switch current_state {
        case STAT_0:
            if current_char == '(' {
                units = append(units, unit{eletype: "ET_SEP", lit:"("})
            } else if current_char == ')' {
                units = append(units, unit{eletype: "ET_SEP", lit:")"})
            } else if current_char == '+' {
                units = append(units, unit{eletype: "ET_OP", lit:"+"})
            } else if current_char == '-' {
                units = append(units, unit{eletype: "ET_OP", lit:"-"})
            } else if current_char == '*' {
                units = append(units, unit{eletype: "ET_OP", lit:"*"})
            } else if current_char == '/' {
                units = append(units, unit{eletype: "ET_OP", lit:"/"})
            } else if current_char >= '0' && current_char <= '9' {
                word += string(current_char)
                current_state = STAT_NUM
            }
        case STAT_NUM:
            if current_char >= '0' && current_char <= '9' {
                word += string(current_char)
            } else {
                units = append(units, unit{eletype: "ET_NUM", lit:word})
                current_state = STAT_0
                word = ""
                i--
            }
        }
        i++
    }
}

// <Exp> ::= <Exp> + <Term> | <Exp> - <Term> | <Term>
// <Term> ::= <Term> * <Factor> | <Term> / <Factor> | <Factor>
// <Factor> ::= x | y | ... | ( <Exp> ) | - <Factor>

type Expresstion struct {
    op string
    left *Term 
    right *Expresstion
}

func (e *Expresstion) value() int {
    if e.op == "+" {
        return e.left.value() + e.right.value()
    } else if e.op == "-" {
        return e.left.value() - e.right.value()
    } else if e.op == "" {
        return e.left.value()
    }
    return 0
}

type Term struct {
    op string
    left *Factor
    right *Term
}

func (t *Term) value() int {
    if t.op == "*" {
        return t.left.value() * t.right.value()
    } else if t.op == "/" {
        return t.left.value() / t.right.value()
    } else if t.op == "" {
        return t.left.value()
    }

    return 0
}

type Factor struct {
    kind int // 0 -> int value, 1 -> expresstion, 2 -> - Factor
    val int
    quotedExpresstion *Expresstion
}

func (f *Factor) value() int {
    if f.kind == 0 {
        return f.val
    } else if f.kind == 1 {
        return f.quotedExpresstion.value()
    }

    return 0
}

func parse_E(units []unit) (*Expresstion, []unit) {
	log.Printf("Enter the parse_E, rest_units = %v\n", units)
    if len(units) == 0 {
        return nil, []unit{}
    }

    var prightExpress *Expresstion
	var rest_units1 []unit

    pterm, rest_units := parse_T(units)

    if len(rest_units) == 1 {
        return &Expresstion{op:"", left:pterm, right: nil}, rest_units
    }

    if rest_units[0].eletype == "ET_OP" && (rest_units[0].lit == "+" || rest_units[0].lit == "-") {
        prightExpress, rest_units1 = parse_E(rest_units[1:])
        log.Printf("Exit parse_E, rest_units = %v\n", rest_units)
        return &Expresstion{op:rest_units[0].lit, left: pterm, right: prightExpress}, rest_units1
    } else {
        log.Printf("Error parse_E, rest_units = %v\n", rest_units)
        return &Expresstion{op:"", left: pterm, right: nil}, rest_units
    }
}

func parse_T(units []unit) (*Term, []unit) {
	log.Printf("Enter the parse_T, rest_units = %v\n", units)
	if len(units) == 0 {
	    return nil, []unit{}
    }

    var prightTerm *Term
	var rest_units1 []unit

    pfactor, rest_units := parse_F(units)
    if len(rest_units) == 1 {
        return &Term{op: "", left: pfactor, right: nil}, rest_units
    }

    if rest_units[0].eletype == "ET_OP" && (rest_units[0].lit == "*" || rest_units[0].lit == "/") {
        prightTerm, rest_units1 = parse_T(rest_units[1:])
        log.Printf("Exit the parse_T, rest_units = %v\n", rest_units)
        return &Term{op: rest_units[0].lit, left: pfactor, right: prightTerm}, rest_units1
    } else {
        log.Printf("Error the parse_T, rest_units = %v\n", rest_units)
        return &Term{op: "", left: pfactor, right: nil}, rest_units
    }
}

func parse_F(units []unit) (*Factor, []unit) {
	log.Printf("Enter parse_F, rest_units = %v\n", units)
	if len(units) == 0 {
	    return nil, []unit{}
    }
    current_unit := units[0]
    if current_unit.eletype == "ET_EOF" {
        return nil, []unit{}
    }

    if current_unit.eletype == "ET_NUM" {
        v, err := strconv.Atoi(current_unit.lit)
        if err != nil {
            log.Fatalf("parse_F error: %s\n", current_unit.lit)
        }

        log.Printf("Exit the parse_F, rest_units = %v\n", units[1:])
        return &Factor{kind: 0, val: v, quotedExpresstion: nil}, units[1:]
    } else if current_unit.eletype == "ET_SEP" && current_unit.lit == "(" {
        pexpresstion, rest_units := parse_E(units[1:])
        if rest_units[0].lit == ")" {
            log.Printf("Exit the parse_F, rest_units = %v\n", units[1:])
            return &Factor{kind: 1, val: 0, quotedExpresstion: pexpresstion}, rest_units[1:]
        } else {
            log.Printf("Exit the parse_F, rest_units = %v\n", units[1:])
            return nil, rest_units
        }
    }

    return nil, units
}


func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    units := generate_units("11-1+1")
    pexpresstion, rest_units := parse_E(units)
    fmt.Printf("%v, %v\n", *pexpresstion, rest_units)
    fmt.Printf("eval result = %v\n", pexpresstion.value())
}
