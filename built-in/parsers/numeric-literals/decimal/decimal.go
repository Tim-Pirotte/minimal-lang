package decimalparser

type number []uint

type WholeParser struct {
    numbers []number
}

type WholeExponentParser struct {
    numbers []wholeExponent
}

type WholeNegativeExponentParser struct {
    numbers []wholeExponent
}

type FractionParser struct {
    numbers []fraction
}

type FractionExponentParser struct {
    numbers []fractionExponent
}

type FractionNegativeExponentParser struct {
    numbers []fractionExponent
}

type fraction struct {
    decimal    number
    fractional number
}

type wholeExponent struct {
    whole    number
    exponent number
}

type fractionExponent struct {
    fraction fraction
    exponent number
}

func NewWholeParser() *WholeParser {
    // These are initialized lazily
    return &WholeParser{nil}
}

func NewWholeExponentParser() *WholeExponentParser {
    return &WholeExponentParser{nil}
}

func NewWholeNegativeExponentParser() *WholeNegativeExponentParser {
    return &WholeNegativeExponentParser{nil}
}

func NewFractionParser() *FractionParser {
    return &FractionParser{nil}
}

func NewFractionExponentParser() *FractionExponentParser {
    return &FractionExponentParser{nil}
}

func NewFractionNegativeExponentParser() *FractionNegativeExponentParser {
    return &FractionNegativeExponentParser{nil}
}
