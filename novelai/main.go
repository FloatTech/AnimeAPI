package novelai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/FloatTech/floatbox/binary"
)

const (
	loginapi = "https://api.novelai.net/user/login"
	genapi   = "https://api.novelai.net/ai/generate-image"
)

type NovalAI struct {
	Tok  string `json:"accessToken"`
	key  string
	conf *Payload
}

func NewNovalAI(key string, config *Payload) *NovalAI {
	return &NovalAI{
		key:  key,
		conf: config,
	}
}

func (nv *NovalAI) Login() error {
	if nv.Tok != "" {
		return nil
	}
	buf := bytes.NewBuffer([]byte(`{"key": "`))
	buf.WriteString(nv.key)
	buf.WriteString(`"}`)
	resp, err := http.Post(loginapi, "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(nv)
}

func (nv *NovalAI) Draw(tags string) (seed int, tagsproceeded string, img []byte, err error) {
	tags = strings.ReplaceAll(tags, "，", ",")
	if !strings.Contains(tags, ",") {
		tags = strings.ReplaceAll(tags, " ", ",")
	}
	if tags == "" {
		return
	}
	config := *nv.conf
	config.Input = tags
	for config.Parameters.Seed == 0 {
		config.Parameters.Seed = rand.Int()
	}
	seed = config.Parameters.Seed
	buf := bytes.NewBuffer(nil)
	err = config.WrtieTo(buf)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", genapi, buf)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+nv.Tok)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var b [8]byte
	for i := 0; i < 2; i++ {
		for b[0] != '\n' {
			_, err = resp.Body.Read(b[:1])
			if err != nil {
				return
			}
		}
		b[0] = 0
	}
	_, err = resp.Body.Read(b[:5])
	if err != nil {
		return
	}
	img, err = io.ReadAll(base64.NewDecoder(base64.StdEncoding, resp.Body))
	return
}

type Para struct {
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	Scale    int     `json:"scale"`
	Sampler  string  `json:"sampler"`
	Steps    int     `json:"steps"`
	Seed     int     `json:"seed"`
	NSamples int     `json:"n_samples"`
	Strength float64 `json:"strength"`
	Noise    float64 `json:"noise"`
	UcPreset int     `json:"ucPreset"`
	Uc       string  `json:"uc"`
}

type Payload struct {
	Input      string `json:"input"`
	Model      string `json:"model"`
	Parameters *Para  `json:"parameters"`
}

func NewDefaultPayload() *Payload {
	return &Payload{
		Model: "safe-diffusion",
		Parameters: &Para{
			Width:    512,
			Height:   768,
			Scale:    12,
			Sampler:  "k_euler_ancestral",
			Steps:    28,
			NSamples: 1,
			Strength: 0.7,
			Noise:    0.2,
			UcPreset: 0,
			Uc:       "lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry",
		},
	}
}

func (p *Payload) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return binary.BytesToString(b)
}

func (p *Payload) WrtieTo(w io.Writer) error {
	return json.NewEncoder(w).Encode(p)
}
