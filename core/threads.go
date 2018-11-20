package core

import (
	"errors"
	"fmt"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

// ErrThreadNotFound indicates thread is not found in the loaded list
var ErrThreadNotFound = errors.New("thread not found")

// ErrThreadLoaded indicates the thread is already loaded from the datastore
var ErrThreadLoaded = errors.New("thread is loaded")

// AddThreadConfig is used to create a new thread model
type AddThreadConfig struct {
	Key       string          `json:"key"`
	Name      string          `json:"name"`
	Schema    mh.Multihash    `json:"schema"`
	Initiator string          `json:"initiator"`
	Type      repo.ThreadType `json:"type"`
	Join      bool            `json:"join"`
}

// AddThread adds a thread with a given name and secret key
func (t *Textile) AddThread(sk libp2pc.PrivKey, conf AddThreadConfig) (*Thread, error) {
	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	skb, err := sk.Bytes()
	if err != nil {
		return nil, err
	}

	var schema string
	if conf.Schema != nil {
		schema = conf.Schema.B58String()
		t.cafeOutbox.Add(schema, repo.CafeStoreRequest)
	}

	threadModel := &repo.Thread{
		Id:        id.Pretty(),
		Key:       conf.Key,
		PrivKey:   skb,
		Name:      conf.Name,
		Schema:    conf.Schema.B58String(),
		Initiator: conf.Initiator,
		Type:      conf.Type,
		State:     repo.ThreadLoaded,
	}
	if err := t.datastore.Threads().Add(threadModel); err != nil {
		return nil, err
	}

	thrd, err := t.loadThread(threadModel)
	if err != nil {
		return nil, err
	}

	// we join here if we're the creator
	if conf.Join {
		if _, err := thrd.joinInitial(); err != nil {
			return nil, err
		}
	}

	if thrd.Schema != nil {
		go t.cafeOutbox.Flush()
	}

	t.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadAdded})

	log.Debugf("added a new thread %s with name %s", thrd.Id, conf.Name)

	return thrd, nil
}

// RemoveThread removes a thread
func (t *Textile) RemoveThread(id string) (mh.Multihash, error) {
	var thrd *Thread
	var index int
	for i, th := range t.threads {
		if th.Id == id {
			thrd = th
			index = i
			break
		}
	}
	if thrd == nil {
		return nil, errors.New("thread not found")
	}

	// notify peers
	addr, err := thrd.leave()
	if err != nil {
		return nil, err
	}

	if err := t.datastore.Threads().Delete(thrd.Id); err != nil {
		return nil, err
	}

	copy(t.threads[index:], t.threads[index+1:])
	t.threads[len(t.threads)-1] = nil
	t.threads = t.threads[:len(t.threads)-1]

	t.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadRemoved})

	log.Infof("removed thread %s with name %s", thrd.Id, thrd.Name)

	return addr, nil
}

// Threads lists loaded threads
func (t *Textile) Threads() []Thread {
	var threads []Thread
	for _, t := range t.threads {
		if t != nil {
			threads = append(threads, *t)
		}
	}
	return threads
}

// Thread get a thread by id from loaded threads
func (t *Textile) Thread(id string) *Thread {
	for _, thrd := range t.threads {
		if thrd.Id == id {
			return thrd
		}
	}
	return nil
}

// ThreadByKey get a thread by key from loaded threads
func (t *Textile) ThreadByKey(key string) *Thread {
	for _, thrd := range t.threads {
		if thrd.Key == key {
			return thrd
		}
	}
	return nil
}

// ThreadInfo gets thread info
func (t *Textile) ThreadInfo(id string) (*ThreadInfo, error) {
	thrd := t.Thread(id)
	if thrd == nil {
		return nil, errors.New(fmt.Sprintf("cound not find thread: %s", id))
	}
	return thrd.Info()
}

// AcceptThreadInvite attemps to download an encrypted thread key from an internal invite,
// add the thread, and notify the inviter of the join
func (t *Textile) AcceptThreadInvite(inviteId string) (mh.Multihash, error) {
	ciphertext, err := ipfs.DataAtPath(t.node, inviteId)
	if err != nil {
		return nil, err
	}
	if err := ipfs.UnpinPath(t.node, inviteId); err != nil {
		log.Warningf("error unpinning path %s: %s", inviteId, err)
	}

	// attempt decrypt w/ own keys
	plaintext, err := crypto.Decrypt(t.node.PrivateKey, ciphertext)
	if err != nil {
		return nil, ErrInvalidThreadBlock
	}
	hash, err := t.handleThreadInvite(plaintext)
	if err != nil {
		return nil, err
	}

	if err := t.datastore.Notifications().DeleteByBlock(inviteId); err != nil {
		return nil, err
	}

	return hash, nil
}

// AcceptExternalThreadInvite attemps to download an encrypted thread key from an external invite,
// add the thread, and notify the inviter of the join
func (t *Textile) AcceptExternalThreadInvite(inviteId string, key []byte) (mh.Multihash, error) {
	ciphertext, err := ipfs.DataAtPath(t.node, fmt.Sprintf("%s", inviteId))
	if err != nil {
		return nil, err
	}

	// attempt decrypt w/ key
	plaintext, err := crypto.DecryptAES(ciphertext, key)
	if err != nil {
		return nil, ErrInvalidThreadBlock
	}
	return t.handleThreadInvite(plaintext)
}

// IgnoreThreadInvite unpins a direct peer-to-peer invite and removes
// the associated notification.
func (t *Textile) IgnoreThreadInvite(inviteId string) error {
	if err := ipfs.UnpinPath(t.node, inviteId); err != nil {
		log.Warningf("error unpinning path %s: %s", inviteId, err)
	}
	return t.datastore.Notifications().DeleteByBlock(inviteId)
}

// handleThreadInvite
func (t *Textile) handleThreadInvite(plaintext []byte) (mh.Multihash, error) {
	block := new(pb.ThreadBlock)
	if err := proto.Unmarshal(plaintext, block); err != nil {
		return nil, err
	}
	if block.Type != pb.ThreadBlock_INVITE {
		return nil, ErrInvalidThreadBlock
	}
	msg := new(pb.ThreadInvite)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	sk, err := libp2pc.UnmarshalPrivateKey(msg.Sk)
	if err != nil {
		return nil, err
	}

	id, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	if thrd := t.Thread(id.Pretty()); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	var sch mh.Multihash
	if msg.Schema != "" {
		sch, err = mh.FromB58String(msg.Schema)
		if err != nil {
			return nil, err
		}
	}
	config := AddThreadConfig{
		Key:       ksuid.New().String(),
		Name:      msg.Name,
		Schema:    sch,
		Initiator: msg.Initiator,
		Type:      repo.OpenThread,
		Join:      false,
	}
	thrd, err := t.AddThread(sk, config)
	if err != nil {
		return nil, err
	}

	// follow parents, update head
	if err := thrd.handleInviteMessage(block); err != nil {
		return nil, err
	}

	// mark any discovered peers as welcomed
	// there's no need to send a welcome because we're about to send a join message
	if err := t.datastore.ThreadPeers().WelcomeByThread(thrd.Id); err != nil {
		return nil, err
	}

	// join the thread
	author, err := peer.IDB58Decode(block.Header.Author)
	if err != nil {
		return nil, err
	}
	hash, err := thrd.join(author)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// addAccountThread adds a thread with seed representing the state of the account
func (t *Textile) addAccountThread() error {
	if t.ThreadByKey(t.config.Account.Address) != nil {
		return nil
	}
	sk, err := t.account.LibP2PPrivKey()
	if err != nil {
		return err
	}

	config := AddThreadConfig{
		Key:       t.account.Address(),
		Name:      "account",
		Initiator: t.account.Address(),
		Type:      repo.PrivateThread,
		Join:      true,
	}
	if _, err := t.AddThread(sk, config); err != nil {
		return err
	}
	return nil
}
