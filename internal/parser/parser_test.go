package parser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestNew test the creation of Parser object
// It's not that important, it's generated automatically as part of gotests
func TestNew(t *testing.T) {
	type args struct {
		client transportClient
		l      *log.Logger
		debug  bool
		p      uint
	}
	tests := []struct {
		name string
		args args
		want *parser
	}{
		{
			name: "TestNew - 1",
			args: args{
				client: &http.Client{},
				l:      log.New(os.Stdout, "logger: ", log.Lshortfile),
				debug:  true,
				p:      2,
			},
			want: &parser{
				client:     &http.Client{},
				log:        log.New(os.Stdout, "logger: ", log.Lshortfile),
				debug:      true,
				concurrent: 2,
				cond:       sync.NewCond(&sync.Mutex{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.client, tt.args.l, tt.args.debug, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

type fakeClient struct {
}

func (f *fakeClient) Get(url string) (*http.Response, error) {
	resp := &http.Response{}
	resp.Header = make(map[string][]string)
	resp.Header.Set("Content-Type", "text/html")
	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(HTMLData)))
	resp.StatusCode = http.StatusOK
	return resp, nil

}

type fakeClientStatusNotFound struct {
}

func (f *fakeClientStatusNotFound) Get(url string) (*http.Response, error) {
	resp := &http.Response{}
	resp.Header = make(map[string][]string)
	resp.Header.Set("Content-Type", "text/html")
	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(HTMLData)))
	resp.StatusCode = http.StatusNotFound
	return resp, nil

}

type fakeClientErr struct {
}

func (f *fakeClientErr) Get(url string) (*http.Response, error) {
	return nil, errors.New("Failed to open page")

}

type fakeClientInvalidContentTypeErr struct {
}

func (f *fakeClientInvalidContentTypeErr) Get(url string) (*http.Response,
	error) {
	resp := &http.Response{}
	resp.StatusCode = http.StatusOK
	resp.Header = make(map[string][]string)
	resp.Header.Set("Content-Type", "text/csv")
	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(HTMLData)))
	return resp, nil
}

func Test_parser_ExtractURLs(t *testing.T) {
	type fields struct {
		client transportClient
	}
	type args struct {
		url string
	}
	want := []string{
		"https://en.wikipedia.org/wiki/H._G._Wells",
		"http://gutenberg.net.au/ebooks13/1303101h.html/",
		"http://gutenberg.net/",
		"http://archive.org",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "Test_parser_ExtractURLs 1 - POS",
			fields:  fields{&fakeClient{}},
			args:    args{"test-url"},
			want:    want,
			wantErr: false,
		},
		{
			name:    "Test_parser_ExtractURLs 2 - NEG",
			fields:  fields{&fakeClientErr{}},
			args:    args{"test-url"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test_parser_ExtractURLs 3 - NEG",
			fields:  fields{&fakeClientInvalidContentTypeErr{}},
			args:    args{"test-url"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test_parser_ExtractURLs 4 - NEG",
			fields:  fields{&fakeClientStatusNotFound{}},
			args:    args{"test-url"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				client:     tt.fields.client,
				log:        log.New(ioutil.Discard, "logger: ", log.Lshortfile),
				debug:      true,
				concurrent: 2,
				cond:       sync.NewCond(&sync.Mutex{}),
			}
			got, err := p.ExtractURLs(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.GetURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.GetURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_ExtractURLsParallel(t *testing.T) {
	type fields struct {
		client transportClient
	}
	type args struct {
		url string
	}

	test := struct {
		name   string
		fields fields
		args   args
	}{
		name:   "Concurrent Clients Test",
		fields: fields{&fakeClient{}},
		args:   args{"test-url"},
	}

	t.Run(test.name, func(t *testing.T) {
		p := &parser{
			client:     test.fields.client,
			log:        log.New(ioutil.Discard, "logger: ", log.Lshortfile),
			debug:      true,
			concurrent: 2,
			cond:       sync.NewCond(&sync.Mutex{}),
		}

		for i := 0; i < 100; i++ {
			go func() {
				p.ExtractURLs(test.args.url)
			}()
		}
		p.cond.L.Lock()
		if p.concurrent >= 2 && runtime.NumGoroutine() > 3 {
			t.Errorf("p.concurrent %d, want %d", p.concurrent, 0)
		}
		p.cond.L.Unlock()

		counter := 0

		lockInFor := func() bool {
			p.cond.L.Lock()
			return true
		}

		for lockInFor() && runtime.NumGoroutine() != 1 && p.concurrent != 2 {
			counter++
			time.Sleep(time.Second)
			p.cond.L.Unlock()
		}

		if counter <= 2 {
			t.Errorf("counter %d got %d want", counter, 2)
		}

	})
}

func Test_parser_linksInBody(t *testing.T) {
	type fields struct {
		client transportClient
	}
	type args struct {
		body io.ReadCloser
	}

	want := []string{
		"https://en.wikipedia.org/wiki/H._G._Wells",
		"http://gutenberg.net.au/ebooks13/1303101h.html/",
		"http://gutenberg.net/",
		"http://archive.org",
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "Test_parser_linksInBody - 1",
			fields: fields{},
			args: args{ioutil.NopCloser(
				bytes.NewReader([]byte(HTMLData)))},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				client: tt.fields.client,
			}
			if got := p.linksInBody(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.linskInHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	HTMLData = `
<!doctype html>
<html lang="en-US">
<head>
</head>

<body>
<article>
	<h1>Sleeper</h1>
	<ul>
		<li>Slepper Awakes is a novel by HG Wells <a href="https://en.wikipedia.org/wiki/H._G._Wells">HG Wells</a></li>
		<li>It's free at gutenberg <a href="http://gutenberg.net.au/ebooks13/1303101h.html/">Sleeper Awakes</a> program</li>
	</ul>
	<h3>Web Design</h3>
	<ul>
		<li><a href="http://gutenberg.net/">Gutenberg</a> contains many great novels in public domains.</li>
		<li><a href="http://archive.org">Archive.org</a> is a great place to get other stuff in public domains.</li>
	</ul>
<!--

<li><a href="http://disqus.org">Discuss</a> is a great place to comment.</li>

-->

</article>

</body>
</html>
`
)
