package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/kbinani/screenshot"

	"github.com/atotto/clipboard"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	slop, err := exec.Command("slop").Output()

	if err != nil {
		panic(err)
	}

	slopOut := string(slop)

	re := regexp.MustCompile(`\d+`)
	slopValues := re.FindAllStringSubmatch(slopOut, -1)

	width, _ := strconv.Atoi(slopValues[0][0])
	height, _ := strconv.Atoi(slopValues[1][0])
	x, _ := strconv.Atoi(slopValues[2][0])
	y, _ := strconv.Atoi(slopValues[3][0])

	rec := image.Rectangle{Max: image.Point{x, y}, Min: image.Point{X: x + width, Y: y + height}}

	img, err := screenshot.CaptureRect(rec)

	if err != nil {
		panic(err)
	}

	imgFile, err := os.Create("screencap.png")

	if err != nil {
		panic(err)
	}

	defer imgFile.Close()

	png.Encode(imgFile, img)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg)

	bucketName := flag.String("b", "my-cheeky-screenshots", "Should be the name of an AWS bucket")
	flag.Parse()

	keyName := fmt.Sprintf("screencap%d_%d-%d.png", width, height, time.Now().Unix())

	freshPtr, err := os.Open(imgFile.Name())
	if err != nil {
		panic(err)
	}

	if _, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(*bucketName),
		Key:    aws.String(keyName),
		Body:   freshPtr,
	}); err != nil {
		panic(err)
	}

	psClient := s3.NewPresignClient(client)

	oneDayFromNow := time.Now().Add(time.Hour * 24)

	resp, err := psClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:                     aws.String(*bucketName),
		Key:                        aws.String(keyName),
		ResponseExpires:            &oneDayFromNow,
		ResponseContentDisposition: aws.String("inline"),
	})

	if err != nil {
		panic(err)
	}

	clipboard.WriteAll(resp.URL)

    exec.Command("notify-send", resp.URL).Run()
}
