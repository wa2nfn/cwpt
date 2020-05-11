//
// Copyright 2019 Bill Lanahan
//
// Read a given number of words from a file, then output the requested number in random order
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	version      = "1.2 04/01/2020"
	maxWordLen   = 25
	minWordLen   = 1
	maxUserWords = 50000
	maxLineLen   = 500
)

var (
	seed    = time.Now().UTC().UnixNano()
	rng     = rand.New(rand.NewSource(seed))
	wordMap = map[string]struct{}{}
)

var (
	flagmax        int
	flaglen        int
	flagmin        int
	flagnum        int64
	flagnumperline int
	flaginput      string
	flagoutput     string
	flagdelimiter  string
	flagversion    bool
	flagsingleword bool
	storedWordCnt  int64
	flaglesson     int
	flagtutor      string
	flaginlist     string
)

func init() {
	flag.IntVar(&flagnumperline, "numPerLine", 0, "Number of strings to print per line (0 means fit to 80 char line).")
	flag.IntVar(&flagmax, "max", 25, "Maximum # characters in a string >= min.")
	flag.IntVar(&flagmin, "min", 1, "Minimum # of characters in a string.")
	flag.IntVar(&flaglen, "len", 80, "Length of characters in output line.")
	flag.Int64Var(&flagnum, "num", 0, fmt.Sprintf("Number of strings to output. Min 0 (0 means ALL input file words), max %d.\n", maxUserWords))
	flag.BoolVar(&flagversion, "version", false, "Display version information. (default false)")
	flag.BoolVar(&flagsingleword, "single", false, "Display single string per line (default false)")
	flag.StringVar(&flaginput, "in", "in.txt", "Input text file name (including extension).")
	flag.StringVar(&flagoutput, "out", "", "Output file name.")
	flag.StringVar(&flagtutor, "tutor", "LCWO", "Only if you use -lessons. Sets order and # of charactersby tutor type.\nChoices: (default LCWO), JustLearnMorseCode, G4FON, MorseElmer")

	flag.IntVar(&flaglesson, "lesson", 0, "Given the Koch lesson number per LCWO, populates options inlist and cglist with appropriate characters. (Default 0)")
	flag.StringVar(&flaginlist, "inlist", "A-Za-z", "Set of characters to define an input word.")
	flag.StringVar(&flagdelimiter, "delimiter", "", "Inter-output string delimiter.")
}

func main() {
	flag.Usage = func() {

	const text =`

random.exe - randomize strings in a file for a morse tutor, sending practice, typing practice, etc.

  random.exe [-in=in.txt] [-out=<stdout>] [-min=1] [-max=25] [-num=0] [-len=80] [-numPerLine=0] [-lesson=0] [-tutor=LCWO] (default options shown)
  
  where:
         -in is the input file to read to obtain character strings (default in.txt)
	 -out is the output file to write to (default is stdout (your screen))
         -min the minimum length of a string to save (default=1)
	 -max the maximum length of a string to save (default=25)
	 -num the number of strings to find in the input file and save for possible printing (default=0, 0 means ALL!)
	 -single or -numPerLine the number of strings to print per line (if line length (len) is not exceeded)
	 -len number of characters in the output line length
	 -lesson X, where X is the lesson number for your tutor (default=0)
	 -tutor X, where X is:  LCWO, JustLearnMorseCode, G4FON, MorseElmer (default=LCWO)
	 -numPerLine number of strings to print on output (if line length (len) is not exceeded)
	 -delimiter a string of characters to delimit output strings (default a space)
	 -version software version
	 -help this help text
	 
   when run:
	The utility will read the input and store strings in random order.
	Report the number of strings read (less than or equal to the number in the input file).
	Prompt you to enter the number of the read strings you want as output (0 means you want them all!).
    `
		fmt.Println(text)

		fmt.Printf("\ne.x. myin.txt contains: \"A short input file for demonstration. 123 <SK>  2+2=4\"\n\nrandom.exe -in=myin.txt -min=4 -max=7\n\nOutput after running:\n\ninput 123 <SK> short file 2+2=4 for\n\nNote: 2 strings did not meet the size criteria so were dropped.\nNote: It is up to your morse tutor or follow on tool to cull out or modify\n      extraneous data.")
		os.Exit(0)
	}

	kochChars := "kmuresnaptlwi.jz=foy,vg5/q92h38b?47c1d60x" // default for LCWO
	var fp *os.File

	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Printf("\nError processing the command line.\nYou may have: forgotten a \"-\" before an option\n  or followed a \"-\" with a space\n  or added extra input.\n")
		os.Exit(1)
	}

	if flagversion {
		fmt.Printf("\nversion: %s\n", version)
		fmt.Println("\nCopyright 2019")
		os.Exit(0)
	}

	if flagmin < minWordLen {
		fmt.Printf("\nError: -min must be >= %d, system min.\n", minWordLen)
		os.Exit(1)
	}

	if flagmin > flagmax {
		fmt.Printf("\nError: -min must <= max(M) <%d>.\n", flagmax)
		os.Exit(1)
	}

	if flagmax < flagmin {
		fmt.Printf("\nError: -max must >= min(m) <%d>.\n", flagmin)
		os.Exit(1)
	}

	if flagmax > maxWordLen {
		fmt.Printf("\nError: -max must <= <%d>, system max.\n", maxWordLen)
		os.Exit(1)
	}

	if flagsingleword && flagnumperline > 1 {
		fmt.Printf("\nError: -single and -numPerLine are mutually exclusive.\n")
		os.Exit(1)
	}

	if flagnumperline < 0 || flagnumperline > 5000 {
		fmt.Printf("\nError: -numPerLine out of range.\n")
		os.Exit(1)
	}

	if flagnum < 0 || flagnum > maxUserWords {
		fmt.Printf("\nError: -num number of output words desired. minimum 0 (0 means all the input words), maximum %d, default 5.\n", maxUserWords)
		os.Exit(1)
	}

	if flaglen < 1 || flagnum > 5000 {
		fmt.Printf("\nError: -len characters per output line >=1 and <= 5000\n")
		os.Exit(1)
	}

	if flagoutput != "" {
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
		fmt.Printf("\n*** Writing to file: %s\n", flagoutput)
	}

	s := 0

	if flaglesson == 0 && flagtutor != "LCWO" {
		fmt.Printf("\nError: Lesson = 0 is invalid for tutor <%s>.\n", flagtutor)
		os.Exit(1)
	}

	if flaglesson >= 1 {

		if flagtutor == "LCWO" {
			kochChars = "kmuresnaptlwi.jz=foy,vg5/q92h38b?47c1d60x"
		} else if flagtutor == "JustLearnMorseCode" {
			kochChars = "kmrsuaptlowi.njef0yv,g5/q9zh38b?427c1d6x@=+"
			s = 0
		} else if flagtutor == "G4FON" {
			kochChars = "kmrsuaptlowi.njef0yv,g5/q9zh38b?427c1d6x"
			s = 0
		} else if flagtutor == "MorseElmer" {
			kochChars = "kmrsuaptlowi.njef0y,vg5/q9zh38b?427c1d6x="
			s = 0
		} else {
			fmt.Printf("\nError: Your tutor name is invalid. Names are case sensitive and without any spaces, see the help.\n")
			os.Exit(1)
		}

		if (flaglesson+1 > len(kochChars)) && flagtutor == "LCWO" {
			fmt.Printf("\nError: Lesson value <%d> exceeds the max <%d>, for tutor <%s>.\n", flaglesson, 40, flagtutor)
			os.Exit(1)
		}

		if flaglesson > len(kochChars) {
			fmt.Printf("\nError: Lesson value <%d> exceeds the max <%d>, for tutor <%s>.\n", flaglesson, len(kochChars), flagtutor)
			os.Exit(1)
		}

		if flagtutor == "LCWO" {
			if flaglesson < len(kochChars) {
				flaginlist = kochChars[0 : flaglesson+1]
			} else {
				flaginlist = kochChars[0:flaglesson]
			}
		} else {
			flaginlist = kochChars[s:flaglesson]
		}

		// now build inlist initailly as lower case
		// now inlist is LC
		temp := strings.ToUpper(flaginlist)

		// now add upper case
		for _, char := range temp {
			if char >= 'A' && char <= 'Z' {
				flaginlist += string(char)
			}
		}
		// now inlist is mixed case
		temp = ""

	}

	readFileMode(fp)
	os.Exit(0)
}

func readFileMode(fp *os.File) {

	file, err := os.Open(flaginput)

	if err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	word := regexp.MustCompile(fmt.Sprintf(`^[%s]{%d,%d}\s*$`, flaginlist, flagmin, flagmax))

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// first way to split the string on spaces
		textWords := strings.FieldsFunc(scanner.Text(), func(r rune) bool {
			if r == ' ' {
				return true
			}
			return false
		})

		for index := 0; index < len(textWords); index++ {
			// every token is now a string of space separated characters

			tmpWord := textWords[index]

			// IF it meets inlist criteria
			if !word.MatchString(tmpWord) {
				continue
			}

			// add to map if not there
			if _, ok := wordMap[tmpWord]; ok != true {
				wordMap[tmpWord] = struct{}{}
				storedWordCnt++

				if storedWordCnt == flagnum {
					break
				}
			}
		}

		if flagnum > 0 && storedWordCnt == flagnum {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	outPut(fp)
	os.Exit(0)

}

func outPut(fp *os.File) {
	var wordsToPrint int64
	numPtr := ""
	strOut := ""
	lineLen := flaglen
	outCnt := 0

	if storedWordCnt == 0 {
		fmt.Println("\nError: Sorry there is nothing to output.\nMake sure your input file has text.")
		os.Exit(0)
	}

	fmt.Printf("\nStrings read in: %d\n\nEnter the number of strings to print (zero prints them all): ", storedWordCnt)
	fmt.Scan(&numPtr)
	fmt.Printf("====================================================================================\n")

	wordsToPrint, err := strconv.ParseInt(numPtr, 10, 64)
	if err != nil {
		fmt.Printf("\nError: invalid input: %T, %v\n", wordsToPrint, wordsToPrint)
	}

	if wordsToPrint == 0 {
		wordsToPrint = storedWordCnt
	}

	if flagnum > 0 && wordsToPrint <= 0 {
		fmt.Printf("\nError: You asked for too little.\n")
		os.Exit(2)
	}

	if wordsToPrint > storedWordCnt {
		fmt.Printf("\nError: You asked for more than what was read.\n")
		os.Exit(2)
	}

	if flagoutput == "" {
		fmt.Printf("\n")
	}

	// first move the map into a slice
	skeys := []string{}
	for key, _ := range wordMap {
		skeys = append(skeys, key)
		// delete key to save space
		delete(wordMap, key)
	}

	// read each string from slice
	newLen := len(skeys) // slice size
	end := len(skeys)    // slice size
	for pos := 0; pos < end; pos++ {
		newPos := rng.Intn(newLen)
		strOut += skeys[newPos]
		outCnt++

		if flagdelimiter == "" {
			strOut += " "
		} else {
			strOut += flagdelimiter
		}

		newLen--                      // "shorten" slice
		skeys[newPos] = skeys[newLen] // fill taken slice value with highest value

		if newLen != 0 {
			skeys[newLen] = "" // reduce slice size
		}

		lineLen -= len(strOut)

		if flagsingleword || flagnumperline == 1 {
			strOut += "\n"
			lineLen = flaglen
		} else {
			// put return if count reached
			if outCnt >= flagnumperline || lineLen <= 0 {
				outCnt = 0
				lineLen = flaglen
				strOut += "\n"
			}
		}

		// print it
		if flagoutput != "" {
			_, err := fp.WriteString(strOut)
			if err != nil {
				fmt.Println(err)
				fp.Close()
				os.Exit(0)
			}
		} else {
			fmt.Printf("%s", strOut)
		}
		strOut = ""

		wordsToPrint--
		if wordsToPrint == 0 {
			fmt.Println()
			os.Exit(0)
		}
	}
}
