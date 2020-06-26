package hash

// uses github.com/tdewolff/hasher
//go:generate hasher -type=Hash -file=hash.go

// Hash defines perfect hashes for a predefined list of strings
type Hash uint32

// Identifiers for the hashes associated with the text in the comments.
const (
	Arccos  Hash = 0x1c06 // arccos
	Arccosh Hash = 0x1c07 // arccosh
	Arcsin  Hash = 0x6    // arcsin
	Arcsinh Hash = 0x7    // arcsinh
	Arctan  Hash = 0x706  // arctan
	Arctanh Hash = 0x707  // arctanh
	Cbrt    Hash = 0xe04  // cbrt
	Cos     Hash = 0x1f03 // cos
	Cosh    Hash = 0x1f04 // cosh
	Erf     Hash = 0x1203 // erf
	Exp     Hash = 0x1503 // exp
	Gamma   Hash = 0x1805 // gamma
	Ln      Hash = 0x2302 // ln
	Log     Hash = 0x2503 // log
	Log10   Hash = 0x2505 // log10
	Log2    Hash = 0x2a04 // log2
	Sin     Hash = 0x303  // sin
	Sinh    Hash = 0x304  // sinh
	Sqrt    Hash = 0x2e04 // sqrt
	Tan     Hash = 0xa03  // tan
	Tanh    Hash = 0xa04  // tanh
)

var HashMap = map[string]Hash{
	"arccos":  Arccos,
	"arccosh": Arccosh,
	"arcsin":  Arcsin,
	"arcsinh": Arcsinh,
	"arctan":  Arctan,
	"arctanh": Arctanh,
	"cbrt":    Cbrt,
	"cos":     Cos,
	"cosh":    Cosh,
	"erf":     Erf,
	"exp":     Exp,
	"gamma":   Gamma,
	"ln":      Ln,
	"log":     Log,
	"log10":   Log10,
	"log2":    Log2,
	"sin":     Sin,
	"sinh":    Sinh,
	"sqrt":    Sqrt,
	"tan":     Tan,
	"tanh":    Tanh,
}

// String returns the text associated with the hash.
func (i Hash) String() string {
	return string(i.Bytes())
}

// Bytes returns the text associated with the hash.
func (i Hash) Bytes() []byte {
	start := uint32(i >> 8)
	n := uint32(i & 0xff)
	if start+n > uint32(len(_Hash_text)) {
		return []byte{}
	}
	return _Hash_text[start : start+n]
}

// ToHash returns a hash Hash for a given []byte. Hash is a uint32 that is associated with the text in []byte. It returns zero if no match found.
func ToHash(s []byte) Hash {
	if 3 < len(s) {
		return HashMap[string(s)]
	}
	h := uint32(_Hash_hash0)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	if i := _Hash_table[h&uint32(len(_Hash_table)-1)]; int(i&0xff) == len(s) {
		t := _Hash_text[i>>8 : i>>8+i&0xff]
		for i := 0; i < len(s); i++ {
			if t[i] != s[i] {
				goto NEXT
			}
		}
		return i
	}
NEXT:
	if i := _Hash_table[(h>>16)&uint32(len(_Hash_table)-1)]; int(i&0xff) == len(s) {
		t := _Hash_text[i>>8 : i>>8+i&0xff]
		for i := 0; i < len(s); i++ {
			if t[i] != s[i] {
				return 0
			}
		}
		return i
	}
	return 0
}

const _Hash_hash0 = 0xf0c5341e
const _Hash_maxLen = 7

var _Hash_text = []byte("" +
	"arcsinharctanhcbrterfexpgammarccoshlnlog10log2sqrt")

var _Hash_table = [1 << 5]Hash{
	0x0:  0x7,    // arcsinh
	0x1:  0x1f04, // cosh
	0x3:  0xa04,  // tanh
	0x4:  0x2302, // ln
	0x5:  0x1503, // exp
	0x7:  0x1805, // gamma
	0x8:  0x6,    // arcsin
	0x9:  0x2505, // log10
	0xd:  0x1c06, // arccos
	0xe:  0x2e04, // sqrt
	0xf:  0x707,  // arctanh
	0x11: 0x1f03, // cos
	0x12: 0x2a04, // log2
	0x14: 0x2503, // log
	0x16: 0x1203, // erf
	0x19: 0xa03,  // tan
	0x1b: 0x303,  // sin
	0x1c: 0x304,  // sinh
	0x1d: 0x706,  // arctan
	0x1e: 0xe04,  // cbrt
	0x1f: 0x1c07, // arccosh
}
