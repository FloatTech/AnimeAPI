package scale

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/FloatTech/zbputils/web"
)

const (
	ModelConservative = "conservative"
	ModelNoDenoise    = "no-denoise"
	ModelDenoise1x    = "denoise1x"
	ModelDenoise2x    = "denoise2x"
	ModelDenoise3x    = "denoise3x"
)

var (
	Models      = [...]string{ModelConservative, ModelNoDenoise, ModelDenoise1x, ModelDenoise2x, ModelDenoise3x}
	ErrInvModel = errors.New("invaild model")
	ErrInvScale = errors.New("invaild scale")
	ErrInvTile  = errors.New("invaild tile")
)

// Get model 0-4, scale 2-4, tile 0-4
func Get(u string, model, scale, tile int) ([]byte, error) {
	if model < 0 || model > 4 || ((scale == 3 || scale == 4) && (model == 2 || model == 3)) {
		return nil, ErrInvModel
	}
	if scale > 4 || scale < 2 {
		return nil, ErrInvScale
	}
	if tile < 0 || tile > 4 {
		return nil, ErrInvTile
	}
	return web.GetData(fmt.Sprintf("https://bilibiliai.azurewebsites.net/api/scale?url=%s&model=%s&scale=%d&tile=%d", url.QueryEscape(u), Models[model], scale, tile))
}

// Post model 0-4, scale 2-4, tile 0-4
func Post(body io.Reader, model, scale, tile int) ([]byte, error) {
	if model < 0 || model > 4 {
		return nil, ErrInvModel
	}
	if scale > 4 || scale < 2 {
		return nil, ErrInvScale
	}
	if tile < 0 || tile > 4 {
		return nil, ErrInvTile
	}
	return web.PostData(
		fmt.Sprintf("https://bilibiliai.azurewebsites.net/api/scale?model=%s&scale=%d&tile=%d", Models[model], scale, tile),
		"application/octet-stream", body,
	)
}
