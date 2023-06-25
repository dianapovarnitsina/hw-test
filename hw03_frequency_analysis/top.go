package main

import (
	"sort"
	"strings"
)

func Top10(text string) []string {

	if text == "" {
		return nil
	} else {
		// Разделение абзаца на слова
		words := strings.Fields(text)

		// Создание карты для подсчета повторений
		wordCount := make(map[string]int)

		// Подсчет повторений слов
		for _, word := range words {
			wordCount[word]++
		}

		// Из мапы получить слова с максимальным количеством повторений 10 штук
		return getTopWords(wordCount, 10)
	}
}

func getTopWords(wordCount map[string]int, n int) []string {

	type wordCountPair struct {
		word  string
		count int
	}

	// Создание среза из структур wordCountPair
	pairs := make([]wordCountPair, len(wordCount))
	i := 0
	for word, count := range wordCount {
		pairs[i] = wordCountPair{word, count}
		i++
	}

	// Сортировка лексикографически
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].word < pairs[j].word
	})

	// Сортировка среза по убыванию частоты
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Получение топ-N слов
	topWords := make([]string, 0, n)
	for _, pair := range pairs[:n] {
		topWords = append(topWords, pair.word)
	}
	return topWords
}
