package decimalparser

type number []uint

type WholeDecimalParser struct {
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
