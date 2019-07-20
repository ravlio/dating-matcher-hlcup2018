package hl

import (
	"bufio"
	"bytes"
	"github.com/ravlio/highloadcup2018/account"
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/idx"
	"github.com/ravlio/highloadcup2018/indexer"
	"github.com/ravlio/highloadcup2018/metrics"
	"github.com/ravlio/highloadcup2018/querier"
	"github.com/ravlio/highloadcup2018/requests"
	"github.com/ravlio/highloadcup2018/requests/likes"
	"github.com/ravlio/highloadcup2018/skiplist"
	"github.com/ravlio/highloadcup2018/skiplist/int32t"
	"github.com/ravlio/highloadcup2018/skiplist/stringt"
	// "github.com/ravlio/highloadcup2018/gojay"
	"github.com/valyala/fastjson"
	// "github.com/ravlio/highloadcup2018/mmon"
	"github.com/ravlio/highloadcup2018/trie"
	"github.com/rs/zerolog/log"
	// "runtime/debug"
	// "archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	// "fmt"
	"time"
)

type Group string
type Order int8

const (
	StartBirthday = -631152000

	GroupBySex       Group = "sex"
	GroupByStatus    Group = "status"
	GroupByInterests Group = "interests"
	GroupByCountry   Group = "country"
	GroupByCity      Group = "city"

	OrderAsc  = 1
	OrderDesc = -1
)

type opType uint8

const (
	AddAccount opType = iota + 1
	UpdateAccount
	AddLikes
)

type opMsg struct {
	opType opType
	msg    *requests.AccountRequest
	cur    *account.Account
	id     uint32
}

type OptionsIndexer struct {
	Enabled     bool
	BufferSize  int
	WorkerCount int
}

type Options struct {
	Indexer OptionsIndexer
	Debug   bool
}

type dbBitmap struct {
	sex         *idx.HashBitmap
	emailDomain *idx.HashBitmap
	status      *idx.HashBitmap
	statusNeq   *idx.HashBitmap
	fname       *idx.HashBitmap
	fnameY      *idx.Bitmap
	fnameN      *idx.Bitmap
	sname       *idx.HashBitmap
	snameY      *idx.Bitmap
	snameN      *idx.Bitmap
	phoneY      *idx.Bitmap
	phoneN      *idx.Bitmap
	phoneCode   *idx.HashBitmap
	country     *idx.HashBitmap
	countryY    *idx.Bitmap
	countryN    *idx.Bitmap
	city        *idx.HashBitmap
	cityY       *idx.Bitmap
	cityN       *idx.Bitmap
	birthYear   *idx.HashBitmap
	joinedYear  *idx.HashBitmap
	interest    *idx.HashBitmap
	premium     *idx.Bitmap
	premiumY    *idx.Bitmap
	premiumN    *idx.Bitmap
	like        *idx.HashBitmap // по вертикали — кого лайкали, по горизонтали — кто лайкал
}

type dbTrie struct {
	sname *trie.Trie
}

type dbHash struct {
	email *idx.Hash
	phone *idx.Hash
}

type dbRbTree struct {
	premium *idx.RBTree
}

type dbSkiplist struct {
	emailLt *skiplist.Skiplist
	emailGt *skiplist.Skiplist
	birthLt *skiplist.Skiplist
	birthGt *skiplist.Skiplist
}

type DB struct {
	curTime   time.Time
	accChan   chan opMsg
	likesChan chan likes.Likes
	likes     map[uint32][]uint32
	bitmap    dbBitmap
	hash      dbHash
	rbtree    dbRbTree
	trie      dbTrie
	skiplist  dbSkiplist
	querier   *querier.Querier
	indexer   *indexer.Indexer
	accounts  *account.Store
	accMx     sync.RWMutex
	lmx       sync.RWMutex
	opts      *Options
}

func NewDB(opts *Options) *DB {
	r := &DB{
		opts:      opts,
		accChan:   make(chan opMsg, opts.Indexer.BufferSize),
		likesChan: make(chan likes.Likes, opts.Indexer.BufferSize),
		likes:     make(map[uint32][]uint32),
		bitmap: dbBitmap{
			sex:         idx.NewSexHashBitmap(),
			emailDomain: idx.NewUint32HashBitmap(),
			status:      idx.NewStatusHashBitmap(),
			statusNeq:   idx.NewStatusHashBitmap(),
			fname:       idx.NewUint32HashBitmap(),
			fnameY:      idx.NewBitmap(),
			fnameN:      idx.NewBitmap(),
			sname:       idx.NewUint32HashBitmap(),
			snameY:      idx.NewBitmap(),
			snameN:      idx.NewBitmap(),
			phoneY:      idx.NewBitmap(),
			phoneN:      idx.NewBitmap(),
			phoneCode:   idx.NewUint32HashBitmap(),
			country:     idx.NewUint32HashBitmap(),
			countryY:    idx.NewBitmap(),
			countryN:    idx.NewBitmap(),
			city:        idx.NewUint32HashBitmap(),
			cityY:       idx.NewBitmap(),
			cityN:       idx.NewBitmap(),
			birthYear:   idx.NewUint32HashBitmap(),
			joinedYear:  idx.NewUint32HashBitmap(),
			interest:    idx.NewUint32HashBitmap(),
			premium:     idx.NewBitmap(),
			premiumY:    idx.NewBitmap(),
			premiumN:    idx.NewBitmap(),
			like:        idx.NewUint32HashBitmap(),
		},
		hash: dbHash{
			email: idx.NewUint32Hash(),
			phone: idx.NewInt64Hash(),
		},
		trie: dbTrie{
			sname: trie.NewTrie(),
		},
		skiplist: dbSkiplist{
			emailLt: skiplist.NewStringSkiplist(stringt.BuiltinLessThan),
			emailGt: skiplist.NewStringSkiplist(stringt.BuiltinGreaterThan),
			birthLt: skiplist.NewInt32Skiplist(int32t.BuiltinLessThan),
			birthGt: skiplist.NewInt32Skiplist(int32t.BuiltinGreaterThan),
		},
		indexer:  indexer.NewIndexer(),
		accounts: account.NewAccounts(),
	}

	if opts.Debug {
		errors.Debug = true
	} else {
		errors.Debug = false
	}
	return r
}

func (d *DB) Start() error {
	mtr := metrics.NewDurarion("index")
	if d.opts.Indexer.Enabled {
		wg := &sync.WaitGroup{}

		for i := 0; i < d.opts.Indexer.WorkerCount; i++ {
			wg.Add(1)
			go d.runIndexWorker(wg, mtr)
		}

		wg.Wait()
	}
	return nil
}

func (d *DB) Stop() error {
	return nil
}

func (d *DB) LoadFromJSON(path string) error {
	op, err := os.OpenFile(path+"options.txt", os.O_RDONLY, 0666)
	if err != nil {
		return errors.Wrap(err, "cant open options.txt file")
	}

	var t int64
	var p int
	_, err = fmt.Fscanln(op, &t)

	if err != nil {
		return errors.Wrap(err, "cant scan timestamp")
	}

	_, err = fmt.Fscanln(op, &p)
	if err != nil {
		return errors.Wrap(err, "cant scan isRating")
	}

	op.Close()
	d.curTime = time.Unix(t, 0)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.Info().Str("path", path).Msg("Starting importing data")

	start := time.Now()
	files, err := ioutil.ReadDir(path + "data/")
	if err != nil {
		return errors.Wrap(err, "cant read data folder")
	}
	// dr, err := zip.OpenReader(path+"data.zip")

	if err != nil {
		return errors.Wrap(err, "cant read data.zip")
	}

	var accs int

	// go mmon.NewMonitor(1)
	// var trie map[string]*roaring.Bitmap
	// for fn, f := range dr.File {

	for fn, f := range files {

		/*if fn>20 {
			 break
		}*/
		log.Info().Str("file", f.Name()).Int("#", fn).Msg("Reading file ...")

		// fd,err:=ioutil.ReadFile(path+"data/"+f.Name())
		fd, err := os.Open(path + "data/" + f.Name())
		// fd,err:=f.Open()
		if err != nil {
			return errors.Wrap(err, "cant open file")
		}

		// bb:=
		// rdr := bufio.NewReaderSize(bytes.NewReader(fd),100000)
		rdr := bufio.NewReader(fd)
		_, err = rdr.Discard(14) // прыгаем на начало массива
		if err != nil {
			return err
		}

		var buf bytes.Buffer

		var c int
		var el int

		// var bm = make(map[uint32][]uint32)
		var fp fastjson.Parser

		for {
			b, err := rdr.ReadByte()

			if err != nil {
				return err
			}
			buf.WriteByte(b)

			if b == '{' {
				c++
			}
			if b == '}' {
				c--
				if c != 0 {
					continue
				}

				// req, err := jsonp.Parse(buf.String())
				// if err != nil {
				// return err
				// }
				// req:=&requests.AccountRequest{}
				req := requests.AccountRequestPoolGet()
				fv, err := fp.ParseBytes(buf.Bytes())
				if err != nil {
					req.ReleaseToPool()
					return err
				}

				err = req.FromFastJson(fv)
				// err := gojay.UnmarshalJSONObject(buf.Bytes(), req)
				if err != nil {
					req.ReleaseToPool()
					return err
				}

				/*if len(req.Likes)>0 {
					for _,l:=range req.Likes {
						m,ok:=bm[l.ID]
					}
				}*/
				if d.opts.Indexer.Enabled {
					err = d.CreateAccount(req)
					if err != nil {
						req.ReleaseToPool()
						return err
					}
				}
				// req.ReleaseToPool()
				el++
				accs++

				buf.Reset()
				n, err := rdr.ReadByte()
				if err != nil {
					return err
				}

				// скипаем ', '
				if n == ',' {
					_, err := rdr.Discard(1)
					if err != nil {
						return err
					}
				} else if n == ']' { // конец
					break
				}
			}
		}

		fd.Close()

	}
	// runtime.GC()
	// debug.FreeOSMemory()

	d.accounts.SortSlice() // на всякий случай сортируем слайс, так как порядок аккаунтов никто не гарантирует
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	log.Info().TimeDiff("spent", time.Now(), start).Int("Accounts count", accs).Uint64("mb_alloc",
		(m2.Sys-m.Sys)/1024/1024).Msg("Importing data finished")

	d.PrintIndexStats()

	return nil
}
