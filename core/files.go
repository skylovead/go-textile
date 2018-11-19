package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/textileio/textile-go/schema"

	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	m "github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"
)

var ErrFileNotFound = errors.New("file not found")

type Keys map[string]string

type Directory map[string]repo.File

type FileInfo struct {
	Path  string     `json:"path"`
	File  *repo.File `json:"file,omitempty"`
	Links Directory  `json:"links,omitempty"`
}

type FilesInfo struct {
	Id       string        `json:"id"`
	Date     time.Time     `json:"date"`
	AuthorId string        `json:"author_id"`
	Username string        `json:"username,omitempty"`
	Caption  string        `json:"caption,omitempty"`
	Files    []FileInfo    `json:"files"`
	Comments []CommentInfo `json:"comments"`
	Likes    []LikeInfo    `json:"likes"`
}

type CommentInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
	Body     string    `json:"body"`
}

type LikeInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
}

const FileLinkName = "f"
const DataLinkName = "d"

func (t *Textile) GetMedia(reader io.Reader, mill m.Mill) (string, error) {
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	media := http.DetectContentType(buffer[:n])

	return media, mill.AcceptMedia(media)
}

func (t *Textile) AddSchema(jsonstr string, name string) (*repo.File, error) {
	var node schema.Node
	if err := json.Unmarshal([]byte(jsonstr), &node); err != nil {
		return nil, err
	}
	data, err := json.Marshal(&node)
	if err != nil {
		return nil, err
	}

	return t.AddFile(&m.Schema{}, AddFileConfig{
		Input: data,
		Media: "application/json",
		Name:  name,
	})
}

type AddFileConfig struct {
	Input []byte `json:"input"`
	Use   string `json:"use"`
	Media string `json:"media"`
	Name  string `json:"name"`
}

func (t *Textile) AddFile(mill m.Mill, conf AddFileConfig) (*repo.File, error) {
	var source string
	if conf.Use != "" {
		source = conf.Use
	} else {
		source = t.checksum(conf.Input)
	}

	opts, err := mill.Options()
	if err != nil {
		return nil, err
	}

	if efile := t.datastore.Files().GetBySource(mill.ID(), source, opts); efile != nil {
		return efile, nil
	}

	res, err := mill.Mill(conf.Input, conf.Name)
	if err != nil {
		return nil, err
	}

	check := t.checksum(res.File)
	if efile := t.datastore.Files().GetByPrimary(mill.ID(), check); efile != nil {
		return efile, nil
	}

	model := &repo.File{
		Mill:     mill.ID(),
		Checksum: check,
		Source:   source,
		Opts:     opts,
		Media:    conf.Media,
		Name:     conf.Name,
		Size:     len(res.File),
		Added:    time.Now(),
		Meta:     res.Meta,
	}

	var reader *bytes.Reader
	if mill.Encrypt() {
		key, err := crypto.GenerateAESKey()
		if err != nil {
			return nil, err
		}
		ciphertext, err := crypto.EncryptAES(res.File, key)
		if err != nil {
			return nil, err
		}
		model.Key = base58.FastBase58Encoding(key)
		reader = bytes.NewReader(ciphertext)
	} else {
		reader = bytes.NewReader(res.File)
	}

	hash, err := ipfs.AddData(t.node, reader, mill.Pin())
	if err != nil {
		return nil, err
	}
	model.Hash = hash.Hash().B58String()

	if err := t.datastore.Files().Add(model); err != nil {
		return nil, err
	}

	// return the model fetched from the datastore to ensure
	// consistent date formatting and therefore consistent
	// directory hashes
	return t.datastore.Files().Get(model.Hash), nil
}

func (t *Textile) AddNodeFromFiles(files []repo.File) (ipld.Node, Keys, error) {
	keys := make(Keys)
	outer := uio.NewDirectory(t.node.DAG)

	for i, file := range files {
		link := strconv.Itoa(i)
		if err := t.fileNode(file, outer, link); err != nil {
			return nil, nil, err
		}
		keys["/"+link+"/"] = file.Key
	}

	node, err := outer.GetNode()
	if err != nil {
		return nil, nil, err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return nil, nil, err
	}
	return node, keys, nil
}

func (t *Textile) AddNodeFromDirs(dirs []Directory) (ipld.Node, Keys, error) {
	keys := make(Keys)
	outer := uio.NewDirectory(t.node.DAG)

	for i, dir := range dirs {
		inner := uio.NewDirectory(t.node.DAG)
		olink := strconv.Itoa(i)

		for link, file := range dir {
			if err := t.fileNode(file, inner, link); err != nil {
				return nil, nil, err
			}
			keys["/"+olink+"/"+link+"/"] = file.Key
		}

		node, err := inner.GetNode()
		if err != nil {
			return nil, nil, err
		}
		if err := ipfs.PinNode(t.node, node, false); err != nil {
			return nil, nil, err
		}

		id := node.Cid().Hash().B58String()
		if err := ipfs.AddLinkToDirectory(t.node, outer, olink, id); err != nil {
			return nil, nil, err
		}
	}

	node, err := outer.GetNode()
	if err != nil {
		return nil, nil, err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return nil, nil, err
	}
	return node, keys, nil
}

func (t *Textile) Files(threadId string, offset string, limit int) ([]FilesInfo, error) {
	list := make([]FilesInfo, 0)

	query := fmt.Sprintf("threadId='%s' and type=%d", threadId, repo.FilesBlock)
	blocks := t.Blocks(offset, limit, query)
	if len(blocks) == 0 {
		return list, nil
	}

	for _, block := range blocks {
		files, err := t.fileAtBlock(block)
		if err != nil {
			return nil, err
		}

		comments, err := t.fileComments(threadId, block.Target)
		if err != nil {
			return nil, err
		}

		likes, err := t.fileLikes(threadId, block.Target)
		if err != nil {
			return nil, err
		}

		list = append(list, FilesInfo{
			Id:       block.Id,
			Date:     block.Date,
			AuthorId: block.AuthorId,
			Username: t.ContactUsername(block.AuthorId),
			Caption:  block.Body,
			Files:    files,
			Comments: comments,
			Likes:    likes,
		})
	}

	return list, nil
}

func (t *Textile) File(threadId string, blockId string) (*FilesInfo, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	files, err := t.fileAtBlock(*block)
	if err != nil {
		return nil, err
	}

	comments, err := t.fileComments(threadId, block.Target)
	if err != nil {
		return nil, err
	}

	likes, err := t.fileLikes(threadId, block.Target)
	if err != nil {
		return nil, err
	}

	return &FilesInfo{
		Id:       block.Id,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: t.ContactUsername(block.AuthorId),
		Caption:  block.Body,
		Files:    files,
		Comments: comments,
		Likes:    likes,
	}, nil
}

func (t *Textile) FilePlaintext(fileId string) (io.ReadSeeker, *repo.File, error) {
	file := t.datastore.Files().Get(fileId)
	if file == nil {
		return nil, nil, errors.New("file not found")
	}
	ciphertext, err := ipfs.DataAtPath(t.node, file.Hash)
	if err != nil {
		return nil, nil, err
	}
	key, err := base58.Decode(file.Key)
	if err != nil {
		return nil, nil, err
	}
	plaintext, err := crypto.DecryptAES(ciphertext, key)
	if err != nil {
		return nil, nil, err
	}
	return bytes.NewReader(plaintext), file, nil
}

func (t *Textile) fileNode(file repo.File, dir uio.Directory, link string) error {
	if t.datastore.Files().Get(file.Hash) == nil {
		return ErrFileNotFound
	}

	// include encrypted file as well
	plaintext, err := json.Marshal(&file)
	if err != nil {
		return err
	}
	key, err := base58.Decode(file.Key)
	if err != nil {
		return err
	}
	ciphertext, err := crypto.EncryptAES(plaintext, key)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(ciphertext)

	pair := uio.NewDirectory(t.node.DAG)
	if _, err := ipfs.AddDataToDirectory(t.node, pair, FileLinkName, reader); err != nil {
		return err
	}
	if err := ipfs.AddLinkToDirectory(t.node, pair, DataLinkName, file.Hash); err != nil {
		return err
	}

	node, err := pair.GetNode()
	if err != nil {
		return err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return err
	}
	return ipfs.AddLinkToDirectory(t.node, dir, link, node.Cid().Hash().B58String())
}

func (t *Textile) fileAtBlock(block repo.Block) ([]FileInfo, error) {
	if block.Type != repo.FilesBlock {
		return nil, ErrBlockNotFile
	}

	links, err := ipfs.LinksAtPath(t.node, block.Target)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, len(links))

	for _, index := range links {
		node, err := ipfs.NodeAtLink(t.node, index)
		if err != nil {
			return nil, err
		}
		fnames := node.Links()

		info := FileInfo{
			Path: block.Target + "/" + index.Name,
		}
		if len(fnames) > 0 {
			// directory of files
			info.Links = make(Directory)
			for _, link := range node.Links() {
				pair, err := ipfs.NodeAtLink(t.node, link)
				if err != nil {
					return nil, err
				}
				file, err := t.fileForPair(pair)
				if err != nil {
					return nil, err
				}
				if file != nil {
					info.Links[link.Name] = *file
				}
			}
		} else {
			// single file
			file, err := t.fileForPair(node)
			if err != nil {
				return nil, err
			}
			info.File = file
		}

		i, err := strconv.Atoi(index.Name)
		if err != nil {
			return nil, err
		}
		files[i] = info
	}

	return files, nil
}

func (t *Textile) checksum(plaintext []byte) string {
	sum := sha256.Sum256(plaintext)
	return base58.FastBase58Encoding(sum[:])
}

func (t *Textile) fileForPair(pair ipld.Node) (*repo.File, error) {
	d, _, err := pair.ResolveLink([]string{DataLinkName})
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	return t.datastore.Files().Get(d.Cid.Hash().B58String()), nil
}

func (t *Textile) fileComments(threadId string, target string) ([]CommentInfo, error) {
	comments := make([]CommentInfo, 0)

	query := fmt.Sprintf("threadId='%s' and type=%d and target='%s'", threadId, repo.CommentBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info := CommentInfo{
			Id:       block.Id,
			Date:     block.Date,
			AuthorId: block.AuthorId,
			Username: t.ContactUsername(block.AuthorId),
			Body:     block.Body,
		}
		comments = append(comments, info)
	}

	return comments, nil
}

func (t *Textile) fileLikes(threadId string, target string) ([]LikeInfo, error) {
	likes := make([]LikeInfo, 0)

	query := fmt.Sprintf("threadId='%s' and type=%d and target='%s'", threadId, repo.LikeBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info := LikeInfo{
			Id:       block.Id,
			Date:     block.Date,
			AuthorId: block.AuthorId,
			Username: t.ContactUsername(block.AuthorId),
		}
		likes = append(likes, info)
	}

	return likes, nil
}
