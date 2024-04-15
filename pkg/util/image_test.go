package util_test

import (
	"bytes"
	"github.com/hankmor/gotools/errs"
	"goboot/pkg/util"
	"image/jpeg"
	"os"
	"testing"
)

var pic = "/Users/sam/Pictures/壁纸/唯美/53733594913dc.jpg"
var dpic = "/Users/sam/Pictures/壁纸/唯美/53733594913dc_1.jpg"
var dpic2 = "/Users/sam/Pictures/壁纸/唯美/53733594913dc_2.jpg"

func TestResize(t *testing.T) {
	f, err := os.Open(pic)
	errs.Throw(err)
	img, err := jpeg.Decode(f)
	errs.Throw(err)
	dimg := util.Resize(200, 200, img)

	errs.Throw(err)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, dimg, &jpeg.Options{Quality: 100})
	err = os.WriteFile(dpic, buf.Bytes(), os.ModePerm)
	errs.Throw(err)
}

func TestThumbnail(t *testing.T) {
	f, err := os.Open(pic)
	errs.Throw(err)
	img, err := jpeg.Decode(f)
	errs.Throw(err)
	dimg := util.Thumbnail(400, 400, img)

	errs.Throw(err)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, dimg, &jpeg.Options{Quality: 100})
	err = os.WriteFile(dpic2, buf.Bytes(), os.ModePerm)
	errs.Throw(err)
}
