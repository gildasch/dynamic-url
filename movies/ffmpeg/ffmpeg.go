package ffmpeg

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var Debug = false

func Duration(video string) (time.Duration, error) {
	// from https://superuser.com/questions/650291/how-to-get-video-duration-in-seconds
	out, err := execCommand(
		fmt.Sprintf(`ffmpeg -i '%s' 2>&1 | awk -F[:.] '/Duration:/ {print $2*3600+$3*60+$4}'`, video))
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(out)
	if err != nil {
		return 0, err
	}

	return time.Duration(i) * time.Second, nil
}

func Capture(video string, after time.Duration, width, height int) (image.Image, error) {
	is, err := Captures(video, after, width, height, 1)
	if err != nil || len(is) < 1 {
		return nil, err
	}
	return is[0], nil
}

func Captures(video string, after time.Duration, width, height, n int) ([]image.Image, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	tmp := "/tmp/" + uuid.String() + "-%04d.jpg"
	fmt.Println("saving", n, tmp)
	defer func() {
		for i := 0; i < n; i++ {
			os.Remove(fmt.Sprintf(tmp, i))
		}
	}()

	resolutionFlag := ""
	if width != 0 || height != 0 {
		resolutionFlag = fmt.Sprintf("-s %dx%d", width, height)
	}

	_, err = execCommand(
		fmt.Sprintf(`ffmpeg -y -ss %f -i '%s' -vframes %d -r 5 %s %s`, after.Seconds(), video, n, resolutionFlag, tmp))
	if err != nil {
		return nil, err
	}

	var images []image.Image
	for i := 1; i <= n; i++ {
		f, err := os.Open(fmt.Sprintf(tmp, i))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		image, err := jpeg.Decode(f)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func GIFCaptures(video string, after time.Duration, width, height, n, framesPerSecond int) ([]*image.Paletted, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	tmp := "/tmp/" + uuid.String() + ".gif"
	fmt.Println("saving", n, tmp)
	defer os.Remove(tmp)

	resolutionFlag := ""
	if width != 0 || height != 0 {
		resolutionFlag = fmt.Sprintf("-s %dx%d", width, height)
	}

	_, err = execCommand(
		fmt.Sprintf(`ffmpeg -y -ss %f -i '%s' -vframes %d -r %d %s %s`,
			after.Seconds(), video, n, framesPerSecond, resolutionFlag, tmp))
	if err != nil {
		return nil, err
	}

	f, err := os.Open(tmp)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, err
	}

	return g.Image, nil
}

// Generates webm sequence. Note ffmpeg must be built with libvpx to encode webm.
func WebM(video string, after time.Duration, width, height, n, framesPerSecond int) ([]byte, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	tmp := "/tmp/" + uuid.String() + ".webm"
	defer os.Remove(tmp)

	_, err = execCommand(
		// -an: Remove audio
		//   This allows autoplay on modern browsers, where
		//   autoplay is disabled due to abuse by ads.
		// -vf 'subtitles=%s': Burn in subtitles (TODO)
		//   Problem, does not take time offset into account.
		//   Possible fix is gen time shifted subs beforehand.
		//   Documentation:
		//   http://ffmpeg.org/ffmpeg-filters.html#subtitles-1
		fmt.Sprintf(`ffmpeg -y -ss %f -i '%s' -an -vframes %d -r %d -s %dx%d %s`,
			after.Seconds(), video, n, framesPerSecond, width, height, tmp))

	f, err := os.Open(tmp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	return ioutil.ReadAll(reader)
}

func execCommand(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)

	if Debug {
		fmt.Println("Executing:", cmd)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return string(out), err
	}
	if Debug {
		fmt.Printf("Output: %s\n", string(out))
	}

	return strings.TrimSuffix(string(out), "\n"), err
}
