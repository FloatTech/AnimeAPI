// Package setu 统一图床
package setu

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/FloatTech/floatbox/file"
	"github.com/corona10/goimagehash"
	base14 "github.com/fumiama/go-base16384"
	"github.com/sirupsen/logrus"
)

type Pool struct {
	folder  string
	rolimg  func(string) (string, error) // (typ) path
	getdat  func(string) ([]byte, error) // (path) imgbytes
	timeout time.Duration
}

var (
	ErrNilFolder  = errors.New("nil folder")
	ErrNoSuchType = errors.New("no such type")
	ErrEmptyType  = errors.New("empty type")
)

func NewPool(folder string, rolimg func(string) (string, error), getdat func(string) ([]byte, error), timeout time.Duration) (*Pool, error) {
	if folder == "" {
		return nil, ErrNilFolder
	}
	if file.IsNotExist(folder) {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return nil, err
		}
	}
	if folder[len(folder)-1] != '/' {
		folder += "/"
	}
	return &Pool{
		folder:  folder,
		rolimg:  rolimg,
		getdat:  getdat,
		timeout: timeout,
	}, nil
}

func (p *Pool) Roll(typ string) (string, error) {
	d := p.folder + typ
	if p.rolimg == nil {
		return p.rollLocal(d)
	}
	var err error
	ch := make(chan string, 1)
	go func() {
		s := ""
		s, err = p.rolimg(typ)
		ch <- s
		close(ch)
	}()
	select {
	case s := <-ch:
		if err != nil {
			logrus.Warnln("[setu.pool] roll img err:", err)
			return p.rollLocal(d)
		}
		ch := make(chan []byte, 1)
		go func() {
			var data []byte
			data, err = p.getdat(s)
			ch <- data
			close(ch)
		}()
		select {
		case data := <-ch:
			if err != nil {
				logrus.Warnln("[setu.pool] get img err:", err)
				return p.rollLocal(d)
			}
			im, ext, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				logrus.Warnln("[setu.pool] decode img err:", err)
				return p.rollLocal(d)
			}
			dh, err := goimagehash.DifferenceHash(im)
			if err != nil {
				logrus.Warnln("[setu.pool] hash img err:", err)
				return p.rollLocal(d)
			}
			var buf [8]byte
			binary.BigEndian.PutUint64(buf[:], dh.GetHash())
			es := base14.EncodeToString(buf[:])
			if len(es) != 6*3 {
				return p.rollLocal(d)
			}
			es = es[:5*3] + "." + ext
			s = d + "/" + es
			if file.IsExist(s) {
				return s, nil
			}
			return s, os.WriteFile(s, data, 0644)
		case <-time.After(p.timeout):
			return p.rollLocal(d)
		}
	case <-time.After(p.timeout):
		return p.rollLocal(d)
	}
}

func (p *Pool) RollLocal(typ string) (string, error) {
	d := p.folder + typ
	if file.IsNotExist(d) {
		return "", ErrNoSuchType
	}
	return p.rollLocal(d)
}

func (p *Pool) rollLocal(d string) (string, error) {
	files, err := os.ReadDir(d)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", ErrEmptyType
	}
	if len(files) == 1 {
		if files[0].IsDir() {
			return "", ErrEmptyType
		}
		return d + "/" + files[0].Name(), nil
	}
	for c := 0; c < 128; c++ {
		f := files[rand.Intn(len(files))]
		if !f.IsDir() {
			return d + "/" + f.Name(), nil
		}
	}
	return "", ErrEmptyType
}
