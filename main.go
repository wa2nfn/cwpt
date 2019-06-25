package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
	"path/filepath"
)

const PROG = "cwpt.exe"
const MAX_WORD_LEN = 15
const MIN_WORD_LEN = 1
const MAX_USER_WORDS = 8000
const MAX_LINE_LEN = 400
const MAX_SUFFIX = 10 // may be perfomance issue

var minWordLen = flag.Int("min", 5, "Minimum word length.\n")
var maxWordLen = flag.Int("max", 5, "Maxmum word length, >= min.\n")
var wordMap = make(map[string]bool)
var numUserWords = flag.Int("num", 100, "Number of output words. Minimum 1, maximum 8000.\n")
var userLineLength = flag.Int("lineLength", 80, "Number characters per line (max 400).\n")
var skipCount = flag.Int("skip", 0, "Number of the first unique words to ignore. Max 2500\n")
var capsBool = flag.Bool("caps", false, "Make output all caps.\n (default lower case)")
var spacePtr = flag.Bool("space", false, "Add a space character to suffixChars.\n (default flase)")
var suffixPtr = flag.Int("suffix", 0, "The max number of suffix characthers to append to words.")
var seed = time.Now().UTC().UnixNano()
var rng = rand.New(rand.NewSource(seed))
var userSuffixChars = flag.String("suffixChars", "", "A set of characters to tack onto a word.\n-suffix X, determines the quanty.\nIntent was for numbers, and punctuation, but can be more.\nUse inside quotes, if a \\\" will be included it must be escaped like: \".,?\\\"ZX\"\n (default \"0-9.,?/=\")")
var alphaChars = flag.String("alphaChars", "A-Za-z", "Set of characters to define an input word.\nUse quotes.\n")

func main() {
	var skipFlag = false

	if filepath.Base(os.Args[0]) != PROG {
		fmt.Printf("\nSorry, the name of the executable must be: %s      - 73, Bill WA2NFN\n\n", PROG)
		os.Exit(0)
	}

	var filePtr = flag.String("file", "afile.txt", "input data file\n")

	flag.Parse()

	if *skipCount  > 0 {
		// we will be skipping some words
		skipFlag = true
	}

	if *minWordLen < MIN_WORD_LEN || *minWordLen > MAX_WORD_LEN {
		fmt.Printf("min must > 0 and <= max\n")
		os.Exit(0)
	}

	if *maxWordLen < *minWordLen || *maxWordLen > MAX_WORD_LEN {
		fmt.Printf("min must > 0 and <= max\n")
		os.Exit(0)
	}

	if *numUserWords < 1 || *numUserWords > MAX_USER_WORDS {
		fmt.Printf("-num number of unique output words desired. minimum 1, maximum %d, default 5\n", MAX_USER_WORDS)
		os.Exit(0)
	}

	if *userLineLength < 0 || *numUserWords > MAX_USER_WORDS {
		fmt.Printf("-num max number of chars per line, \"0\" means single column of words, default 80\n", MAX_LINE_LEN)
		os.Exit(0)
	}

	if *suffixPtr < 0 || *suffixPtr > MAX_SUFFIX {
			fmt.Printf("-suffix, 0=no suffix, max number of characters is %s", MAX_SUFFIX)
			os.Exit(0)
	}

	file, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if *alphaChars == "" {
		fmt.Printf("\n%s-alphaChars can't be empty or nothing gets matched\n")
		os.Exit(0)
	}
	word := regexp.MustCompile(fmt.Sprintf("^[%s]{%d,%v}$", *alphaChars,*minWordLen, *maxWordLen))

	for scanner.Scan() {
		text := scanner.Text()
		for _, text := range regexp.MustCompile(fmt.Sprintf("[^%s]",*alphaChars)).Split(text, -1) {
			// every token is now a string of space separated characters
			if word.MatchString(text) {

				// set case before storing
				if *capsBool {
					text = strings.ToUpper(text)
				} else {
					text = strings.ToLower(text)
				}

				// add to map if not there
				if _, ok := wordMap[text]; ok != true {
					wordMap[text] = true
					if skipFlag && *skipCount > 0 {
						*skipCount--
					}

					if skipFlag &&  *skipCount == 0 {
						// clear the map and start saving again
						//wordMap = nil
						wordMap = make(map[string]bool)
						skipFlag = false
					}

					if len(wordMap) == *numUserWords {
						break
					}
				}
			}
		}

		if len(wordMap) == *numUserWords {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(wordMap) == 0 {
		fmt.Println("Sorry there is nothing to output.\nVerify your options (or the defaults) will find words in your input file.")
		os.Exit(0)
	}

	fillArray()
}

// fill the array from the word map but might need to stuff more values
func fillArray() {
	var wordArray = make([]string, 0, *numUserWords)
	factor := 0

	for key := range wordMap {
		// make first population of slice
		wordArray = append(wordArray, key)
	}

	// see if initial array satisfies the number of words the user wanted
	// if less, we will reuse words from map to grow the array (or slice)
	factor-- // one pass of map already in array
	factor = (*numUserWords - len(wordArray)) / len(wordMap)
	for ; factor > 0; factor-- {
		for key := range wordMap {
			wordArray = append(wordArray, key)
		}
	}

	remainder := *numUserWords % len(wordArray)
	for key := range wordMap {
		if remainder == 0 {
			break
		}
		wordArray = append(wordArray, key)
		remainder--
	}

	// trash the map to conserve memory
	wordMap = nil

	doOutput(wordArray)
}

// ready to print the users practice word
func doOutput(words []string) {
	lineLength := 0
	wordLength := 0
	wordsYetToSend := len(words) // always shrinking
	wordOut := "" // always changing

	for index := 0; index < len(words); index++ {
		// select words randomly from array, then lower high water mark so all words get used
		rn :=  rng.Intn(wordsYetToSend)
		wordOut = words[rn]
		words[rn] = words[wordsYetToSend-1]
		words[wordsYetToSend-1] = "" //save memory
		wordsYetToSend--

		if *userLineLength == 0 { // print in a column
			suffix, sufLen := sufStr()

			if sufLen != 0  {
				fmt.Printf("%s%s\n", wordOut, suffix)
			} else {
				fmt.Printf("%s\n", wordOut)
			}
		} else {
			// words in a line at users length

			wordLength = len(wordOut)
			if lineLength + wordLength > *userLineLength {
				// output will be too long
				fmt.Printf("\n%s", wordOut)
				lineLength = 0
			} else {
				fmt.Printf("%s", wordOut)
			}

			lineLength += wordLength

			// see if we need a suffix or just a space
			suffix, sufLen := sufStr()
			if sufLen != 0 {
				fmt.Printf("%s ", suffix)
				lineLength =  lineLength + sufLen + 1
			} else {
				fmt.Printf(" ")
				lineLength++
			}
		}

	}
	fmt.Println()
}

// returns a single random suffix from list to add to output word
func sufStr() (string,int) {
	suf := "1234567890=,./?"
	retStr := ""

	if *suffixPtr == 0 {
		return "",0
	} else {
		if len(*userSuffixChars) >= 1 {
			suf = *userSuffixChars
		}
		if *spacePtr {
			suf  +=" "
		}
		ll := len(suf)

		// user wants a suffix
		for count := 1;  count <= rng.Intn(*suffixPtr)+1; count++ {
			retStr += string(suf[rng.Intn(ll)])
		}
		return retStr, len(retStr)
	}
}
