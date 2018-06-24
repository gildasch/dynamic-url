package ffmpeg

import (
	"fmt"
	"image"
	"image/jpeg"
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
		fmt.Sprintf(`ffmpeg -i %s 2>&1 | grep "Duration"| cut -d ' ' -f 4 | sed s/,// | sed 's@\..*@@g' | awk '{ split($1, A, ":"); split(A[3], B, "."); print 3600*A[1] + 60*A[2] + B[1] }'`, video))
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
	tmp := "/tmp/" + uuid.NewV4().String() + ".jpg"

	resolutionFlag := ""
	if width != 0 || height != 0 {
		resolutionFlag = fmt.Sprintf("-s %dx%d", width, height)
	}

	_, err := execCommand(
		fmt.Sprintf(`ffmpeg -y -ss %f -i %s -vframes 1 %s %s`, after.Seconds(), video, resolutionFlag, tmp))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp)

	f, err := os.Open(tmp)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	return image, nil
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
