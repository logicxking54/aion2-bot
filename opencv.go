package main

import (
	"fmt"
	"gitlab.logicxking.com/core/brain/errmsgs"
	"gocv.io/x/gocv"
	"image"
	"io/fs"
	"path/filepath"
	"strings"
)

var (
	buildScales = []float64{0.50, 0.60, 0.70, 0.80, 0.90, 1.00, 1.15, 1.30, 1.45}
)

func buildMonster(a *App) error {
	err := filepath.WalkDir("./monsters", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		name := strings.ToLower(d.Name())
		if !strings.HasSuffix(name, ".png") {
			return nil
		}

		m := gocv.IMRead(path, gocv.IMReadColor)
		if m.Empty() {
			m.Close()
			return fmt.Errorf("IMRead failed: %s", path)
		}

		for i := range buildScales {
			final := gocv.NewMat()

			err := gocv.Resize(m, &final, image.Point{}, buildScales[i], buildScales[i], gocv.InterpolationLinear)
			if err != nil {
				return a.appCtx.NewError(err, errmsgs.NewError(err.Error()))
			}

			a.monsterMats = append(a.monsterMats, MonsterMat{
				Path: path,
				Mat:  final,
			})
		}

		return nil
	})
	if err != nil {
		return a.appCtx.NewError(err, errmsgs.NewError(err.Error()))
	}

	return nil
}
