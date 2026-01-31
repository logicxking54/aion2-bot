package main

import (
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/screenshot"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
)

func FindPosition(displayIndex int, targetImagePath string) (bool, int64, int64, int64, int64) {
	// 1. รับค่าขอบเขตของหน้าจอ
	displayRect := robotgo.GetDisplayRect(displayIndex)

	// 2. Screenshot หน้าจอเต็ม
	img, err := screenshot.CaptureDisplay(displayIndex)
	if err != nil {
		log.Printf("Error capturing display: %v", err)
		return false, 0, 0, 0, 0
	}

	// 3. แปลงเป็น GoCV Mat
	sceneMat, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		log.Printf("Error converting to Mat: %v", err)
		return false, 0, 0, 0, 0
	}
	defer sceneMat.Close()

	// 4. โหลด Template (รูปเล็ก) แบบ GrayScale
	template := gocv.IMRead(targetImagePath, gocv.IMReadGrayScale)
	if template.Empty() {
		log.Printf("Error: Could not load template at %s", targetImagePath)
		return false, 0, 0, 0, 0
	}
	defer template.Close()

	// เก็บขนาดของรูปเล็กไว้ส่งกลับ
	w := int64(template.Cols())
	h := int64(template.Rows())

	// เตรียม Scene เป็น GrayScale สำหรับ SIFT
	sceneGray := gocv.NewMat()
	defer sceneGray.Close()
	gocv.CvtColor(sceneMat, &sceneGray, gocv.ColorRGBAToGray)

	// 5. เริ่มใช้ SIFT
	sift := gocv.NewSIFT()
	defer sift.Close()

	// หา Keypoints และ Descriptors
	_, desc1 := sift.DetectAndCompute(template, gocv.NewMat())
	defer desc1.Close()

	kp2, desc2 := sift.DetectAndCompute(sceneGray, gocv.NewMat())
	defer desc2.Close()

	// 6. Matching จุดเด่น
	matcher := gocv.NewBFMatcherWithParams(gocv.NormL2, false)
	defer matcher.Close()

	// ใช้ KnnMatch เพื่อกรองผลลัพธ์ด้วย Lowe's ratio test
	matches := matcher.KnnMatch(desc1, desc2, 2)

	var goodMatches []gocv.DMatch
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		// กรองจุดที่ใกล้เคียงกันเกินไปเพื่อให้มั่นใจว่าใช่จุดที่ถูกต้อง
		if m[0].Distance < 0.7*m[1].Distance {
			goodMatches = append(goodMatches, m[0])
		}
	}

	// 7. ตรวจสอบว่าเจอจุดมากพอหรือไม่ (ขั้นต่ำ 4 จุด)
	if len(goodMatches) > 4 {
		var sumX, sumY float64
		for _, m := range goodMatches {
			sumX += float64(kp2[m.TrainIdx].X)
			sumY += float64(kp2[m.TrainIdx].Y)
		}

		// คำนวณค่าเฉลี่ยของจุดที่เจอ (พิกัดภายในรูป Mat)
		avgX := int(sumX / float64(len(goodMatches)))
		avgY := int(sumY / float64(len(goodMatches)))

		// วาดจุด Debug ลงในรูปภาพ
		gocv.Circle(&sceneMat, image.Pt(avgX, avgY), 30, color.RGBA{0, 255, 0, 255}, 3)
		gocv.IMWrite("./debug_result.png", sceneMat) // บันทึกรูปเพื่อตรวจความแม่นยำ

		// คำนวณพิกัดจริงบนหน้าจอ (Display X/Y + Found X/Y)
		screenX := int64(displayRect.X + avgX)
		screenY := int64(displayRect.Y + avgY)

		log.Printf("Found target at Screen X: %d, Y: %d", screenX, screenY)
		return true, screenX, screenY, w, h
	}

	return false, 0, 0, 0, 0
}
