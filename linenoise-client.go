// SPDX-License-Identifier: MIT
//
// Copyright (c) 2017-2021 Mark Cornick
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Noise struct {
	Text string
}

type Error struct {
	Message string
}

func btos(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	var (
		length = 16
		upper  = true
		lower  = true
		digit  = true
	)
	flag.IntVar(&length, "length", 16, "Length to generate")
	flag.BoolVar(&upper, "upper", true, "Include uppercase letters")
	flag.BoolVar(&lower, "lower", true, "Include lowercase letters")
	flag.BoolVar(&digit, "digit", true, "Include digits")
	flag.Parse()

	endpoint := fmt.Sprintf("http://127.0.0.1:8080/v1/noise?length=%d&upper=%s&lower=%s&digit=%s",
		length,
		btos(upper),
		btos(lower),
		btos(digit),
	)
	res, err := http.Get(endpoint)
	check(err)
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	check(err)
	if res.StatusCode > 299 {
		var errorObj Error
		err := json.Unmarshal(body, &errorObj)
		check(err)
		if res.StatusCode == 422 {
			log.Fatalf("%s\n", errorObj.Message)
			os.Exit(1)
		}
		log.Fatalf("%d %s\n", res.StatusCode, errorObj.Message)
		os.Exit(1)
	}
	var noise Noise
	err = json.Unmarshal(body, &noise)
	check(err)
	fmt.Printf("%s\n", noise.Text)
}
