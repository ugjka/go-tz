// Code generation tool for embedding the timezone shapefile in the gotz package
// run "go generate" in the parent directory after changing the -release flag in gen.go
// You need mapshaper to be installed and it must be in your $PATH
// More info on mapshaper: https://github.com/mbloch/mapshaper
package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/machinebox/progress"
)

const dlURL = "https://github.com/evansiroky/timezone-boundary-builder/releases/download/%s/timezones.geojson.zip"

func main() {
	_, err := exec.LookPath("mapshaper")
	if err != nil {
		log.Fatalln("Error: mapshaper executable not found in $PATH")
	}

	release := flag.String("release", "", "timezone boundary builder release version")
	flag.Parse()

	resp, err := http.Get(fmt.Sprintf(dlURL, *release))
	if err != nil {
		log.Fatalf("Error: could not download tz shapefile: %v\n", err)
	}
	defer resp.Body.Close()

	respBody := progress.NewReader(resp.Body)
	size := resp.ContentLength
	wg := sync.WaitGroup{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		progressChan := progress.NewTicker(ctx, respBody, size, time.Second)
		fmt.Println("Downloading timezone shape file", *release)
		for p := range progressChan {
			fmt.Printf("\r%v  Remaining...", p.Remaining().Round(time.Second))
		}
		fmt.Println("")
	}()

	buffer := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buffer, respBody)
	if err != nil {
		cancel()
		wg.Wait()
		log.Printf("Download failed: %v\n", err)
		return
	}
	wg.Wait()
	cancel()

	bufferReader := bytes.NewReader(buffer.Bytes())
	zipReader, err := zip.NewReader(bufferReader, size)
	if err != nil {
		log.Printf("Could not access zipfile: %v\n", err)
		return
	}
	if len(zipReader.File) == 0 {
		log.Println("Error: release zip file have no files!")
		return
	} else if zipReader.File[0].Name != "combined.json" {
		log.Println("Error: first file in zip file is not combined.json")
		return
	}

	geojsonData, err := zipReader.File[0].Open()
	if err != nil {
		log.Printf("Error: could not read from zip file: %v\n", err)
		return
	}

	currDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error: could not get current dir: %v\n", err)
		return
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Printf("Error: could not create tmp dir: %v\n", err)
		return
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		log.Printf("Error: could not switch to tmp dir: %v\n", err)
		return
	}

	geojsonFile, err := os.Create("./combined.json")
	if err != nil {
		log.Printf("Error: could not create combinedJSON file: %v\n", err)
		return
	}

	_, err = io.Copy(geojsonFile, geojsonData)
	if err != nil {
		geojsonFile.Close()
		log.Printf("Error: could not copy from zip to combined.json: %v\n", err)
		return
	}
	geojsonFile.Close()

	fmt.Println("*** RUNNING MAPSHAPER ***")
	mapshaper := exec.Command("mapshaper", "-i", "combined.json", "-simplify", "visvalingam", "20%", "-o", "reduced.json")
	mapshaper.Stdout = os.Stdout
	mapshaper.Stderr = os.Stderr
	err = mapshaper.Run()
	if err != nil {
		log.Printf("Error: could not run mapshaper: %v\n", err)
		return
	}
	fmt.Println("*** MAPSHAPER FINISHED ***")

	fmt.Println("*** CREATING COMPRESSED SHAPEFILE ***")
	reducedFile, err := os.Open("reduced.json")
	if err != nil {
		log.Printf("Error: could not open file: %v\n", err)
		return
	}
	defer reducedFile.Close()

	buffer = bytes.NewBuffer([]byte{})
	gzipper, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	if err != nil {
		log.Printf("Error: could not create gzip writer: %v\n", err)
		return
	}

	_, err = io.Copy(gzipper, reducedFile)
	if err != nil {
		log.Printf("Error: could not copy data: %v\n", err)
		return
	}
	if err := gzipper.Close(); err != nil {
		log.Printf("Error: could not flush/close gzip: %v\n", err)
		return
	}

	err = os.Chdir(currDir)
	if err != nil {
		log.Printf("Error: could not switch to previous dir: %v", err)
		return
	}

	err = os.WriteFile("shapefile.gz", buffer.Bytes(), 0644)
	if err != nil {
		log.Printf("Error: could not write compressed shapefile: %v", err)
		return
	}

	os.RemoveAll(tmpDir)
	fmt.Println("*** ALL DONE, YAY ***")
}
