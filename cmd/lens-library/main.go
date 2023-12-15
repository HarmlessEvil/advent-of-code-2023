package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"log/slog"
	"os"
	"slices"
)

type Hash struct {
	State byte
}

func (h *Hash) Write(p []byte) (n int, err error) {
	for _, symbol := range p {
		h.State = byte((int(h.State) + int(symbol)) * 17)
	}

	return len(p), nil
}

func (h *Hash) Sum64() uint64 {
	return uint64(h.State)
}

func (h *Hash) Sum(b []byte) []byte {
	return append(b, h.State)
}

func (h *Hash) Reset() {
	h.State = 0
}

func (h *Hash) Size() int {
	return 1
}

func (h *Hash) BlockSize() int {
	return 1
}

type HashMap struct {
	Buckets [256]list.List
}

type MapEntry struct {
	Key   []byte
	Value int
}

func (m *HashMap) bucketIndex(key []byte) uint64 {
	var hash Hash
	_, _ = hash.Write(key)

	return hash.Sum64()
}

func (m *HashMap) find(bucket *list.List, key []byte) *list.Element {
	for e := bucket.Front(); e != nil; e = e.Next() {
		mapEntry := e.Value.(*MapEntry)
		if slices.Equal(mapEntry.Key, key) {
			return e
		}
	}

	return nil
}

func (m *HashMap) Set(key []byte, value int) {
	bucket := &m.Buckets[m.bucketIndex(key)]

	if e := m.find(bucket, key); e != nil {
		mapEntry := e.Value.(*MapEntry)
		mapEntry.Value = value
	} else {
		bucket.PushBack(&MapEntry{
			Key:   append([]byte{}, key...),
			Value: value,
		})
	}
}

func (m *HashMap) Delete(key []byte) {
	bucket := &m.Buckets[m.bucketIndex(key)]

	if e := m.find(bucket, key); e != nil {
		bucket.Remove(e)
	}
}

func runMain() error {
	hashMap, err := parseHashMap()
	if err != nil {
		return fmt.Errorf("parse hash map: %w", err)
	}

	sum := 0

	for i, bucket := range hashMap.Buckets {
		if bucket.Len() == 0 {
			continue
		}

		fmt.Printf("Box %d:", i)
		for e, j := bucket.Front(), 0; e != nil; e, j = e.Next(), j+1 {
			mapEntry := e.Value.(*MapEntry)
			fmt.Printf(" [%s %d]", mapEntry.Key, mapEntry.Value)

			sum += (i + 1) * (j + 1) * mapEntry.Value
		}
		fmt.Println()
	}

	fmt.Println(sum)

	return nil
}

func parseHashMap() (HashMap, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return HashMap{}, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := bytes.IndexRune(data, ','); i != -1 {
			return i + 1, data[:i], nil
		}

		if !atEOF {
			return 0, nil, nil
		}

		return 0, bytes.TrimRight(data, "\r\n"), bufio.ErrFinalToken
	})

	var hashMap HashMap

	for j := 0; scanner.Scan(); j++ {
		line := scanner.Bytes()

		label, i := parseLabel(line)
		switch line[i] {
		case '=':
			hashMap.Set(label, parseInt(line[i+1:]))
		case '-':
			hashMap.Delete(label)
		}
	}

	if err := scanner.Err(); err != nil {
		return HashMap{}, fmt.Errorf("scan: %w", err)
	}

	return hashMap, nil
}

func parseLabel(line []byte) ([]byte, int) {
	i := bytes.IndexRune(line, '=')
	if i == -1 {
		i = len(line) - 1
	}

	return line[:i], i
}

func parseInt(data []byte) int {
	res := 0
	for _, item := range data {
		res = res*10 + int(item) - '0'
	}

	return res
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
