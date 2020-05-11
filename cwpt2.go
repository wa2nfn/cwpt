//
// Copyright 2019, 2020 Bill Lanahan -  WA2NFN
//
// Only my second GO program. Lots could be more GO like
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"runtime"
)

const (
	prog          = "cwpt2"
	program       = "cwpt2.exe"
	version       = "1.3.0 04/25/2020"
	maxWordLen    = 40
	maxUserWords  = 10000
	maxLineLen    = 500
	maxSuffix     = 20
	maxPrefix     = 20
	maxDelimChars = 20
	maxRepeat     = 20
	maxSkips      = 5000
	maxMixedMode  = 20
	inListStr     = "A-Za-z%C0%E0%C4%E4%C9%E9%C8%E8%C7%E7%D1%F1%D6%F6%DC%FC"
	wordCount	= 5
	//inListStr     = "A-Za-z\u00C0\u00E0\u00C4\u00E4\u00C9\u00E9\u00C8\u00E8\u00C7\u00E7\u00D1\u00F1\u00D6\u00F6\u00DC\u00FC"
)

var (
	seed           = time.Now().UTC().UnixNano()
	rng            = rand.New(rand.NewSource(seed))
	wordMap        = make(map[string]struct{})
	wordArray      = make([]string, 0, 0)
	delimiterSlice []string
	effDelta       int
	proSign        []string
	runeMap        = make(map[rune]struct{})
	runeMapInt     = make(map[string]rune)
)

var (
	flagmax         int
	flagcgmax       int
	flaglen         int
	flagmin         int
	flagcgmin       int
	flagrepeat      int
	flagnum         int
	flagskip        int
	flagsuf         int
	flagpre         int
	flagDM          int
	flagDR          bool
	flaglesson      int
	flagMixedMode   int
	flagwordcount	int
	flagEBsf        string
	flagEBfs        string
	flagEBstep      int
	flagEBnum       int
	flagEBslow      int
	flagEBlow       int
	flagEBfast      int
	flagEBrepeat    int
	flagEBeff       int
	flagEBramp      bool
	flagEBefframp   bool
	flagheader      string
	flagprelist     string
	flagPrelistRune []rune
	flagsuflist     string
	flagSuflistRune []rune
	flaginlist      string
	flagcglist      string
	flagCglistRune  []rune
	flaginput       string
	flagoutput      string
	flagopt         string
	flagprosign     string
	flagdelimit     string
	flagtutor       string
	flagcaps        bool
	flagversion     bool
	flagrandom      bool
	flagunique      bool
	flagNR          bool
	flagMMR         bool
	flagCG          bool
	flagreverse     bool
	flaghelp      string
)

func init() {
	flag.IntVar(&flagmax, "max", 10, "Maximum # characters in a word >= min. (Default 10)")
	flag.IntVar(&flagcgmax, "cgmax", 5, "Maximum # characters in a code group >= cgmin. (Default 5)")
	flag.IntVar(&flagmin, "min", 1, "Minimum # of characters in a word (or code group). (Default 1)")
	flag.IntVar(&flagcgmin, "cgmin", 5, "Minimum # of characters in a code group. (Default 5)")
	flag.IntVar(&flagrepeat, "repeat", 1, "Number of times to repeat word sequentially. (Default 1)")
	flag.IntVar(&flagnum, "num", 100, fmt.Sprintf("Number of words (or code groups) to output. Min 1, max %d.\n", maxUserWords))
	flag.IntVar(&flaglen, "len", 80, fmt.Sprintf("Length characters in an output line (max %d).", maxLineLen))
	flag.IntVar(&flagskip, "skip", 0, fmt.Sprintf("Number of the first unique words in the input to skip. Max %d", maxSkips))
	flag.IntVar(&flagsuf, "suffix", 0, "The max number of suffix characters to append to words.")
	flag.IntVar(&flagpre, "prefix", 0, "The max number of prefix characters to affix to words.")
	flag.BoolVar(&flagcaps, "caps", false, "Print output in all capitals. (default lower case)")
	flag.BoolVar(&flagversion, "version", false, "Display version information. (default false)")
	flag.BoolVar(&flagrandom, "random", false, "If prefix/suffix is used, will determine if either is used on a\nword-by-word basis. (default false)")
	flag.StringVar(&flagsuflist, "suflist", "0-9,.?/=", "Characters to append to a word. Suffix X, sets the quantity.")
	flag.StringVar(&flagprelist, "prelist", "0-9,.?/=", "Characters to insert before a word. Prefix X, sets the quantity.")
	flag.StringVar(&flaginlist, "inlist", inListStr, "Set of characters to define an input word.")
	flag.StringVar(&flaginput, "in", "", "Input text file name (including extension).")
	flag.StringVar(&flagoutput, "out", "", "Output file name.")
	flag.StringVar(&flagopt, "opt", "", "Specify an options file name")
	flag.StringVar(&flagprosign, "prosign", "", "ProSign file name. 1-4 TWO letter ProSigns per line.\n No space in between, as in \"<BT> <AR>\".\n<SOS> is the only 3 letter ProSign.")
	flag.StringVar(&flagdelimit, "delimiter", "", "Output an inter-word delimiter string. A \"^\" separates delimiters e.g. <SK>^abc^123.\nA blank field e.g. aa^ ^bb, is valid to NOT get a delimiter. (default \"\"). ")
	flag.BoolVar(&flagunique, "unique", false, "Each output word is sent only once (num option quantity may be reduced).\n (default false)")
	flag.StringVar(&flagtutor, "tutor", "LCWO", "Only if you use -lessons. Sets order and # of charactersby tutor type.\nChoices: (default LCWO), JustLearnMorseCode, G4FON, MorseElmer, MorseCodeNinja, HamMorse, LockdownMorse\nUse -help=tutors for more info.")
	flag.IntVar(&flagDM, "DM", 0, fmt.Sprintf("Delimiter multiple, (if delimiter is used.) Between 1 and DM delimiter\nstrings are concatenated. (min 0, max %d)", maxDelimChars))
	flag.IntVar(&flaglesson, "lesson", 0, "Given the Koch lesson number per LCWO, populates options inlist and cglist with appropriate characters. (Default 0)")
	flag.BoolVar(&flagDR, "DR", false, "Delimiter randomness, (if DM > 0) DR=true makes a delimiter randomly print on an instance-by-instance basis")
	flag.IntVar(&flagMixedMode, "mixedMode", 0, fmt.Sprintf("mixedMode X, If X gt 1 & le %d, a code group will print every X words.", maxMixedMode))
	flag.BoolVar(&flagreverse, "reverse", false, "Reverses the spelling of words from inlist file (ignored for codeGroups_. (default false)")
	flag.BoolVar(&flagCG, "codeGroups", false, "Random code groups from cglist characters.")
	flag.BoolVar(&flagNR, "NR", false, "Non-Randomized output words read from input.")
	flag.BoolVar(&flagMMR, "MMR", false, "Mixed-Mode-Random, randomizes the output of a code group in mixed mode.")
	flag.StringVar(&flagcglist, "cglist", "a-z0-9.,?/=", "Set of characters to make random code groups.")
	flag.StringVar(&flagheader, "header", "", "string copied verbatim to head of output")
	flag.IntVar(&flagEBlow, "EB_LOW", 15, "ebook2cw low character speed wpm setting.")
	flag.IntVar(&flagEBstep, "EB_STEP", 5, "ebook2cw wpm and/or effectie speed change increment.")
	flag.IntVar(&flagEBslow, "EB_SLOW", 0, "ebook2cw number of words to send at slower speed.")
	flag.IntVar(&flagEBfast, "EB_FAST", 0, "ebook2cw number of words to send at faster speed.")
	flag.IntVar(&flagEBnum, "EB_NUM", 0, "ebook2cw number of speed change steps.")
	flag.BoolVar(&flagEBramp, "EB_RAMP", false, "ebook2cw ramps speed up in steps (default false).")
	flag.IntVar(&flagEBrepeat, "EB_REPEAT", 0, "ebook2cw times to repeat each word with increasing speed.")
	flag.IntVar(&flagEBeff, "EB_EFFECTIVE", 0, "ebook2cw effective (aka Farnsworth) speed must be < EB_LOW.")
	flag.BoolVar(&flagEBefframp, "EB_EFFECTIVE_RAMP", false, "ebook2cw ramp effective speed (char speed constant) must be < EB_LOW.")
	flag.StringVar(&flagEBsf, "EB_SF", "", "to alert transition from EB_LOW to EB_LOW+EB_STEP for plain text in mixedMode\nor EB_SLOW text to EB_FAST text,")
	flag.StringVar(&flagEBfs, "EB_FS", "", "to alert transition from EB_LOW+EB_STEP speed for plain text to EB_LOW for codeGroup mixedMode\nor EB_FAST text to EB_SLOW text.")
	flag.StringVar(&flaghelp, "help", "", "[EBOOK|INTERNATIONAL|TUTORS] more help of given topics.")
	flag.IntVar(&flagwordcount, "wordCount", 0, "Number of words to link as a phrase IF <repeat> option is also used.(Max 5)")

	// fill the rune map which is used to validate option string like: cglist, prelist, delimiter

	runeMap['a'] = struct{}{}
	runeMap['b'] = struct{}{}
	runeMap['c'] = struct{}{}
	runeMap['d'] = struct{}{}
	runeMap['e'] = struct{}{}
	runeMap['f'] = struct{}{}
	runeMap['g'] = struct{}{}
	runeMap['h'] = struct{}{}
	runeMap['i'] = struct{}{}
	runeMap['j'] = struct{}{}
	runeMap['k'] = struct{}{}
	runeMap['l'] = struct{}{}
	runeMap['m'] = struct{}{}
	runeMap['n'] = struct{}{}
	runeMap['o'] = struct{}{}
	runeMap['p'] = struct{}{}
	runeMap['q'] = struct{}{}
	runeMap['r'] = struct{}{}
	runeMap['s'] = struct{}{}
	runeMap['t'] = struct{}{}
	runeMap['u'] = struct{}{}
	runeMap['v'] = struct{}{}
	runeMap['w'] = struct{}{}
	runeMap['x'] = struct{}{}
	runeMap['y'] = struct{}{}
	runeMap['z'] = struct{}{}
	runeMap['A'] = struct{}{}
	runeMap['B'] = struct{}{}
	runeMap['C'] = struct{}{}
	runeMap['D'] = struct{}{}
	runeMap['E'] = struct{}{}
	runeMap['F'] = struct{}{}
	runeMap['G'] = struct{}{}
	runeMap['H'] = struct{}{}
	runeMap['I'] = struct{}{}
	runeMap['J'] = struct{}{}
	runeMap['K'] = struct{}{}
	runeMap['L'] = struct{}{}
	runeMap['M'] = struct{}{}
	runeMap['N'] = struct{}{}
	runeMap['O'] = struct{}{}
	runeMap['P'] = struct{}{}
	runeMap['Q'] = struct{}{}
	runeMap['R'] = struct{}{}
	runeMap['S'] = struct{}{}
	runeMap['T'] = struct{}{}
	runeMap['U'] = struct{}{}
	runeMap['V'] = struct{}{}
	runeMap['W'] = struct{}{}
	runeMap['X'] = struct{}{}
	runeMap['Y'] = struct{}{}
	runeMap['Z'] = struct{}{}
	runeMap['0'] = struct{}{}
	runeMap['1'] = struct{}{}
	runeMap['2'] = struct{}{}
	runeMap['3'] = struct{}{}
	runeMap['4'] = struct{}{}
	runeMap['5'] = struct{}{}
	runeMap['6'] = struct{}{}
	runeMap['7'] = struct{}{}
	runeMap['8'] = struct{}{}
	runeMap['9'] = struct{}{}
	runeMap[','] = struct{}{}
	runeMap['.'] = struct{}{}
	runeMap['/'] = struct{}{}
	runeMap['?'] = struct{}{}
	runeMap['='] = struct{}{}
	runeMap['+'] = struct{}{}
	runeMap['!'] = struct{}{}      // added at bottom of LCWO
	runeMap['"'] = struct{}{}      // added at bottom of LCWO
	runeMap['\''] = struct{}{}     // added at bottom of LCWO
	runeMap['('] = struct{}{}      // added at bottom of LCWO
	runeMap[')'] = struct{}{}      // added at bottom of LCWO
	runeMap['-'] = struct{}{}      // added at bottom of LCWO
	runeMap[':'] = struct{}{}      // added at bottom of LCWO
	runeMap[';'] = struct{}{}      // added at bottom of LCWO
	runeMap['\u00C0'] = struct{}{} // cap A grave
	runeMap['\u00E0'] = struct{}{} // low a grave
	runeMap['\u00C4'] = struct{}{} // cap A diaeresis
	runeMap['\u00E4'] = struct{}{} // low a diaeresis
	runeMap['\u00C9'] = struct{}{} // cap E acute
	runeMap['\u00E9'] = struct{}{} // low e acute
	runeMap['\u00C8'] = struct{}{} // cap E grave
	runeMap['\u00E8'] = struct{}{} // cap E acute
	runeMap['\u00C7'] = struct{}{} // cap C cedilla
	runeMap['\u00E7'] = struct{}{} // low c cedilla
	runeMap['\u00D1'] = struct{}{} // cap N tilde
	runeMap['\u00F1'] = struct{}{} // low n tilde
	runeMap['\u00D6'] = struct{}{} // cap O diaeresis
	runeMap['\u00F6'] = struct{}{} // low o diaeresis
	runeMap['\u00DC'] = struct{}{} // cap U diaeresis
	runeMap['\u00FC'] = struct{}{} // low u diaeresis
	runeMap['*'] = struct{}{} // DUMMY value for delimiter and ebook users

	runeMapInt["C0"] = '\u00C0' // cap A grave
	runeMapInt["E0"] = '\u00E0' // low a grave
	runeMapInt["C4"] = '\u00C4' // cap A diaeresis
	runeMapInt["E4"] = '\u00E4' // low a diaeresis
	runeMapInt["C9"] = '\u00C9' // cap E acute
	runeMapInt["E9"] = '\u00E9' // low e acute
	runeMapInt["C8"] = '\u00C8' // cap E grave
	runeMapInt["E8"] = '\u00E8' // cap E acute
	runeMapInt["C7"] = '\u00C7' // cap C cedilla
	runeMapInt["E7"] = '\u00E7' // low c cedilla
	runeMapInt["D1"] = '\u00D1' // cap N tilde
	runeMapInt["F1"] = '\u00F1' // low n tilde
	runeMapInt["D6"] = '\u00D6' // cap O diaeresis
	runeMapInt["F6"] = '\u00F6' // low o diaeresis
	runeMapInt["DC"] = '\u00DC' // cap U diaeresis
	runeMapInt["FC"] = '\u00FC' // low u diaeresis
}

func main() {
	kochChars := "kmuresnaptlwi.jz=foy,vg5/q92h38b?47c1d60x" // default for LCWO
	var fp *os.File
	localSkipFlag := false
	localSkipCount := 0

	if strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe") != prog {
		fmt.Printf("\nThe executable must be named: %s or %s\n73\nBill, WA2NFN\n", prog, program)
		os.Exit(1)
	}

	flag.Parse() // first parse to see if we had -opt

	if flagopt != "" {
		_, err := os.Stat(flagopt)
		if os.IsNotExist(err) {
			fmt.Printf("\nError: Can't find options file=<%s>.\n", flagopt)
			os.Exit(1)
		}

		optFile, err := os.Open(flagopt)
		if err != nil {
			fmt.Printf("\n%s File name <%s>.\n", err, flagopt)
			os.Exit(1)
		}

		// do file parse
		doOptFile(optFile)
		optFile.Close()

		flag.Parse() // second parse since options read
	}

	//
	// verify valid options
	//
	if flag.NArg() > 0 {
		fmt.Printf("\nError processing the command line.\n\nYou may have:\n   forgotten a \"-\" before an option\n   or followed a \"-\" with a space\n   or added extra input\n   or put spaces arount the \"=\"\n")
		os.Exit(1)
	}

	if flagversion {
		fmt.Printf("\n%s version: %s\nCopyright 2019, 2020", program, version)
		os.Exit(0)
	}

	if flaghelp != "" {
		flaghelp = strings.ToUpper(flaghelp)

		if flaghelp == "EBOOK" {
			fmt.Println("\nEBOOK Help Info")
			fmt.Printf("\nRelationships and compatibilities for EB_ Options\n\nOption names starting with EB_ are only for users that will input the generated practice text\nto LCWO's \"Convert text to CW\" screen, or for input to ebook2cw (see Note).\n\nOther non-EB_ options can be used as well.\n\nEB_SLOW: (alternating slow/fast word groups)\n\t\trequires: EB_SLOW, EB_LOW, EB_FAST, EB_STEP\n\t\toptional: EB_SF, EB_FS\n\t\tnon-compatible: EB_RAMP, EB_EFFECTIVE_RAMP, EB_NUM, EB_REPEAT\n")

			fmt.Printf("\nEB_RAMP: (steady increase of character speed sections of words)\n\t\trequires: EB_RAMP, EB_LOW, EB_STEP, EB_NUM\n\t\toptional: EB_SF, EB_REPEAT, EB_EFFECTIVE\n\t\tnon-compatible: EB_EFFECTIVE_RAMP\n")

			fmt.Printf("\nEB_EFFECTIVE_RAMP: (steady increase of effective speed sections, with fixed character speed)\n\t\trequires: EB_EFFECTIVE_RAMP, EB_LOW, EB_STEP, EB_NUM, EB_EFFECTIVE\n\t\toptional: EB_SF\n\t\tnon-compatible: EB_RAMP, EB_REPEAT\n")

			fmt.Printf("\nEB_REPEAT: (each word repeated at increasing character speeds)\n\t\trequires: EB_REPEAT, EB_LOW, EB_STEP\n\t\toptional: EB_EFFECTIVE, EB_RAMP (EB_NUM if EB_RAMP)\n\t\tnon-compatible: EB_EFFECTIVE_RAMP\n")

			fmt.Printf("\nEB_LOW, EB_NUM, EB_STEP (no other EB_ options): (\"bounce\" each word has random character speed)\n\t\trequires: EB_LOW, EB_STEP, EB_NUM\n\t\toptional:\n\t\tnon-compatible: EB_EFFECTIVE_RAMP, EB_RAMP, EB_REPEAT, EB_SLOW, EB_FAST, EB_SF, EB_FS\n")

			fmt.Printf("\nNote: EB_SF and EB_FS MAY be used outside of LCWO/ebook2cw. If used with option mixedMode,\nthese will insert the specified string immediately before and/or after the codeGroup.\n")
			os.Exit(1)
		} else if flaghelp == "TUTORS" {
			fmt.Println("\nTUTORS Help Info")
			fmt.Printf(`
The options <tutor> and <lesson> are for user convience. By choosing the pair, you are prepopluating the
option <inlist> which reads words from the <in> file, and <cglist> which is used to create code groups.

The generated practice text therefore can be given to ANY tutor. In some cases, a tutor will teach a ProSign,
but these two options only function with single characters.

The option <inlist> will be populated with both upper and lower case for each alpha character.

The term "lesson" may not be used in each tutor, but its just the order that the character was taught.

Lesson is cummulative, that is if you enter lesson=5, all the characters from lesson 1 through 5 are used.

Lesson  LCWO  JustLearnMorseCode  G4FON  MorseElmer  MorseCodeNinja  HamMorse  LockdownMorse
------  ----  ------------------  -----  ----------  --------------  --------  -------------
1       km            k           k          t          k               e           e
2       u             m           m          m          a               m           o
3       r             r           r          r          e               r           a
4       e             s           s          s          n               s           i
5       s             u           u          u          o               u           u
6       n             a           a          a          i               a           y
7       a             p           p          p          s               p           z
8       P             t           t          t          l               t           q
9       t             l           l          l          4               l           j
10      l             o           o          o          r               o           x

Lesson  LCWO  JustLearnMorseCode  G4FON  MorseElmer  MorseCodeNinja  HamMorse  LockdownMorse
------  ----  ------------------  -----  ----------  --------------  --------  -------------
11      w            w            w          w          h              w            k
12      i            i            i          i          d              i            v
13      .            .            .          .          1              .            b
14      j            n            n          n          2              n            p
15      z            j            j          j          5              j            + <AR> (1)
16      = <BT>       e            e          e          c              g            g
17      f            f            f          f          u              f            w
18      o            0            0          0          m              0            f
19      y            y            y          y          w              y            c
20      ,            v            v          ,          3              ,            l

Lesson  LCWO  JustLearnMorseCode  G4FON  MorseElmer  MorseCodeNinja  HamMorse  LockdownMorse
------  ----  ------------------  -----  ----------  --------------  --------  -------------
21      v            ,            ,          v          6              v            d
22      g            g            g          g          ?              g            m
23      5            5            5          5          f              5            h
24      /            /            /          /          y              /            r
25      q            q            q          q          p              q            s
26      9            9            9          9          g              9            n
27      2            z            z          z          7              z            t
28      h            h            h          h          9              h
29      3            3            3          3          /              3
30      8            8            8          8          b              8

Lesson  LCWO  JustLearnMorseCode  G4FON  MorseElmer  MorseCodeNinja  HamMorse LockdownMorse
------  ----  ------------------  -----  ----------  --------------  -------- -------------
31      b            b            b          b          v              b
32      ?            ?            ?          ?          k              ?
33      4            4            4          4          j              4
34      7            2            2          2          8              2 
35      c            7            7          7          0              7
36      1            c            c          c          x              c
37      d            1            1          1          q              1
38      6            d            d          d          z              6
39      0            6            6          6          =              x
40      x            x            x          x          .              =

Lesson  LCWO  JustLearnMorseCode  G4FON  MorseElmer  MorseCodeNinja HamMorse LockdownMorse
------  ----  ------------------  -----  ----------  -------------- -------- -------------
41                   @            = <BT>                + <AR> (2) 
42                   = <BT>       + <AR> (2)
43                   + <AR> 

Notes: 
1- <KA> also introduced but must be handled with the prosign option
2- <SK> (same as above)
			`)  // end of table

			os.Exit(1)
		} else if flaghelp == "INTERNATIONAL" {
			fmt.Println("\nINTERNATIONAL Characters  Help Info")
			fmt.Printf(`
If the code tutor you intend to use knows how to send morse for international characters, you can specify them
in the "list" options. All of the these are already included in the option <inlist>, used to find words in
your <in> file. They are in both upper and lower case.

You may include them in the following options: prelist, suflist, cglist, or delimiter. To do that, they MUST
be entered as %%XX, as shown below. No spaces are used, and the alphas MUST be uppercase. You only need to
inlude the leter you want in EITHER upper or lower case (cwpt2 will make the case correct based on the <caps>
option. To include the "grave a" and "acute e" to english vowels for use as word suffixes, as an example,
you would set the option like: suflist="aeiou%%C0%%C9".

%%C0 to represent '\u00C0'  Upper Case A grave
%%E0 to represent '\u00E0'  Lower Case a grave
%%C4 to represent '\u00C4'  Upper Case A diaeresis
%%E4 to represent '\u00E4'  Lower Case a diaeresis
%%C9 to represent '\u00C9'  Upper Case E acute
%%E9 to represent '\u00E9'  Lower Case e acute
%%C8 to represent '\u00C8'  Upper Case E grave
%%E8 to represent '\u00E8'  Upper Case E acute
%%C7 to represent '\u00C7'  Upper Case C cedilla
%%E7 to represent '\u00E7'  Lower Case c cedilla
%%D1 to represent '\u00D1'  Upper Case N tilde
%%F1 to represent '\u00F1'  Lower Case n tilde
%%D6 to represent '\u00D6'  Upper Case O diaeresis
%%F6 to represent '\u00F6'  Lower Case o diaeresis
%%DC to represent '\u00DC'  Upper Case U diaeresis
%%FC to represent '\u00FC'  Lower Case u diaeresis

`)

			os.Exit(1)
		} else {
			fmt.Printf("\nError: Invalid value for option <help>, choices are (case insensitive): TUTOR, EBOOK, or INTERNATIONAL.\n")
			os.Exit(1)
		}
	}

	if flagNR == false {
		wordArray = nil // save space we're using the map not array
	}

	//
	// out of range checks
	if flagDM < 0 || flagDM > maxDelimChars {
		fmt.Printf("\nError: DM, delimiter multiple min >=0, max <= %d.\n", maxDelimChars)
		os.Exit(1)
	} else if flagDM >= 1 {
		// ok DM is in range

		// split into fields if any

		// if prosigns are in a delimiter field
		runeMap['<'] = struct{}{}
		runeMap['>'] = struct{}{}
		runeMap[' '] = struct{}{}
		runeMap['*'] = struct{}{}

		// first make sure any prosign is valid format
		m := regexp.MustCompile(`<[a-zA-Z]{2}>`)
		tStr := flagdelimit

		if m.MatchString(tStr) {
			tStr = m.ReplaceAllString(tStr, "")
		}

		if strings.Contains(tStr,"<") || strings.Contains(tStr, ">" ) {
			fmt.Printf("\nError: option <delimiter> contains invalid prosign format.\n")
			os.Exit(77)
		}

		for _, field := range strings.Split(flagdelimit, "^") {
			processDelimiter(field)
		}

		delete(runeMap, ' ')
		delete(runeMap, '<')
		delete(runeMap, '>')
		delete(runeMap, '*')
	}

	if flagmin < 1 {
		fmt.Printf("\nError: min must be >= 1.\n")
		os.Exit(1)
	}

	if flagmin > flagmax {
		fmt.Printf("\nError: min must <= max <%d>.\n", flagmax)
		os.Exit(1)
	}

	if flagmax < flagmin {
		fmt.Printf("\nError: max must >= min <%d>.\n", flagmin)
		os.Exit(1)
	}

	if flagmax > maxWordLen {
		fmt.Printf("\nError: max must <= <%d>, system max.\n", maxWordLen)
		os.Exit(1)
	}

	if flagcgmax < flagcgmin {
		fmt.Println("\nError: cgmax must >= cgmin.")
		os.Exit(1)
	}

	if flagMixedMode < 0 || flagMixedMode == 1 || flagMixedMode > maxMixedMode {
		fmt.Printf("\nError: mixedMode X Where X  minimum 2, maximum %d, default 0=off.\n", maxMixedMode)
		os.Exit(1)
	}

	if flagskip < 0 || flagskip > maxSkips {
		fmt.Printf("\nError: skip x  minimum 0, maximum %d, default 0.\n", maxSkips)
		os.Exit(1)
	}

	if flagskip >= 1 {
		// we will be skipping some words
		localSkipFlag = true
		localSkipCount = flagskip
	}

	if flagnum < 1 || flagnum > maxUserWords {
		fmt.Printf("\nError: num, number of output words desired. minimum 1, maximum %d, default 100.\n", maxUserWords)
		os.Exit(1)
	}

	if flaglen < 1 || flaglen > maxUserWords {
		fmt.Printf("\nError: len max output line length, default 80, maximum %d.\n", maxLineLen)
		os.Exit(1)
	}

	if flagsuf < 0 || flagsuf > maxSuffix {
		fmt.Printf("\nError: suffix, 0=no suffix, max number of characters is %d.\n", maxSuffix)
		os.Exit(1)
	}

	if flagMixedMode > 0 && flagCG == true {
		fmt.Printf("\nError: mixedMode is mutually exclusive with codeGroups option.\n")
		os.Exit(1)
	}

	if flagsuf > 0 {
		if flagsuflist != "" {
			flagSuflistRune = ckValidInString(flagsuflist, "suflist")
		} else {
			fmt.Printf("\nError: if suffix > 0, the suflist must contain characters, its empty.\n")
			os.Exit(1)
		}
	}

	if flagpre < 0 || flagpre > maxPrefix {
		fmt.Printf("\nError: prefix, 0=no prefix, max number of characters is %d.\n", maxPrefix)
		os.Exit(1)
	}

	if flagpre > 0 {
		if flagprelist != "" {
			// return expanded
			flagPrelistRune = ckValidInString(flagprelist, "prelist")
		} else {
			fmt.Printf("\nError: if prefix > 0, the prelist must contain characters, its empty.\n")
			os.Exit(1)
		}
	}

	if flagrandom && (flagsuf == 0 && flagpre == 0) {
		fmt.Printf("\nError: random requires either prefix > 0 or suffix > 0.\n")
		os.Exit(1)
	}

	if flagrepeat < 1 || flagrepeat > maxRepeat {
		fmt.Printf("\nError: repeat (default 1) must be between 1 and %d.\n", maxRepeat)
		os.Exit(1)
	}

	if flagwordcount < 0 || flagwordcount > 5 {
		fmt.Printf("\nError: wordCount (default 0) must be between 1 and %d, to link words as a phrase.\n", wordCount)
		os.Exit(1)
	}

	if flagwordcount >= 1 && flagrepeat < 2 {
		fmt.Printf("\nError: wordCount requires the repeat option >= 2\n")
		os.Exit(1)
	}

	if flagNR == true && (flagunique || flagCG) {
		fmt.Printf("\nError: NR is mutually exclusive with unique and codeGroups options.\n")
		os.Exit(1)
	}

	if flagoutput != "" {
		if flagoutput == flaginput {
			fmt.Printf("\nError: -out can't equal -in, or the input file would be over written.\n")
			os.Exit(1)
		}

		// check for existance first
		_, err := os.Stat(flagoutput)
		if err == nil {
			fmt.Printf("\nWarning: out file: <%s> exists!\n\nEnter \"y\" to overwrite it: ", flagoutput)
			ans := ""
			fmt.Scanf("%s", &ans)
			if ans != "y" {
				fmt.Printf("\nNo output as requested.\n")
				os.Exit(0)
			}
		}

		fp, err = os.Create(flagoutput)
		if err != nil {
			fmt.Println(err)
			os.Exit(9)
		}
		fmt.Printf("\nWriting to file: %s\n", flagoutput)
	}

	// ebook options
	// hard code some values since they are arbitrary

	if flagEBnum < 0 || flagEBnum > 30 {
		fmt.Printf("\nError: EB_NUM number of speed values must be >= 0 and <= 30.\n")
		os.Exit(0)
	}

	if flagEBnum > 0 {
		if flagEBslow > 0 || flagEBfast > 0 {
			fmt.Printf("\nError: EB_NUM must = 0(off) if EB_SLOW > 0.\n")
			os.Exit(0)
		}

		if flagEBrepeat > 1 && !flagEBramp {
			fmt.Printf("\nError: EB_REPEAT is mutually exclusive with EB_NUM unless also with EB_RAMP.\n")
			os.Exit(0)
		}
	}

	if flagEBlow < 5 {
		fmt.Printf("\nError: EB_LOW lowest speed must be at least 5 wpm.\n")
		os.Exit(0)
	}

	if flagEBeff > 0 {
		if flagEBeff >= flagEBlow {
			fmt.Printf("\nError: EB_EFFECTIVE speed must be < EB_LOW wpm.\n")
			os.Exit(0)
		}

		// set delta since eblow and eff are set
		effDelta = flagEBlow - flagEBeff
	}

	if flagEBstep < 0 || flagEBstep > 20 {
		fmt.Printf("\nError: EB_STEP speed incremental step must be >= 0 and <= 30 wpm, 0(off).\n")
		os.Exit(0)
	}

	if flagEBslow < 0 || flagEBfast < 0 {
		fmt.Printf("\nError: EB_FAST and EB_SLOW must be >= 0, 0(off).\n")
		os.Exit(0)
	}

	if flagEBfast > 0 && flagEBslow == 0 {
		fmt.Printf("\nError: EB_FAST > 0 requires EB_SLOW to be specified.\n")
		os.Exit(0)
	}

	// we want EB, lots of exclusions to try
	if flagEBslow > 0 {
		if flagEBfast == 0 {
			flagEBfast = flagEBslow
		}

		if flagEBstep < 1 {
			fmt.Printf("\nError: EB_STEP must be >=1 with EB_SLOW/EB_FAST.\n")
			os.Exit(0)
		}

		if flagEBrepeat >= 1 {
			fmt.Printf("\nError: EB_REPEAT is mutually exclusive with EB_SLOW and EB_FAST options.\n")
			os.Exit(1)
		}
	}

	if flagEBrepeat < 0 || flagEBrepeat > 30 {
		fmt.Printf("\nError: EB_REPEAT must be >=2 and <= 30 for word speed repeat, 0(off).\n")
		os.Exit(0)
	}

	if flagEBramp {
		if flagEBfast > 0 {
			fmt.Printf("\nError: EB_RAMP is mutually exclusive with EB_SLOW and EB_FAST options.\n")
			os.Exit(1)
		}

		if flagEBnum == 0 {
			fmt.Printf("\nError: EB_RAMP requires EB_NUM > 0.\n")
			os.Exit(1)
		}

		if flagEBstep == 0 {
			fmt.Printf("\nError: EB_RAMP requires EB_STEP > 0.\n")
			os.Exit(1)
		}
	}

	if flagEBefframp {
		if flagEBrepeat > 0 {
			fmt.Printf("\nError: EB_EFFECTIVE_RAMP is mutually exclusive with EB_REPEAT.\n")
			os.Exit(1)
		}

		if flagEBnum == 0 {
			fmt.Printf("\nError: EB_EFFECTIVE_RAMP requires EB_NUM > 0.\n")
			os.Exit(1)
		}

		if flagEBstep < 1 {
			fmt.Printf("\nError: EB_EFFECTIVE_RAMP requires EB_STEP >= 1 and its less than EB_LOW.\n")
			os.Exit(1)
		}

		if flagEBramp {
			fmt.Printf("\nError: EB_EFFECTIVE_RAMP is mutually exclusive with EB_RAMP.\n")
			os.Exit(1)
		}

		if flagEBlow < 1 {
			fmt.Printf("\nError: EB_EFFECTIVE_RAMP requires EB_LOW >= 5.\n")
			os.Exit(1)
		}
	}

	s := 0

	flagtutor = strings.ToUpper(flagtutor)

	if flaglesson == 0 && flagtutor != "LCWO" {
		fmt.Printf("\nError: Lesson = 0 is invalid for tutor <%s>.\n", flagtutor)
		os.Exit(1)
	}

	if flaglesson >= 1 {

		if flagtutor == "LCWO" {
			kochChars = "kmuresnaptlwi.jz=foy,vg5/q92h38b?47c1d60x"
		} else if flagtutor == "JUSTLEARNMORSECODE" {
			kochChars = "kmrsuaptlowi.njef0yv,g5/q9zh38b?427c1d6x@=+"
			s = 0
		} else if flagtutor == "G4FON" {
			kochChars = "kmrsuaptlowi.njef0yv,g5/q9zh38b?427c1d6x"
			s = 0
		} else if flagtutor == "MORSEELMER" {
			kochChars = "kmrsuaptlowi.njef0y,vg5/q9zh38b?427c1d6x=+"
			s = 0
		} else if flagtutor == "MORSECODENINJA" {
			kochChars = "taenois14rhdl25cumw36?fypg79/bvkj80xqz=."
			s = 0
		} else if flagtutor == "HAMMORSE" {
			kochChars = "kmrsuaptlowi.njef0y,vg5/q9zh38b?427c1d6x=+"
			s = 0
		} else if flagtutor == "LOCKDOWNMORSE" {
			kochChars = "eoaiuyzqjxkvbp+gwfcldmhrsnt"
			s = 0
		} else {
			fmt.Printf("\nError: Your tutor name is invalid. Names are NOT case sensitive,  and without any spaces, see the help.\n")
			os.Exit(1)
		}

		if (flaglesson+1 > len(kochChars)) && flagtutor == "LCWO" {
			fmt.Printf("\nError: Lesson value <%d> exceeds the max <%d>, for tutor <%s>.\n", flaglesson, 40, flagtutor)
			os.Exit(1)
		}

		if flaglesson > len(kochChars) {
			os.Exit(1)
		}

		if flagtutor == "LCWO" {
			if flaglesson < len(kochChars) {
				flagcglist = kochChars[0 : flaglesson+1]
			} else {
				flagcglist = kochChars[0:flaglesson]
			}
		} else {
			flagcglist = kochChars[s:flaglesson]
		}

		// now build inlist initailly as lower case
		// cglist is LC
		flaginlist = flagcglist
		// now inlist is LC
		temp := strings.ToUpper(flagcglist)

		// now add upper case
		for _, char := range temp {
			if char >= 'A' && char <= 'Z' {
				flaginlist += string(char)
			}
		}
		// now inlist is mixed case
		temp = ""

		if flagcaps {
			flagcglist = strings.ToUpper(flagcglist)
		}
	}

	// check inlist for %XX codes
	if flaglesson == 0 {
		out := ""
		m := regexp.MustCompile("%[C-F][0146789C]")

		for {
			if m.MatchString(flaginlist) {
				s := m.FindString(flaginlist)
				out = strings.TrimLeft(s, "%")
				out = string(runeMapInt[out])
				flaginlist = strings.Replace(flaginlist, string(s), out, 1)
				out = ""
			} else {

				break
			}
		}
	}

	// must follow other cglist manipulation
	// either case lets get cglist expanded now
	if flagCG || flagMixedMode > 0 {
		// make sure we have chars to work with
		if len(flagcglist) < 2 {
			fmt.Printf("Error: you requested codeGroups or mixedMode, so cglist must have at least 2 characters.\n")
			os.Exit(1)
		} else {
			if flagcglist != "" {
				// return expanded
				flagCglistRune = ckValidInString(flagcglist, "cglist")
			}
		}
	}


	// no longer needed save space
	runeMap = nil
	runeMapInt = nil

	//
	// major flow decision - WORD_MODE or CODE_GROUPS ?
	//
	if flagCG {
		makeGroups(fp)
	} else {
		readFileMode(localSkipFlag, localSkipCount, fp)
		doOutput(wordArray, fp)
	}
	os.Exit(0) // program done
}

//
// buildCharSlice - create a byte slice to use for codeGroups
//
func buildCharSlice() []rune {

	// make slice of chars for MAY NEED use later
	// if word mode, we only need for MAX possible codeGroups
	numChars := 0

	if flagMixedMode == 0 {
		numChars = flagcgmax * flagnum // may be extra
	} else {
		numChars = flagcgmax * (flagnum / flagMixedMode) // may be extra
	}

	charSlice := make([]rune, 0, numChars)
	cgSlice := flagCglistRune

	//charSlice = append(charSlice, cgSlice...)
	charSlice = cgSlice

	// charSlice now has the user given list of chars
	// then just copy cgSlice into charSlice as needed
	if len(cgSlice) < numChars {
		// flush out the charSlice to max we may need
		factor := numChars / len(cgSlice)
		factor-- // we have the original already

		// only does FULL slice
		for ; factor > 0; factor-- {
			charSlice = append(charSlice, cgSlice...)
		}

		// may still be a partial shortage
		howShort := numChars - len(charSlice)

		for _, key := range cgSlice {
			charSlice = append(charSlice, key)

			if howShort == 0 {
				break
			}
			howShort--
		}
	}

	return charSlice
}

// make random code groups
// uses the presaved chars in charSlice based on uniform distribution
func makeGroups(fp *os.File) {
	strBuf := ""
	var tmpOut []rune

	charSlice := buildCharSlice()

	// make the code groups
	for i := 0; i < flagnum; i++ {

		// tmpOut is our code group
		tmpOut, charSlice = makeSingleGroup(charSlice)

		// text repeat!
		if flagrepeat > 0 {
			// we need to repeat
			temp := tmpOut

			for cnt := 1; cnt < flagrepeat; cnt++ {
				// wordOut is the word plus trailing space already
				temp = append(temp, tmpOut...)
			}
			 strBuf += string(temp) 

		} else {
			// non repeat case
			strBuf += string(tmpOut)
		}

		if flagDM > 0 && (flagDR == false || (flagDR == true && flipFlop())) {
			for i := 1; i <= (1 + rng.Intn(flagDM)); i++ {
				strBuf += delimiterSlice[rng.Intn(len(delimiterSlice))]
			}
			strBuf += " "
		}
	}


	printStrBuf(strBuf, fp)
}

/*
** make a code group of random length
** character pulled from byte slice that even distribution
** of characters.
 */
func makeSingleGroup(charSlice []rune) ([]rune, []rune) {
	var cg []rune
	var tmp rune
	gl := flagcgmin

	// choose random grp len from min to max
	if flagcgmax != flagcgmin {
		gl = rng.Intn(flagcgmax-flagcgmin) + flagcgmin
	}

	for i := 0; i < gl; i++ {
		if len(charSlice) < gl {
			break
		}
		tmp, charSlice = getRandomChar(charSlice)
		cg = append(cg, tmp)
	}

	cg = append(cg, ' ')
	return cg, charSlice
}

// used for codeGroup
func getRandomChar(randCharSlice []rune) (rune, []rune) {
	sLen := len(randCharSlice)

	index := rng.Intn(sLen)
	newChar := randCharSlice[index] // to be returned

	// eat the value used
	sLen--
	randCharSlice[index] = randCharSlice[sLen]
	randCharSlice = randCharSlice[:sLen]

	return newChar, randCharSlice
}

// fill the array from the word map but might need to stuff more values
func fillArray(fp *os.File) {
	var wordArray = make([]string, 0, flagnum)

	for key := range wordMap {
		// make first population of slice
		wordArray = append(wordArray, key)
	}

	// see if initial array satisfies the number of words the user wanted
	// if less, we will reuse words from map to grow the array (or slice)
	if !flagunique && len(wordArray) < flagnum {

		howShort := flagnum - len(wordArray)
		factor := flagnum / len(wordArray)
		factor-- // we have the original already

		// only does FULL maps
		for ; factor > 0; factor-- {
			for key := range wordMap {
				wordArray = append(wordArray, key)
			}
		}

		// may still be a partial shortage
		howShort = flagnum - len(wordArray)
		i := 1
		for key := range wordMap {
			wordArray = append(wordArray, key)
			if i == howShort {
				break
			}
			i++
		}
	}

	// trash the map to conserve memory
	wordMap = nil
	doOutput(wordArray, fp)
}

// returns a single random prefix/suffix from list to add to output word
func ixStr(ps string) string {
	retStr := ""

	if ps == "s" {
		// user wants a suffix
		for count := 1; count <= rng.Intn(flagsuf)+1; count++ {
			retStr += string(flagSuflistRune[rng.Intn(len(flagSuflistRune))])
		}

	} else {
		// user wants a prefix
		for count := 1; count <= rng.Intn(flagpre)+1; count++ {
			retStr += string(flagPrelistRune[rng.Intn(len(flagPrelistRune))])
		}
	}

	return retStr
}

//
// make sure the string can be expanded into visable ASCII since all morse is limited to that
func ckValidInString(ck string, whoAmI string) []rune {
	str := strRangeExpand(ck, whoAmI)

	// check each rune to make sure its in the runeMap
	// also build a new string that is in the proper case
	newRune := []rune{}

	gotEscSymbol := false
	strInternational := ""

	for _, runeRead := range []rune(str) {
		if whoAmI != "delimiter" && runeRead == '*'  {
			fmt.Printf("\nError: Invalid character <%s>, in option <%s>.\nOnly used in delimiter option, as a special case delay for LCWO/ebook2cw users.\n", string(runeRead), whoAmI)
			os.Exit(98)
		}

		if runeRead == '%' && gotEscSymbol == false {

			// we potenially have an international
			gotEscSymbol = true
			continue
		}

		if gotEscSymbol == true && len(strInternational) < 2 {
			if runeRead == '%' {
				// two %s NG
				fmt.Printf("\nError: Invalid <%%>, in string <%s> for option <%s>, not followed by appropriate 2 upper case letters.\n", str, whoAmI)
				os.Exit(98)
			}

			// need to get two proper upper case
			strInternational += string(runeRead)
			if len(strInternational) < 2 {
				continue
			}

			runeRead = runeMapInt[strInternational]
			// reset these two
			gotEscSymbol = false
			strInternational = ""
		}

		if _, ok := runeMap[runeRead]; ok {

			if flagcaps {
				newRune = append(newRune, unicode.ToUpper(runeRead))
			} else {
				newRune = append(newRune, unicode.ToLower(runeRead))
			}
		} else {
			fmt.Printf("\nError: Invalid entry <%v>, in string <%s> for option <%s>.\n", string(runeRead), str, whoAmI)
			os.Exit(99)
		}
	}

	if whoAmI == "delimiter" {
		s := string(newRune)
		s = strings.ReplaceAll(s, "*", " |S500 ")
		delimiterSlice = append(delimiterSlice, s)
		newRune = nil
	}

	return newRune
}

//
// process the file of prosigns, check their validity
//
func doProSigns(file *os.File) {
	ps := ""

	scanner := bufio.NewScanner(file)

	word := regexp.MustCompile("^\\s*(<[A-Za-z][A-Za-z]>){1,4}\\s*$|^\\s*<[Ss][Oo][Ss]>{1,4}\\s*$")

	for scanner.Scan() {
		ps = strings.TrimSpace(scanner.Text())

		if word.MatchString(ps) {
			if flagcaps {
				ps = strings.ToUpper(ps)
			}

			if flagNR {
				proSign = append(proSign, ps)
			} else {
				// add to map if not there
				if _, ok := wordMap[ps]; ok != true {
					wordMap[ps] = struct{}{}
				}
			}

		} // ignore non matching ProSigns
	}
}

func doOptFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	ignore := regexp.MustCompile("^\\s*#|^\\s*$")
	doneEnd := regexp.MustCompile("^\\s*#\\s*(END|DONE)$")
	blkStart := regexp.MustCompile("^\\s*/\\*")
	blkEnd := regexp.MustCompile("^\\s*\\*")
	lineNum := 1
	inBlk := false

	for scanner.Scan() {
		str := scanner.Text()

		// start block comment
		if inBlk == false && blkStart.MatchString(str) {
			lineNum++
			inBlk = true
			continue
		}

		// end block comment
		if inBlk == true {
			lineNum++

			if blkEnd.MatchString(str) {
				inBlk = false
				continue
			} else {
				continue
			}
		}

		if doneEnd.MatchString(str) {
			return
		}

		if ignore.MatchString(str) {
			lineNum++
			continue
		}

		str = strings.TrimLeft(str, "-")

		// cuts EOL comments off
		dex := strings.Index(str, "#")
		if dex != -1 {
			str = str[:dex]
		}

		// we assume its an option
		// trim down to get the string at the end
		str = strings.TrimSpace(str)

		if str[len(str)-1] == '=' {
			fmt.Printf("\nError: Invalid format for option <%v> on line <%d> of file <%s>. Appears to be missing a value after \"=\".\n", str, lineNum, flagopt)
			os.Exit(7)

		}

		// = sep?
		arr := strings.SplitN(str, "=", 2)
		if len(arr) == 2 {
			if flag.Lookup(arr[0]) == nil {
				fmt.Printf("\nError: Invalid option <%s> on line <%d> of file <%s>.\n", arr[0], lineNum, flagopt)
				os.Exit(7)
			}

			arr[1] = strings.TrimLeft(arr[1], "'\"")
			arr[1] = strings.TrimRight(arr[1], "'\"")

			if arr[0] == "opt" {
				fmt.Printf("\nWarning: option \"opt\" can't be reset in the options file <%s> on line <%d>. Ignoring it and continuing.\n", flagopt, lineNum)
				continue
			}

			flag.Set(arr[0], arr[1])
			continue
		}

		// space sep?
		arr = strings.SplitN(str, " ", 2)
		if len(arr) == 2 {
			if flag.Lookup(arr[0]) == nil {
				fmt.Printf("\n2Error: Invalid option <%s> on line <%d> of file <%s>.\n", arr[0], lineNum, flagopt)
				os.Exit(7)
			}

			arr[1] = strings.TrimLeft(arr[1], "'\"")
			arr[1] = strings.TrimRight(arr[1], "'\"")

			if arr[0] == "opt" {
				fmt.Printf("\nWarning: option \"opt\" can't be reset in the options file <%s> on line <%d>. Ignoring it and continuing.\n", flagopt, lineNum)
				continue
			}

			flag.Set(arr[0], arr[1])
		} else if len(arr) == 1 {
			flag.Set(arr[0], "true")
			continue
		} else {
			if flag.Lookup(arr[0]) == nil {
				fmt.Printf("\nError: Invalid option <%s> on line <%d> of file <%s>.\n", arr[0], lineNum, flagopt)
				os.Exit(7)
			}
		}

		flag.Set(arr[0], "")

	}
}

/*
** walk the users string for prelist or suflist to see if we need to expand char ranges
 */
func strRangeExpand(inStr string, whoAmI string) string {
	outStr := ""
	last := ""
	gotDash := false

	// take care of special case of dash in beginning
	if inStr[0] == '-' {
		fmt.Printf("\nError: \"-\" character(s) at start of list option: %s=%s\n", whoAmI, inStr)
		os.Exit(7)
	}

	for _, char := range strings.Split(inStr, "") {

		if char == "-" {
			if gotDash == true {
				fmt.Printf("\nError: sequencial \"-\" characters in a list option: %s=%s\n", whoAmI, inStr)
				os.Exit(7)
			}

			gotDash = true
			continue
		}

		if gotDash == false {
			last = char
			outStr += last
		} else {
			outStr += expandIt(last, char, whoAmI)
			gotDash = false
			last = ""
		}
	}

	//  detect error
	if gotDash && len(inStr) > 1 {
		fmt.Printf("\nError: trailing \"-\" characters in a list option: %s=%s\n", whoAmI, inStr)
		os.Exit(7)
	}

	return outStr
}

// expandIt expands a char range into the individual chars
func expandIt(lower string, upper string, whoAmI string) string {
	outStr := ""

	low := rune(lower[0])
	up := rune(upper[0])

	if up < low {
		fmt.Printf("\nError: range in an option list is not in ASCII/UTF-8 order: i.e. C-A (invalid) vs. A-C (correct) \n")
		fmt.Printf("       Delimiters support ONLY a single range in a field. i.e. ^[A-D]^ or ^[0-3]^\n")
		os.Exit(7)
	}

	if whoAmI != "delimiter_simple" {
		low++ // we did low already for all other list options
	}

	for i := low; i <= up; i++ {
		if whoAmI == "delimiter_simple" {
			delimiterSlice = append(delimiterSlice, string(i))
		} else {
			outStr += string(i)
		}
	}
	return outStr
}

/*
** added for Ordered write of input "NR"
 */

// read input file and create words. vs do code groups
func readFileMode(localSkipFlag bool, localSkipCount int, fp *os.File) {
	done := false
	discarded := false

	if flaginput == "" {
		fmt.Printf("\nError: an input file must be given to -in, unless -codeGroups is used.\n")
		os.Exit(0)
	}

	if flaginlist == "" {
		fmt.Printf("\nError: inlist can't be empty or nothing gets matched.\n")
		os.Exit(0)
	}

	file, err := os.Open(flaginput)
	if err != nil {
		fmt.Printf("\n%s File name <%s>.\n", err, flaginput)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// to match what user wants
	s := fmt.Sprintf(`^[%s]{%d,%d}$|^(<[A-Za-z]{2}>){1,}$`, flaginlist, flagmin, flagmax)
	word := regexp.MustCompile(s)

	if !flagNR {
		wordMap = make(map[string]struct{})
	}

	// read ProSigns
	if flagprosign != "" {
		psfile, err := os.Open(flagprosign)
		if err != nil {
			fmt.Printf("\n%s File name <%s>.\n", err, flagprosign)
			os.Exit(1)
		}

		// fill the map with prosigns
		doProSigns(psfile)
		psfile.Close()
	}

	for scanner.Scan() {
		// first way to split the string on spaces

		textWords := strings.FieldsFunc(scanner.Text(), func(r rune) bool {
			if r == ' ' {
				return true
			}
			return false
		})

		for index := 0; done == false && index < len(textWords); index++ {
			// every token is now a string of space separated characters
			tmpWord := strings.TrimRight(textWords[index], ".\",?!")
			tmpWord = strings.TrimLeft(tmpWord, "\"")

			if word.MatchString(tmpWord) {

				// skip only viable matching words
				if localSkipFlag {
					if localSkipCount > 0 {
						localSkipCount--
						continue
					} else {
						localSkipFlag = false
					}
				}

				// set case before storing
				if flagcaps {
					tmpWord = strings.ToUpper(tmpWord)
				} else {
					tmpWord = strings.ToLower(tmpWord)
				}

				// reverse the string
				if flagreverse {
					tmpWord = reverse(tmpWord)
				}

				/*
				** if -NR words are ordered so we store and retrieve from an array
				** else we use a map
				 */
				if flagNR {
					wordArray = append(wordArray, tmpWord)
					if len(wordArray) == flagnum {
						done = true

					}

				} else {
					// add to map if not there
					if _, ok := wordMap[tmpWord]; ok != true {
						wordMap[tmpWord] = struct{}{}
					}
				}
			} else {
				discarded = true
			}
		}

		// proSigns for NR = false done differently
		if flagNR && flagprosign != "" {
			replaceIndex := make(map[int]struct{})

			for i := 0; i < len(proSign); {
				rand := rng.Intn(len(wordArray))
				if _, ok := replaceIndex[rand]; ok != true {
					replaceIndex[rand] = struct{}{}
					i++
				}
			}

			// now do the substitions
			j := 0
			for index := range replaceIndex {
				temp := append([]string{proSign[j]}, wordArray[index:]...)
				wordArray = append(wordArray[:index], temp...)
				j++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	msg := "\nSorry there is nothing to output.\nMake sure the options (or defaults) are not to restrictive (min, max, inlist).\nVerify your input file is sufficiently populated with matchable text.\n"

	if flagNR {

		if len(wordArray) == 0 {
			fmt.Println(msg)

			if discarded {
				fmt.Printf("Your input file DID have some text.\n")
				if localSkipFlag {
					fmt.Printf("Your -skip X option maybe too aggresive.\n")
				}
			}
			os.Exit(0)
		}

		ct := len(wordArray)
		if ct < flagnum {
			ct = flagnum - ct
			// we need to append more words in order
			for i := 0; i < ct; i++ {
				wordArray = append(wordArray, wordArray[i])
			}
		}

	} else {
		if len(wordMap) == 0 {
			fmt.Println(msg)
			if discarded {
				fmt.Printf("Your input file DID have some text.\n")
				if localSkipFlag {
					fmt.Printf("Your -skip X option maybe too aggresive.\n")
				}
			}
			os.Exit(0)
		}

		// trim the saved wordMap to save time and memory later
		if len(wordMap) > flagnum {
			cntr := 0
			m := make(map[string]struct{})
			for v, _ := range wordMap {
				if cntr == flagnum {
					wordMap = m
					m = nil
					break
				} else {
					m[v] = struct{}{}
					cntr++
				}
			}
		}
		runtime.GC()

		fillArray(fp)
	}
}

// ready to print the users practice word
func doOutput(words []string, fp *os.File) {
	strBuf := ""
	strOut := ""
	// for eb options
	firstSlowFast := true
	lastSpeed := flagEBlow
	lastSpeedEff := flagEBeff
	counter := 1
	sectionSize := 0
	EBspeeds := []int{}
	EBspeedsRepeat := []int{}
	ebslowcnt := 0
	ebfastcnt := 0
	ebinslow := false
	fBOUNCE := false
	fSLOWFAST := false
	fRAMP := false
	fEFFRAMP := false
	fREPEAT := false
	speedCount := 0
	var charSlice []rune

	if flagMixedMode > 0 {
		charSlice = buildCharSlice()
	}

	// header
	if flagheader != "" {
		strOut += fmt.Sprintf("%s\n", flagheader)
	}

	// for runtime eff
	if flagEBramp {
		fRAMP = true
	} else if flagEBefframp {
		fEFFRAMP = true
	} else if flagEBslow > 0 {
		fSLOWFAST = true
	} else if flagEBnum > 0 && flagEBlow > 0 && flagEBstep >= 1 && !fRAMP && !fSLOWFAST && !fEFFRAMP && !fREPEAT {
		fBOUNCE = true
	}

	// for runtime eff
	if flagEBrepeat >= 1 && !fBOUNCE {
		fREPEAT = true
	}

	///////////////////////////////////
	////// EB handling - intial setup
	///////////////////////////////////
	// seed array with EB_LOW, and fill as appropriate
	if flagEBnum >= 1 && flagEBstep > 0 {
		for i := 0; i < flagEBnum; i++ {
			EBspeeds = append(EBspeeds, flagEBlow+(i*flagEBstep))
		}
	}

	if flagEBrepeat >= 1 && flagEBstep > 0 {
		for i := 0; i < flagEBrepeat; i++ {
			EBspeedsRepeat = append(EBspeedsRepeat, flagEBlow+(i*flagEBstep))
		}
	}

	// EB fRAMP how many words per ramp section
	if fRAMP && !fEFFRAMP {
		sectionSize = flagnum / flagEBnum
		if sectionSize < 1 {
			fmt.Printf("\nError: EB_NUM is too large for the -num value.\nThere would not be any words in each speed change section.\n")
			os.Exit(1)
		}

		lastSpeed = EBspeeds[0]

		if !fREPEAT {
			if flagEBeff > 0 {
				strOut += fmt.Sprintf("|w%d |e%d ", EBspeeds[0], flagEBeff)
			} else {
				strOut += fmt.Sprintf("|w%d ", EBspeeds[0])
			}
			speedCount++
		}
	}

	// EB EFFECTIVE_RAMP how many words per ramp section
	if fEFFRAMP && !fRAMP {
		sectionSize = flagnum / flagEBnum

		if sectionSize < 1 {
			fmt.Printf("\nError: EB_NUM is too large for the -num value.\nThere would not be any words in each speed change section.\n")
			os.Exit(1)
		}

		lastSpeedEff = flagEBeff
		strOut += fmt.Sprintf("|w%d |e%d ", flagEBlow, flagEBeff)
		counter = 0
	}

	////////////////////////////////////////
	/// setup done, now process input words
	/////////////////////////////////////////
	// select words from array, then lower high water mark so all words get used

	wcFlg := false
	cnt := 0
	dWord := ""

	if flagwordcount > 1 {
		wcFlg = true
	}

	for index, wordOut := range words {

		// if true we are linking words together to treat as a entity
		if wcFlg {

			if cnt < flagwordcount {
				dWord += wordOut
				dWord += "_"
				cnt++
				continue
			} else {
				wordOut = strings.Trim(dWord, "_")
				dWord = ""
				cnt = 0
			}


			/*
			if wcCnt == 1 {
				wcCnt++
				continue
			} else if wcCnt <= flagwordcount {
				dWord += wordOut
				dWord += "_"
				wcCnt++

				if wcCnt > flagwordcount {
					wcCnt = 1
				}
			}
			*/
		}

		//////////////
		// EB fBOUNCE
		//////////////
		if fBOUNCE {
			speed := EBspeeds[0]

			if len(EBspeeds) >= 1 {

				for {
					speed = EBspeeds[rng.Intn(len(EBspeeds))]

					if speed != lastSpeed {
						lastSpeed = speed
						break
					}
				}
			}

			if flagEBeff > 0 {
				strOut += fmt.Sprintf("|w%d |e%d ", speed, speed-effDelta)
			} else {
				strOut += fmt.Sprintf("|w%d ", speed)
			}
		}

		/////////////////
		/// EB FAST_SLOW
		/////////////////
		if fSLOWFAST {
			s := flagEBlow

			if ebinslow {
				if ebslowcnt >= flagEBslow {
					s = flagEBlow + flagEBstep // now fast
					// slow words are done
					ebfastcnt = 0

					// keep eff same
					strOut += fmt.Sprintf("%s|w%d ", flagEBsf, s)
					ebinslow = false
				}
			} else {
				if ebfastcnt >= flagEBfast || firstSlowFast {
					firstSlowFast = false
					s := flagEBlow // now slow

					// fast words are done
					ebslowcnt = 0

					// set up slow section
					if index == 0 {
						if flagEBeff > 0 {
							strOut += fmt.Sprintf("|e%d |w%d ", flagEBeff, s)
						} else {
							strOut += fmt.Sprintf("|w%d ", s)
						}
					} else {
						strOut += fmt.Sprintf("%s|w%d ", flagEBfs, s)
					}
					ebinslow = true
				}
			}
		}

		// end raw word, and get back word to print
		wordOut, charSlice = prepWord(wordOut, lastSpeed, index, charSlice)

		///////////////////////////////////
		// EB CHECK FOR SPEED MARKERS
		///////////////////////////////////

		if fRAMP {

			if counter >= sectionSize && speedCount < len(EBspeeds) {
				sf := ""
				if flagEBsf != "" {
					sf = " " + flagEBsf
				}

				if index+flagEBnum <= flagnum {
					if flagEBeff > 0 {
						strOut += fmt.Sprintf("%s%s|e%d |w%d ", wordOut, sf, EBspeeds[speedCount]-effDelta, EBspeeds[speedCount])
					} else {
						strOut += fmt.Sprintf("%s%s|w%d ", wordOut, sf, EBspeeds[speedCount])
					}
					wordOut = ""
					speedCount++
				}

				if speedCount < len(EBspeeds) {
					lastSpeed = EBspeeds[speedCount]
				}
				counter = 1
			} else {
				counter++
			}
		}

		////////////
		// fEFFRAMP
		///////////
		if fEFFRAMP {
			if counter == sectionSize {
				// ck if eff is going to over take word speed
				if lastSpeedEff+flagEBstep <= flagEBlow {
					// cap the eff speed
					lastSpeedEff += flagEBstep
				}
				strOut += fmt.Sprintf("%s|e%d ", wordOut, lastSpeedEff)
				counter = 1
				strOut += wordOut
				wordOut = ""
			} else {
				counter++
			}

		}

		//////////////
		/// fSLOWFAST
		//////////////
		if fSLOWFAST {
			if ebinslow {
				ebslowcnt++
			} else {
				ebfastcnt++
			}
		}

		if wcFlg {
			// get back to individual words
			wordOut = strings.ReplaceAll(wordOut, "_", " ")
		}

		// this is the processed word to be used
		strBuf += strOut + wordOut
		strOut = ""
	}

	printStrBuf(strBuf, fp)
}

//
// prints the bufStr adjusting the length per flaglen
//
func printStrBuf(strBuf string, fp *os.File) {
	// done processing now output it
	res := ""
	index := 0
	for _, r := range strBuf {

		if index <= flaglen {
			res = res + string(r)
			index++
			continue
		}

		if index >= flaglen {
			if r != ' ' && r != '\n' {
				res = res + string(r)
				index++
				continue
			} else {
				res = res + "\n"
				index = 0
			}
		}
	}

	if flagoutput == "" {
		fmt.Printf("%s", res)
	} else {

		_, err := fp.WriteString(res)
		if err != nil {
			fmt.Println(err)
			fp.Close()
			os.Exit(0)
		}
	}
} // end

// simple random true or false
//func flipFlop(s string) bool {
func flipFlop() bool {
	if rng.Intn(2) == 1 {
		return true
	}
	return false
}

/*
** take in a raw word from input file and tack on: prefix, suffix
** repeat if necessay,do mixedMode
 */
func prepWord(wordOut string, lastSpeed int, index int, charSlice []rune) (string, []rune) {
	strOut := ""
	rand := 3

	if flagrandom {
		if flagsuf >= 1 || flagpre >= 1 {
			// 0 - neither ix, 1 prefix,2 do suffix, 3 both
			rand = rng.Intn(4)
		}
	}

	// end raw word, and get back word to print
	// do we need prefix?
	if flagpre >= 1 && (rand == 3 || rand == 1) {
		wordOut = ixStr("p") + wordOut
	}

	// do we need a suffix or just a space
	if flagsuf >= 1 && (rand == 3 || rand == 2) {
		wordOut += ixStr("s")
	}

	// text repeat!
	if flagrepeat > 0 {
		// we need to repeat
		wordOut += " "
		temp := wordOut

		for cnt := 1; cnt < flagrepeat; cnt++ {
			// wordOut is the word plus trailing space already
			wordOut += temp
		}
	}

	// EB_REPEAT
	if flagEBrepeat > 1 {
		for i := 0; i < flagEBrepeat; i++ {
			// if we ALSO have fRAMP we must offset speed
			spd := lastSpeed + (i * flagEBstep)

			if flagEBeff > 0 {
				strOut += fmt.Sprintf("|w%d |e%d %s", spd, spd-effDelta, wordOut)
			} else {
				strOut += fmt.Sprintf("|w%d %s", spd, wordOut)
			}

			if flagEBramp {
				spd += lastSpeed
			}
		}
	}

	// why WDL
	/*
		if flagEBramp && flagEBrepeat == 0 {
			strOut = wordOut
			wordOut = ""
		}
	*/

	// mixedMode put out code Group
	if flagMixedMode > 1 && (flagMMR == false || (flagMMR == true && flipFlop())) {
		g := []rune{}

		if index%flagMixedMode == 0 {
			if flagEBsf != "" {
				strOut += flagEBsf + " "
			}

			g, charSlice = makeSingleGroup(charSlice)
			strOut += string(g)
			if flagEBfs != "" {
				strOut += flagEBfs + " "
			}
		}
	}

	// this means NOTHING was done to the word
	if strOut == "" {
		strOut = wordOut
	}

	// use delimiter
	if flagDM > 0 && (flagDR == false || (flagDR == true && flipFlop())) {
		d := ""
		for i := 1; i <= (1 + rng.Intn(flagDM)); i++ {
			d += delimiterSlice[rng.Intn(len(delimiterSlice))]
		}
		strOut += d
		strOut += " "
	}

	return strOut, charSlice
}

// reverse a string
func reverse(s string) string {
	rs := []rune(s)

	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}

	return string(rs)
}

func findPercent(in string) (string, string) {
	out := ""

	m := regexp.MustCompile("%[C-F][0146789C]")
	if m.MatchString(in) {
		s := m.FindString(in)
		in = strings.Replace(in, string(s), "", 1)
		out = strings.TrimLeft(s, "%")
		out = string(runeMapInt[out])
	} else {
		fmt.Printf("Error: invalid character in delimiter option\n")
		os.Exit(88)
	}

	return in, out
}

//
// called for each field of a delimiter option
func processDelimiter(inStr string) {
	// eliminate special case of simple range
	m := regexp.MustCompile("^([0-9]-[0-9])|([a-z]-[a-z])|([A-Z]-[A-Z])$")
	if m.MatchString(inStr) {
		expandIt(string(inStr[0]), string(inStr[2]), "delimiter_simple")
		return
	}

	ckValidInString(inStr, "delimiter")
}


